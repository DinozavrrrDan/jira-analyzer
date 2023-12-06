package main

import (
	"flag"
	"github.com/DinozvrrDan/jira-analyzer/connector/config"
	"github.com/DinozvrrDan/jira-analyzer/connector/internal/app"
	"github.com/DinozvrrDan/jira-analyzer/connector/pkg/logger"
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

	defer newApp.Close()
	newApp.Run()
}
