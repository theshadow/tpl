package plugins

import (
	"fmt"
	"plugin"
)

const (
	SDKVersion = "1.0.0"

	PluginConstSDKVersion = "SDKVersion"
	PluginConstName = "Name"
	PluginConstVersion = "Version"
	PluginConstConstructor = "New"
)

type Manager struct {
	Plugins  map[string]Plugin
	Versions map[string]string
}

func New() *Manager {
	return &Manager{
		Plugins:  make(map[string]Plugin),
		Versions: make(map[string]string),
	}
}

func (m *Manager) Load(paths ...string) error {
	for _, path := range paths {
		p, err := loadPlugin(path)
		if err != nil {
			return fmt.Errorf("unable to load plugin from '%s': %s", path, err)
		}
		m.Plugins[p.Name()] = p
		m.Versions[p.Name()] = p.Version()
	}
	return nil
}

func loadPlugin(path string) (Plugin, error) {
	pn, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load plugin %s", err)
	}

	s, err := pn.Lookup(PluginConstSDKVersion)
	if err != nil {
		return nil, fmt.Errorf("unable to lookup plugin var '%s' because %s", PluginConstSDKVersion, err)
	}
	if *s.(*string) != SDKVersion {
		return nil, fmt.Errorf("incompatible SDK version %s, expected %s", *s.(*string), SDKVersion)
	}

	s, err = pn.Lookup(PluginConstName)
	if err != nil {
		return nil, fmt.Errorf("unable to lookup plugin var '%s' %s", PluginConstName, err)
	}

	name := *s.(*string)
	if len(name) == 0 {
		return nil, fmt.Errorf("invalid plugin %s, must be a valid non-zero length string", PluginConstName)
	}

	s, err = pn.Lookup(PluginConstVersion)
	if err != nil {
		return nil, fmt.Errorf("unable to lookup plugin %s: %s", PluginConstVersion, err)
	}

	version := *s.(*string)
	if len(version) == 0 {
		return nil, fmt.Errorf("invalid plugin %s, must be a valid non-zero length string", PluginConstVersion)
	}

	s, err = pn.Lookup(PluginConstConstructor)
	if err != nil {
		return nil, fmt.Errorf("unable to lookup plugin constructor %s: %s", PluginConstConstructor, err)
	}

	p := s.(func() Plugin)()
	return p, nil
}

