package sdkplugin

type PluginOption func(iPlugin)

type iPlugin interface {
	Name() string
	SetName(name string)
}

type plugin struct {
	name string
}

func (s *plugin) SetName(name string) {
	s.name = name
}

func (s *plugin) Name() string {
	return s.name
}

func WithName(name string) PluginOption {
	return func(p iPlugin) {
		p.SetName(name)
	}
}
