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
	soapUser := soapservice.GetUserByToken(token, config.Soap)
	if soapUser.ErrorMsg == "" {
		user := db.GetUserByWSOLogin(soapUser.WsoLogin)
		if user.Name == "" {
			println(response.GenerateErrorJson(soapUser.WsoLogin + " not found"))
		} else {
			println(response.GenerateSuccessJson(user))
		}
	} else {
		println(response.GenerateErrorJson(soapUser.ErrorMsg))
	}
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
