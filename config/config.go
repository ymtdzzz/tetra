package config

import "github.com/BurntSushi/toml"

const (
	DBTypeMySQL DBType = "mysql"
)

type DBType string

type DBConfig struct {
	Type     DBType
	Host     string
	Port     int
	User     string
	Password string
}

type Config map[string]DBConfig

func LoadConfig(path string) (Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
