package auth

import (
	"github.com/hanson/coze-go/client"
	"github.com/patrickmn/go-cache"
	"time"
)

var c *cache.Cache

func init() {
	c = cache.New(5*time.Minute, 10*time.Minute)
}

type Auth interface {
	GetClient() (cli *client.Client, err error)
}
