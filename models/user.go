package models

import (
	"encoding/json"
	"time"
)

type User struct {
	Id        int `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	PartnerId int `json:"partner_id,omitempty"`
	Type      string `json:"type,omitempty"`
	Role      int `json:"role,omitempty"`
	Expire    time.Time `json:"expire,omitempty"`
}

func (user *User) Serialize() string {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func (d *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		*Alias
		Expire string `json:"expire"`
	}{
		Alias: (*Alias)(d),
		Expire: d.Expire.Format("2006-01-02 15:04:05"),
	})
}
