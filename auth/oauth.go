package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hanson/coze-go/client"
	"github.com/hanson/go-toolbox/utils"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Oauth struct {
	appId string
	kid   string
	pem   string

	key   []byte
	token *jwt.Token

	accessToken *OauthTokenResp
}

func NewOauth(appId, kid string) (oauth *Oauth) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": appId,                    // OAuth 应用的 ID
		"aud": "api.coze.cn",            //扣子 API 的Endpoint
		"iat": time.Now().Unix(),        // JWT开始生效的时间，秒级时间戳
		"exp": time.Now().Unix() + 3600, // JWT过期时间，秒级时间戳
		"jti": utils.RandStr(16, 0),     // 随机字符串，防止重放攻击
	})
	token.Header["kid"] = kid

	return &Oauth{
		token: token,
	}
}

func (o *Oauth) WithPemByte(pem []byte) *Oauth {
	o.key = pem

	return o
}

func (o *Oauth) WithPemFile(path string) *Oauth {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	o.key = b

	return o
}

func (o *Oauth) sign() (string, error) {
	if len(o.key) == 0 {
		return "", fmt.Errorf("key is empty")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(o.key)
	if err != nil {
		log.Printf("err: %+v", err)
		return "", err
	}

	sign, err := o.token.SignedString(key)
	if err != nil {
		log.Printf("err: %+v", err)
		return "", err
	}

	return sign, nil
}

type OauthTokenResp struct {
	ExpiresIn   int64
	AccessToken string
}

func (o *Oauth) GetClient() (*client.Client, error) {
	rsp, err := o.GetToken()
	if err != nil {
		return nil, err
	}

	return client.NewClient(rsp.AccessToken, rsp.ExpiresIn), nil
}

func (o *Oauth) GetToken() (*OauthTokenResp, error) {
	sign, err := o.sign()
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	req, _ := http.NewRequest("POST", "https://api.coze.cn/api/permission/oauth2/token", bytes.NewReader([]byte(`{"duration_seconds": 86399,"grant_type": "urn:ietf:params:oauth:grant-type:jwt-bearer"}`)))

	req.Header.Add("Authorization", "Bearer "+sign)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("err: %+v", err)
		return nil, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rsp OauthTokenResp
	err = json.Unmarshal(b, &rsp)
	if err != nil {
		return nil, err
	}

	o.accessToken = &rsp

	return &rsp, nil
}
