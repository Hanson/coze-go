本项目持续更新中，你可以参考 oauth 中如何进行 jwt 授权，也欢迎催更（催一个接口更一个）

# Install

```
go get github.com/hanson/coze-go
```

# Usage

## 授权
```
// 获取授权对象
oauth := auth.NewOauth("appid", "kid").WithPemByte([]byte("xxx"))

// 获取授权token
oauth.GetToken()
```

## 会话
```
// 创建会话对象
c := NewConversation(o)

// 创建会话
c.CreateConversation(&CreateConversationReq{})
```

## 对话
```
// 创建对话对象
c := chat.NewChat(o)

// 发起对话
c.Chat(conversationId, &chat.ChatReq{
    BotId:  "xxx",
    UserId: "test",
    AdditionalMessages: []*common.Message{
        {
            Role: "user",
            Content: `你是谁？`,
        },
    },
    CustomVariables: map[string]string{},
    Stream:          false,
    AutoSaveHistory: true,
})
// 查看对话详情
c.Retrieve(conversationId, chatRsp.Data.Id)
// 查看对话消息详情
c.MessageList(conversationId, chatRsp.Data.Id)
```

## 工作流	
```
// 创建工作流对象
w := workflow.NewWorkflow(o)
// 执行工作流
w.WorkflowRun(&workflow.WorkflowRunReq{
    WorkflowId: "",
    Parameters: map[string]interface{}{
        "USER_INPUT":        `你好`,
        "CONVERSATION_NAME": "Default",
    },
})
```
