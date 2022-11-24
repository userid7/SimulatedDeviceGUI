package main

import (
	dm "SimulatedDeviceGUI/deviceManager"
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

func (a *App) GetReader() []dm.Device {
	devices := dm.GetActiveDevices()
	return devices
}

func (a *App) SetReaderCardPresent() {
}
func (a *App) SetReaderConnection() {
}
func (a *App) SetReaderEpc() {
}
