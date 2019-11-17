package push

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	//Recieve chan []byte
	Send chan []byte
	sync.Mutex
	client   mqtt.Client
	conState uint32
}

func NewClient(url string, user string, pwd string) *MqttClient {
	//设置连接参数, 断线自动重连, 初次连接直到成功
	clinetOptions := mqtt.NewClientOptions().AddBroker(url).SetUsername(user).SetPassword(pwd).SetAutoReconnect(true).SetConnectRetry(true)
	//设置客户端ID
	clinetOptions.SetClientID(fmt.Sprintf("app-api：%d", time.Now().Unix()))
	//设置handler
	clinetOptions.SetDefaultPublishHandler(messagePubHandler)
	//设置连接超时
	clinetOptions.SetConnectTimeout(time.Duration(60) * time.Second)
	//创建客户端连接
	client := mqtt.NewClient(clinetOptions)

	ret := &MqttClient{
		//Recieve: make(chan []byte),
		Send:   make(chan []byte),
		client: client,
	}

	/*
		i := 0

		for {
			i++
			time.Sleep(time.Duration(3) * time.Second)
			text := fmt.Sprintf("this is test msg #%d ! from task :%d", i, taskId)
			//fmt.Printf("start publish msg to mqtt broker, taskId: %d, count: %d \n", taskId, i)
			//发布消息
			token := client.Publish("go-test-topic", 1, false, text)
			fmt.Printf("[Pub] end publish msg to mqtt broker, taskId: %d, count: %d, token : %s \n", taskId, i, token)
			token.Wait()
		}

		client.Disconnect(250)
		fmt.Println("[Pub] task is ok")
	*/

	return ret
}

func (c *MqttClient) Connect() {
	c.Lock()
	defer c.Unlock()
	//客户端连接判断
	if token := c.client.Connect(); token.WaitTimeout(time.Duration(60)*time.Second) && token.Wait() && token.Error() != nil {
		return
	}

}

func messagePubHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Pub Client Topic : %s \n", msg.Topic())
	fmt.Printf("Pub Client msg : %s \n", msg.Payload())
}

func convertDataToBytes(payload interface{}) ([]byte, error) {
	switch p := payload.(type) {
	case string:
		return []byte(p), nil
	case []byte:
		return []byte(p), nil
	case bytes.Buffer:
		return p.Bytes(), nil
	default:
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

func (c *MqttClient) Publish(topic string, data interface{}, qos int, retain bool) error {
	if !c.client.IsConnected() {
		return errors.New("mqtt client disconnect")
	}

	payload, err := convertDataToBytes(data)
	if err != nil {
		return err
	}

	token := c.client.Publish(topic, byte(qos), retain, payload)
	if token.Wait() && token.Error() != nil {
		log.Println("publish message error:%s", token.Error().Error())
		return token.Error()
	}

	return nil
}
