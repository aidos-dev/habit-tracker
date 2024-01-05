package models

import "errors"

type User struct {
	Id         int    `json:"userId" db:"id"`
	Username   string `json:"userName" db:"user_name"`
	TgUsername string `json:"tg_user_name" db:"tg_user_name"`
	FirstName  string `json:"firstName" db:"first_name"`
	LastName   string `json:"lastName" db:"last_name"`
	Email      string `json:"eMail" db:"email"`
	Password   string `json:"password" db:"password_hash"`
	Role       string `json:"role" db:"role" `
}

func (u *User) Validate() error {
	if u.Username == "" && u.TgUsername == "" {
		return errors.New("user structure has no values")
	}

	if u.TgUsername == "" {
		if u.Username == "" || u.Email == "" || u.Password == "" {
			return errors.New("user structure has no values")
		}
	}

	return nil
}

/*
GetUser struct is userd when admin gets the list of all users.
The difference with User struct is that GetUser has Id field and
doesn not include Password field
*/
type GetUser struct {
	Id         int    `json:"userId" db:"id"`
	Username   string `json:"userName" db:"user_name" binding:"required"`
	TgUsername string `json:"tg_user_name" db:"tg_user_name"`
	FirstName  string `json:"firstName" db:"first_name" binding:"required"`
	LastName   string `json:"lastName" db:"last_name" binding:"required"`
	Email      string `json:"eMail" db:"email" binding:"required"`
	Role       string `json:"role" db:"role" `
}

type TgUser struct {
	TgUsername string `json:"tg_user_name" db:"tg_user_name"`
}

const Empty = "NULL"

const (
	WebClintPath      = "/web"
	TelegramClintPath = "/telegram"
)
