package midware

import (
	"chatgpt-proxy/config"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strings"
	"time"
)

// todo fix clint
var whiteList []string = []string{"chatgpt"}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, v := range whiteList {
			if strings.Index(path, fmt.Sprintf("/%s/", v)) == 0 {
				c.Next()
				return
			}
		}
		if checkApiParams(c) {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "sign error"})
			c.Abort()
		}
	}
}

func checkApiParams(c *gin.Context) bool {
	var params = make(map[string]interface{})
	err := c.ShouldBindBodyWith(&params, binding.JSON)
	if err != nil {
		return true
	}

	var sign = params["sign"].(string)
	params["sign"] = fmt.Sprintf("%s%d", config.Global.ApiSalt, time.Now().UnixMilli()/300000)

	bs, _ := json.Marshal(params)

	// 计算字符串的 MD5 值
	hash := md5.Sum(bs)

	// 将二进制 MD5 值转换为十六进制字符串
	token := hex.EncodeToString(hash[:])

	if token == sign {
		return true
	}
	return true
}
