package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Timeout Timeout `yaml:"timeout"`
	Server  Server  `yaml:"handler"`
	Gateway Gateway `yaml:"gateway"`
}

type Gateway struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Server struct {
	ApiServer     ApiServer     `yaml:"api_server"`
	ResourceHTTP  ResourceHTTP  `yaml:"resource_http"`
	AnalyticsHTTP AnalyticsHTTP `yaml:"analytics_http"`
	ConnectorHTTP ConnectorHTTP `yaml:"connector_http"`
}

type Timeout struct {
	Resource  int64 `yaml:"resource"`
	Analytics int64 `yaml:"analytics"`
}

type ApiServer struct {
	ApiServerHost string `yaml:"host_server"`
	ApiServerPort string `yaml:"port_server"`
	ApiPrefix     string `yaml:"api_pref"`
}

type ResourceHTTP struct {
	ResourceHost   string `yaml:"resource_host"`
	ResourcePort   string `yaml:"resource_port"`
	ResourcePrefix string `yaml:"resource_pref"`
}

type AnalyticsHTTP struct {
	AnalyticsHost   string `yaml:"analytics_host"`
	AnalyticsPort   string `yaml:"analytics_port"`
	AnalyticsPrefix string `yaml:"analytics_pref"`
}

type ConnectorHTTP struct {
	ConnectorHost   string `yaml:"connector_host"`
	ConnectorPort   string `yaml:"connector_port"`
	ConnectorPrefix string `yaml:"connector_pref"`
}

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	f, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error while reading configuration: %w", err)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
