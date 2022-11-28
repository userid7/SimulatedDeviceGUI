package pm

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type PM struct {
	Id          uint `gorm:"primaryKey"`
	PMGatewayId uint
	Post        string
	Code        string
	Kw          float32
	Kwh         float32
	IsOk        bool
	IsRandom    bool
}

type PMGateway struct {
	Id          uint `gorm:"primaryKey"`
	Line        string
	Code        string
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
	willTopic string
	isActive  bool
}

func (pmg *ActivePMGateway) Setup() {
	pmg.mu.Lock()
	defer pmg.mu.Unlock()

	fmt.Println("starting device ", pmg.PMGateway.Id, "loop")
	fmt.Println(pmg.PMGateway)

	opts := mqtt.NewClientOptions()

	pmg.clientId = (pmg.PMGateway.Line + "-" + pmg.PMGateway.Code)
	pmg.willTopic = "status/PM/" + pmg.PMGateway.Line + "/" + pmg.PMGateway.Code

	opts.AddBroker(pmg.PMGateway.TargetUrl)
	opts.SetClientID(pmg.clientId)
	opts.SetWill(pmg.willTopic, "Offline", 0, true)
	opts.SetKeepAlive(5 * time.Second)

	pmg.c = mqtt.NewClient(opts)

	pmg.isActive = true

}

func (pmg *ActivePMGateway) Loop() {
	for {
		if !pmg.isActive {
			return
		}

		pmg.mu.Lock()
		d := pmg.PMGateway
		pmg.mu.Unlock()

		fmt.Println("device", d.Code, "loop")
		fmt.Println(d)

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

				fmt.Println("mqtt Topic : ", pmg.willTopic)
				fmt.Println("mqtt payload :", "Online")

				if token := pmg.c.Publish(pmg.willTopic, 1, true, "Online"); token.Wait() && token.Error() != nil {
					fmt.Println(token.Error())
				}

				pmg.NextSend = time.Now().Add(5 * time.Second)
			} else {
				pmg.c.Disconnect(100)
			}
		}

		if pmg.c.IsConnected() {
			fmt.Println("IsConnected!")

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

					data := Data{Kw: pm.Kw, Kwh: pm.Kwh}

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
			fmt.Println("IsNotConnected")
			pmg.mu.Lock()
			pmg.PMGateway.IsConnected = false
			pmg.mu.Unlock()
		}

		time.Sleep(1 * time.Second)

	}
}

func (pmg *ActivePMGateway) Destroy() {
	fmt.Println("Destroy device with id", pmg.id)
	pmg.c.Disconnect(100)
	fmt.Println("is device with id", pmg.id, "connected : ", pmg.c.IsConnected())
	pmg.c = nil
	pmg.isActive = false
}

func (pmg *ActivePMGateway) sendMqttPayload(url string, line string, post string, code string, status string, message string, data Data, c mqtt.Client) error {
	payload := &Payload{Status: status, Message: message, Data: data}

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
	Data    Data   `json:"data"`
}
