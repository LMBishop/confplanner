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
	AcceptRegistrations bool   `yaml:"acceptRegistrations"`
	BaseURL             string `yaml:"baseURL"`
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
