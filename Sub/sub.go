package main

/*
	a simple sub client
*/

import (
	Iot "YKIot/MQTT"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var qt_chan chan bool = make(chan bool, 1)

func main() {
	myNoOpStore := &Iot.NoOpStore{}
	opts := MQTT.NewClientOptions()
	opts.AddBroker(Iot.BrokerAddr)
	opts.SetClientID(Iot.Subcriber)
	opts.SetStore(myNoOpStore) //default ":memory"
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetUsername(Iot.UserName)
	opts.SetPassword(Iot.Password)
	opts.SetCleanSession(Iot.ClearSession)
	opts.SetAutoReconnect(Iot.AutoConnect)

	opts.SetConnectionLostHandler(func(client MQTT.Client, err error) {
		fmt.Println("warn Disconnect!")
	})

	var callback MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("Sub recv TOPIC: %s\n", msg.Topic())
		fmt.Printf("Sub recv MSG: %s\n", msg.Payload())
		if cmd := fmt.Sprintf("%s", msg.Payload()); strings.EqualFold(cmd, "exit") {
			qt_chan <- true
			//os.Exit(1)
		}
	}

	opts.SetOnConnectHandler(func(c MQTT.Client) {
		fmt.Println("yeah Connected!")

		if token := c.Subscribe(Iot.Topic, Iot.Qos, callback); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
	})

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	/*
			if token := c.Subscribe("$SYS/broker/load/#", 0, Iot.BrokerLoadHandler); token.Wait() && token.Error() != nil {
				fmt.Println(token.Error())
				os.Exit(1)
			}

			if token := c.Subscribe("$SYS/broker/connection/#", 0, Iot.BrokerConnectionHandler); token.Wait() && token.Error() != nil {
				fmt.Println(token.Error())
				os.Exit(1)
			}

			if token := c.Subscribe("$SYS/broker/clients/#", 0, Iot.BrokerClientsHandler); token.Wait() && token.Error() != nil {
				fmt.Println(token.Error())
				os.Exit(1)
			}

			loadCount := 0
			connectionCount := 0
			clientsCount := 0

			for {
				select {
				case <-Iot.BrokerLoad:
					loadCount++
				case <-Iot.BrokerConnection:
					connectionCount++
				case <-Iot.BrokerClients:
					clientsCount++
				case <-qt_chan:
					fmt.Println("recv server exit")
					break
				}
			}
		fmt.Printf("Received %3d Broker Load messages\n", loadCount)
		fmt.Printf("Received %3d Broker Connection messages\n", connectionCount)
		fmt.Printf("Received %3d Broker Clients messages\n", clientsCount)

	*/
	<-qt_chan

	if token := c.Unsubscribe(Iot.Topic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(50)

	qt := make(chan os.Signal, 1)
	signal.Notify(qt, os.Interrupt, os.Kill)
	<-qt
}
