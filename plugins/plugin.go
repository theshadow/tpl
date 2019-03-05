package plugins

type Plugin interface {
	Name() string
	Functions() map[string]interface{}
	Version() string
}