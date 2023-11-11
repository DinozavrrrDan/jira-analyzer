package configreader

import (
	"github.com/spf13/viper"
)

type ConfigRaeder struct {
	viperConfigReader *viper.Viper
}

func newConfigReader() *ConfigRaeder {
	configReader := ConfigRaeder{}
	configReader.viperConfigReader = viper.New()
	configReader.viperConfigReader.SetConfigName("") //придумать имя конфига
	configReader.viperConfigReader.SetConfigType("yaml")
	configReader.viperConfigReader.AddConfigPath("") //соответственно добавить путь

	return &configReader
}
