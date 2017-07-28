package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"time"
)

const defaultDbConnectionTimeout = 5

var (
	connection *sql.DB
	dbTimeout  = defaultDbConnectionTimeout // timeout in minutes
)

func getConnection() *sql.DB {
	if connection == nil {
		openConnection(config.Db)
	}
	dbTimeout = defaultDbConnectionTimeout

	return connection
}

func formatDbConnectString(dbConfig dbConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name)
}

func openConnection(config dbConfig) {
	conn, err := sql.Open("mysql", formatDbConnectString(config))
	checkError(err, "DB connection", fatalLogLevel)
	connection = conn
	logMessage("DB connection open")

	go func() {
		for {
			time.Sleep(time.Minute)
			dbTimeout--
			if dbTimeout == 0 {
				connection.Close()
				connection = nil
				break
				logMessage("DB connection closed by timeout")
			}
		}
	}()
}
