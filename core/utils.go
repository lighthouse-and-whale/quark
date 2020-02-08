package core

import (
	"bytes"
	"crypto/tls"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	//"github.com/chromedp/chromedp"
	//"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------

//func WebElementRender(Url string, sleep int) (res string, err error) {
//	err = chromedp.Run(
//		context.Background(),
//		chromedp.Navigate(Url),
//		chromedp.OuterHTML("html", &res),
//		chromedp.Sleep(time.Second*sleep),
//	)
//	return
//}

func REQ(header http.Header, method, url string, body io.Reader) (http.Header, []byte, error) {
	cli := http.Client{
		// 忽略证书校验
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	if req, err := http.NewRequest(method, url, body); err != nil {
		return nil, nil, err
	} else {
		req.Header = header
		res, err := cli.Do(req)
		if err != nil {
			return nil, nil, err
		} else {
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			return res.Header, body, err
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------

func DefaultPagination(limit, page int) (int, int) {
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}
	return limit, page
}

//----------------------------------------------------------------------------------------------------------------------

func Json(data interface{}) {
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(js))
}

func JsonLine(data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(js))
}

func GetJsonBytes(data interface{}) []byte {
	js, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return js
}

func JsonValue(data []byte, key string) (value interface{}, err error) {
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}

	switch x := value.(type) {
	default:
		return nil, errors.New("json: not found")

	case map[string]interface{}:
		if subMap, _ := x[key]; subMap == nil {
			return nil, errors.New("json: not found")
		} else {
			value = subMap
		}

	case []interface{}:
		var iKey int
		i, err := strconv.Atoi(key)
		if err != nil {
			return nil, errors.New("json: not found")
		}
		iKey = i
		if iKey < 0 || iKey >= len(x) {
			return nil, errors.New("json: not found")
		}
		value = x[iKey]
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------

// ASCII 33~126
func Rand33126(n int) (output []byte) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := make([]byte, 0)
	for i := 33; i <= 126; i++ {
		seed = append(seed, byte(i))
	}
	for {
		i := r.Intn(len(seed))
		output = append(output, seed[i])
		if len(output) == n {
			return output
		}
	}
}

func RandInt(n int) (output []byte) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := []byte{
		48, 49, 50, 51, 52, 53, 54, 55, 56, 57}
	for {
		i := r.Intn(len(seed))
		output = append(output, seed[i])
		if len(output) == n {
			return output
		}
	}
}

func RandAz09(n int) (output []byte) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := []byte{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	}
	for {
		i := r.Intn(len(seed))
		output = append(output, seed[i])
		if len(output) == n {
			return output
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------

// 计算两个坐标的距离(米)
func GeographicDistance(x1, y1, x2, y2 float64) float64 {
	rad := math.Pi / 180
	x1 *= rad
	y1 *= rad
	x2 *= rad
	y2 *= rad
	dist := math.Acos(
		math.Sin(x1)*math.Sin(x2) +
			math.Cos(x1)*math.Cos(x2)*math.Cos(y2-y1),
	)
	return dist * 6371000
}

//----------------------------------------------------------------------------------------------------------------------

// 启动浏览器
func OpenWeb(url string) {
	system := runtime.GOOS
	switch {
	case system == "windows":
		if err := exec.Command("cmd", "/c", "start", url).Run(); err != nil {
			panic(err)
		}
	case system == "linux":
		if err := exec.Command("xdg-open", url).Run(); err != nil {
			panic(err)
		}
	case system == "":
		if err := exec.Command("start", url).Run(); err != nil {
			panic(err)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------

// 给文本中的链接包装<a>标签
// 该方法比较占用资源
func TextareaFormat(src string) string {
	src = strings.ReplaceAll(src, "<", "&lt;")
	src = strings.ReplaceAll(src, ">", "&gt;")
	src = strings.ReplaceAll(src, "\n", "<br>")
	var urls = regexp.MustCompile(`(ftp|http|https)://[-A-Za-z0-9_./?=&#%]+[\w]`).FindAllString(src, 1000)
	// 去除重复
	var keys = make(map[string]bool)
	for i := range urls {
		if _, has := keys[urls[i]]; !has {
			keys[urls[i]] = true
			src = strings.ReplaceAll(src, urls[i], fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, urls[i], urls[i]))
		}
	}
	return src
}

//----------------------------------------------------------------------------------------------------------------------

func SaveGobStorage(data interface{}, name string) {
	buffer := new(bytes.Buffer)
	if err := gob.NewEncoder(buffer).Encode(data); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(name, buffer.Bytes(), 0600); err != nil {
		panic(err)
	}
}

func LoadGobStorage(data interface{}, name string) {
	raw, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewBuffer(raw)
	if err := gob.NewDecoder(buffer).Decode(data); err != nil {
		panic(err)
	}
}

//----------------------------------------------------------------------------------------------------------------------

// 检测数据是否是二进制
func IsBinary(data []byte) (b bool) {
	if len(data) > 100 {
		data = data[0:100]
	}
	for i := range data {
		if data[i] == 0 {
			return true
		}
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------

func Timekeeper(name string) func() {
	start := time.Now()
	log.Printf("[  run  ] %s\n", name)
	return func() {
		log.Printf("[  end  ] %s T:%s\n", name, time.Since(start))
	}
}
