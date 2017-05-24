package main

import (
	"flag"
	"errors"
	"sauth/utils"
	"sauth/configuration"
	"sauth/soapservice"
	"encoding/xml"
	"os"
	"log"
	"sauth/response"
	"sauth/models"
)

func main() {
	token, err := getToken()
	utils.CheckError(err)
	setLogFile()

	config := configuration.GetConfig()
	auth := soapservice.BasicAuth{Login: config.Soap.User, Password: config.Soap.Password}
	soapClient := soapservice.NewSOAPClient(config.Soap.Host, true, &auth)
	request := soapservice.FindOAuthConsumerIfTokenIsValid{
		XMLName: xml.Name{},
		ValidationReqDTO: &soapservice.OAuth2TokenValidationRequestDTO{
			XMLName: xml.Name{},
			AccessToken: &soapservice.OAuth2TokenValidationRequestDTOOAuth2AccessToken{
				XMLName:    xml.Name{},
				Identifier: token,
				TokenType:  "bearer"}}}
	wsoResponse := soapservice.FindOAuthConsumerIfTokenIsValidResponse{}
	err = soapClient.Call("FindOAuthConsumerIfTokenIsValid", &request, &wsoResponse)
	utils.CheckError(err)

	validationResponse := wsoResponse.Return_.AccessTokenValidationResponse
	if validationResponse.ErrorMsg != "" {
		println(generateErrorJson(validationResponse))
	} else {
		println(generateSuccessJson(validationResponse))
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

func generateErrorJson(dto *soapservice.OAuth2TokenValidationResponseDTO) string {
	errorResponse := response.NewErrorResponse("fail", dto.ErrorMsg)
	return errorResponse.Serialize()
}

func generateSuccessJson(dto *soapservice.OAuth2TokenValidationResponseDTO) string {
	user := models.GetUserByWSOLogin(dto.AuthorizedUser)
	if user.Name == "" {
		errorResponse := response.NewErrorResponse("fail", dto.AuthorizedUser+" not found")
		return errorResponse.Serialize()
	}

	successResponse := response.NewSuccessResponse("success", user)
	return successResponse.Serialize()
}
