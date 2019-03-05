package main

import (
	"fmt"
	"github.com/theshadow/tpl/plugins"
)

var (
	// Name is the name of the plugin and used to prefix all of the functions.
	Name = "example"
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
	m["greeting"] = greeting
	return m
}

func (p *plugin) Version() string {
	return Version
}

func (p *plugin) Name() string {
	return Name
}

func greeting(fname, lname string) string {
	return fmt.Sprintf("Hello, %s %s!", fname, lname)
}
