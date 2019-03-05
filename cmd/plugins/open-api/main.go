package plugins

import "github.com/theshadow/tpl/plugins"

const (
	// Name is the name of the plugin and used to prefix all of the functions.
	Name = "openapi"
	// SDKVersion is the version of the Plugin SDK this plugin is targeting.
	SDKVersion = "1.0.0"
	// Version is the version of this plugin
	Version = "1.0.0"
)

func New() plugins.Plugin {
	return &plugin{}
}

type plugin struct {
}

func (p *plugin) Functions() map[string]interface{} {
	m := make(map[string]interface{})
	m["definition_from_model"] = definitionFromModel
	return m
}

func (p *plugin) Version() string {
	return Version
}

func (p *plugin) Name() string {
	return Name
}

func definitionFromModel(model map[string]interface{}) string {
	return ""
}
