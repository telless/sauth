package main

import (
	"net"
	"fmt"
	"time"
	"strconv"
	"errors"
	"strings"
)

var (
	soapUsersCache = make(map[string]soapUser)
	dbUsersCache   = make(map[string]apiUser)
)

func startServer(listener net.Listener) (error) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			return err
		}

		go handler(connection)
	}

	return nil;
}

func handler(connection net.Conn) {
	defer connection.Close()

	buf := make([]byte, 512)

	n, err := connection.Read(buf)
	checkError(err, "Read inc message", fatalLogLevel)
	requestData := string(buf[0:n])
	logMessage(fmt.Sprintln("Server got:", requestData))
	responseData, err := handleRequest(strings.TrimRight(requestData, " \r\n\t"))
	if err != nil {
		checkError(err, "Handle inc request "+requestData, minorLogLevel)
		responseData = err.Error()
	}
	_, err = connection.Write([]byte(responseData))
	checkError(err, "Write response "+responseData, fatalLogLevel)

	logMessage(fmt.Sprintln("Server send:", responseData))
}

func handleRequest(parsedToken string) (string, error) {
	soapUser := soapUsersCache[parsedToken]
	if soapUser.WsoLogin == "" || soapUser.ExpireTime.Before(time.Now()) {
		soapUser = getUserByToken(parsedToken, config.Soap)
		soapUsersCache[parsedToken] = soapUser
	} else {
		logMessage(fmt.Sprintf("Get WSO soapUser %s from cache, expire in %s seconds",
			soapUser.WsoLogin,
			strconv.FormatFloat(soapUser.ExpireTime.Sub(time.Now()).Seconds(), 'f', 0, 64)))
	}
	if soapUser.ErrorMsg != "" {
		return "", errors.New(fmt.Sprintf("WSO soapUser not found, error: %s, token: %s", soapUser.ErrorMsg, parsedToken))
	}
	if soapUser.ExpireTime.Before(time.Now()) {
		return "", errors.New(fmt.Sprintf("WSO token %s expired", parsedToken))
	}
	apiUser := dbUsersCache[soapUser.WsoLogin]
	if apiUser.Name == "" || apiUser.Expire.Before(time.Now()) {
		var err error = nil
		apiUser, err = getUserByWSOLogin(soapUser.WsoLogin)
		if err != nil {
			return "", errors.New(fmt.Sprintf("User %s not found in DB", soapUser.WsoLogin))
		}
		apiUser.Expire = time.Now().Add(time.Duration(config.Db.Expire) * time.Second)
		apiUser.TokenExpire = soapUser.ExpireTime
		dbUsersCache[soapUser.WsoLogin] = apiUser
	} else {
		logMessage(fmt.Sprintf("Get API soapUser %s from cache, expire in %s seconds",
			soapUser.WsoLogin,
			strconv.FormatFloat(apiUser.Expire.Sub(time.Now()).Seconds(), 'f', 0, 64)))
	}
	return apiUser.Serialize(), nil
}
