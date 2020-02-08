package core

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------

func (c *Core) writeHeader(w http.ResponseWriter, code int) {
	c.code = code
	w.WriteHeader(200)
}

func (c *Core) toJson(w http.ResponseWriter, code int, comment string, v interface{}) {
	c.comment = comment
	if res, err := json.Marshal(ResponseJson{Code: code, Comment: comment, Data: v}); err == nil {
		c.writeHeader(w, code)
		_, _ = w.Write(res)
	} else {
		panic(err)
	}
}

func (c Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		return
	}
	now := time.Now()

	defer func(now time.Time) {
		// 消化日志
		loggerQueue <- requestLoggerQueue{c, r, time.Now().Sub(now)}
		_ = r.Body.Close()
	}(now)
	defer internalServerError(&c, w) // 错误追踪

	// 全局请求头
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Add("Access-Control-Allow-Origin", "*") // 跨域
	w.Header().Add("Access-Control-Max-Age", "3600")   // OPTIONS 缓存时间

	if r.URL.Path == "/" || r.URL.Path == "" {
		r.URL.Path = "/main"
	}

	var prefix string
	var rup = strings.Split(r.URL.Path, "/")

	for i := range rup {
		if i == 0 {
			continue
		}
		prefix += "/" + rup[i]
		if h, ok := multiplexer[prefix]; ok {
			if len(h.handles) != 0 {
				for j := range h.handles {
					if err := h.handles[j](Context{
						core: &c, Response: w, Request: r, User: r.Header.Get(AccessTokenHeader),
					}); err == nil {
						return
					}
				}
			}
		}
	}

	c.toJson(w, 0, "No such API", map[int]int{})
}

//----------------------------------------------------------------------------------------------------------------------

// example:
// Add("tag", "/url", callback)
//
func (c *Core) Add(tag, url string, handles ...handleFunc) *Core {
	if len(url) == 0 || url == "/" {
		println(`[  ERROR  ] Add(): url cannot be set to "" or "/"`)
		os.Exit(0)
	}
	prefix := c.prefix + url
	multiplexer[prefix] = multiplexerHandles{handles: handles}
	if len(tag) != 0 {
		apis = append(apis, interfaceDoc{Comment: tag, Path: prefix})
	}
	return &Core{prefix: prefix}
}

//----------------------------------------------------------------------------------------------------------------------

// 日志记录
func requestLogger(c *Core, r *http.Request, delay time.Duration) {
	var contentType = r.Header.Get("Content-Type")
	var content string
	if contentType != "application/json" {
		content = "{}"
	} else {
		content = strings.ReplaceAll(string(c.body), `"`, `\"`)
	}
	Line := `{"Proto":"%s","code":%d,"comment":"%s","url":"%s","delay":"%s","user":"%s","ip":"%s","device":"%s","body":"%s"}` + "\n"
	log.Printf(
		Line, r.Proto, c.code, c.comment, r.URL, delay, r.Header.Get(AccessTokenHeader),
		ClientIP(r), __deviceUserAgent(r.UserAgent()), content)
}

func __deviceUserAgent(v string) (o string) {
	s := strings.Split(regexp.MustCompile(`\(.+\)`).FindString(v), ")")[0]
	if len(s) != 0 {
		o = s[1:]
	}
	if len(o) == 0 {
		return v
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------

// 错误追踪
func internalServerError(c *Core, w http.ResponseWriter) {
	if e := recover(); e != nil {
		__panicTrace(e)
		if res, e := json.Marshal(ResponseJson{500, "busy server", map[int]int{}}); e == nil {
			c.writeHeader(w, 200)
			_, _ = w.Write(res)
		}
	}
}

func __panicTrace(err interface{}) {
	var p string
	trace := make([]byte, 1<<16)
	now := time.Now().Local().Format(time.RFC3339Nano)
	v := strings.Split(string(trace[:int(math.Min(float64(runtime.Stack(trace, true)), 5000))]), "\n")
	for i := range v {
		if strings.HasPrefix(v[i], "	") {
			if !strings.Contains(v[i], "/go.mongodb.org/") {
				if !strings.Contains(v[i], "/github.com/") {
					if !strings.Contains(v[i], "/usr/local/") {
						p += v[i] + "\n"
					}
				}
			}
		}
	}
	ret := fmt.Sprintf(">>> %s\npanic: %v\n%s\n", now, err, p)
	fmt.Println(ret)
	var f *os.File
	f, err = os.OpenFile(now[:10]+".panic.log", 0x2|0x400|0x40, 0600)
	if err == nil {
		_, _ = f.WriteString(ret)
		_ = f.Close()
	}
}

//----------------------------------------------------------------------------------------------------------------------

// 获取本地IP
func LanIp() string {
	addr, _ := net.InterfaceAddrs()
	return addr[1].(*net.IPNet).IP.String()
}

// 尽最大努力实现获取客户端 IP 的算法
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作
func ClientIP(r *http.Request) (ip string) {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip = strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return
	}
	ip, _, _ = net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	return
}
