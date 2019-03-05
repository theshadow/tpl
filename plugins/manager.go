package plugins

import (
	"fmt"
	"plugin"
)

const (
	SDKVersion = "1.0.0"

	PluginConstSDKVersion = "PluginSDKVersion"
	PluginConstName = "Name"
	PluginConstVersion = "Version"
	PluginConstConstructor = "New"
)

type ConstructorFn func() Plugin

type Manager struct {
	plugins map[string]Plugin
	versions map[string]string
}

func (m Manager) Load(paths ...string) error {
	for _, path := range paths {
		p, err := loadPlugin(path)
		if err != nil {
			return fmt.Errorf("unable to load plugin '%s' from '%s': %s", p.Name(), path, err)
		}
		m.plugins[p.Name()] = p
	}
	return nil
}

func loadPlugin(path string) (Plugin, error) {
	pn, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load plugin")
	}

	s, err := pn.Lookup(PluginConstSDKVersion)
	if err != nil {
		return nil, fmt.Errorf("unable to lookup plugin %s", PluginConstSDKVersion)
	}
	if *s.(*string) != SDKVersion {
		return nil, fmt.Errorf("incompatible SDK version %s, expected %s", *s.(*string), SDKVersion)
	}

	s, err = pn.Lookup(PluginConstName)
	if err != nil {
		return nil, fmt.Errorf("unable to lookup plugin %s", PluginConstName)
	}

	name := *s.(*string)
	if len(name) == 0 {
		return nil, fmt.Errorf("invalid plugin %s, must be a valid non-zero length string", PluginConstName)
	}

	s, err = pn.Lookup(PluginConstVersion)
	if err != nil {
		return nil, fmt.Errorf("unable to lookup plugin %s", PluginConstVersion)
	}

	version := *s.(*string)
	if len(version) == 0 {
		return nil, fmt.Errorf("invalid plugin %s, must be a valid non-zero length string", PluginConstVersion)
	}

	s, err = pn.Lookup(PluginConstConstructor)
	if err != nil {
		return nil, fmt.Errorf("unable to lookup plugin constructor %s", PluginConstConstructor)
	}

	p := s.(ConstructorFn)()
	return p, nil
}

