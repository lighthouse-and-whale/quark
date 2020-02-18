package core

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------

func newSignature(key, data []byte) string {
	h := hmac.New(sha256.New, key)
	if _, e := h.Write(data); e != nil {
		panic(e)
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

//----------------------------------------------------------------------------------------------------------------------

//easyjson:json
type AccessToken struct {
	// 生效时间
	N int64
	// 失效时间
	E int64
	// 携带数据
	D string
}

//----------------------------------------------------------------------------------------------------------------------

// 用户令牌签发
func NewAccessToken(key []byte, data string, nbf, exp time.Duration) string {
	var now = time.Now()
	var T, _ = json.Marshal(AccessToken{
		N: now.Add(nbf).Unix(),
		E: now.Add(exp).Unix(),
		D: data,
	})
	return base64.StdEncoding.EncodeToString([]byte(string(T) + "." + newSignature(key, T)))
}

// 用户令牌校验
func VerifyAccessToken(token string) (user string, b bool) {
	var vbs [][]byte
	if str, e := base64.StdEncoding.DecodeString(token); e != nil {
		if e != nil {
			panic(e)
		}
	} else {
		vbs = bytes.Split(str, []byte("."))
	}
	if len(vbs) != 2 {
		log.Println("if len(vbs) != 2")
		return
	}
	// 签名校验
	if string(vbs[1]) != newSignature([]byte(ATK), vbs[0]) {
		return
	}
	// 抽取数据
	var a AccessToken
	if e := a.UnmarshalJSON(vbs[0]); e != nil {
		panic(e)
	}
	// 是否生效
	if time.Now().Unix() < a.N {
		return
	}
	// 是否过期
	if time.Now().Unix() > a.E {
		return
	}
	return a.D, true
}

//----------------------------------------------------------------------------------------------------------------------

func (c *Context) GetAccessToken() string {
	return c.Request.Header.Get("Authorization")
}
