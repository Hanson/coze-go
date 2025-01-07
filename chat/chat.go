package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hanson/coze-go/auth"
	"github.com/hanson/coze-go/common"
)

type Chat struct {
	//client *client.Client
	auth auth.Auth
}

func NewChat(auth auth.Auth) *Chat {
	return &Chat{auth: auth}
}

type ChatReq struct {
	BotId              string            `json:"bot_id"`
	UserId             string            `json:"user_id"`
	AdditionalMessages []*common.Message `json:"additional_messages"`
	Stream             bool              `json:"stream"`
	CustomVariables    map[string]string `json:"custom_variables"`
	AutoSaveHistory    bool              `json:"auto_save_history"`
	MetaData           map[string]string `json:"meta_data"`
	ExtraParams        map[string]string `json:"extra_params"`
}

type ChatResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Id             string `json:"id"`
		ConversationId string `json:"conversation_id"`
		BotId          string `json:"bot_id"`
		CreatedAt      int64  `json:"created_at"`
		CompletedAt    int64  `json:"completed_at"`
		Status         string `json:"status"`
		Usage          struct {
			TokenCount  int64 `json:"token_count"`
			OutputCount int64 `json:"output_count"`
			InputCount  int64 `json:"input_count"`
		}
	}
}

func (c *Chat) Chat(conversationId string, req *ChatReq) (resp *ChatResp, err error) {
	cli, err := c.auth.GetClient()
	if err != nil {
		return
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return
	}

	respBody, err := cli.Request("POST", "https://api.coze.cn/v3/chat?conversation_id="+conversationId, bytes.NewReader(reqBody))
	if err != nil {
		return
	}

	fmt.Println(string(reqBody))
	fmt.Println(string(respBody))

	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return
	}

	return
}

type RetrieveResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Id             string `json:"id"`
		ConversationId string `json:"conversation_id"`
		BotId          string `json:"bot_id"`
		CreatedAt      int64  `json:"created_at"`
		CompletedAt    int64  `json:"completed_at"`
		Status         string `json:"status"`
		Usage          struct {
			TokenCount  int64 `json:"token_count"`
			OutputCount int64 `json:"output_count"`
			InputCount  int64 `json:"input_count"`
		}
	}
}

func (c *Chat) Retrieve(conversationId, chatId string) (resp *RetrieveResp, err error) {
	cli, err := c.auth.GetClient()
	if err != nil {
		return
	}

	respBody, err := cli.Request("GET", fmt.Sprintf("https://api.coze.cn/v3/chat/retrieve?conversation_id=%s&chat_id=%s", conversationId, chatId), nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return
	}

	return
}

type MessageListResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Id             string `json:"id"`
		ConversationId string `json:"conversation_id"`
		ChatId         string `json:"chat_id"`
		BotId          string `json:"bot_id"`
		Role           string `json:"role"`
		Content        string `json:"content"`
		Type           string `json:"type"`
	}
}

func (c *Chat) MessageList(conversationId, chatId string) (resp *MessageListResp, err error) {
	cli, err := c.auth.GetClient()
	if err != nil {
		return
	}

	respBody, err := cli.Request("GET", fmt.Sprintf("https://api.coze.cn/v3/chat/message/list?conversation_id=%s&chat_id=%s", conversationId, chatId), nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return
	}

	return
}
