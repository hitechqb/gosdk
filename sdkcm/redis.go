package sdkcm

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"time"
)

type RedisClient interface {
	Connect() error
	Ping() error
	Do(command string, args ...interface{}) (interface{}, error)
	GetByKey(key string) (interface{}, error)
	GetStruct(out interface{}, cmd string, params ...interface{}) error
	GetStructByKey(key string, out interface{}) error
	Set(key string, value interface{}, expire time.Duration) error
}

type redisClient struct {
	logger *logrus.Logger
	pool   *redis.Pool
	cf     SDKConfig
}

func NewRedisClient(logger *logrus.Logger, cf SDKConfig) *redisClient {
	return &redisClient{
		logger: logger,
		cf:     cf,
		pool: &redis.Pool{
			MaxIdle:   80,
			MaxActive: 12000,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", cf.RedisURL())
			},
		},
	}
}

func (s *redisClient) Connect() error {
	if err := s.Ping(); err != nil {
		return err
	}

	logger.Infof(`🎉 Connected to Redis server on "%s" !`, s.cf.RedisURL())

	return nil
}

func (s *redisClient) Ping() error {
	_, err := s.Do("PING")
	if err != nil {
		return err
	}

	return nil
}

func (s *redisClient) Do(command string, args ...interface{}) (interface{}, error) {
	c := s.pool.Get()
	defer c.Close()
	return c.Do(command, args...)
}

// Get redis item with key
func (s *redisClient) GetByKey(key string) (interface{}, error) {
	return s.Do("GET", key)
}

// Get redis item and parse into struct
func (s *redisClient) GetStruct(out interface{}, cmd string, params ...interface{}) error {
	data, err := s.Do(cmd, params...)
	if err != nil {
		return err
	}

	if data == nil {
		return nil
	}

	bData, err := redis.Bytes(data, err)
	if err != nil {
		return err
	}

	return json.Unmarshal(bData, out)
}

// Get redis item with key and parse into struct
func (s *redisClient) GetStructByKey(key string, out interface{}) error {
	return s.GetStruct(out, "GET", key)
}

func (s *redisClient) Set(key string, value interface{}, expire time.Duration) error {
	_, err := s.Do("SETEX", key, int64(expire.Seconds()), value)
	if err != nil {
		return err
	}

	return nil
}
