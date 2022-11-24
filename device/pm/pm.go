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
	pmg.willTopic = "status/HF/" + pmg.PMGateway.Line + "/" + pmg.PMGateway.Code

	opts.AddBroker(pmg.PMGateway.TargetUrl)
	opts.SetClientID(pmg.clientId)
	opts.SetWill(pmg.willTopic, "Offline", 1, true)
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
		d := pmg.HFReader
		pmg.mu.Unlock()

		fmt.Println("device", d.Code, "loop")
		fmt.Println(d)
		// fmt.Println(pmg.NextSend)

		if !pmg.c.IsConnected() {
			if d.IsConnected {
				if token := pmg.c.Connect(); token.Wait() && token.Error() != nil {
					fmt.Println("device with id", pmg.HFReader.Id, "failed to connect")
					fmt.Println(token.Error())
					pmg.mu.Lock()
					pmg.HFReader.IsConnected = false
					pmg.mu.Unlock()
					continue
				} else {
					pmg.mu.Lock()
					pmg.HFReader.IsConnected = true
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

			if d.IsCardPresent && d.UidBuffer != "" {
				if pmg.CurrentUid != d.UidBuffer {
					// TODO: sent first in to server
					err := pmg.sendMqttPayload(d.TargetUrl, d.Line, d.Post, d.Code, d.UidBuffer, "IN", pmg.c)
					if err == nil {
						pmg.mu.Lock()
						pmg.CurrentUid = d.UidBuffer
						pmg.NextSend = time.Now().Add(30 * time.Second)
						pmg.mu.Unlock()
					}
				} else {
					if time.Now().After(pmg.NextSend) {
						// TODO: sent in to server periodically
						err := pmg.sendMqttPayload(d.TargetUrl, d.Line, d.Post, d.Code, pmg.CurrentUid, "IN", pmg.c)
						if err != nil {
							continue
						}
						pmg.mu.Lock()
						pmg.NextSend = time.Now().Add(30 * time.Second)
						pmg.mu.Unlock()
					}
				}
			} else {
				if pmg.CurrentUid == "" {
					if time.Now().After(pmg.NextSend) {
						// TODO: sent empty to server periodically
						err := pmg.sendMqttPayload(d.TargetUrl, d.Line, d.Post, d.Code, pmg.CurrentUid, "EMPTY", pmg.c)
						if err != nil {
							continue
						}
						pmg.mu.Lock()
						pmg.NextSend = time.Now().Add(30 * time.Second)
						pmg.mu.Unlock()
					}
				} else {
					// TODO: sent out to server
					err := pmg.sendMqttPayload(d.TargetUrl, d.Line, d.Post, d.Code, pmg.CurrentUid, "OUT", pmg.c)
					if err == nil {
						pmg.mu.Lock()
						pmg.CurrentUid = ""
						pmg.NextSend = time.Now().Add(30 * time.Second)
						pmg.mu.Unlock()
					}
				}
			}
		} else {
			fmt.Println("IsNotConnected")
			pmg.mu.Lock()
			pmg.HFReader.IsConnected = false
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

func (pmg *ActivePMGateway) sendMqttPayload(url string, line string, post string, code string, uid string, state string, c mqtt.Client) error {
	fmt.Printf("%s is %s\n", code, state)

	payload := &PayloadData{Uid: uid, Status: state}

	b, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to send")
		fmt.Println(err)
		return err
	}
	fmt.Println(string(b))

	topic := "data/HF/" + line + "/" + post + "/" + code

	fmt.Println("mqtt topic : ", topic)
	fmt.Println("mqtt payload :", string(b))

	if token := c.Publish(topic, 1, false, string(b)); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return token.Error()
	}

	return nil
}

type PayloadData struct {
	Uid    string `json:"uid"`
	Status string `json:"operation"`
}
