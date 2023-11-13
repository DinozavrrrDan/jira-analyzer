package apiserver

import (
	"Jira-analyzer/jiraConnector/configReader"
	"Jira-analyzer/jiraConnector/connector"
	"Jira-analyzer/jiraConnector/dbPusher"
	"Jira-analyzer/jiraConnector/logger"
)

//Как я понял он создает просто все части

type apiServer struct {
	configReader   *configReader.ConfigRaeder
	logger         *logger.JiraLogger
	jiraConnector  *connector.Connector
	databasePusher *dbPusher.DatabasePusher
	//transformer *transformer.ne будет
}

func CreateNewApiServer() *apiServer {
	newReader := configReader.CreateNewConfigReader()
	return &apiServer{
		configReader:   newReader,
		logger:         logger.CreateNewLogger(),
		jiraConnector:  connector.CreateNewJiraConnector(),
		databasePusher: dbPusher.CreateNewDatabasePusher(),
	}
}
