package dependency

type Configuration struct {
	VendorsFolder string `json:"folder"`
}

var defaultConfig = Configuration{
	VendorsFolder: "vendor",
}
