package main

import (
	dm "SimulatedDeviceGUI/deviceManager"
	d "SimulatedDeviceGUI/deviceMqtt"
	"context"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	dm.Main()
	go dm.Loop()
}

func (a *App) CreateReader(line string, post string, code string, targetUrl string) {
	dm.CreateDevice(line, post, code, targetUrl)
}

func (a *App) DeleteReader(id uint) {
	dm.DeleteDevice(id)
}

func (a *App) GetReader() []d.Device {
	// fmt.Println("Trying to get Reader")
	devices := dm.GetActiveDevices()
	// fmt.Println(devices)
	return devices
}

func (a *App) SetReaderEpc(id uint, epc string) {
	dm.SetDeviceEpc(id, epc)
}

func (a *App) SetReaderCardPresent(id uint, isCardPresent bool) {
	dm.SetDeviceIsCardPresent(id, isCardPresent)
}

func (a *App) SetReaderConnection(id uint, isConnected bool) {
	dm.SetDeviceIsConnected(id, isConnected)
}
