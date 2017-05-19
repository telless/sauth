package configuration

import (
	"encoding/json"
	"sauth/utils"
)

const config_path string = "configuration/parameters.json"

type Config struct {
	DB   DBConfig
	Soap SoapConfig
}

type DBConfig struct {
	Host     string `json:"db.host"`
	Port     int `json:"db.port"`
	Name     string `json:"db.name"`
	User     string `json:"db.user"`
	Password string `json:"db.password"`
}
type SoapConfig struct {
	Host     string `json:"soap.host"`
	User     string `json:"soap.user"`
	Password string `json:"soap.password"`
}

func GetConfig() Config {
	config := Config{}
	file := utils.GetFileContent(config_path)

	dbConfig := DBConfig{}
	json.Unmarshal(file, &dbConfig)
	config.DB = dbConfig

	soapConfig := SoapConfig{}
	json.Unmarshal(file, &soapConfig)
	config.Soap = soapConfig

	return config
}
