package device

// import (
// 	"fmt"
// 	"sync"
// 	"time"

// 	mqtt "github.com/eclipse/paho.mqtt.golang"
// )

// type MqttDevice struct {
// 	NextSend   time.Time
// 	c          mqtt.Client
// 	mu         sync.Mutex
// 	clientId   string
// 	connTopic  string
// 	isActive   bool
// }

// func (md *MqttDevice) Setup() {
// 	md.mu.Lock()
// 	defer md.mu.Unlock()

// 	fmt.Println("starting device ", md.HFReader.Id, "loop")
// 	fmt.Println(md.HFReader)

// 	opts := mqtt.NewClientOptions()

// 	md.clientId = (md.HFReader.Line + "-" + md.HFReader.Post + "-" + md.HFReader.Code)
// 	md.connTopic = "status/HF/" + md.HFReader.Line + "/" + md.HFReader.Post + "/" + md.HFReader.Code

// 	opts.AddBroker(md.HFReader.TargetUrl)
// 	opts.SetClientID(md.clientId)
// 	opts.SetWill(md.connTopic, "Offline", 2, true)
// 	opts.SetKeepAlive(5 * time.Second)

// 	md.c = mqtt.NewClient(opts)

// 	md.isActive = true
// }

// func (md *MqttDevice) Loop() {
// 	for {
// 		if !md.isActive {
// 			md.Disconnect()
// 			return
// 		}

// 		md.mu.Lock()
// 		d := md.HFReader
// 		md.mu.Unlock()

// 		fmt.Println("device", d.Code, "loop")
// 		fmt.Println(d)

// 		if !md.c.IsConnected() {
// 			if d.IsConnected {
// 				if token := md.c.Connect(); token.Wait() && token.Error() != nil {
// 					fmt.Println("device with id", md.HFReader.Id, "failed to connect")
// 					fmt.Println(token.Error())
// 					md.mu.Lock()
// 					md.HFReader.IsConnected = false
// 					md.mu.Unlock()
// 					continue
// 				} else {
// 					md.mu.Lock()
// 					md.HFReader.IsConnected = true
// 					md.mu.Unlock()
// 				}

// 				fmt.Println("mqtt Topic : ", md.connTopic)
// 				fmt.Println("mqtt payload :", "Online")

// 				if token := md.c.Publish(md.connTopic, 1, true, "Online"); token.Wait() && token.Error() != nil {
// 					fmt.Println(token.Error())
// 				}

// 				md.NextSend = time.Now().Add(5 * time.Second)
// 			} else {
// 				md.c.Disconnect(100)
// 			}
// 		} else {
// 			if d.IsConnected {
// 				fmt.Println("IsConnected!")

// 				if d.IsCardPresent && d.UidBuffer != "" {
// 					if md.CurrentUid != d.UidBuffer {
// 						// TODO: sent first in to server
// 						err := md.sendMqttPayload(d.TargetUrl, d.Line, d.Post, d.Code, d.UidBuffer, "IN", md.c)
// 						if err == nil {
// 							md.mu.Lock()
// 							md.CurrentUid = d.UidBuffer
// 							md.NextSend = time.Now().Add(30 * time.Second)
// 							md.mu.Unlock()
// 						}
// 					} else {
// 						if time.Now().After(md.NextSend) {
// 							// TODO: sent in to server periodically
// 							err := md.sendMqttPayload(d.TargetUrl, d.Line, d.Post, d.Code, md.CurrentUid, "IN", md.c)
// 							if err != nil {
// 								continue
// 							}
// 							md.mu.Lock()
// 							md.NextSend = time.Now().Add(30 * time.Second)
// 							md.mu.Unlock()
// 						}
// 					}
// 				} else {
// 					if md.CurrentUid == "" {
// 						if time.Now().After(md.NextSend) {
// 							// TODO: sent empty to server periodically
// 							err := md.sendMqttPayload(d.TargetUrl, d.Line, d.Post, d.Code, md.CurrentUid, "EMPTY", md.c)
// 							if err != nil {
// 								continue
// 							}
// 							md.mu.Lock()
// 							md.NextSend = time.Now().Add(30 * time.Second)
// 							md.mu.Unlock()
// 						}
// 					} else {
// 						// TODO: sent out to server
// 						err := md.sendMqttPayload(d.TargetUrl, d.Line, d.Post, d.Code, md.CurrentUid, "OUT", md.c)
// 						if err == nil {
// 							md.mu.Lock()
// 							md.CurrentUid = ""
// 							md.NextSend = time.Now().Add(30 * time.Second)
// 							md.mu.Unlock()
// 						}
// 					}
// 				}
// 			} else {
// 				md.Disconnect()
// 			}
// 		}
// 		time.Sleep(1 * time.Second)
// 	}
// }

// func (md *MqttDevice) Disconnect() {
// 	fmt.Println("Diconnecting from broker...")
// 	if token := md.c.Publish(md.connTopic, 1, true, "Offline"); token.Wait() && token.Error() != nil {
// 		fmt.Println(token.Error())
// 	}
// 	md.c.Disconnect(100)
// 	md.mu.Lock()
// 	md.HFReader.IsConnected = false
// 	md.mu.Unlock()
// 	fmt.Println("Done!")
// }

// func (md *MqttDevice) Destroy() {
// 	fmt.Println("Destroy device with id", md.id)
// 	md.c.Disconnect(100)
// 	fmt.Println("is device with id", md.id, "connected : ", md.c.IsConnected())
// 	md.c = nil
// 	md.isActive = false
// }