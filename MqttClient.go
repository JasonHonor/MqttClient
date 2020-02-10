package MqttClient

import (
	"fmt"
	"math/rand"
	"os"
	"sync"

	"github.com/gogf/gf/frame/g"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"crypto/tls"
	"log"
	"strings"
	"time"
)

type MqttClient struct {
	c                mqtt.Client
	ClientOnce       sync.Once
	gGetHostNameOnce sync.Once
	gHostName        string
}

func Subscribe(client mqtt.Client, topic string, callback mqtt.MessageHandler) bool {
	log.Printf("Doing Subscribe Topic %s Handler %v\r\n", topic, callback)
	//topic string, qos byte, retained bool, payload interface{}
	//QoS 0(At most once)」、「QoS 1(At least once)」、「QoS 2(Exactly once」
	token := client.Subscribe(topic, 1, callback)

	if token.Wait() && token.Error() == nil {
		log.Printf("MqttClient Subscribed!\r\n")
		return true
	} else {
		log.Printf("MqttClient Subscribe Error:%s\r\n", token.Error())
	}

	return false
}

func (client *MqttClient) Publish(topic string, msg *string) bool {

	log.Println("Publishing....")

	for {
		fmt.Printf(".")

		if client.c == nil {
			fmt.Printf(" nil\n")
			continue
		}

		fmt.Printf(". %v", client.c.IsConnected())

		if client.c.IsConnected() {
			//topic string, qos byte, retained bool, payload interface{}
			//QoS 0(At most once)」、「QoS 1(At least once)」、「QoS 2(Exactly once」
			token := client.c.Publish(topic, 1, true, *msg)

			if token.Wait() && token.Error() == nil {
				//return true
				log.Println("MsgClient Published.")
				return true
			} else {
				log.Printf("MqttClient Publish Error:%s\r\n", token.Error())
			}
			//return false
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func (mqttClient *MqttClient) newMqttClient(serverUri, clientId string, subscribes map[string]mqtt.MessageHandler) {
	log.Println("MqttClient.newMqttClient")

	// set the protocol, ip and port of the broker. tcp://localhost:1883,ssl://127.0.0.1:883
	opts := mqtt.NewClientOptions().AddBroker(serverUri)

	// set the id to the client.
	opts.SetClientID(clientId)

	opts.SetAutoReconnect(true)

	opts.SetKeepAlive(60 * time.Second)

	opts.SetMaxReconnectInterval(30 * time.Second)

	//on connection lost.
	opts.SetConnectionLostHandler(func(mqtt.Client, error) {
		log.Printf("Connection Lost.\n")
	})

	//on connected
	opts.SetOnConnectHandler(func(mc mqtt.Client) {
		mqttClient.c = mc

		log.Printf("Connection created. TopicList %d\n", len(subscribes))

		//subscribe topics
		for topic, callback := range subscribes {
			Subscribe(mc, topic, callback)
			log.Println("Topic " + topic + " subscribed.")
		}
	})

	opts.SetReconnectingHandler(func(mqtt.Client, *mqtt.ClientOptions) {
		log.Printf("Connection recreating.\n")
	})

	//support ssl protocol.
	if strings.HasPrefix(serverUri, "ssl://") {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}

		opts.SetTLSConfig(tlsConfig)
	}

	c := mqtt.NewClient(opts)

	// create a new client
	token := c.Connect()
	if token == nil {
		log.Printf("mqttClient token =nil\n")
	} else {
		if token.Wait() && token.Error() == nil {
			return
		} else {
			log.Printf("MqttClient Connect Error:%s\r\n", token.Error())
		}
	}
}

func (mc MqttClient) getHostName(NeedRand bool) string {
	mc.gGetHostNameOnce.Do(func() {
		var err error
		mc.gHostName, err = os.Hostname()
		if err != nil {
			mc.gHostName = ""

			rand.Seed(time.Now().UnixNano())

			var NamePart = make([]string, 8)
			//create a random number.
			for i := 0; i < 7; i++ {
				NamePart[i] = fmt.Sprintf("%d", rand.Intn(10))
			}

			mc.gHostName = strings.Join(NamePart, "")

		} else {
			if NeedRand {
				rand.Seed(time.Now().UnixNano())

				var NamePart = make([]string, 8)

				//create a random number.
				for i := 0; i < 7; i++ {
					NamePart[i] = fmt.Sprintf("%d", rand.Intn(10))
				}

				mc.gHostName += strings.Join(NamePart, "")
			}
		}
	})
	return mc.gHostName
}

func (mqttClient *MqttClient) autoConnect(subscribes map[string]mqtt.MessageHandler, clientIds ...string) {

	log.Println("MqttClient.autoConnect")

	mqttClient.ClientOnce.Do(func() {

		log.Println("MqttClient.ConfigOnce")

		vServerUrl := g.Cfg().Get("client.mqtt-server")

		if vServerUrl == nil {
			log.Println("No client.mqtt-server defined.")
			return
		}

		log.Printf("mqtt-server:%v\n", vServerUrl)

		if vServerUrl != nil {

			serverUrl := vServerUrl.(string)

			var clientId string
			if len(clientIds) == 0 {
				clientId = mqttClient.getHostName(true)
			} else {
				clientId = clientIds[0]
			}

			mqttClient.newMqttClient(serverUrl, clientId, subscribes)
		}
	})
}

func NewMqttClient(subscribes map[string]mqtt.MessageHandler, clientIds ...string) *MqttClient {

	client := new(MqttClient)

	client.autoConnect(subscribes, clientIds...)

	return client
}
