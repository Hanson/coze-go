package client

import (
	"io"
	"net/http"
	"time"
)

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

func (c *Client) SetToken(token string) {
	c.token = token
}

func (c *Client) IsExpired() bool {
	if c.expiredIn <= 0 {
		return true
	}

	if c.expiredIn < time.Now().Unix() {
		return true
	}

	return false
}

func (c *Client) Request(method, url string, body io.Reader) (b []byte, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.token)

	rsp, err := c.http.Do(req)
	if err != nil {
		return
	}

	defer rsp.Body.Close()

	b, err = io.ReadAll(rsp.Body)
	if err != nil {
		return
	}

	return
}
