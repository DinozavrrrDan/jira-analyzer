package main

import (
	"connector/config"
	"connector/internal/app"
	"connector/pkg/logger"
	"flag"
)

func main() {
	configPath := flag.String("configPath", "config/config-connector.yaml", "Path to the config file")
	log := logger.CreateNewLogger()

	cfg, err := config.NewConfig(*configPath)

	if err != nil {
		log.Log(logger.ERROR, err.Error())
		panic(err)
	}

	newApp := app.NewApp(cfg, log)
	newApp.Run()
}
