package sdkplugin

import (
	"github.com/hitech/gosdk"
	"github.com/hitech/gosdk/sdkcm"
	"github.com/sirupsen/logrus"
)

// type NatPluginOption PluginOption

type natPlugin struct {
	*plugin
	logger    *logrus.Logger
	cf        sdkcm.SDKConfig
	natClient sdkcm.NatClient
	router    func(service sdk.Service)
}

func NewNatPlugin(options ...PluginOption) *natPlugin {
	logger := sdkcm.GetLogger()
	cf := sdkcm.GetConfig()
	client := sdkcm.NewNatClient(logger, cf)
	p := &natPlugin{
		logger:    logger,
		cf:        cf,
		natClient: client,
		plugin:    &plugin{},
	}

	for _, o := range options {
		o(p)
	}

	return p
}

func (s *natPlugin) GetNat() sdkcm.NatClient {
	return s.natClient
}

func (s *natPlugin) Name() string {
	return s.name
}

func (s *natPlugin) Run(service sdk.Service) error {
	if err := s.natClient.Connect(); err != nil {
		return err
	}

	if s.router != nil {
		s.router(service)
	}

	return nil
}

func (s *natPlugin) Stop() error {
	return nil
}

func WithNatPluginRouter(router func(service sdk.Service)) PluginOption {
	return func(i iPlugin) {
		n := i.(*natPlugin)
		n.router = router
	}
	//return func(plg *natPlugin) {
	//	plg.router = router
	//}
}
