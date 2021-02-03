package config

import (
	"encoding/json"
	"fmt"

	"github/demo/env"
	"github/demo/utils/log"
)

type Logger struct {
	Env      string `json:"env"`
	Filename string `json:"filename"`
	Level    string `json:"level"`
}

type Database struct {
	Dialect  string `json:"dialect"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Config struct {
	Logger   *Logger   `json:"logger"`
	Database *Database `json:"database"`
}

func (c *Config) Init(v env.Variables) bool {
	for _, key := range v.Keys() {
		switch key {
		case env.LogEnv:
			c.Logger.Env = fmt.Sprintf("%v", v[env.LogEnv])
		case env.LogFile:
			c.Logger.Filename = fmt.Sprintf("%v", v[env.LogFile])
		case env.LogLevel:
			c.Logger.Level = fmt.Sprintf("%v", v[env.LogLevel])
		case env.DBDialect:
			c.Database.Dialect = fmt.Sprintf("%v", v[env.DBDialect])
		case env.DBHost:
			c.Database.Host = fmt.Sprintf("%v", v[env.DBHost])
		case env.DBPort:
			c.Database.Port = fmt.Sprintf("%v", v[env.DBPort])
		case env.DBName:
			c.Database.Name = fmt.Sprintf("%v", v[env.DBName])
		case env.DBUser:
			c.Database.User = fmt.Sprintf("%v", v[env.DBUser])
		case env.DBPassword:
			c.Database.Password = fmt.Sprintf("%v", v[env.DBPassword])
		}
	}

	return true
}

func (c *Config) Watch() (string, error) {
	cJson, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Infof("config watch:\n%s", string(cJson))
	return string(cJson), nil
}

func NewConfig() *Config {
	c := &Config{
		Logger: &Logger{
			Env:   "development",
			Level: "debug",
		},
		Database: &Database{},
	}

	return c
}
