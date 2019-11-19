package redis

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type RedisCacheProvider struct {
	Addr     string
	Password string
	DB       int
	client   *redis.Client
	once     sync.Once
}

func (r *RedisCacheProvider) getClient() *redis.Client {
	if r.client == nil {
		r.once.Do(func() {
			r.client = redis.NewClient(&redis.Options{
				Addr:     r.Addr,
				Password: r.Password,
				DB:       r.DB,
			})
		})
	}
	return r.client
}

func (r *RedisCacheProvider) GetJsonObject(key string, data interface{}) error {
	client := r.getClient()
	val, err := client.Get(key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), data)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisCacheProvider) GetString(key string) (string, error) {
	client := r.getClient()
	val, err := client.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisCacheProvider) Set(key string, data interface{}, expiration time.Duration) error {
	client := r.getClient()
	var err error
	switch data.(type) {
	case string:
		err = client.Set(key, data, expiration).Err()
	default:
		val, err2 := json.Marshal(data)
		if err2 == nil {
			err = client.Set(key, val, expiration).Err()
		} else {
			err = err2
			log.Printf("Set Key [%s] Fail , %s", key, err2.Error())
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisCacheProvider) GetSetString(key string, setFunc func() (data string, expiration time.Duration, err error)) (string, error) {
	val, err := r.GetString(key)
	if err == redis.Nil {
		log.Printf("Not Exists Key [%s] Then Setting", key)
		v, exp, err := setFunc()
		if err != nil {
			log.Printf("Call SetFunc Fail , %s", key, err.Error())
			return "", err
		}
		err = r.Set(key, v, exp)
		if err != nil {
			log.Printf("Set Key [%s] Fail , %s", key, err.Error())
		}
		return v, nil
	} else if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisCacheProvider) GetSetJsonObject(key string, data interface{}, setFunc func() (data interface{}, expiration time.Duration, err error)) (interface{}, error) {
	err := r.GetJsonObject(key, data)
	if err == redis.Nil {
		log.Printf("Not Exists Key [%s] Then Setting", key)
		v, exp, err := setFunc()
		if err != nil {
			log.Printf("Call SetFunc Fail , %s", key, err.Error())
			return nil, err
		}
		err = r.Set(key, v, exp)
		if err != nil {
			log.Printf("Set Key [%s] Fail , %s", key, err.Error())
		}
		return v, nil
	} else if err != nil {
		return nil, err
	}

	return data, nil
}
