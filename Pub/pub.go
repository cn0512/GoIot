package main

/*
	a simple pub svr
*/

import (
	Iot "YKIot/MQTT"
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	myNoOpStore := &Iot.NoOpStore{}
	opts := MQTT.NewClientOptions()
	opts.AddBroker(Iot.BrokerAddr)
	opts.SetClientID(Iot.Publisher)
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

	opts.SetOnConnectHandler(func(client MQTT.Client) {
		fmt.Println("yeah Connected!")
	})

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//use input to pub
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		payload := scanner.Text()
		token := c.Publish(Iot.Topic, Iot.Qos, Iot.Retain, payload)
		token.Wait()

		if payload == "exit" {
			fmt.Println("ctrl+c to exit")
			break
		}
	}

	c.Disconnect(50) //250

	qt := make(chan os.Signal, 1)
	signal.Notify(qt, os.Interrupt, os.Kill)
	<-qt
}
