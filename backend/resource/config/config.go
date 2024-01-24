package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	DB     DB     `yaml:"database"`
	Server Server `yaml:"handler"`
}

type DB struct {
	HostDB     string `yaml:"host_db"`
	PortDB     string `yaml:"port_db"`
	NameDB     string `yaml:"name_db"`
	PasswordDB string `yaml:"password_db"`
	UserDB     string `yaml:"user_db"`
}

type Server struct {
	ApiServer    ApiServer    `yaml:"api_server"`
	ResourceHTTP ResourceHTTP `yaml:"resource_http"`
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
