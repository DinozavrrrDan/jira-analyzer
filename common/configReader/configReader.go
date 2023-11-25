package configReader

import (
	"Jira-analyzer/common/logger"

	"github.com/spf13/viper"
)

/*
Конфигурационный файл должен содержать следующую информацию:
● Настройки базы данных (Жду инфе от Бори)
	o имя пользователя PostgreSQL
	o пароль пользователя PostgreSQL
	o хост PostgreSQL
	o порт PostgreSQL
	o название базы данных
● Настройки модуля
	o Порт для локального http сервера
	o URL сайта JIRA
	o Количество Issues, загружаемые из Jira за 1 запрос
	o Количество потоков при выгрузке из Jira
	o Максимальное время ожидания для повторной отправки запроса в Jira
	o Начальное время ожидания для повторной отправки запроса в Jira
*/

type ConfigReader struct {
	viperConfigReader *viper.Viper
	viperLogger       *logger.JiraLogger
}

func CreateNewConfigReader() *ConfigReader {
	configReader := ConfigReader{}
	configReader.viperConfigReader = viper.New()
	configReader.viperConfigReader.SetConfigName("config")
	configReader.viperConfigReader.SetConfigType("yaml")
	configReader.viperConfigReader.AddConfigPath("././configs")
	configReader.viperConfigReader.AddConfigPath("../../configs")
	viperLogger := logger.CreateNewLogger()
	if err := configReader.viperConfigReader.ReadInConfig(); err != nil {
		viperLogger.Log(logger.ERROR, "error in configReader"+err.Error())
	}
	return &configReader
}

func (configReader *ConfigReader) GetJiraUrl() string {
	return configReader.viperConfigReader.GetString("jiraUrl")
}

func (configReader *ConfigReader) GetHostDB() string {
	return configReader.viperConfigReader.GetString("hostDB")
}

func (configReader *ConfigReader) GetDatabaseName() string {
	return configReader.viperConfigReader.GetString("NameDB")
}

func (configReader *ConfigReader) GetPasswordDB() string {
	return configReader.viperConfigReader.GetString("passwordDB")
}

func (configReader *ConfigReader) GetUserDb() string {
	return configReader.viperConfigReader.GetString("userDB")
}

func (configReader *ConfigReader) GetPortDB() int {
	return configReader.viperConfigReader.GetInt("portDB")
}

func (configReader *ConfigReader) GetThreadCount() int {
	return configReader.viperConfigReader.GetInt("threadCount")
}

func (configReader *ConfigReader) GetIssusOnOneRequest() int {
	return configReader.viperConfigReader.GetInt("issueInOneRequest")
}

func (configReader *ConfigReader) GetMaxTimeSleep() int {
	return configReader.viperConfigReader.GetInt("maxTimeSleep")
}

func (configReader *ConfigReader) GetMinTimeSleep() int {
	return configReader.viperConfigReader.GetInt("maxTimeSleep")
}

func (configReader *ConfigReader) GetHostServer() string {
	return configReader.viperConfigReader.GetString("hostServer")
}

func (configReader *ConfigReader) GetPortServer() int {
	return configReader.viperConfigReader.GetInt("portServer")
}

func (configReader *ConfigReader) GetApiPrefix() string {
	return configReader.viperConfigReader.GetString("apiPref")
}

func (configReader *ConfigReader) GetResourcePrefix() string {
	return configReader.viperConfigReader.GetString("resourcePref")
}

func (configReader *ConfigReader) GetResourceHost() string {
	return configReader.viperConfigReader.GetString("resourceHost")
}

func (configReader *ConfigReader) GetResourcePort() string {
	return configReader.viperConfigReader.GetString("resourcePort")
}

func (configReader *ConfigReader) GetConnectorHost() string {
	return configReader.viperConfigReader.GetString("connectorHost")
}

func (configReader *ConfigReader) GetConnectorPort() string {
	return configReader.viperConfigReader.GetString("connectorPort")
}
