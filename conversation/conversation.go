package conversation

import (
	"bytes"
	"encoding/json"
	"github.com/hanson/coze-go/auth"
	"github.com/hanson/coze-go/common"
)

type Conversation struct {
	//client *client.Client
	auth auth.Auth
}

func NewConversation(auth auth.Auth) *Conversation {
	return &Conversation{auth: auth}
}

//func NewConversation(client *client.Client) *Conversation {
//	return &Conversation{client: client}
//}

type CreateConversationReq struct {
	Messages []*common.Message `json:"messages"`
	MetaData map[string]string `json:"meta_data"`
}

type CreateConversationResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Id        string            `json:"id"`
		CreatedAt int64             `json:"created_at"`
		MetaData  map[string]string `json:"meta_data"`
	} `json:"data"`
}

func (c *Conversation) CreateConversation(req *CreateConversationReq) (resp *CreateConversationResp, err error) {
	cli, err := c.auth.GetClient()
	if err != nil {
		return
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return
	}

	respBody, err := cli.Request("POST", "https://api.coze.cn/v1/conversation/create", bytes.NewReader(reqBody))
	if err != nil {
		return
	}

	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return
	}

	return
}
