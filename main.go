/*
# @Time : 2019-01-09 10:50
# @Author : smallForest
# @SoftWare : GoLand
*/
package main

import (
	"fmt"
	"sfwechat/jssdk"
)

func main() {
	jssdk := jssdk.New("wx02d6c9061******", "7acfb40fb2f70cd331a*******", "http://a.com")
	config := jssdk.GetWechatConfig()
	fmt.Println(config)
}
