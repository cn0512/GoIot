package MQTT

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var BrokerLoad = make(chan bool)
var BrokerConnection = make(chan bool)
var BrokerClients = make(chan bool)

func BrokerLoadHandler(client MQTT.Client, msg MQTT.Message) {
	BrokerLoad <- true
	fmt.Printf("BrokerLoadHandler         ")
	fmt.Printf("[%s]  ", msg.Topic())
	fmt.Printf("%s\n", msg.Payload())
}

func BrokerConnectionHandler(client MQTT.Client, msg MQTT.Message) {
	BrokerConnection <- true
	fmt.Printf("BrokerConnectionHandler   ")
	fmt.Printf("[%s]  ", msg.Topic())
	fmt.Printf("%s\n", msg.Payload())
}

func BrokerClientsHandler(client MQTT.Client, msg MQTT.Message) {
	BrokerClients <- true
	fmt.Printf("BrokerClientsHandler      ")
	fmt.Printf("[%s]  ", msg.Topic())
	fmt.Printf("%s\n", msg.Payload())
}
