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
	//DATABASE_DSN = "host=localhost port=5432 user=postgres password=jkmrf1905 dbname=postgres sslmode=disable"
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
	ps := os.Getenv("DATABASE_DSN")

	// используем init-функцию
	flag.StringVar(&Configs.ServiceURL, "a", ":8081", "default server adress")
	flag.StringVar(&Configs.BaseAddressURL, "b", "http://localhost:8000/", "base adress of short URL")
	flag.StringVar(&Configs.FilePath, "f", "short-url-db.json", "It's a FilePATH") //"short-url-db.json"
	if ps == "" {
		flag.StringVar(&Configs.DBDSN, "d", "", "It's DSN string for DB") //"host=localhost port=5432 user=postgres password=jkmrf1905 dbname=postgres sslmode=disable"
	} else {
		Configs.DBDSN = ps
	}

}
