package cache

import (
	"errors"
	"sync"
	"time"
)

type SetStringFunction func() (data string, expiration time.Duration)
type SetObjectFunction func() (data interface{}, expiration time.Duration)

type Cache interface {
	GetString(key string) (string, error)
	GetJsonObject(key string, data interface{}) error
	Set(key string, data interface{}, expiration time.Duration) error
	GetSetString(key string, setFunc func() (data string, expiration time.Duration)) (string, error)
	GetSetJsonObject(key string, data interface{}, setFunc func() (data interface{}, expiration time.Duration)) (interface{}, error)
}

var provider Cache
var once sync.Once

func RegisterProvider(cache Cache) {
	if provider == nil {
		once.Do(func() {
			provider = cache
		})
	}
}

func GetString(key string) (string, error) {
	if provider != nil {
		return provider.GetString(key)
	}
	return "", errors.New("cache provider is nil")
}

func GetJsonObject(key string, data interface{}) error {
	if provider != nil {
		return provider.GetJsonObject(key, data)
	}
	return errors.New("cache provider is nil")
}

func Set(key string, data interface{}, expiration time.Duration) error {
	if provider != nil {
		return provider.Set(key, data, expiration)
	}
	return errors.New("cache provider is nil")
}

func GetSetString(key string, setFunc SetStringFunction) (string, error) {
	if provider != nil {
		return provider.GetSetString(key, setFunc)
	}
	return "", errors.New("cache provider is nil")
}

func GetSetJsonObject(key string, data interface{}, setFunc SetObjectFunction) (interface{}, error) {
	if provider != nil {
		return provider.GetSetJsonObject(key, data, setFunc)
	}
	return nil, errors.New("cache provider is nil")
}
