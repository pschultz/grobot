package grobot

var modules = []Module{}

type Module interface {
	Name() string
	LoadConfiguration(conf *Configuration) error
}

func RegisterModule(module Module) {
	modules = append(modules, module)
}
