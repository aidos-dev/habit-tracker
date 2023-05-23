package models

type User struct {
	Id        int    `json:"-" db:"id"`
	Username  string `json:"userName" db:"user_name" binding:"required"`
	FirstName string `json:"firstName" db:"first_name" binding:"required"`
	LastName  string `json:"lastName" db:"last_name" binding:"required"`
	Email     string `json:"eMail" db:"email" binding:"required"`
	Password  string `json:"password" db:"password_hash" binding:"required"`
	Role      string `json:"role" db:"role" `
}
