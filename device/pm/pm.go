package pm

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
)

type PM struct {
	Id          uint `gorm:"primaryKey"`
	PMGatewayId uint
	Post        string `validate:"required"`
	Code        string `validate:"required"`
	Kw          float32
	Kwh         float32
	IsOk        bool
	IsRandom    bool
}

type PMGateway struct {
	Id          uint   `gorm:"primaryKey"`
	Line        string `validate:"required"`
	Code        string `validate:"required"`
	PMs         []PM
	IsConnected bool `gorm:"-"`
	TargetUrl   string
}

type ActivePMGateway struct {
	id        uint
	PMGateway PMGateway
	NextSend  time.Time
	c         mqtt.Client
	mu        sync.Mutex
	clientId  string
	connTopic string
	isActive  bool
	db        *gorm.DB
}

func (pmg *ActivePMGateway) Setup() {
	pmg.mu.Lock()
	defer pmg.mu.Unlock()

	fmt.Println("starting device ", pmg.PMGateway.Id, "loop")
	fmt.Println(pmg.PMGateway)

	opts := mqtt.NewClientOptions()

	pmg.clientId = (pmg.PMGateway.Line + "-" + pmg.PMGateway.Code)
	pmg.connTopic = "status/PM/" + pmg.PMGateway.Line + "/" + pmg.PMGateway.Code

	opts.AddBroker(pmg.PMGateway.TargetUrl)
	opts.SetClientID(pmg.clientId)
	opts.SetWill(pmg.connTopic, "Offline", 0, true)
	opts.SetKeepAlive(5 * time.Second)

	pmg.c = mqtt.NewClient(opts)

	pmg.isActive = true

	rand.Seed(time.Now().UnixNano())

	go func() {
		for {
			for _, pm := range pmg.PMGateway.PMs {
				if pm.IsRandom {
					min := -3
					max := 3
					offset := rand.Intn(max-min+1)  + min
					newKw := pm.Kw + float32(offset)
					pmg.SetPMKw(pm.Id, newKw)
				}
				if pm.IsOk {
					fmt.Println("NewKwh")
					newKwh := pm.Kwh + (pm.Kw / 360)
					fmt.Println(newKwh)
					pmg.SetPMKwh(pm.Id, newKwh)
				}
			}
			time.Sleep(10 * time.Second)
		}
	}()

}

func (pmg *ActivePMGateway) Loop() {
	for {
		if !pmg.isActive {
			pmg.Disconnect()
			time.Sleep(1 * time.Second)
			return
		}

		pmg.mu.Lock()
		d := pmg.PMGateway
		pmg.mu.Unlock()

		// fmt.Println("device", d.Code, "loop")
		// fmt.Println(d)

		if !pmg.c.IsConnected() {
			if d.IsConnected {
				if token := pmg.c.Connect(); token.Wait() && token.Error() != nil {
					fmt.Println("device with id", pmg.PMGateway.Id, "failed to connect")
					fmt.Println(token.Error())
					pmg.mu.Lock()
					pmg.PMGateway.IsConnected = false
					pmg.mu.Unlock()
					continue
				} else {
					pmg.mu.Lock()
					pmg.PMGateway.IsConnected = true
					pmg.mu.Unlock()
				}

				fmt.Println("mqtt Topic : ", pmg.connTopic)
				fmt.Println("mqtt payload :", "Online")

				if token := pmg.c.Publish(pmg.connTopic, 1, true, "Online"); token.Wait() && token.Error() != nil {
					fmt.Println(token.Error())
				}

				pmg.NextSend = time.Now().Add(5 * time.Second)
			} else {
				// fmt.Println("IsNotConnected")
				pmg.c.Disconnect(100)
			}
		} else {
			if d.IsConnected {
				// fmt.Println("IsConnected!")

				if time.Now().After(pmg.NextSend) {
					for _, pm := range pmg.PMGateway.PMs {
						var status, message string

						if pm.IsOk {
							status = "OK"
							message = ""
						} else {
							status = "ERROR"
							message = "Unable to get data"
						}

						data := &Data{Kw: pm.Kw, Kwh: pm.Kwh}

						err := pmg.sendMqttPayload(d.TargetUrl, d.Line, pm.Post, pm.Code, status, message, data, pmg.c)
						if err != nil {
							fmt.Println("Failed to send pm data with id", pm.Id)
							fmt.Println(err)
							continue
						}

						time.Sleep(1 * time.Second)
					}

					pmg.NextSend = time.Now().Add(30 * time.Second)
				}
			} else {
				pmg.Disconnect()
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (pmg *ActivePMGateway) Disconnect() {
	fmt.Println("Diconnecting from broker...")
	if token := pmg.c.Publish(pmg.connTopic, 1, true, "Offline"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	pmg.c.Disconnect(100)
	pmg.mu.Lock()
	pmg.PMGateway.IsConnected = false
	pmg.mu.Unlock()
	fmt.Println("Done!")
}

func (pmg *ActivePMGateway) Destroy() {
	fmt.Println("Destroy device with id", pmg.id)
	pmg.c.Disconnect(100)
	fmt.Println("is device with id", pmg.id, "connected : ", pmg.c.IsConnected())
	// pmg.c = nil
	pmg.isActive = false
}

func (pmg *ActivePMGateway) SetPMKw(id uint, kw float32) {
	fmt.Println("Set PM kw with id", id)

	selectedPM := &PM{Id: id}

	if result := pmg.db.Model(&selectedPM).Update("Kw", kw); result.Error != nil {
		fmt.Println("Failed to update PM with id", id)
		return
	}
}

func (pmg *ActivePMGateway) SetPMKwh(id uint, kw float32) {
	fmt.Println("Set PM kwh with id", id)

	selectedPM := &PM{Id: id}

	if result := pmg.db.Model(&selectedPM).Update("Kwh", kw); result.Error != nil {
		fmt.Println("Failed to update PM with id", id)
		return
	}
}

func (pmg *ActivePMGateway) sendMqttPayload(url string, line string, post string, code string, status string, message string, data *Data, c mqtt.Client) error {
	var payload *Payload

	if status == "OK" {
		payload = &Payload{Status: status, Message: message, Data: data}
	} else {
		payload = &Payload{Status: status, Message: message}
	}

	b, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal")
		fmt.Println(err)
		return err
	}
	fmt.Println(string(b))

	topic := "data/PM/" + line + "/" + post + "/" + code

	fmt.Println("mqtt topic : ", topic)
	fmt.Println("mqtt payload :", string(b))

	if token := c.Publish(topic, 1, false, string(b)); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return token.Error()
	}

	return nil
}

type Data struct {
	Kw  float32 `json:"kw"`
	Kwh float32 `json:"kwh"`
}

type Payload struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    *Data  `json:"data"`
}
