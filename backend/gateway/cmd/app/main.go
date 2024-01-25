package main

import (
	"flag"
	"github.com/DinozvrrDan/jira-analyzer/backend/gateway/config"
	"github.com/DinozvrrDan/jira-analyzer/backend/gateway/internal/app"
	"github.com/DinozvrrDan/jira-analyzer/backend/gateway/pkg/logger"
)

func main() {
	configPath := flag.String("configPath", "backend/gateway/config/config-gateway.yaml", "Path to the config file")
	flag.Parse()

	log := logger.CreateNewLogger()

	cfg, err := config.NewConfig(*configPath)

	if err != nil {
		log.Log(logger.ERROR, err.Error())
		panic(err)
	}

	newApp, err := app.NewApp(cfg, log)

	if err != nil {
		log.Log(logger.ERROR, err.Error())
		panic(err)
	}

	defer func(newApp *app.App) {
		err := newApp.Close()
		if err != nil {
			log.Log(logger.ERROR, err.Error())
			panic(err)
		}
	}(newApp)

	if err = newApp.Run(); err != nil {
		log.Log(logger.ERROR, err.Error())
		panic(err)
	}
}
