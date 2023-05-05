package models

type User struct {
	Id        int    `json:"-" db:"id"`
	Username  string `json:"userName" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"eMail" binding:"required"`
	Password  string `json:"password" binding:"required"`
}
