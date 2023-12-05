package config

import (
	"connector/pkg/logger"
	"fmt"
	"github.com/spf13/viper"
)

type Reader struct {
	viperConfigReader *viper.Viper
}

func CreateNewConfigReader(configPath string, log *logger.Logger) (*Reader, error) {
	configReader := Reader{}
	configReader.viperConfigReader = viper.New()
	configReader.viperConfigReader.SetConfigName("config")
	configReader.viperConfigReader.SetConfigType("yaml")
	configReader.viperConfigReader.AddConfigPath("./config")
	configReader.viperConfigReader.AddConfigPath("../config")
	configReader.viperConfigReader.AddConfigPath("../config/config.yaml")

	//Как только разберусь с приколом viper добавлю
	//configReader.viperConfigReader.AddConfigPath(configPath)
	if err := configReader.viperConfigReader.ReadInConfig(); err != nil {
		fmt.Print("AAAAAAAA")
		log.Log(logger.ERROR, "error in configReader"+err.Error())
		return nil, err
	}
	return &configReader, nil
}

func (configReader *Reader) GetJiraUrl() string {
	return configReader.viperConfigReader.GetString("jiraUrl")
}

func (configReader *Reader) GetHostDB() string {
	return configReader.viperConfigReader.GetString("database.hostDB")
}

func (configReader *Reader) GetDatabaseName() string {
	return configReader.viperConfigReader.GetString("database.nameDB")
}

func (configReader *Reader) GetPasswordDB() string {
	return configReader.viperConfigReader.GetString("database.passwordDB")
}

func (configReader *Reader) GetUserDb() string {
	return configReader.viperConfigReader.GetString("database.userDB")
}

func (configReader *Reader) GetPortDB() int {
	return configReader.viperConfigReader.GetInt("database.portDB")
}

func (configReader *Reader) GetThreadCount() int {
	return configReader.viperConfigReader.GetInt("threadCount")
}

func (configReader *Reader) GetIssuesOnOneRequest() int {
	return configReader.viperConfigReader.GetInt("issueInOneRequest")
}

func (configReader *Reader) GetMaxTimeSleep() int {
	return configReader.viperConfigReader.GetInt("maxTimeSleep")
}

func (configReader *Reader) GetMinTimeSleep() int {
	return configReader.viperConfigReader.GetInt("maxTimeSleep")
}

func (configReader *Reader) GetHostServer() string {
	return configReader.viperConfigReader.GetString("apiServer.hostServer")
}

func (configReader *Reader) GetPortServer() int {
	return configReader.viperConfigReader.GetInt("apiServer.portServer")
}

func (configReader *Reader) GetApiPrefix() string {
	return configReader.viperConfigReader.GetString("apiServer.apiPref")
}

func (configReader *Reader) GetResourcePrefix() string {
	return configReader.viperConfigReader.GetString("resource.resourcePref")
}

func (configReader *Reader) GetResourceHost() string {
	return configReader.viperConfigReader.GetString("resource.resourceHost")
}

func (configReader *Reader) GetResourcePort() string {
	return configReader.viperConfigReader.GetString("resource.resourcePort")
}

func (configReader *Reader) GetConnectorPref() string {
	return configReader.viperConfigReader.GetString("connector.connectorPref")
}

func (configReader *Reader) GetConnectorHost() string {
	return configReader.viperConfigReader.GetString("connector.connectorHost")
}

func (configReader *Reader) GetConnectorPort() string {
	return configReader.viperConfigReader.GetString("connector.connectorPort")
}

func (configReader *Reader) GetAnalyticPref() string {
	return configReader.viperConfigReader.GetString("analytics.analyticPref")
}

func (configReader *Reader) GetAnalyticHost() string {
	return configReader.viperConfigReader.GetString("analytics.analyticHost")
}

func (configReader *Reader) GetAnalyticPort() string {
	return configReader.viperConfigReader.GetString("analytics.analyticPort")
}
