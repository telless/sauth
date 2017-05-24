package models

import (
	"strings"
	"sauth/db"
	"fmt"
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
	sqlQuery := fmt.Sprintf(`SELECT
  u.id,
  CONCAT_WS(' ', u.last_name, u.first_name, u.middle_name) AS name,
  u.email,
  u.partner_id,
  u.type,
  u.role
FROM user AS u
WHERE email = %q
LIMIT 1`, email)

	connection := db.GetConnection()
	defer connection.Close()

	result := db.GetConnection().QueryRow(sqlQuery)
	err := result.Scan(&apiUser.Id,
		&apiUser.Name,
		&apiUser.Email,
		&apiUser.PartnerId,
		&apiUser.Type,
		&apiUser.Role)

	utils.CheckError(err)

	return apiUser
}
