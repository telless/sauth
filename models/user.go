package models

type User struct {
	Id        int `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	PartnerId int `json:"partner_id"`
	Type      string `json:"type"`
	Role      int `json:"role"`
}
