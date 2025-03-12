package config

import "flag"

type myConfigs struct {
	ServiceURL     string
	BaseAddressURL string
}

var MyConfigs myConfigs

func init() {

	// используем init-функцию
	flag.StringVar(&MyConfigs.ServiceURL, "a", "localhost:8080", "default server adress")
	flag.StringVar(&MyConfigs.BaseAddressURL, "b", "http://localhost:8000/", "base adress of short URL")

}
