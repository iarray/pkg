package nsqclient

import (
	"encoding/json"
	"sync"

	"github.com/iarray/pkg/ddd/infrastruct/mqclient"
	"github.com/nsqio/go-nsq"
)

type NsqPublisher struct {
	LookupHost string
	producer   *nsq.Producer
	lock       sync.Mutex
}

func NewPublisher(lookuphost string) mqclient.IMqClient {
	return &NsqPublisher{LookupHost: lookuphost}
}

func (n *NsqPublisher) Connect() error {
	n.lock.Lock()
	defer n.lock.Unlock()
	if n.producer != nil {
		n.producer.Stop()
	}
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(n.LookupHost, config)
	if err != nil {
		return err
	}
	n.producer = producer
	return nil
}

func (n *NsqPublisher) Publish(topic string, data interface{}, qos int, retain bool) error {
	if n.producer == nil {
		n.Connect()
	}
	buf, err := json.Marshal(data)
	if err != nil {
		//重连
		n.Connect()
		return err
	}
	messageBody := buf
	topicName := topic

	// Synchronously publish a single message to the specified topic.
	// Messages can also be sent asynchronously and/or in batches.
	err = n.producer.Publish(topicName, messageBody)
	if err != nil {
		return err
	}

	return nil
	//producer.Stop()
}
