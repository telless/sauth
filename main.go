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
	println(generateJsonResponse(validationResponse))
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

func generateJsonResponse(dto *soapservice.OAuth2TokenValidationResponseDTO) string {
	if dto.ErrorMsg != "" {
		return generateErrorJson(dto.ErrorMsg)
	} else {
		user := models.GetUserByWSOLogin(dto.AuthorizedUser)
		if user.Name == "" {
			return generateErrorJson(dto.AuthorizedUser + " not found")
		}
		return generateSuccessJson(user)
	}
}

func generateErrorJson(errorMsg string) string {
	errorResponse := response.NewErrorResponse("fail", errorMsg)
	return errorResponse.Serialize()
}

func generateSuccessJson(user models.User) string {
	successResponse := response.NewSuccessResponse("success", user)
	return successResponse.Serialize()
}
