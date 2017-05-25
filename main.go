package main

import (
	"flag"
	"errors"
	"sauth/utils"
	"sauth/configuration"
	"sauth/soapservice"
	"os"
	"log"
	"sauth/response"
	"sauth/db"
)

func main() {
	setLogFile()

	token, err := getToken()
	utils.CheckError(err)

	connection := db.GetConnection()
	defer connection.Close()

	config := configuration.GetConfig()
	wsoLogin, err := soapservice.GetUserByToken(token, config.Soap)
	println(response.GenerateJsonResponse(wsoLogin, err))
}

func getToken() (string, error) {
	tokenPtr := flag.String("token", "default", "user token")
	flag.Parse()
	if *tokenPtr == "default" {
		return "", errors.New("token: default token is not acceptable")
	}
	return *tokenPtr, nil
}

func setLogFile() error {
	f, err := os.OpenFile("logs/sauth.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		return err
	}
	log.SetOutput(f)
	return nil
}
