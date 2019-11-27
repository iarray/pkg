package mqclient

type IMqClient interface {
	Connect() error
	Publish(topic string, data interface{}, qos int, retain bool) error
}

var instance IMqClient

func Register(c IMqClient) {
	instance = c
}

func I() IMqClient {
	return instance
}
