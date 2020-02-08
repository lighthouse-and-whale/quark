package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/profile"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------

// example:
// core.NewError(err, "package.func()")
//
func NewError(e error, m string) {
	text := fmt.Sprintf("%s [  ERROR  ] %s: %s\n",
		time.Now().Local().Format("2006/01/02 15:04:05"), m, e.Error())
	fmt.Println(text)
	var f *os.File
	f, e = os.OpenFile("errors.log", 0x2|0x400|0x40, 0600)
	if e == nil {
		_, _ = f.WriteString(text)
		_ = f.Close()
	}
}

//----------------------------------------------------------------------------------------------------------------------

func init() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		err := errors.New("open config.json: no such file or directory")
		NewError(err, "core.NewHTTP()")
		os.Exit(1)
	} else {
		if err := json.Unmarshal(data, &CONFIG); err != nil {
			err := errors.New("open config.json: does not conform to json syntax")
			NewError(err, "core.NewHTTP()")
			os.Exit(1)
		}
		ATK = CONFIG["accessToken_key"].(string)
	}
}

// example:
// HTTP := core.NewHTTP(":2000")
// HTTP.Server.ListenAndServe()
//
func NewHTTP(addr string) Core {
	go requestLoggerQueueGo()
	go blockerRestore()
	return Core{
		Server: &http.Server{Addr: addr, Handler: new(Core)},
	}
}

func requestLoggerQueueGo() {
	go func() {
		for {
			select {
			case v := <-loggerQueue:
				requestLogger(&v.c, v.r, v.d)
			}
		}
	}()
}

func blockerRestore() {
	timer := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-timer.C:
			for k, it := range online {
				it.Limit = 0
				onlineMx.Lock()
				online[k] = it
				onlineMx.Unlock()
			}
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------

func ShowAPIs(b bool) {
	log.Printf("---- SERVER STARTED ----\n")
	if !b {
		return
	}
	fmt.Println()
	for i, item := range apis {
		fmt.Printf("| %3d | [%s](#%s)| %s\n", i+1, item.Path, item.Comment, item.Comment)
	}
	fmt.Println()
}

func GetOnlineCount() int {
	return len(online)
}

func PProfCPU() interface{ Stop() } {
	return profile.Start(profile.ProfilePath("."))
}
func PProfMEM() interface{ Stop() } {
	return profile.Start(profile.MemProfile, profile.ProfilePath("."))
}

//----------------------------------------------------------------------------------------------------------------------

type ResponseJson struct {
	Code    int         `json:"code"`
	Comment string      `json:"comment"`
	Data    interface{} `json:"data"`
}

type Core struct {
	Server  *http.Server
	prefix  string
	code    int
	comment string
	body    []byte
}

type Context struct {
	core     *Core
	Request  *http.Request
	Response http.ResponseWriter
	User     string
}

type (
	handleFunc         func(Context) error
	multiplexerHandles struct {
		handles []handleFunc
	}
)

type interfaceDoc struct {
	Comment string
	Path    string
}

type requestLoggerQueue struct {
	c Core
	r *http.Request
	d time.Duration
}

type onlineStruct struct {
	WsConn *websocket.Conn
	Limit  int
}

type (
	RegexVerifyItem struct {
		V interface{}
		P string
	}
	RegexVerifyMap map[string]RegexVerifyItem
)

//----------------------------------------------------------------------------------------------------------------------

const AccessTokenHeader = "Atu"

//----------------------------------------------------------------------------------------------------------------------

// 工具类句柄
var (
	Next = errors.New("NEXT")
	xss  = bluemonday.UGCPolicy()
)

//----------------------------------------------------------------------------------------------------------------------

// WebSocket句柄
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//----------------------------------------------------------------------------------------------------------------------

const ChannelQueueMax = 10000 * 10

// 日志队列
var loggerQueue = make(chan requestLoggerQueue, ChannelQueueMax)

//----------------------------------------------------------------------------------------------------------------------

// 全局路由缓存
var (
	multiplexer = make(map[string]multiplexerHandles)
	apis        = make([]interfaceDoc, 0)
)

//----------------------------------------------------------------------------------------------------------------------

// 配置文件缓存
var (
	CONFIG = make(map[string]interface{})
	ATK    string
)

//----------------------------------------------------------------------------------------------------------------------

// 超频用户黑名单缓存
var (
	interceptor   = make(map[string]error)
	interceptorMx sync.Mutex
)

//----------------------------------------------------------------------------------------------------------------------

// 在线用户缓存
var (
	online   = make(map[string]onlineStruct)
	onlineMx sync.Mutex
)

const (
	_       = iota
	online1 // 用户上线
	online2 // 用户重新上线
	online3 // 心跳异常被下线
)

const MBx1 = 1024 * 1024
