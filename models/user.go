package models

import "encoding/json"

type User struct {
	Id        int `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	PartnerId int `json:"partner_id"`
	Type      string `json:"type"`
	Role      int `json:"role"`
}

func (user *User) Serialize() string {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
