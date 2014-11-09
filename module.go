package gobot

import "encoding/json"

var modules = []Module{}

type Module interface {
	Name() string
	LoadConfiguration(map[string]*json.RawMessage) error
}

func RegisterModule(module Module) {
	modules = append(modules, module)
}
