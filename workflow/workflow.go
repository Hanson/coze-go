package workflow

import (
	"bytes"
	"encoding/json"
	"github.com/hanson/coze-go/auth"
)

type Workflow struct {
	//client *client.Client
	auth auth.Auth
}

func NewWorkflow(auth auth.Auth) *Workflow {
	return &Workflow{auth: auth}
}

type WorkflowRunReq struct {
	WorkflowId string                 `json:"workflow_id"`
	Parameters map[string]interface{} `json:"parameters"`
	BotId      string                 `json:"bot_id"`
	Ext        map[string]string      `json:"ext"`
	IsAsync    bool                   `json:"is_async"`
	AppId      string                 `json:"app_id"`
}

type WorkflowRunResp struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Data      string `json:"data"`
	DebugUrl  string `json:"debug_url"`
	ExecuteId string `json:"execute_id"`
}

func (w *Workflow) WorkflowRun(req *WorkflowRunReq) (resp *WorkflowRunResp, err error) {
	cli, err := w.auth.GetClient()
	if err != nil {
		return
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return
	}

	respBody, err := cli.Request("POST", "https://api.coze.cn/v1/workflow/run", bytes.NewReader(reqBody))
	if err != nil {
		return
	}

	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return
	}

	return
}
