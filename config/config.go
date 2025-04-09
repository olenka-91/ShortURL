package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type MyConfigs struct {
	ServiceURL     string
	BaseAddressURL string
	FilePath       string
	DBDSN          string
}

var Configs MyConfigs

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
	ps := os.Getenv("DATABASE_DSN")

	// используем init-функцию
	flag.StringVar(&Configs.ServiceURL, "a", ":8081", "default server adress")
	flag.StringVar(&Configs.BaseAddressURL, "b", "http://localhost:8000/", "base adress of short URL")
	flag.StringVar(&Configs.FilePath, "f", "short-url-db.json", "It's a FilePATH")
	if ps == "" {
		flag.StringVar(&Configs.DBDSN, "d", "host=localhost port=5432 user=postgres password=jkmrf1905 dbname=postgres sslmode=disable", "It's DSN string for DB")
	} else {
		Configs.DBDSN = ps
	}

}
