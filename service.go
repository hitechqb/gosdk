package sdk

import (
	"github.com/hitech/gosdk/sdkcm"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/color"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Option func(*service)

type Service interface {
	Start()
	GetPlugin(pluginName string) interface{}
	GetDb(pluginName string) *gorm.DB
	GetRedis(pluginName string) sdkcm.RedisClient
	GetHttpServer(pluginName string) *echo.Echo
	GetNAT(pluginName string) sdkcm.NatClient
	Info()
}

type Plugin interface {
	Name() string
	Run(Service) error
	Stop() error
}

type ConfigPlugin interface {
	GetConfig() sdkcm.SDKConfig
}

type GormPlugin interface {
	GetDb() *gorm.DB
}

type HttpServerPlugin interface {
	GetServer() interface{}
}

type NatPlugin interface {
	GetNat() sdkcm.NatClient
}

type RedisPlugin interface {
	GetRedis() sdkcm.RedisClient
}

type Commander interface {
	AddCommand(use, desc string, run func(Service))
	Execute()
}

type service struct {
	plugins   map[string]Plugin
	logger    *logrus.Logger
	stopChan  chan error
	commander Commander
}

func NewService(options ...Option) *service {
	logger := sdkcm.GetLogger()

	s := &service{
		plugins: make(map[string]Plugin),
		logger:  logger,
	}

	s.commander = NewCommander(s)

	for _, o := range options {
		o(s)
	}

	return s
}

func (s *service) run() {
	// load env
	if err := godotenv.Load(); err != nil {
		s.logger.Fatalln(err)
	}

	var httpServers, normalPlugins []Plugin
	for _, p := range s.plugins {
		if _, ok := p.(HttpServerPlugin); ok {
			httpServers = append(httpServers, p)
			continue
		}
		normalPlugins = append(normalPlugins, p)
	}

	// start plugins
	for _, sub := range normalPlugins {
		if err := sub.Run(s); err != nil {
			s.logger.Fatalf("Plugin %s has error when start: %s", sub.Name(), err)
		}
	}

	// start servers
	for _, sub := range httpServers {
		if err := sub.Run(s); err != nil {
			s.logger.Fatalf("Plugin %s has error when start: %s", sub.Name(), err)
		}
	}

	// stop
	for _, sub := range s.plugins {
		if err := sub.Stop(); err != nil {
			s.logger.Fatalf("Plugin %s has error when start: %s", sub.Name(), err)
		}
	}
}

func (s *service) Start() {
	s.commander.AddCommand(
		"serve",
		"Start service", func(Service) {
			s.run()
		})
	s.commander.Execute()
}

func (s *service) GetPlugin(pluginName string) interface{} {
	return s.plugins[pluginName]
}

func (s *service) GetDb(pluginName string) *gorm.DB {
	plg, ok := s.plugins[pluginName].(GormPlugin)
	if ok {
		return plg.GetDb()
	}

	s.logger.Warnf(`🔥 Gorm plugin "%s" does not exist !`, pluginName)

	return nil
}

func (s *service) GetRedis(pluginName string) sdkcm.RedisClient {
	plg, ok := s.plugins[pluginName].(RedisPlugin)
	if ok {
		return plg.GetRedis()
	}

	s.logger.Warnf(`🔥 Redis plugin "%s" does not exist !`, pluginName)

	return nil
}

func (s *service) GetHttpServer(pluginName string) *echo.Echo {
	plg, ok := s.plugins[pluginName].(HttpServerPlugin)
	if ok {
		server := plg.GetServer()

		if e, ok := server.(*echo.Echo); ok {
			return e
		}

		if _, ok := server.(*http.Server); ok || server == nil {
			ec := echo.New()
			ec.Server = server.(*http.Server)
			return ec
		}

		s.logger.Fatalln(`Unknown server type, can not process`)
	}

	s.logger.Warnf(`🔥 HTTP Server plugin "%s" does not exist !`, pluginName)

	return nil
}

func (s *service) GetNAT(pluginName string) sdkcm.NatClient {
	plg, ok := s.plugins[pluginName].(NatPlugin)
	if ok {
		return plg.GetNat()
	}

	s.logger.Warnf(`🔥 NAT plugin "%s" does not exist !`, pluginName)

	return nil
}

func (s *service) Info() {
	colored := color.New()
	colored.Println(`
=================================================================================
=                                     Info                                      =
=================================================================================`)
	colored.Println("📦 Plugins:")
	index := 0
	for k, _ := range s.plugins {
		index++
		colored.Printf("%d. 👉 %s \n", index, k)
	}
	colored.Println("=================================================================================")

}

func WithPlugin(plg Plugin) Option {
	return func(s *service) {
		for _, m := range s.plugins {
			if _, ok := s.plugins[plg.Name()]; ok {
				s.logger.Fatalf(`Plugin "%s" has already existed`, m.Name())
			}
		}
		s.plugins[plg.Name()] = plg
	}
}

func WithCommand(use, desc string, run func(Service)) Option {
	return func(s *service) {
		s.commander.AddCommand(use, desc, run)
	}
}
