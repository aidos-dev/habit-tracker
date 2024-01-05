package models

type Client struct {
	ClientType string `json:"client_type"`
}

const (
	WebClient      = "web"
	TelegramClient = "telegram"
)
