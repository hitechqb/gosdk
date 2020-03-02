package sdkplugin

import (
	"github.com/hitechqb/gosdk"
	"github.com/hitechqb/gosdk/sdkcm"
	"github.com/sirupsen/logrus"
)

type redisPlugin struct {
	*plugin
	logger      *logrus.Logger
	cf          sdkcm.SDKConfig
	redisClient sdkcm.RedisClient
}

func NewRedisPlugin(options ...PluginOption) *redisPlugin {
	logger := sdkcm.GetLogger()
	cf := sdkcm.GetConfig()
	client := sdkcm.NewRedisClient(logger, cf)

	s := &redisPlugin{
		plugin:      &plugin{},
		logger:      logger,
		cf:          cf,
		redisClient: client,
	}

	for _, o := range options {
		o(s)
	}

	return s
}

func (s *redisPlugin) GetRedis() sdkcm.RedisClient {
	return s.redisClient
}

func (s *redisPlugin) Run(service sdk.Service) error {
	return s.redisClient.Connect()
}

func (s *redisPlugin) Stop() error {
	return nil
}
