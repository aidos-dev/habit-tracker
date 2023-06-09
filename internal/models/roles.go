package models

const (
	UserGeneral   = "user_basic"
	Administrator = "admin"
)

type UpdateRoleInput struct {
	Role string `json:"role"`
}
