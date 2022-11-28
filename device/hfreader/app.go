package hfreader

import (
	"context"
	"fmt"
	"time"

	util "SimulatedDeviceGUI/util"

	"gorm.io/gorm"
)

// App struct
type App struct {
	ctx             context.Context
	db              *gorm.DB
	activeHFReaders []*ActiveHFReader
}

// NewApp creates a new App application struct
func NewReaderApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context, db *gorm.DB) {
	a.ctx = ctx
	a.db = db

	db.AutoMigrate(&HFReader{})

	go func() {
		for {
			fmt.Println(a.activeHFReaders)
			a.SyncReaderWithDB()
			time.Sleep(2 * time.Second)
		}
	}()
}

func (a *App) CreateReader(line string, post string, code string, targetUrl string) {
	if targetUrl == "" {
		targetUrl = "tcp://localhost:1883"
	}
	hfReader := HFReader{Line: line, Post: post, Code: code, TargetUrl: targetUrl}
	if result := a.db.Create(&hfReader); result.Error != nil {
		fmt.Println("Failed to create device")
	}
}

func (a *App) DeleteReader(id uint) {
	if result := a.db.Delete(&HFReader{}, id); result.Error != nil {
		fmt.Println("Failed to delete device")
	}
}

func (a *App) GetAllDeviceFromDB() []HFReader {
	var devices []HFReader
	if result := a.db.Find(&devices); result.Error != nil {
		fmt.Println("Failed to get all device")
	}
	return devices
}

func (a *App) GetAllReader() []HFReader {
	var hfr []HFReader

	for _, activeDevice := range a.activeHFReaders {
		hfr = append(hfr, activeDevice.HFReader)
	}
	return hfr
}

func (a *App) SetReaderCardPresent(id uint, isCardPresent bool) {
	a.db.Model(&HFReader{Id: id}).Update("IsCardPresent", isCardPresent)
}
func (a *App) SetReaderConnection(id uint, isConnected bool) {
	fmt.Println("Connection change")
	// a.db.Model(&HFReader{Id: id}).Update("IsConnected", isConnected)
	for i, activeHFReader := range a.activeHFReaders{
		if(activeHFReader.id == id){
			a.activeHFReaders[i].mu.Lock()
			a.activeHFReaders[i].HFReader.IsConnected = isConnected
			a.activeHFReaders[i].mu.Unlock()
		}
	}
}
func (a *App) SetReaderEpc(id uint, epc string) {
	a.db.Model(&HFReader{Id: id}).Update("UidBuffer", epc)
}

func (a *App) SyncReaderWithDB() {
	devices := a.GetAllDeviceFromDB()
	var unExistIndex []int

	for i, activeDevice := range a.activeHFReaders {
		isExistInDB := false

		for j, device := range devices {
			if activeDevice.id == device.Id {
				a.activeHFReaders[i].mu.Lock()
				a.activeHFReaders[i].HFReader.UidBuffer = device.UidBuffer
				a.activeHFReaders[i].HFReader.IsCardPresent = device.IsCardPresent
				a.activeHFReaders[i].mu.Unlock()
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
		id := a.activeHFReaders[i].id
		fmt.Println("active device with id", id, ", not exist in db, removing...")
		a.activeHFReaders[i].Destroy()
		a.activeHFReaders = util.RemoveFromSlice(a.activeHFReaders, i)
		fmt.Println("active device with id", id, ", has removed")
	}

	for _, device := range devices {
		device.IsConnected = true
		activeHFReader := &ActiveHFReader{id: device.Id, HFReader: device, isActive: true}
		activeHFReader.Setup()
		go activeHFReader.Loop()
		a.activeHFReaders = append(a.activeHFReaders, activeHFReader)
	}
}
