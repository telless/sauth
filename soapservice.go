package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"net"
	"net/http"
	"time"
	"encoding/json"
	"strings"
)

func (b *soapBody) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
				b.Fault = &soapFault{}
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

// Inc token struct
type authToken struct {
	tokenType        string
	tokenCredentials string
}

// Request struct

type soapUser struct {
	WsoLogin   string `json:"wso_login"`
	ExpireTime time.Time `json:"expire_time"`
	Token      string `json:"token"`
	Valid      bool `json:"valid"`
	ErrorMsg   string `json:"error_msg"`
}

func (user *soapUser) Serialize() string {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return ""
	}

	return string(jsonData)
}

type findOAuthConsumerIfTokenIsValid struct {
	XMLName xml.Name `xml:"http://org.apache.axis2/xsd findOAuthConsumerIfTokenIsValid"`

	ValidationReqDTO *oAuth2TokenValidationRequestDTO `xml:"validationReqDTO,omitempty"`
}

type oAuth2TokenValidationRequestDTO struct {
	XMLName xml.Name `xml:"http://dto.oauth2.identity.carbon.wso2.org/xsd validationReqDTO"`

	AccessToken       *oAuth2TokenValidationRequestDTOOAuth2AccessToken             `xml:"accessToken,omitempty"`
	Context           []*oAuth2TokenValidationRequestDTOTokenValidationContextParam `xml:"context,omitempty"`
	RequiredClaimURIs []string                                                      `xml:"requiredClaimURIs,omitempty"`
}

type oAuth2TokenValidationRequestDTOOAuth2AccessToken struct {
	XMLName xml.Name `xml:"http://dto.oauth2.identity.carbon.wso2.org/xsd accessToken"`

	Identifier string `xml:"identifier,omitempty"`
	TokenType  string `xml:"tokenType,omitempty"`
}

// Response struct

type oAuth2TokenValidationRequestDTOTokenValidationContextParam struct {
	XMLName xml.Name `xml:"http://dto.oauth2.identity.carbon.wso2.org/xsd OAuth2TokenValidationRequestDTO_TokenValidationContextParam"`

	Key   string `xml:"key,omitempty"`
	Value string `xml:"value,omitempty"`
}

type findOAuthConsumerIfTokenIsValidResponse struct {
	XMLName xml.Name `xml:"http://org.apache.axis2/xsd findOAuthConsumerIfTokenIsValidResponse"`

	Return_ *oAuth2ClientApplicationDTO `xml:"return,omitempty"`
}

type oAuth2ClientApplicationDTO struct {
	XMLName xml.Name `xml:"http://org.apache.axis2/xsd return"`

	AccessTokenValidationResponse *oAuth2TokenValidationResponseDTO `xml:"accessTokenValidationResponse,omitempty"`
	ConsumerKey                   string                            `xml:"consumerKey,omitempty"`
}

type oAuth2TokenValidationResponseDTO struct {
	XMLName xml.Name `xml:"http://dto.oauth2.identity.carbon.wso2.org/xsd accessTokenValidationResponse"`

	AuthorizationContextToken *oAuth2TokenValidationResponseDTOAuthorizationContextToken `xml:"authorizationContextToken,omitempty"`
	AuthorizedUser            string                                                     `xml:"authorizedUser,omitempty"`
	ErrorMsg                  string                                                     `xml:"errorMsg,omitempty"`
	ExpiryTime                int64                                                      `xml:"expiryTime,omitempty"`
	Scope                     []string                                                   `xml:"scope,omitempty"`
	Valid                     bool                                                       `xml:"valid,omitempty"`
}

// Client struct

type oAuth2TokenValidationResponseDTOAuthorizationContextToken struct {
	XMLName xml.Name `xml:"http://dto.oauth2.identity.carbon.wso2.org/xsd authorizationContextToken"`

	TokenString string `xml:"tokenString,omitempty"`
	TokenType   string `xml:"tokenType,omitempty"`
}

type oAuth2TokenValidationServicePortType struct {
	client *soapClient
}

type soapClient struct {
	url  string
	tls  bool
	auth *basicAuth
}

type basicAuth struct {
	Login    string
	Password string
}

type soapEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`

	Body soapBody
}

type soapBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault   *soapFault  `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

type soapFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string `xml:"faultcode,omitempty"`
	String string `xml:"faultstring,omitempty"`
	Actor  string `xml:"faultactor,omitempty"`
	Detail string `xml:"detail,omitempty"`
}

var soapTimeout = time.Duration(30 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, soapTimeout)
}

func (f *soapFault) Error() string {
	return f.String
}

func newSoapClient(url string, tls bool, auth *basicAuth) *soapClient {
	return &soapClient{
		url:  url,
		tls:  tls,
		auth: auth,
	}
}

func (s *soapClient) call(soapAction string, request, response interface{}) error {
	envelope := soapEnvelope{}

	envelope.Body.Content = request
	buffer := new(bytes.Buffer)

	encoder := xml.NewEncoder(buffer)

	if err := encoder.Encode(envelope); err != nil {
		return err
	}

	if err := encoder.Flush(); err != nil {
		return err
	}

	logMessage("WSO request: " + buffer.String())

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

	req.Header.Set("apiUser-Agent", "sauth/1.0")
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
		return nil
	}

	logMessage("WSO response: " + string(rawBody))
	respEnvelope := new(soapEnvelope)
	respEnvelope.Body = soapBody{Content: response}
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

func getUserByToken(tokenString string, config soapConfig) (soapUser) {
	token:=prepareToken(tokenString)
	auth := basicAuth{Login: config.User, Password: config.Password}
	soapClient := newSoapClient(config.Host, true, &auth)
	request := findOAuthConsumerIfTokenIsValid{
		XMLName: xml.Name{},
		ValidationReqDTO: &oAuth2TokenValidationRequestDTO{
			XMLName: xml.Name{},
			AccessToken: &oAuth2TokenValidationRequestDTOOAuth2AccessToken{
				XMLName:    xml.Name{},
				Identifier: token.tokenCredentials,
				TokenType:  strings.ToLower(token.tokenType)}}}
	wsoResponse := findOAuthConsumerIfTokenIsValidResponse{}
	err := soapClient.call("findOAuthConsumerIfTokenIsValid", &request, &wsoResponse)
	checkError(err, "soap call error", fatalLogLevel)
	validationResponse := wsoResponse.Return_.AccessTokenValidationResponse
	return soapUser{
		WsoLogin:   validationResponse.AuthorizedUser,
		ExpireTime: time.Now().Add(time.Duration(validationResponse.ExpiryTime) * time.Second),
		Token:      tokenString,
		Valid:      validationResponse.Valid,
		ErrorMsg:   validationResponse.ErrorMsg,
	}
}

func prepareToken(tokenString string) authToken {
	data := strings.Split(tokenString, " ");

	return authToken{
		tokenType: data[0],
		tokenCredentials: data[1]}
}
