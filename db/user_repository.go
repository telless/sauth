package db

import (
	"sauth/models"
	"strings"
	"sauth/utils"
)

func GetUserByWSOLogin(wsoLogin string) models.User {
	apiUser := models.User{}
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

	connection := GetConnection()
	defer connection.Close()

	stmt, err := GetConnection().Prepare(sqlQuery)
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
