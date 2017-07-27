package db

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"sauth/configuration"
	"sauth/utils"
)

var connection *sql.DB

func GetConnection() *sql.DB {
	if connection == nil {
		openConnection(configuration.GetConfig().DB)
	}
	return connection
}

func RestartConnection() {
	if connection != nil {
		connection.Close()
	}
	openConnection(configuration.GetConfig().DB)
}
func formatDbConnectString(dbConfig configuration.DBConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name)
}

func openConnection(config configuration.DBConfig) {
	conn, err := sql.Open("mysql", formatDbConnectString(config))
	utils.CheckError(err, "db connection", utils.FatalLogLevel)

	connection = conn
}
