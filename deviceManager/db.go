package deviceManager

import (
	d "SimulatedDeviceGUI/deviceMqtt"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

var activeDevices []*Device

type DeviceMethod interface {
	Main()
	Loop()
	Destroy()
}

type Device struct {
	Id       int
	Type     string
	IsActive bool
	DeviceMethod
}

func Main() {
	var err error
	db, err = gorm.Open(sqlite.Open("deviceManager.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&d.Device{})
}

func Loop() {
	for {
		syncActiveDevices3()
		time.Sleep(2 * time.Second)
	}
}

func syncActiveDevices3() {
	devices := GetAllDeviceFromDB()

	for i, activeDevice := range activeDevices {
		isExistInNewData := false

		for j, device := range devices {
			if activeDevice.Id == device.Id {
				devices = removeFromSlice(devices, j)
				isExistInNewData = true
				break
			}
		}
		if !isExistInNewData {
			fmt.Println("active device with id", activeDevices[i].Id, ", not exist in db, removing...")
			activeDevices[i].Destroy()
			activeDevices = removeFromSlice(activeDevices, i)
		}
	}

	for _, device := range devices {
		// device.Mu = &sync.Mutex{}
		activeDevices = append(activeDevices, device)
	}

	for i, activeDevice := range activeDevices {
		if !activeDevice.IsActive {
			activeDevices[i].Main()
			go activeDevices[i].Loop()
		}
	}

}

// Device Factory
func GetAllDeviceFromDB() []*Device {
	var devices []*Device
	if result := db.Find(&devices); result.Error != nil {
		fmt.Println("Failed to get all device")
	}
	return devices
}

func CreateDevice(line string, post string, code string, targetUrl string) {
	var device = d.Device{Line: line, Post: post, Code: code, TargetUrl: targetUrl}
	if result := db.Create(&device); result.Error != nil {
		fmt.Println("Failed to create device")
	}
}

func DeleteDevice(id uint) {
	if result := db.Delete(&d.Device{}, id); result.Error != nil {
		fmt.Println("Failed to delete device")
	}
}

func GetActiveDevices() []Device {
	var d []Device

	for _, activeDevice := range activeDevices {
		d = append(d, *activeDevice)
	}

	fmt.Println("GetActiveDevices")
	fmt.Println(activeDevices)
	fmt.Println(d)
	return d
}

func removeFromSlice[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}
