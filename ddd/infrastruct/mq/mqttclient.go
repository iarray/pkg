package mq

type IMqttClient interface {
	Connect()
	Publish(topic string, data interface{}, qos int, retain bool) error
}

var instance IMqttClient

func Register(c IMqttClient) {
	instance = c
}

func I() IMqttClient {
	return instance
}
