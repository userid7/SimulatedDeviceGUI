package deviceMqtt

import (
	"sync"
	"time"
)

type Device struct {
	Id            uint `gorm:"primaryKey"`
	Line          string
	Post          string
	Code          string
	UidBuffer     string
	IsCardPresent bool
	CurrentUid    string
	IsConnected   bool
	IsActive      bool
	NextSend      time.Time
	TargetUrl     string
	Mu            *sync.Mutex `gorm:"-"`
}

// type Device struct {
// 	Id            uint `gorm:"primaryKey"`
// 	Line          string
// 	Post          string
// 	Code          string
// 	UidBuffer     string
// 	IsCardPresent bool
// 	// CurrentUid    string
// 	IsConnected bool
// 	// IsActive      bool
// 	// NextSend      time.Time
// 	TargetUrl string
// }

type DeviceRunner struct {
	Device        Device
	NextSend      time.Time
	CurrentUid    string
	ExitChannel   chan (bool)
	DeviceChannel chan (Device)
}
