package nsqclient

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/iarray/pkg/ddd/infrastruct/mqclient"
	"github.com/nsqio/go-nsq"
)

type NsqPublisher struct {
	LookupHost string
	producer   *nsq.Producer
	once       sync.Once
	err        error
}

func NewPublisher(lookuphost string) mqclient.IMqClient {
	return &NsqPublisher{LookupHost: lookuphost}
}

func (n *NsqPublisher) Connect() error {
	n.once.Do(func() {
		config := nsq.NewConfig()
		producer, err2 := nsq.NewProducer(n.LookupHost, config)
		if err2 != nil {
			n.err = err2
		}
		n.producer = producer
	})
	return n.err
}

func (n *NsqPublisher) Publish(topic string, data interface{}, qos int, retain bool) error {
	if n.producer == nil {
		return errors.New("Publisher not connect")
	}
	buf, err := json.Marshal(data)
	if err != nil {
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
