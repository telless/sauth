package main

import (
	"encoding/json"
)

const config_path string = "parameters.json"

type baseConfig struct {
	Db   dbConfig `json:"db"`
	Soap soapConfig `json:"soap"`
}

type dbConfig struct {
	Host     string `json:"host"`
	Port     int `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Expire   int `json:"cache_expire"`
}
type soapConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func getConfig() baseConfig {
	config := baseConfig{}
	file := getFileContent(config_path)
	json.Unmarshal(file, &config)

	return config
}
