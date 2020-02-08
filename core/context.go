package core

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"regexp"
)

//----------------------------------------------------------------------------------------------------------------------

func (c *Context) Redirect(url string) (err error) {
	http.Redirect(c.Response, c.Request, url, 308)
	return
}

func (c *Context) ToString(s string) (err error) {
	_, _ = c.Response.Write([]byte(s))
	return
}

func (c *Context) ToHTML(dir string) (err error) {
	tpl, _ := template.ParseFiles(dir)
	_ = tpl.Execute(c.Response, nil)
	return
}

func (c *Context) ToJson(code int, comment string, v ...interface{}) (err error) {
	if len(v) == 0 {
		c.core.toJson(c.Response, code, comment, map[int]int{})
	} else {
		c.core.toJson(c.Response, code, comment, v[0])
	}
	return
}

func (c *Context) To401() (err error) {
	c.core.toJson(c.Response, 401, "unauthorized", map[int]int{})
	return
}

func (c *Context) TO403(comment string, v ...interface{}) (err error) {
	c.core.toJson(c.Response, 403, "非法:"+comment, map[int]int{})
	return
}

func (c *Context) ToStream(r io.Reader, name string, length float64) (err error) {
	c.Response.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
	c.Response.Header().Add("Content-Length", fmt.Sprintf("%d", int64(length)))
	c.core.writeHeader(c.Response, 200)
	c.core.comment = "dl/" + name
	if len(name) == 0 {
		c.core.comment = "dl/stream"
	}
	_, _ = io.Copy(c.Response, r)
	return
}

//----------------------------------------------------------------------------------------------------------------------

func (c *Context) Query(v string) string {
	return xss.Sanitize(c.Request.URL.Query().Get(v))
}

//----------------------------------------------------------------------------------------------------------------------

// example:
// input := struct {
//     Name string `json:"name"`
// }{}
// ct.Bind(&input).RegexVerify(core.RegexVerifyMap{
//     "title": {V: input.Name, P: `^.{1,20}$`},
// })
//
func (c *Context) Bind(v interface{}) *Context {
	if err := json.NewDecoder(c.Request.Body).Decode(&v); err != nil {
		_ = c.ToJson(400, "json syntax error")
		panic(nil)
	}
	c.core.body, _ = json.Marshal(v)
	if fmt.Sprintf("%T", v) == "*map[string]string" {
		println("error: Bind() do not use map")
		os.Exit(0)
	}
	// XSS
	var x = reflect.ValueOf(v).Elem()
	for i := 0; i < x.NumField(); i++ {
		if x.Field(i).Type().String() == "string" {
			x.Field(i).SetString(xss.Sanitize(x.Field(i).String()))
		}
	}
	return c
}

func (c *Context) RegexVerify(fields map[string]RegexVerifyItem) {
	for i := range fields {
		T := reflect.TypeOf(fields[i].V).String()
		has := true

		switch T {
		case "string":
			if len(fields[i].P) != 0 {
				has, _ = regexp.MatchString(fields[i].P, fields[i].V.(string))
			}
		case "[]string":
			for _, field := range fields[i].V.([]string) {
				if len(fields[i].P) != 0 {
					has, _ = regexp.MatchString(fields[i].P, field)
					if !has {
						break
					}
				}
			}
		}

		if !has {
			_ = c.ToJson(400, fmt.Sprintf("字段 '%s' 未通过正则检查 '%s'", i, fields[i].P))
			panic(nil)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------

// 获取上传的文件
// example:
// file, filename, err := ct.GetRequestFile("file", 10)
// defer file.Close()
//
func (c *Context) GetRequestFile(name string, fileSizeMB int64,
) (output multipart.File, filename string, err error) {
	// r.FormFile 缓冲区大小
	_ = c.Request.ParseMultipartForm(1024 * 10) // 10KB

	var fileHeader *multipart.FileHeader
	_, fileHeader, err = c.Request.FormFile(name)
	if err != nil {
		_ = c.ToJson(400, "multipartForm ("+name+") is null")
		panic(nil)
	} else {
		if fileHeader.Size > fileSizeMB*1024*1024 {
			_ = c.ToJson(400, fmt.Sprintf("文件体积:%dMB > %dMB", fileHeader.Size/1024/1024, fileSizeMB))
			panic(nil)
		}
		filename = fileHeader.Filename
	}
	output, err = fileHeader.Open()
	return
}
