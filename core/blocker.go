package core

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------

// 用户请求拦截器
func UserRequestInterceptor(ct Context) error {
	token := ct.GetAccessToken()

	// 令牌校验
	if len(token) != 0 {
		user, has := VerifyAccessToken(token)
		if !has {
			return ct.To401()
		}
		if user != "" {
			ct.User = user
			ct.Request.Header.Add(AccessTokenHeader, ct.User)
		} else {
			return Next
		}
	} else {
		return Next
	}

	// 超频拦截
	if _, has := interceptor[ct.User]; has {
		return ct.ToJson(1102, "请求超频")
	}

	// 在线状态检测
	if conn, has := online[ct.User]; !has {
		//return ct.ToJson(1100, "当前状态为离线")
	} else {
		if conn.Limit > 25 {
			// 请求频繁锁定账户1分钟，5秒内请求超过25次
			__interceptorSet(ct.User, time.Minute)
		} else {
			// 请求+1
			conn.Limit = conn.Limit + 1
			onlineMx.Lock()
			online[ct.User] = conn
			onlineMx.Unlock()
		}
	}

	return Next
}

func __interceptorSet(k string, expire time.Duration) {
	if _, has := interceptor[k]; !has {
		interceptorMx.Lock()
		interceptor[k] = nil
		interceptorMx.Unlock()
		if expire != 0 {
			go func(k string, expire time.Duration) {
				time.Sleep(expire)
				delete(interceptor, k)
			}(k, expire)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------

// 用户在线状态机
func UserOnlineStateMachine(ct Context) error {
	user, has := VerifyAccessToken(ct.Query("at"))
	if !has {
		return ct.To401()
	}
	if user != "" {
		__userOnlineStateMachine(ct.Response, ct.Request, user)
	}
	return ct.ToJson(200, "长连接结束")
}

func __userOnlineStateMachine(w http.ResponseWriter, r *http.Request, user string) {
	conn, e := upGrader.Upgrade(w, r, nil)
	if e != nil {
		panic(e)
	} else {
		defer conn.Close()

		if this, has := online[user]; !has {

			// 用户上线
			online[user] = onlineStruct{WsConn: conn}
			onlineStateMachineLogger(online1, user)

		} else {

			// 发送下线通知
			if e = online[user].WsConn.WriteMessage(1, []byte("1100")); e != nil {
				log.Printf("[ error ] %s\n", e.Error())
				return
			}
			_ = online[user].WsConn.Close()

			// 用户重新上线
			this.WsConn = conn
			online[user] = this
			onlineStateMachineLogger(online2, user)
		}
	}

	// 心跳
	for {
		time.Sleep(time.Second * 1)

		if e = conn.WriteMessage(1, []byte("200")); e != nil {
			break
		}
		_, msg, _ := conn.ReadMessage()
		if string(msg) != "200" {

			// 心跳异常，下线处理
			_ = online[user].WsConn.Close()
			delete(online, user)
			onlineStateMachineLogger(online3, user)
			break
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------

func onlineStateMachineLogger(event int, user string) {
	fmt.Printf(`ONLINE %s (%d) %s`+"\n", time.Now().Format(time.RFC3339), event, user)
}
