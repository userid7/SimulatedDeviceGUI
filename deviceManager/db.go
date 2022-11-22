package deviceManager

import (
	d "SimulatedDeviceGUI/deviceMqtt"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

var activeDevices []d.Device

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
		syncActiveDevices2()
		time.Sleep(2 * time.Second)
	}
}

func syncActiveDevices2() {
	devices := GetAllDevice()

	for i, activeDevice := range activeDevices {
		isExistInNewData := false
		for j, device := range devices {
			if activeDevice.Id == device.Id {
				activeDevices[i].Mu.Lock()
				activeDevices[i].UidBuffer = device.UidBuffer
				activeDevices[i].IsCardPresent = device.IsCardPresent
				activeDevices[i].IsConnected = true
				activeDevices[i].TargetUrl = device.TargetUrl
				activeDevices[i].Mu.Unlock()

				devices = removeFromSlice(devices, j)
				isExistInNewData = true
				break
			}
		}
		if !isExistInNewData {
			fmt.Println("active device with id", activeDevices[i].Id, ", not exist in db, removing...")
			// TODO : send signal via channel to inactive deviceRunner loop
			activeDevices[i].Mu.Lock()
			activeDevices[i].IsActive = false
			activeDevices[i].Mu.Unlock()
			activeDevices = removeFromSlice(activeDevices, i)
		}
	}

	fmt.Println("devices")
	fmt.Println(devices)
	for _, device := range devices {
		device.Mu = &sync.Mutex{}
		activeDevices = append(activeDevices, device)
	}

	for i, activeDevice := range activeDevices {
		if !activeDevice.IsActive {
			activeDevices[i].Mu.Lock()
			activeDevices[i].IsActive = true
			activeDevices[i].IsConnected = true
			activeDevices[i].Mu.Unlock()
			// go activeDevices[i].Loop()
			go activeDevices[i].Main()
		}
	}
}

func syncActiveDevices(devices []d.Device) {
	var newActiveDevices []d.Device
	// fmt.Println("SyncActiveDevices")

	// newActiveDevices := devices

	// fmt.Println(devices)

	// // sync with data in db
	for _, device := range devices {
		isAlreadyExist := false
		for _, activeDevice := range activeDevices {
			if device.Id == activeDevice.Id {
				device.IsActive = activeDevice.IsActive
				newActiveDevices = append(newActiveDevices, device)
				isAlreadyExist = true
				break
			}
		}
		if !isAlreadyExist {
			newActiveDevices = append(newActiveDevices, device)
		}

	}

	// sync data with memory
	for _, activeDevice := range activeDevices {
		isAlreadyExist := false
		for _, newActiveDevice := range newActiveDevices {
			if newActiveDevice.Id == activeDevice.Id {
				isAlreadyExist = true
			}
		}

		if !isAlreadyExist {
			activeDevice.IsActive = false
		}
	}

	activeDevices = newActiveDevices

	// fmt.Println(activeDevices)

	for i, activeDevice := range activeDevices {
		fmt.Println(&activeDevices[i])
		if !activeDevice.IsActive {
			fmt.Println("Starting new device loop", activeDevice.Id)
			activeDevices[i].IsActive = true
			go activeDevices[i].Loop()
		}
	}
}

func GetAllDevice() []d.Device {
	var devices []d.Device
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
	// var device = d.Device{Line: line, Post: post, Code: code, TargetUrl: targetUrl}
	if result := db.Delete(&d.Device{}, id); result.Error != nil {
		fmt.Println("Failed to create device")
	}
}

func GetActiveDevices() []d.Device {
	fmt.Println("GetActiveDevices")
	fmt.Println(activeDevices)
	return activeDevices
}

func SetDeviceEpc(id uint, epc string) {
	// var device = d.Device{Line: line, Post: post, Code: code, TargetUrl: targetUrl}
	device := d.Device{Id: id}
	if result := db.Model(&device).Update("UidBuffer", epc); result.Error != nil {
		fmt.Println("Failed to set device Epc")
	}
}

func SetDeviceIsCardPresent(id uint, isCardPresent bool) {
	// var device = d.Device{Line: line, Post: post, Code: code, TargetUrl: targetUrl}
	device := d.Device{Id: id}
	if result := db.Model(&device).Update("IsCardPresent", isCardPresent); result.Error != nil {
		fmt.Println("Failed to set device isCardPresent")
	}
}

func SetDeviceIsConnected(id uint, isConnected bool) {
	// var device = d.Device{Line: line, Post: post, Code: code, TargetUrl: targetUrl}
	device := d.Device{Id: id}
	if result := db.Model(&device).Update("IsConnected", isConnected); result.Error != nil {
		fmt.Println("Failed to set device isConnected")
	}
}

func removeFromSlice[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}
