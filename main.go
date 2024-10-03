package main

import (
	"log"
	"transfer/api"
	"transfer/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config")
	}

	server, err := api.NewServer()
	if err != nil {
		log.Fatal("Failed to create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}

}
