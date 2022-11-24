package hfreader

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type HFReader struct {
	Id            uint `gorm:"primaryKey"`
	Line          string
	Post          string
	Code          string
	UidBuffer     string
	IsCardPresent bool
	IsConnected   bool `gorm:"-"`
	TargetUrl     string
}

type ActiveHFReader struct {
	id         uint
	HFReader   HFReader
	CurrentUid string
	NextSend   time.Time
	c          mqtt.Client
	mu         sync.Mutex
	clientId   string
	willTopic  string
	isActive   bool
}

func (hf *ActiveHFReader) Setup() {
	hf.mu.Lock()
	defer hf.mu.Unlock()

	fmt.Println("starting device ", hf.HFReader.Id, "loop")
	fmt.Println(hf.HFReader)

	opts := mqtt.NewClientOptions()

	hf.clientId = (hf.HFReader.Line + "-" + hf.HFReader.Post + "-" + hf.HFReader.Code)
	hf.willTopic = "status/HF/" + hf.HFReader.Line + "/" + hf.HFReader.Post + "/" + hf.HFReader.Code

	opts.AddBroker(hf.HFReader.TargetUrl)
	opts.SetClientID(hf.clientId)
	opts.SetWill(hf.willTopic, "Offline", 1, true)
	opts.SetKeepAlive(5 * time.Second)

	hf.c = mqtt.NewClient(opts)

	hf.isActive = true

}

func (hf *ActiveHFReader) Loop() {
	for {
		if !hf.isActive {
			return
		}

		hf.mu.Lock()
		d := hf.HFReader
		hf.mu.Unlock()

		fmt.Println("device", d.Code, "loop")
		fmt.Println(d)
		// fmt.Println(hf.NextSend)

		if !hf.c.IsConnected() {
			if d.IsConnected {
				if token := hf.c.Connect(); token.Wait() && token.Error() != nil {
					fmt.Println("device with id", hf.HFReader.Id, "failed to connect")
					fmt.Println(token.Error())
					hf.mu.Lock()
					hf.HFReader.IsConnected = false
					hf.mu.Unlock()
					continue
				} else {
					hf.mu.Lock()
					hf.HFReader.IsConnected = true
					hf.mu.Unlock()
				}

				fmt.Println("mqtt Topic : ", hf.willTopic)
				fmt.Println("mqtt payload :", "Online")

				if token := hf.c.Publish(hf.willTopic, 1, true, "Online"); token.Wait() && token.Error() != nil {
					fmt.Println(token.Error())
				}

				hf.NextSend = time.Now().Add(5 * time.Second)
			} else {
				hf.c.Disconnect(100)
			}
		}

		if hf.c.IsConnected() {
			fmt.Println("IsConnected!")

			if d.IsCardPresent && d.UidBuffer != "" {
				if hf.CurrentUid != d.UidBuffer {
					// TODO: sent first in to server
					err := hf.sendMqttPayload(d.TargetUrl, d.Line, d.Post, d.Code, d.UidBuffer, "IN", hf.c)
					if err == nil {
						hf.mu.Lock()
						hf.CurrentUid = d.UidBuffer
						hf.NextSend = time.Now().Add(30 * time.Second)
						hf.mu.Unlock()
					}
				} else {
					if time.Now().After(hf.NextSend) {
						// TODO: sent in to server periodically
						err := hf.sendMqttPayload(d.TargetUrl, d.Line, d.Post, d.Code, hf.CurrentUid, "IN", hf.c)
						if err != nil {
							continue
						}
						hf.mu.Lock()
						hf.NextSend = time.Now().Add(30 * time.Second)
						hf.mu.Unlock()
					}
				}
			} else {
				if hf.CurrentUid == "" {
					if time.Now().After(hf.NextSend) {
						// TODO: sent empty to server periodically
						err := hf.sendMqttPayload(d.TargetUrl, d.Line, d.Post, d.Code, hf.CurrentUid, "EMPTY", hf.c)
						if err != nil {
							continue
						}
						hf.mu.Lock()
						hf.NextSend = time.Now().Add(30 * time.Second)
						hf.mu.Unlock()
					}
				} else {
					// TODO: sent out to server
					err := hf.sendMqttPayload(d.TargetUrl, d.Line, d.Post, d.Code, hf.CurrentUid, "OUT", hf.c)
					if err == nil {
						hf.mu.Lock()
						hf.CurrentUid = ""
						hf.NextSend = time.Now().Add(30 * time.Second)
						hf.mu.Unlock()
					}
				}
			}
		} else {
			fmt.Println("IsNotConnected")
			hf.mu.Lock()
			hf.HFReader.IsConnected = false
			hf.mu.Unlock()
		}

		time.Sleep(1 * time.Second)

	}
}

func (hf *ActiveHFReader) Destroy() {
	fmt.Println("Destroy device with id", hf.id)
	hf.c.Disconnect(100)
	fmt.Println("is device with id", hf.id, "connected : ", hf.c.IsConnected())
	hf.c = nil
	hf.isActive = false
}

func (hf *ActiveHFReader) sendMqttPayload(url string, line string, post string, code string, uid string, state string, c mqtt.Client) error {
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
