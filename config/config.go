package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type myConfigs struct {
	ServiceURL     string
	BaseAddressURL string
	FileName       string
	DBDSN          string
}

var MyConfigs myConfigs

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
	ps := os.Getenv("DATABASE_DSN")

	// используем init-функцию
	flag.StringVar(&MyConfigs.ServiceURL, "a", ":8081", "default server adress")
	flag.StringVar(&MyConfigs.BaseAddressURL, "b", "http://localhost:8000/", "base adress of short URL")
	flag.StringVar(&MyConfigs.FileName, "f", "short-url-db.json", "It's a FilePATH")
	if ps == "" {
		flag.StringVar(&MyConfigs.DBDSN, "d", "host=localhost port=5432 user=postgres password=jkmrf1905 dbname=postgres sslmode=disable", "It's DSN string for DB")
	} else {
		MyConfigs.DBDSN = ps
	}

}
