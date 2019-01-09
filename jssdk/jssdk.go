/*
# @Time : 2019-01-09 10:49
# @Author : smallForest
# @SoftWare : GoLand
# 微信jssdk生成
*/
package jssdk

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var appid string
var appsecret string

type wxconfig struct {
	appId     string
	appSecret string
	url       string
}
type json1 struct {
	Expire_time  int
	Jsapi_ticket string
}

type json2 struct {
	Expire_time  int
	Access_token string
}
type token struct {
	Access_token string
	Expires_in   int
}

type ticket struct {
	Errcode    int
	Errmsg     string
	Ticket     string
	Expires_in int
}

type signPackage struct {
	AppId     string   `json:"appId"`
	NonceStr  string   `json:"nonceStr"`
	Timestamp int      `json:"timestamp"`
	Url       string   `json:"url"`
	Signature string   `json:"signature"`
	JsApiList []string `json:"jsApiList"`
}

func New(appId string, appSecret string, url string) wxconfig {
	wx := wxconfig{appId, appSecret, url}
	return wx
}

// 创建类方法获取jssdk配置项目 返回json串
func (wx wxconfig) GetWechatConfig() string {
	appid = wx.appId
	appsecret = wx.appSecret
	fmt.Printf("running appid: %s appsecret: %s\n", wx.appId, wx.appSecret)
	ticket := getJsApiTicket()
	timestamp := timestamp()
	nonceStr := getRandomString(16)
	url := wx.url
	my_string := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", ticket, nonceStr, timestamp, url)
	signature := my_sha1(my_string)
	signPackage := signPackage{wx.appId, nonceStr, timestamp, url, signature, []string{"onMenuShareTimeline", "onMenuShareAppMessage"}}
	fmt.Println(signPackage)
	data, err := json.Marshal(signPackage)
	if err != nil {
		fmt.Println(err)
	}
	return string(data)
}

// 函数首字母小写 表示私有
func getJsApiTicket() string {
	fmt.Printf("获取getJsApiTicket\n")
	// 获取ticket内容
	json_str := readFile("./jssdk/ticket.txt")
	// 解析json串
	j := json1{}
	err := json.Unmarshal([]byte(json_str), &j)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	expire_time := j.Expire_time
	jsapi_ticket := j.Jsapi_ticket
	if expire_time < timestamp() {
		// ticket过期需要从新获取
		// 获取access_token生成ticket
		access_token := getAccessToken()
		fmt.Println(access_token)
		if access_token == "" {
			return ""
		}
		url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token=%s", access_token)
		resp, _ := http.Get(url)

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		// 解析body
		ticket := ticket{}
		err := json.Unmarshal([]byte(body), &ticket)
		if err != nil {
			fmt.Println(err.Error())
			return ""
		}
		if ticket.Errcode == 0 {
			jj := json1{timestamp() + 7000, ticket.Ticket}
			data, err := json.Marshal(jj)
			if err != nil {
				fmt.Println(err)
			}
			setFile(string(data), "./jssdk/ticket.txt")
			return ticket.Ticket
		} else {
			return ""
		}
	}

	return jsapi_ticket
}

// 写文件"./jssdk/ticket.txt"
func setFile(ticket string, filename string) bool {
	d1 := []byte(ticket)
	err := ioutil.WriteFile(filename, d1, 0644)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// 读取文件内容
func readFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("文件读取出错", err)
		return fmt.Sprintf("%s", err)
	}
	fmt.Println("文件读取成功", string(data))
	return string(data)
}

// 获取access token
func getAccessToken() string {
	// 获取ticket内容
	json_str := readFile("./jssdk/token.txt")
	// 解析json串
	j := json2{}
	err := json.Unmarshal([]byte(json_str), &j)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	expire_time := j.Expire_time
	access_token := j.Access_token
	if expire_time < timestamp() {
		// 获取新token
		url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appid, appsecret)
		resp, _ := http.Get(url)

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		// 解析body
		token := token{}
		err := json.Unmarshal([]byte(body), &token)
		if err != nil {
			fmt.Println(err.Error())
			return ""
		}
		if len(token.Access_token) > 0 {
			jj := json2{timestamp() + 7000, token.Access_token}
			data, err := json.Marshal(jj)
			if err != nil {
				fmt.Println(err)
			}
			setFile(string(data), "./jssdk/token.txt")
			return token.Access_token
		} else {
			return ""
		}
	}
	return access_token
}

func timestamp() int {
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	my_timestamp, _ := strconv.Atoi(timestamp)
	return my_timestamp
}

// 获取随机字符串
func getRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// sha1加密
func my_sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
