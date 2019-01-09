### 使用go写的微信jssdk
使用方法参考main.go
```go
package main

import (
	"fmt"
	"sfwechat/jssdk"
)

func main() {
	jssdk := jssdk.New("wx02d6c9061******", "7acfb40fb2f70cd331a*******","http://a.com")
	config := jssdk.GetWechatConfig()
	fmt.Println(config)
}

```