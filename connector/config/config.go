package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Connector Connector `yaml:"connector"`
	DB        DB        `yaml:"database"`
	Server    Server    `yaml:"server"`
}

type Connector struct {
	JiraUrl        string `yaml:"jira_url"`
	ThreadCount    int    `yaml:"thread_count"`
	IssueInRequest int    `yaml:"issue_in_request"`
	MaxTimeSleep   int    `yaml:"max_time_sleep"`
	MinTimeSleep   int    `yaml:"min_time_sleep"`
}

type DB struct {
	HostDB     string `yaml:"host_db"`
	PortDB     string `yaml:"port_db"`
	NameDB     string `yaml:"name_db"`
	PasswordDB string `yaml:"password_db"`
	UserDB     string `yaml:"user_db"`
}

type Server struct {
	ApiServer     ApiServer     `yaml:"api_server"`
	ConnectorHTTP ConnectorHTTP `yaml:"connector_http"`
}

type ApiServer struct {
	ApiServerHost string `yaml:"host_server"`
	ApiServerPort string `yaml:"port_server"`
	ApiPrefix     string `yaml:"api_pref"`
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
