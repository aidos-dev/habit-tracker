package models

type Auth struct {
	Client Client `json:"client_data"`
	User   User   `json:"user_data"`
}
