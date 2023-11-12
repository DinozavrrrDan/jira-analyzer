package configreader

import (
	"Jira-analyzer/jiraConnector/logger"

	"github.com/spf13/viper"
)

type ConfigRaeder struct {
	viperConfigReader *viper.Viper
	viperLogger       *logger.JiraLogger
}

func CreateNewConfigReader() *ConfigRaeder {
	configReader := ConfigRaeder{}
	configReader.viperConfigReader = viper.New()
	configReader.viperConfigReader.SetConfigName("config") //придумать имя конфига
	configReader.viperConfigReader.SetConfigType("yaml")
	configReader.viperConfigReader.AddConfigPath("analyzer/configs") //соответственно добавить путь
	viperLogger := logger.CreateNewLogger()
	if err := configReader.viperConfigReader.ReadInConfig(); err != nil {
		viperLogger.Log(logger.ERROR, err.Error())
	}
	return &configReader
}

func (configReader *ConfigRaeder) GetJiraUrl() string {
	return configReader.viperConfigReader.GetString("jiraUrl")
}
