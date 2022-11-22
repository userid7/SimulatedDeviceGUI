package deviceMqtt

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (device *Device) Setup() mqtt.Client {
	device.Mu.Lock()
	defer device.Mu.Unlock()

	fmt.Println("starting device ", device.Id, "loop")
	fmt.Println(device)

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	clientId := (device.Line + "-" + device.Post + "-" + device.Code)
	opts.SetClientID(clientId)

	willTopic := "status/HF/" + device.Line + "/" + device.Post + "/" + device.Code
	opts.SetWill(willTopic, "Offline", 1, true)

	opts.SetKeepAlive(3 * time.Second)

	c := mqtt.NewClient(opts)
	return c
}

func (device *Device) Main() {
	c := device.Setup()

	for {
		device.Mu.Lock()
		d := device
		device.Mu.Unlock()

		fmt.Println("device", d.Code, "loop")
		fmt.Println(d)

		if !c.IsConnected() {
			if d.IsConnected {
				if token := c.Connect(); token.Wait() && token.Error() != nil {
					fmt.Println(token.Error())
					continue
				}
				willTopic := "status/HF/" + device.Line + "/" + device.Post + "/" + device.Code
				fmt.Println("mqtt Topic : ", willTopic)
				fmt.Println("mqtt payload :", "Online")

				if token := c.Publish(willTopic, 1, true, "Online"); token.Wait() && token.Error() != nil {
					fmt.Println(token.Error())
				}

				device.NextSend = time.Now().Add(5 * time.Second)
			} else {
				c.Disconnect(100)
			}
		}

		if c.IsConnected() {
			fmt.Println("IsConnected!")

			if d.IsCardPresent && d.UidBuffer != "" {
				if d.CurrentUid != d.UidBuffer {
					// TODO: sent first in to server
					err := sendMqttPayload(d.TargetUrl, d.Line, d.Post, d.Code, d.UidBuffer, "IN", c)
					if err == nil {
						device.Mu.Lock()
						device.CurrentUid = device.UidBuffer
						device.NextSend = time.Now().Add(30 * time.Second)
						device.Mu.Unlock()
					}
				} else {
					if time.Now().After(d.NextSend) {
						// TODO: sent in to server periodically
						err := sendMqttPayload(d.TargetUrl, d.Line, d.Post, d.Code, d.CurrentUid, "IN", c)
						if err != nil {
							continue
						}
						device.Mu.Lock()
						device.NextSend = time.Now().Add(30 * time.Second)
						device.Mu.Unlock()
					}
				}
			} else {
				if device.CurrentUid == "" {
					if time.Now().After(device.NextSend) {
						// TODO: sent empty to server periodically
						err := sendMqttPayload(device.TargetUrl, device.Line, device.Post, device.Code, device.CurrentUid, "EMPTY", c)
						if err != nil {
							continue
						}
						device.Mu.Lock()
						device.NextSend = time.Now().Add(30 * time.Second)
						device.Mu.Unlock()
					}
				} else {
					// TODO: sent out to server
					err := sendMqttPayload(device.TargetUrl, device.Line, device.Post, device.Code, device.CurrentUid, "OUT", c)
					if err == nil {
						device.Mu.Lock()
						device.CurrentUid = ""
						device.NextSend = time.Now().Add(30 * time.Second)
						device.Mu.Unlock()
					}
				}
			}
		} else {
			device.Mu.Lock()
			device.IsConnected = false
			device.Mu.Unlock()
		}
		time.Sleep(1 * time.Second)
	}
}

func (device *Device) Loop() {
	isMqttConnect := false

	// TODO: add device loop task
	fmt.Println("starting device ", device.Id, "loop")
	fmt.Println(device)

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	clientId := (device.Line + "-" + device.Post + "-" + device.Code)
	opts.SetClientID(clientId)

	willTopic := "status/HF/" + device.Line + "/" + device.Post + "/" + device.Code
	opts.SetWill(willTopic, "Offline", 1, true)

	opts.SetKeepAlive(3 * time.Second)

	c := mqtt.NewClient(opts)

	for {
		if device.IsConnected {
			if !isMqttConnect {
				if token := c.Connect(); token.Wait() && token.Error() != nil {
					fmt.Println(token.Error())
					continue
				}

				isMqttConnect = true

				fmt.Println("mqtt Topic : ", willTopic)
				fmt.Println("mqtt payload :", "Online")

				if token := c.Publish(willTopic, 1, true, "Online"); token.Wait() && token.Error() != nil {
					fmt.Println(token.Error())
				}

				device.NextSend = time.Now().Add(5 * time.Second)

			} else {
				fmt.Println(device)

				if !device.IsActive {
					fmt.Println("stopping device ", device.Id, "loop")
					return
				}

				if device.IsCardPresent {
					if device.CurrentUid == "" {
						// TODO: sent first in to server
						err := sendMqttPayload(device.TargetUrl, device.Line, device.Post, device.Code, device.UidBuffer, "IN", c)
						if err == nil {
							device.CurrentUid = device.UidBuffer
							device.NextSend = time.Now().Add(30 * time.Second)
						}
					} else {
						if time.Now().After(device.NextSend) {
							// TODO: sent in to server periodically
							err := sendMqttPayload(device.TargetUrl, device.Line, device.Post, device.Code, device.CurrentUid, "IN", c)
							if err != nil {
								device.IsConnected = false
							}
							device.NextSend = time.Now().Add(30 * time.Second)
						}
					}
				} else {
					if device.CurrentUid == "" {
						if time.Now().After(device.NextSend) {
							// TODO: sent empty to server periodically
							err := sendMqttPayload(device.TargetUrl, device.Line, device.Post, device.Code, device.CurrentUid, "EMPTY", c)
							if err != nil {
								device.IsConnected = false
							}
							device.NextSend = time.Now().Add(30 * time.Second)
						}
					} else {
						// TODO: sent out to server
						err := sendMqttPayload(device.TargetUrl, device.Line, device.Post, device.Code, device.CurrentUid, "OUT", c)
						if err == nil {
							device.CurrentUid = ""
							device.NextSend = time.Now().Add(30 * time.Second)
						}
					}
				}

				// time.Sleep(500 * time.Millisecond)
				// fmt.Println("Done 1 device MQTT cycle")
			}
		} else {
			c.Disconnect(1000)
		}
	}
}

func (device *Device) CardIn(uid string) {
	device.UidBuffer = uid
	device.IsCardPresent = true
}

func (device *Device) CardOut() {
	device.IsCardPresent = false
}

func (device *Device) DeleteDevice() {
	device.IsActive = false
}

func sendMqttPayload(url string, line string, post string, code string, uid string, state string, c mqtt.Client) error {
	fmt.Printf("%s is %s\n", code, state)

	payload := &PayloadData{Uid: uid, Status: state}

	b, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to send")
		fmt.Println(err)
		return err
	}
	fmt.Println(string(b))

	// if state == "IN" {
	// 	s = fmt.Sprintf("%s/1/%s/IN", uid, code)
	// } else if state == "OUT" {
	// 	s = fmt.Sprintf("%s/%s/OUT", uid, code)
	// } else if state == "EMPTY" {
	// 	s = fmt.Sprintf("0")
	// }

	// topic := "data/device/reader/" + code

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
