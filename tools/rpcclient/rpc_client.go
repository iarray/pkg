package rpcclient

import (
	"net/rpc"
	"sync"
)

var client *rpc.Client
var nw, remoteAddr string
var lock sync.Mutex

func setRemoteAddress(netword, addr string) {
	lock.Lock()
	defer lock.Unlock()
	remoteAddr = addr
	nw = netword
	if client != nil {
		client.Close()
	}
}

func connect() {
	if client == nil {
		lock.Lock()
		defer lock.Unlock()
		if client == nil {
			client, _ = rpc.Dial(nw, remoteAddr)
		}
	}
}

func run(method string, args interface{}, result interface{}) error {
	connect()
	return client.Call(method, args, result)
}
