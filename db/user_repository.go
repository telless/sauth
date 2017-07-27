package db

import (
	"sauth/models"
	"strings"
	"sauth/utils"
)

func GetUserByWSOLogin(wsoLogin string) models.User {
	user := models.User{}
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
	stmt, err := GetConnection().Prepare(sqlQuery)
	utils.CheckError(err, "prepare sql query", utils.FatalLogLevel)
	result := stmt.QueryRow(email)
	err = result.Scan(&user.Id,
		&user.Name,
		&user.Email,
		&user.PartnerId,
		&user.Type,
		&user.Role)

	utils.CheckError(err, "parse sql response", utils.FatalLogLevel)

	return user
}
