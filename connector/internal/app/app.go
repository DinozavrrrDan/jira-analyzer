package app

import (
	"connector/config"
	"connector/internal/server"
	"connector/internal/service"
	"connector/pkg/logger"
)

func Run(configPath string) {

	log := logger.CreateNewLogger()

	cfg, err := config.CreateNewConfigReader(configPath, log)
	if err != nil {
		log.Log(logger.ERROR, "Config err: %s"+err.Error())
	}
	//Инициализация репозиториев ?? это с DB

	//Инициализация зависимостей
	deps := service.ServicesDependencies{
		JiraRepositoryUrl: cfg.GetJiraUrl(),
	}

	//Инициализация сервисов
	services := service.NewServices(deps, log, cfg)
	//Cтарт сервера

	apiServer := server.NewApiServer(services, log, cfg)

	apiServer.StartServer()
}
