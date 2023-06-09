package telegram

import "net/http"

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host, toke string) {
}
