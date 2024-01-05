package models

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

const (
	TgBotHost   = "api.telegram.org"
	CtxUsername = "username"
)

type TgUserName struct {
	Username string `json:"tg_user_name"`
}

type Event struct {
	ChatId   int
	UserName string
	Text     string
}
