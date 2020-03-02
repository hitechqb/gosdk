package sdkplugin

import (
	"fmt"
	"github.com/hitechqb/gosdk"
	"github.com/hitechqb/gosdk/sdkcm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type httpServer struct {
	*plugin
	logger *logrus.Logger
	cf     sdkcm.SDKConfig
	router func(*echo.Echo, sdk.Service)
	engine *echo.Echo
}

func NewDefaultHttpServerPlugin(options ...PluginOption) *httpServer {
	s := NewHttpServerPlugin(options...)

	// add standard middleware
	s.engine.Use(middleware.Gzip())
	s.engine.Use(middleware.Recover())
	s.engine.Pre(middleware.RemoveTrailingSlash())
	//...

	return s
}

func NewHttpServerPlugin(options ...PluginOption) *httpServer {
	logger := sdkcm.GetLogger()
	cf := sdkcm.GetConfig()
	engine := echo.New()
	engine.HideBanner = true

	engine.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(sdk.NewSDKContext(c, logger, cf))
		}
	})

	// ...
	s := &httpServer{
		plugin: &plugin{},
		engine: engine,
		logger: logger,
		cf:     cf,
	}

	for _, o := range options {
		o(s)
	}

	return s
}

func (s *httpServer) GetServer() interface{} {
	return s.engine
}

func (s *httpServer) Run(service sdk.Service) error {
	if s.router != nil {
		s.router(s.engine, service)
	}

	return s.engine.Start(fmt.Sprintf(`:%s`, s.cf.AppPort()))
}

func (s *httpServer) Stop() error {
	return nil
}

func WithHttpServerRouter(router func(*echo.Echo, sdk.Service)) PluginOption {
	return func(p iPlugin) {
		s, ok := p.(*httpServer)
		if !ok {
			sdkcm.GetLogger().Fatalf(`HttpServer plugin can't cast: iPlugin -> *httpServer`)
		}
		s.router = router
	}
}
