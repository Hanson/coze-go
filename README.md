本项目持续更新中，你可以参考 oauth 中如何进行 jwt 授权，也欢迎催更（催一个接口更一个）

# Install

```
go get github.com/hanson/coze-go
```

# Usage

```
  // 获取 token
  oauth := NewOauth("appid", "kid").WithPemByte([]byte("xxx"))
  oauth.GetToken()
```
