package client

import "net/http"

type Client struct {
	token     string
	expiredIn int64
	http      *http.Client
}

func NewClient(token string, expiredIn int64) *Client {
	return &Client{
		token:     token,
		expiredIn: expiredIn,
		http:      http.DefaultClient,
	}
}
