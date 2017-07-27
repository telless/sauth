package main

import (
	"sauth/db"
	"net"
	"os"
	"fmt"
	"os/signal"
	"syscall"
	"sauth/utils"
	"sauth/models"
	"sauth/services"
	"strings"
	"errors"
	"time"
	"strconv"
)

var (
	soapUsersCache = make(map[string]services.SOAPUser)
	dbUsersCache   = make(map[string]models.User)
)

func serve(sockFile string) (error) {
	os.Remove(sockFile)
	listener, err := net.Listen("unix", sockFile)
	if err != nil {
		return err
	}
	utils.Log(fmt.Sprintf("Serve started on %s", sockFile))

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func(ln net.Listener, c chan os.Signal) {
		sig := <-c
		ln.Close()
		os.Remove(sockFile)
		utils.Die(fmt.Sprintf("Caught signal %s: shutting down.", sig))
	}(listener, signals)

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
	dbCheck()
	buf := make([]byte, 512)

	n, err := connection.Read(buf)
	utils.CheckError(err, "read inc message", utils.FatalLogLevel)
	requestData := string(buf[0:n])
	utils.Log(fmt.Sprintln("Server got:", requestData))
	responseData, err := handleRequest(requestData)
	if err != nil {
		utils.CheckError(err, "handle inc request "+requestData, utils.LowLogLevel)
		responseData = err.Error()
	}
	_, err = connection.Write([]byte(responseData))
	utils.CheckError(err, "write response "+responseData, utils.FatalLogLevel)

	utils.Log(fmt.Sprintln("Server send:", responseData))
}

func dbCheck() {
	err := db.GetConnection().Ping()
	if err != nil {
		db.RestartConnection()
	}
}

func handleRequest(requestData string) (string, error) {
	responseData, err := parseToken(requestData)
	if err != nil {
		return "", err
	}

	return responseData, nil
}

func parseToken(token string) (string, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(token, prefix) {
		return "", errors.New("Bearer token required for authorization")
	}

	return getUser(strings.TrimRight(token[len(prefix):], " \r\n\t"))
}

func getUser(parsedToken string) (string, error) {
	user := soapUsersCache[parsedToken]
	if user.WsoLogin == "" || user.ExpireTime.Before(time.Now()) {
		user = services.GetUserByToken(parsedToken, config.Soap)
		soapUsersCache[parsedToken] = user
	} else {
		logMsg := fmt.Sprintf("Get WSO user %s from cache, expire in %s seconds",
			user.WsoLogin,
			strconv.FormatFloat(user.ExpireTime.Sub(time.Now()).Seconds(), 'f', 0, 64))
		utils.Log(logMsg)
	}
	if user.ErrorMsg != "" {
		emptyUser := models.User{}
		return emptyUser.Serialize(), errors.New(fmt.Sprintf("WSO user not found, error: %s, token: %s", user.ErrorMsg, parsedToken))
	}

	dbUser := dbUsersCache[user.WsoLogin]
	if dbUser.Name == "" || dbUser.Expire.Before(time.Now()) {
		dbUser = db.GetUserByWSOLogin(user.WsoLogin)
		dbUser.Expire = time.Now().Add(time.Duration(config.DB.Expire) * time.Second)
		dbUsersCache[user.WsoLogin] = dbUser
	} else {
		logMsg := fmt.Sprintf("Get API user %s from cache, expire in %s seconds",
			user.WsoLogin,
			strconv.FormatFloat(dbUser.Expire.Sub(time.Now()).Seconds(), 'f', 0, 64))
		utils.Log(logMsg)
	}

	return dbUser.Serialize(), nil
}
