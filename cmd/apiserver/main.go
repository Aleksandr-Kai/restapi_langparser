package main

import (
	"log"
	"restapi_langparser/internal/apiserver"
	"restapi_langparser/internal/config"
)

func main() {
	cfg := config.New()
	//cfg.Type = config.MemStore
	cfg.DatabaseURL = "user=kai dbname=postgres password=hitomi sslmode=disable"
	if err := apiserver.Start(cfg); err != nil {
		log.Println(err)
	}
}
