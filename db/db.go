package db

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"sauth/configuration"
	"sauth/utils"
)

var connection *sql.DB

func OpenConnection(config configuration.DBConfig) {
	conn, err := sql.Open("mysql", formatDbConnectString(config))
	utils.CheckError(err)

	connection = conn
}

func GetConnection() *sql.DB {
	return connection
}

func formatDbConnectString(dbConfig configuration.DBConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name)
}
