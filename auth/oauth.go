package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hanson/coze-go/client"
	"github.com/hanson/go-toolbox/utils"
	"log"
	"os"
	"time"
)

type Oauth struct {
	appId string
	kid   string
	pem   string

	key   []byte
	token *jwt.Token

	client *client.Client
}

func NewOauth(appId, kid string) (oauth *Oauth) {
	oauth = &Oauth{
		appId:  appId,
		kid:    kid,
		client: client.NewClient("", 0),
	}

	oauth.NewJwtToken()

	return
}

func (o *Oauth) NewJwtToken() *Oauth {
	o.token = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": o.appId,                  // OAuth 应用的 ID
		"aud": "api.coze.cn",            //扣子 API 的Endpoint
		"iat": time.Now().Unix(),        // JWT开始生效的时间，秒级时间戳
		"exp": time.Now().Unix() + 3600, // JWT过期时间，秒级时间戳
		"jti": utils.RandStr(16, 0),     // 随机字符串，防止重放攻击
	})
	o.token.Header["kid"] = o.kid

	return o
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

func (o *Oauth) Sign() (string, error) {
	return o.sign()
}

func (o *Oauth) GetClient() (cli *client.Client, err error) {
	tokenCache, found := c.Get("token_" + o.kid)
	var rsp *OauthTokenResp
	if found {
		rsp = tokenCache.(*OauthTokenResp)
	} else {
		rsp, err = o.GetToken()
		if err != nil {
			return nil, err
		}

	}

	o.client = client.NewClient(rsp.AccessToken, rsp.ExpiresIn)

	return o.client, nil
}

type OauthTokenResp struct {
	ExpiresIn   int64  `json:"expires_in,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}

func (o *Oauth) GetToken() (resp *OauthTokenResp, err error) {

	sign, err := o.sign()
	if err != nil {
		return nil, err
	}

	o.client.SetToken(sign)

	b, err := o.client.Request("POST", "https://api.coze.cn/api/permission/oauth2/token", bytes.NewReader([]byte(`{"duration_seconds": 86399,"grant_type": "urn:ietf:params:oauth:grant-type:jwt-bearer"}`)))
	if err != nil {
		return
	}

	fmt.Println(string(b))

	err = json.Unmarshal(b, &resp)
	if err != nil {
		return
	}

	c.Set("token_"+o.kid, resp, 86000*time.Second)

	return
}

//func (o *Oauth) getToken() (*OauthTokenResp, error) {
//
//	sign, err := o.sign()
//	if err != nil {
//		return nil, err
//	}
//
//	req, _ := http.NewRequest("POST", "https://api.coze.cn/api/permission/oauth2/token", bytes.NewReader([]byte(`{"duration_seconds": 86399,"grant_type": "urn:ietf:params:oauth:grant-type:jwt-bearer"}`)))
//
//	req.Header.Add("Authorization", "Bearer "+sign)
//	req.Header.Add("Content-Type", "application/json")
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		log.Printf("err: %+v", err)
//		return nil, err
//	}
//
//	defer resp.Body.Close()
//
//	b, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	fmt.Println(string(b))
//
//	var rsp OauthTokenResp
//	err = json.Unmarshal(b, &rsp)
//	if err != nil {
//		return nil, err
//	}
//
//	c.Set("token_"+o.kid, &rsp, 86000*time.Second)
//
//	return &rsp, nil
//}
