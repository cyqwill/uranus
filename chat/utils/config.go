package utils

import (
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Port string
	Addr string
}

type DbConfig struct {
	DbType     string `toml:"db_type"`
	DbUserName string `toml:"db_username"`
	DbPassword string `toml:"db_password"`
	DbName     string `toml:"db_name"`
}

// more config to go
type AppConfig struct {
	Server   ServerConfig `toml:"server"`
	Database DbConfig `toml:"db"`
}

// This func should do only once in whole program
var (
	cfg  *AppConfig
	once sync.Once
)

func Config(c string) *AppConfig {
	once.Do(func() {
		filePath, err := filepath.Abs(c)
		if err != nil {
			panic(err)
		}
		log.Infof("parse toml file once. filePath: %s.", filePath)
		if _, err := toml.DecodeFile(filePath, &cfg); err != nil {
			panic(err)
		}
	})
	return cfg
}
