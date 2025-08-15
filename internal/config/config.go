package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		ConnString string `yaml:"connString"`
	} `yaml:"database"`
	Conference struct {
		ScheduleURL string `yaml:"scheduleURL"`
	} `yaml:"conference"`
	Auth struct {
		EnableBasicAuth bool           `yaml:"enableBasicAuth"`
		AuthProviders   []AuthProvider `yaml:"authProviders"`
	}
	AcceptRegistrations bool   `yaml:"acceptRegistrations"`
	BaseURL             string `yaml:"baseURL"`
}

type AuthProvider struct {
	Identifier               string   `yaml:"identifier"`
	Name                     string   `yaml:"name"`
	ClientID                 string   `yaml:"clientID"`
	ClientSecret             string   `yaml:"clientSecret"`
	Endpoint                 string   `yaml:"endpoint"`
	LoginFilter              string   `yaml:"loginFilter"`
	LoginFilterAllowedValues []string `yaml:"loginFilterAllowedValues"`
	UserSyncFilter           string   `yaml:"userSyncFilter"`
}

func ReadConfig(configPath string, dst *Config) error {
	config, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(config, dst); err != nil {
		return err
	}
	return nil
}
