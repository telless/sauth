package configuration

import (
	"encoding/json"
	"sauth/utils"
)

const config_path string = "parameters.json"

type Config struct {
	DB   DBConfig `json:"db"`
	Soap SoapConfig `json:"soap"`
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     int `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Expire   int `json:"cache_expire"`
}
type SoapConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func GetConfig() Config {
	config := Config{}
	file := utils.GetFileContent(config_path)
	json.Unmarshal(file, &config)
	return config
}
