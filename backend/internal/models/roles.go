package models

import "errors"

const (
	UserGeneral   = "user_basic"
	Administrator = "admin"
)

type UpdateRoleInput struct {
	Role *string `json:"role"`
}

func (i UpdateRoleInput) Validate() error {
	if i.Role == nil || *i.Role == "" {
		return errors.New("user role update structure has no values")
	}
	return nil
}
