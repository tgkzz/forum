package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Host string `json:"Host"`
	Port string `json:"Port"`
	DB   DB
}

type DB struct {
	DriverName     string `json:"DriverName"`
	DataSourceName string `json:"DataSourceName"`
	Sql            string `json:"Sql"`
}

func OpenConfig() (Config, error) {
	configFile, err := os.Open("./config/config.json")
	if err != nil {
		return Config{}, err
	}
	defer configFile.Close()
	var config Config
	if err = json.NewDecoder(configFile).Decode(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}
