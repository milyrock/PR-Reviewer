package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `yaml:"postgres"`
}

type DatabaseConfig struct {
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Network  string `yaml:"network"`
	Port     string `yaml:"port"`
}

func ReadConfig(path string) (*Config, error) {
	var config Config

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(file))
	expanded_yaml := os.ExpandEnv(string(file))

	err = yaml.Unmarshal([]byte(expanded_yaml), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
