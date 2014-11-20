package grobot

var modules = []Module{}

type Module interface {
	Name() string
	LoadConfiguration(conf *Configuration) error
}

func RegisterModule(module Module) {
	modules = append(modules, module)
}

func GetModule(name string) Module {
	for _, m := range modules {
		if m.Name() == name {
			return m
		}
	}
	return nil
}
