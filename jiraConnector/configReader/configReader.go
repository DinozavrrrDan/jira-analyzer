package configReader

import (
	"Jira-analyzer/jiraConnector/logger"

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

func (configReader *ConfigRaeder) GetHostDB() string {
	return configReader.viperConfigReader.GetString("hostDB")
}

func (configReader *ConfigRaeder) GetDatabaseName() string {
	return configReader.viperConfigReader.GetString("databaseName")
}

func (configReader *ConfigRaeder) GetPasswordDB() string {
	return configReader.viperConfigReader.GetString("passwordDB")
}

func (configReader *ConfigRaeder) GetUserDb() string {
	return configReader.viperConfigReader.GetString("userDb")
}

func (configReader *ConfigRaeder) GetPortDB() int {
	return configReader.viperConfigReader.GetInt("portDB")
}

func (configReader *ConfigRaeder) GetThreadCount() int {
	return configReader.viperConfigReader.GetInt("threadCount")
}

func (configReader *ConfigRaeder) GetIssusOnOneRequest() int {
	return configReader.viperConfigReader.GetInt("issueInOneRequest")
}

func (configReader *ConfigRaeder) GetMaxTimeSleep() int {
	return configReader.viperConfigReader.GetInt("maxTimeSleep")
}

func (configReader *ConfigRaeder) GetMinTimeSleep() int {
	return configReader.viperConfigReader.GetInt("maxTimeSleep")
}

func (configReader *ConfigRaeder) GetHostServer() string {
	return configReader.viperConfigReader.GetString("hostServer")
}

func (configReader *ConfigRaeder) GetPortServer() int {
	return configReader.viperConfigReader.GetInt("portServer")
}
