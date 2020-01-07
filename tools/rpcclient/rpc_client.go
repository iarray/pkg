package rpcclient

import (
	"errors"
	"log"
	"net/rpc"
	"sync"
)

var client *rpc.Client
var nw, remoteAddr string
var lock sync.Mutex

func SetRemoteAddress(netword, addr string) {
	lock.Lock()
	defer lock.Unlock()
	remoteAddr = addr
	nw = netword
	if client != nil {
		client.Close()
	}
}

func Connect() error {
	if client == nil {
		lock.Lock()
		defer lock.Unlock()
		if client == nil {
			var err error
			if client, err = rpc.Dial(nw, remoteAddr); err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

func Close() {
	lock.Lock()
	defer lock.Unlock()
	if client != nil {
		client.Close()
		client = nil
	}
}

func Run(method string, args interface{}, result interface{}) error {
	Connect()
	if client == nil {
		return errors.New("rpc客户端未连接")
	}
	if err := client.Call(method, args, result); err != nil {
		if err == rpc.ErrShutdown {
			//重连
			Close()
			Connect()
		}
		return err
	}
	return nil
}
