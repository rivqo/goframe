package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Driver   string `yaml:"driver"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"database"`
	Auth struct {
		Secret   string        `yaml:"secret"`
		Duration time.Duration `yaml:"duration"`
	} `yaml:"auth"`
	RateLimit struct {
		Requests int           `yaml:"requests"`
		Period   time.Duration `yaml:"period"`
	} `yaml:"rateLimit"`
	App struct {
		Name string           `yaml:"name"`
		Version  string		  `yaml:"version"`
	} `yaml:"app"`
}

// Load reads the configuration from the provided file path
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

