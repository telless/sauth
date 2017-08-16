package main

import (
	"encoding/json"
	"time"
	"strings"
)

type apiUser struct {
	Id          int `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	PartnerId   int `json:"partnerId,omitempty"`
	Type        string `json:"type,omitempty"`
	Role        int `json:"role,omitempty"`
	Expire      time.Time `json:"dataExpire,omitempty"`
	TokenExpire time.Time `json:"tokenExpire,omitempty"`
}

func (user *apiUser) Serialize() string {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return ""
	}

	return string(jsonData)
}

// Override default time format
func (d *apiUser) MarshalJSON() ([]byte, error) {
	type Alias apiUser

	return json.Marshal(&struct {
		*Alias
		Expire      string `json:"dataExpire"`
		TokenExpire string `json:"tokenExpire"`
	}{
		Alias:       (*Alias)(d),
		Expire:      d.Expire.Format(projectTimeFormat),
		TokenExpire: d.TokenExpire.Format(projectTimeFormat),
	})
}

// Get user data from DB
func getUserByWSOLogin(wsoLogin string) (apiUser, error) {
	user := apiUser{}
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
	stmt, err := getConnection().Prepare(sqlQuery)
	checkError(err, "prepare sql query", fatalLogLevel)
	result := stmt.QueryRow(email)
	err = result.Scan(&user.Id,
		&user.Name,
		&user.Email,
		&user.PartnerId,
		&user.Type,
		&user.Role)

	return user, err
}
