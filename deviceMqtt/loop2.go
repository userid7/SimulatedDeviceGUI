package deviceMqtt

// import (
// 	"encoding/json"
// 	"fmt"
// 	"time"

// 	mqtt "github.com/eclipse/paho.mqtt.golang"
// )

// func (deviceRunner DeviceRunner) Loop() {
// 	device := deviceRunner.D

// 	// TODO: add device loop task
// 	fmt.Println("starting device ", device.Id, "loop")
// 	fmt.Println(device)

// 	opts := mqtt.NewClientOptions()
// 	opts.AddBroker("tcp://localhost:1883")
// 	clientId := (device.Line + "-" + device.Post + "-" + device.Code)
// 	opts.SetClientID(clientId)

// 	willTopic := "status/HF/" + device.Line + "/" + device.Post + "/" + device.Code
// 	// opts.SetWill(willTopic, "OFF", 1, true)

// 	c := mqtt.NewClient(opts)

// 	if token := c.Connect(); token.Wait() && token.Error() != nil {
// 		panic(token.Error())
// 	}

// 	fmt.Println("mqtt Topic : ", willTopic)
// 	fmt.Println("mqtt payload :", "Online")

// 	if token := c.Publish(willTopic, 1, true, "Online"); token.Wait() && token.Error() != nil {
// 		fmt.Println(token.Error())
// 	}

// 	deviceRunner.NextSend = time.Now().Add(5 * time.Second)

// 	for {
// 		fmt.Println(device)

// 		if !deviceRunner.IsActive {
// 			fmt.Println("stopping device ", device.Id, "loop")
// 			return
// 		}

// 		if device.IsCardPresent {
// 			if deviceRunner.CurrentUid == "" {
// 				// TODO: sent first in to server
// 				err := sendPostRequest(device.TargetUrl, device.Line, device.Post, device.Code, device.UidBuffer, "IN", c)
// 				if err == nil {
// 					deviceRunner.CurrentUid = device.UidBuffer
// 					deviceRunner.NextSend = time.Now().Add(30 * time.Second)
// 				}
// 			} else {
// 				if time.Now().After(deviceRunner.NextSend) {
// 					// TODO: sent in to server periodically
// 					err := sendPostRequest(device.TargetUrl, device.Line, device.Post, device.Code, deviceRunner.CurrentUid, "IN", c)
// 					if err != nil {
// 						device.IsConnected = false
// 					}
// 					deviceRunner.NextSend = time.Now().Add(30 * time.Second)
// 				}
// 			}
// 		} else {
// 			if deviceRunner.CurrentUid == "" {
// 				if time.Now().After(deviceRunner.NextSend) {
// 					// TODO: sent empty to server periodically
// 					err := sendPostRequest(device.TargetUrl, device.Line, device.Post, device.Code, deviceRunner.CurrentUid, "EMPTY", c)
// 					if err != nil {
// 						device.IsConnected = false
// 					}
// 					deviceRunner.NextSend = time.Now().Add(30 * time.Second)
// 				}
// 			} else {
// 				// TODO: sent out to server
// 				err := sendPostRequest(device.TargetUrl, device.Line, device.Post, device.Code, deviceRunner.CurrentUid, "OUT", c)
// 				if err == nil {
// 					deviceRunner.CurrentUid = ""
// 					deviceRunner.NextSend = time.Now().Add(30 * time.Second)
// 				}
// 			}
// 		}

// 		time.Sleep(3 * time.Second)
// 	}
// }

// func (deviceRunner DeviceRunner) CardIn(uid string) {
// 	deviceRunner.D.UidBuffer = uid
// 	deviceRunner.D.IsCardPresent = true
// }

// func (deviceRunner DeviceRunner) CardOut() {
// 	deviceRunner.D.IsCardPresent = false
// }

// func (deviceRunner DeviceRunner) DeleteDeviceRunner() {
// 	deviceRunner.IsActive = false
// }

// func sendPostRequest(url string, line string, post string, code string, uid string, state string, c mqtt.Client) error {
// 	fmt.Printf("%s is %s\n", code, state)

// 	payload := &PayloadData{Uid: uid, Status: state}

// 	b, err := json.Marshal(payload)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	fmt.Println(string(b))

// 	// if state == "IN" {
// 	// 	s = fmt.Sprintf("%s/1/%s/IN", uid, code)
// 	// } else if state == "OUT" {
// 	// 	s = fmt.Sprintf("%s/%s/OUT", uid, code)
// 	// } else if state == "EMPTY" {
// 	// 	s = fmt.Sprintf("0")
// 	// }

// 	// topic := "data/device/reader/" + code

// 	topic := "data/HF/" + line + "/" + post + "/" + code

// 	fmt.Println("mqtt topic : ", topic)
// 	fmt.Println("mqtt payload :", string(b))

// 	if token := c.Publish(topic, 1, false, string(b)); token.Wait() && token.Error() != nil {
// 		fmt.Println(token.Error())
// 		return token.Error()
// 	}

// 	return nil
// }

// type PayloadData struct {
// 	Uid    string `json:"uid"`
// 	Status string `json:"operation"`
// }
