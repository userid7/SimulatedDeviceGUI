package pm

import (
	"context"
	"fmt"
	"time"

	util "SimulatedDeviceGUI/util"

	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
)

// App struct
type App struct {
	ctx              context.Context
	db               *gorm.DB
	activePMGateways []*ActivePMGateway
	validate         *validator.Validate
}

func NewPMApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context, db *gorm.DB) {
	a.ctx = ctx
	a.db = db
	a.validate = validator.New()

	db.AutoMigrate(&PMGateway{})
	db.AutoMigrate(&PM{})

	go func() {
		for {
			fmt.Println(a.activePMGateways)
			a.SyncPMGatewayWithDB()
			time.Sleep(2 * time.Second)
		}
	}()
}

func (a *App) CreatePMGateway(line string, code string, targetUrl string) {
	if targetUrl == "" {
		targetUrl = "tcp://localhost:1883"
	}

	pmGateway := PMGateway{Line: line, Code: code, TargetUrl: targetUrl}

	if err := a.validate.Struct(pmGateway); err != nil {
		fmt.Println(err)
		return
	}

	if result := a.db.Create(&pmGateway); result.Error != nil {
		fmt.Println("Failed to create device")
	}
}

func (a *App) CreatePM(post string, code string, pmGatewayId uint) {
	pm := PM{Post: post, Code: code, PMGatewayId: pmGatewayId}
	
	if err := a.validate.Struct(pm); err != nil {
		fmt.Println(err)
		return
	}

	if result := a.db.Create(&pm); result.Error != nil {
		fmt.Println("Failed to create device")
	}
}

func (a *App) DeletePMGateway(id uint) {
	if result := a.db.Delete(&PMGateway{}, id); result.Error != nil {
		fmt.Println("Failed to delete PMGateway with id", id)
		return
	}
	if result := a.db.Where("PMGatewayId = ?", id).Delete(&PM{}); result.Error != nil {
		fmt.Println("Failed to delete PMs of PMGateway with id", id)
		return
	}
}

func (a *App) DeletePM(id uint) {
	if result := a.db.Delete(&PM{}, id); result.Error != nil {
		fmt.Println("Failed to delete PM with id", id)
		return
	}
}

func (a *App) GetAllPMGatewayFromDB() []PMGateway {
	var pmGateway []PMGateway
	if result := a.db.Preload("PMs").Find(&pmGateway); result.Error != nil {
		fmt.Println("Failed to get all PMGateway")
	}
	return pmGateway
}

func (a *App) GetAllActivePMGateway() []PMGateway {
	var pmGateway []PMGateway

	for _, activeDevice := range a.activePMGateways {
		pmGateway = append(pmGateway, activeDevice.PMGateway)
	}
	return pmGateway
}

func (a *App) SetPMGatewayConnection(id uint, isConnected bool) {
	for i, activeHFReader := range a.activePMGateways {
		if activeHFReader.id == id {
			a.activePMGateways[i].mu.Lock()
			a.activePMGateways[i].PMGateway.IsConnected = isConnected
			a.activePMGateways[i].mu.Unlock()
		}
	}
}

func (a *App) SetPM(id uint, pm PM) {
	fmt.Println("Set PM with id", pm.Id)
	fmt.Println(pm)

	selectedPM := &PM{Id: id}

	// TODO : Updates not receive zero value (false also not detect)
	if result := a.db.Model(&selectedPM).Updates(pm); result.Error != nil {
		fmt.Println("Failed to update PM with id", id)
		return
	}
}

func (a *App) SetPMIsOk(id uint, isOk bool) {
	fmt.Println("Set PM IsOk with id", id)

	selectedPM := &PM{Id: id}

	if result := a.db.Model(&selectedPM).Update("IsOk", isOk); result.Error != nil {
		fmt.Println("Failed to update PM with id", id)
		return
	}
}

func (a *App) SetPMKw(id uint, kw float32) {
	fmt.Println("Set PM kw with id", id)

	selectedPM := &PM{Id: id}

	if result := a.db.Model(&selectedPM).Update("Kw", kw); result.Error != nil {
		fmt.Println("Failed to update PM with id", id)
		return
	}
}

func (a *App) SetPMKwh(id uint, kw float32) {
	fmt.Println("Set PM kwh with id", id)

	selectedPM := &PM{Id: id}

	if result := a.db.Model(&selectedPM).Update("Kwh", kw); result.Error != nil {
		fmt.Println("Failed to update PM with id", id)
		return
	}
}

func (a *App) SetPMIsRandom(id uint, isRandom bool) {
	fmt.Println("Set PM IsRandom with id", id)

	selectedPM := &PM{Id: id}

	if result := a.db.Model(&selectedPM).Update("IsRandom", isRandom); result.Error != nil {
		fmt.Println("Failed to delete PM with id", id)
		return
	}
}

func (a *App) SyncPMGatewayWithDB() {
	devices := a.GetAllPMGatewayFromDB()
	var unExistIndex []int

	for i, activeDevice := range a.activePMGateways {
		isExistInDB := false

		for j, device := range devices {
			if activeDevice.id == device.Id {
				a.activePMGateways[i].mu.Lock()
				a.activePMGateways[i].PMGateway.PMs = device.PMs
				a.activePMGateways[i].mu.Unlock()

				devices = util.RemoveFromSlice(devices, j)
				isExistInDB = true
				break
			}
		}
		if !isExistInDB {
			unExistIndex = append(unExistIndex, i)
		}
	}

	unExistIndex = util.Reverse(unExistIndex)

	for _, i := range unExistIndex {
		id := a.activePMGateways[i].id
		fmt.Println("active device with id", id, ", not exist in db, removing...")
		a.activePMGateways[i].Destroy()
		a.activePMGateways = util.RemoveFromSlice(a.activePMGateways, i)
		fmt.Println("active device with id", id, ", has removed")
	}

	for _, device := range devices {
		device.IsConnected = true
		activePMGateway := &ActivePMGateway{id: device.Id, PMGateway: device, isActive: true, db: a.db}
		activePMGateway.Setup()
		go activePMGateway.Loop()
		a.activePMGateways = append(a.activePMGateways, activePMGateway)
	}
}
