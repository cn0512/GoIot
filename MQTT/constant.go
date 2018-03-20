package MQTT

import (
	_ "fmt"
)

const (
	//BrokerAddr = "tcp://iot.eclipse.org:1883"
	BrokerAddr = "tcp://127.0.0.1:36799"
	Publisher  = "Iot-ykit-pub"
	Subcriber  = "Iot-yit-sub"

	Topic = "Iot-ykit"

	UserName = "ykit"
	Password = "pwd"

	Qos = 1
	/*在MQTT协议中，PUBLISH消息固定头部RETAIN标记，只有为1才要求服务器需要持久保存此消息，除非新的PUBLISH覆盖。*/
	Retain = false

	/*
		如果为false(flag=0)，Client断开连接后，Server应该保存Client的订阅信息。
		如果为true(flag=1)，表示Server应该立刻丢弃任何会话状态信息。
	*/
	ClearSession = false

	AutoConnect = true
)
