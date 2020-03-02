package sdkplugin

import (
	"github.com/hitech/gosdk"
	"github.com/hitech/gosdk/sdkcm"
)

type config struct {
	*plugin
}

func NewConfigPlugin(options ...PluginOption) *config {
	p := &plugin{}
	s := &config{
		plugin: p,
	}

	for _, o := range options {
		o(p)
	}

	return s
}

func (s *config) GetConfig() sdkcm.SDKConfig {
	return sdkcm.GetConfig()
}

func (s *config) Run(sdk.Service) error {
	return nil
}

func (s *config) Stop() error {
	return nil
}

func WithConfig(cf sdkcm.SDKConfig) PluginOption {
	return func(p iPlugin) {
		sdkcm.SetConfig(cf)
	}
}
