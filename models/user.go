package models

import (
	"strings"
	"sauth/db"
	"sauth/utils"
)

type User struct {
	Id        int `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	PartnerId int `json:"partner_id"`
	Type      string `json:"type"`
	Role      int `json:"role"`
}

func GetUserByWSOLogin(wsoLogin string) User {
	apiUser := User{}
	email := strings.TrimSuffix(wsoLogin, "@carbon.super") // remove @carbon.super in the end
	sqlQuery := `SELECT
  u.id,
  CONCAT_WS(' ', u.last_name, u.first_name, u.middle_name) AS name,
  u.email,
  u.partner_id,
  u.type,
  u.role
FROM user AS u
WHERE email = ?
LIMIT 1`

	connection := db.GetConnection()
	defer connection.Close()

	stmt, err := db.GetConnection().Prepare(sqlQuery)
	utils.CheckError(err)
	result := stmt.QueryRow(email)
	err = result.Scan(&apiUser.Id,
		&apiUser.Name,
		&apiUser.Email,
		&apiUser.PartnerId,
		&apiUser.Type,
		&apiUser.Role)

	utils.CheckError(err)

	return apiUser
}
