package soapservice

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"net"
	"net/http"
	"time"
	"log"
	"sauth/configuration"
	"sauth/utils"
	"errors"
)

// Request struct

type FindOAuthConsumerIfTokenIsValid struct {
	XMLName xml.Name `xml:"http://org.apache.axis2/xsd findOAuthConsumerIfTokenIsValid"`

	ValidationReqDTO *OAuth2TokenValidationRequestDTO `xml:"validationReqDTO,omitempty"`
}

type OAuth2TokenValidationRequestDTO struct {
	XMLName xml.Name `xml:"http://dto.oauth2.identity.carbon.wso2.org/xsd validationReqDTO"`

	AccessToken       *OAuth2TokenValidationRequestDTOOAuth2AccessToken             `xml:"accessToken,omitempty"`
	Context           []*OAuth2TokenValidationRequestDTOTokenValidationContextParam `xml:"context,omitempty"`
	RequiredClaimURIs []string                                                      `xml:"requiredClaimURIs,omitempty"`
}

type OAuth2TokenValidationRequestDTOOAuth2AccessToken struct {
	XMLName xml.Name `xml:"http://dto.oauth2.identity.carbon.wso2.org/xsd accessToken"`

	Identifier string `xml:"identifier,omitempty"`
	TokenType  string `xml:"tokenType,omitempty"`
}

type OAuth2TokenValidationRequestDTOTokenValidationContextParam struct {
	XMLName xml.Name `xml:"http://dto.oauth2.identity.carbon.wso2.org/xsd OAuth2TokenValidationRequestDTO_TokenValidationContextParam"`

	Key   string `xml:"key,omitempty"`
	Value string `xml:"value,omitempty"`
}

// Response struct

type FindOAuthConsumerIfTokenIsValidResponse struct {
	XMLName xml.Name `xml:"http://org.apache.axis2/xsd findOAuthConsumerIfTokenIsValidResponse"`

	Return_ *OAuth2ClientApplicationDTO `xml:"return,omitempty"`
}

type OAuth2ClientApplicationDTO struct {
	XMLName xml.Name `xml:"http://org.apache.axis2/xsd return"`

	AccessTokenValidationResponse *OAuth2TokenValidationResponseDTO `xml:"accessTokenValidationResponse,omitempty"`
	ConsumerKey                   string                            `xml:"consumerKey,omitempty"`
}

type OAuth2TokenValidationResponseDTO struct {
	XMLName xml.Name `xml:"http://dto.oauth2.identity.carbon.wso2.org/xsd accessTokenValidationResponse"`

	AuthorizationContextToken *OAuth2TokenValidationResponseDTOAuthorizationContextToken `xml:"authorizationContextToken,omitempty"`
	AuthorizedUser            string                                                     `xml:"authorizedUser,omitempty"`
	ErrorMsg                  string                                                     `xml:"errorMsg,omitempty"`
	ExpiryTime                int64                                                      `xml:"expiryTime,omitempty"`
	Scope                     []string                                                   `xml:"scope,omitempty"`
	Valid                     bool                                                       `xml:"valid,omitempty"`
}

type OAuth2TokenValidationResponseDTOAuthorizationContextToken struct {
	XMLName xml.Name `xml:"http://dto.oauth2.identity.carbon.wso2.org/xsd authorizationContextToken"`

	TokenString string `xml:"tokenString,omitempty"`
	TokenType   string `xml:"tokenType,omitempty"`
}

// Client struct

type OAuth2TokenValidationServicePortType struct {
	client *SOAPClient
}

type SOAPClient struct {
	url  string
	tls  bool
	auth *BasicAuth
}

type BasicAuth struct {
	Login    string
	Password string
}

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`

	Body SOAPBody
}

type SOAPBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault   *SOAPFault  `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string `xml:"faultcode,omitempty"`
	String string `xml:"faultstring,omitempty"`
	Actor  string `xml:"faultactor,omitempty"`
	Detail string `xml:"detail,omitempty"`
}

var timeout = time.Duration(30 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

func (b *SOAPBody) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if b.Content == nil {
		return xml.UnmarshalError("Content must be a pointer to a struct")
	}

	var (
		token    xml.Token
		err      error
		consumed bool
	)

Loop:
	for {
		if token, err = d.Token(); err != nil {
			return err
		}

		if token == nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			if consumed {
				return xml.UnmarshalError("Found multiple elements inside SOAP body; not wrapped-document/literal WS-I compliant")
			} else if se.Name.Space == "http://schemas.xmlsoap.org/soap/envelope/" && se.Name.Local == "Fault" {
				b.Fault = &SOAPFault{}
				b.Content = nil

				err = d.DecodeElement(b.Fault, &se)
				if err != nil {
					return err
				}

				consumed = true
			} else {
				if err = d.DecodeElement(b.Content, &se); err != nil {
					return err
				}

				consumed = true
			}
		case xml.EndElement:
			break Loop
		}
	}

	return nil
}

func (f *SOAPFault) Error() string {
	return f.String
}

func NewSOAPClient(url string, tls bool, auth *BasicAuth) *SOAPClient {
	return &SOAPClient{
		url:  url,
		tls:  tls,
		auth: auth,
	}
}

func (s *SOAPClient) Call(soapAction string, request, response interface{}) error {
	envelope := SOAPEnvelope{}

	envelope.Body.Content = request
	buffer := new(bytes.Buffer)

	encoder := xml.NewEncoder(buffer)

	if err := encoder.Encode(envelope); err != nil {
		return err
	}

	if err := encoder.Flush(); err != nil {
		return err
	}

	log.Println("request: " + buffer.String())

	req, err := http.NewRequest("POST", s.url, buffer)
	if err != nil {
		return err
	}
	if s.auth != nil {
		req.SetBasicAuth(s.auth.Login, s.auth.Password)
	}

	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	if soapAction != "" {
		req.Header.Add("SOAPAction", soapAction)
	}

	req.Header.Set("User-Agent", "gowsdl/0.1")
	req.Close = true

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: s.tls,
		},
		Dial: dialTimeout,
	}

	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rawBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(rawBody) == 0 {
		log.Println("empty response")
		return nil
	}

	log.Println("response" + string(rawBody))
	respEnvelope := new(SOAPEnvelope)
	respEnvelope.Body = SOAPBody{Content: response}
	err = xml.Unmarshal(rawBody, respEnvelope)
	if err != nil {
		return err
	}

	fault := respEnvelope.Body.Fault
	if fault != nil {
		return fault
	}

	return nil
}

func GetUserByToken(token string, config configuration.SoapConfig) (string, error) {
	auth := BasicAuth{Login: config.User, Password: config.Password}
	soapClient := NewSOAPClient(config.Host, true, &auth)
	request := FindOAuthConsumerIfTokenIsValid{
		XMLName: xml.Name{},
		ValidationReqDTO: &OAuth2TokenValidationRequestDTO{
			XMLName: xml.Name{},
			AccessToken: &OAuth2TokenValidationRequestDTOOAuth2AccessToken{
				XMLName:    xml.Name{},
				Identifier: token,
				TokenType:  "bearer"}}}
	wsoResponse := FindOAuthConsumerIfTokenIsValidResponse{}
	err := soapClient.Call("FindOAuthConsumerIfTokenIsValid", &request, &wsoResponse)
	utils.CheckError(err)
	validationResponse := wsoResponse.Return_.AccessTokenValidationResponse
	if validationResponse.ErrorMsg != "" {
		return "", errors.New(validationResponse.ErrorMsg)
	}
	return validationResponse.AuthorizedUser, nil
}
