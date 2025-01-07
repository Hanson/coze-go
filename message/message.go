package message

import (
	"bytes"
	"encoding/json"
	"github.com/hanson/coze-go/auth"
	"log"
)

type Message struct {
	//client *client.Client
	auth auth.Auth
}

func NewMessage(auth auth.Auth) *Message {
	return &Message{auth: auth}
}

type CreateMessageReq struct {
	Role        string            `json:"role"`
	Content     string            `json:"content"`
	ContentType string            `json:"content_type"`
	MetaData    map[string]string `json:"meta_data"`
}

type CreateMessageResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Id        string            `json:"id"`
		CreatedAt int64             `json:"created_at"`
		MetaData  map[string]string `json:"meta_data"`
	} `json:"data"`
}

func (m *Message) CreateMessage(conversationId string, req *CreateMessageReq) (resp *CreateMessageResp, err error) {
	cli, err := m.auth.GetClient()
	if err != nil {
		return
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return
	}

	respBody, err := cli.Request("POST", "https://api.coze.cn/v1/conversation/message/create?conversation_id="+conversationId, bytes.NewReader(reqBody))
	if err != nil {
		return
	}
	log.Printf("11111111 err: %+v", string(respBody))

	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return
	}

	return
}
