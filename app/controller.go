package app

import (
	"chatgpt-proxy/component"
	"chatgpt-proxy/component/chatgpt"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func ask(c *gin.Context) {
	var params map[string]string
	component.ParsePostMap(c, &params)
	if params == nil {
		return
	}

	var message = params["message"]
	var nickname = params["conversationId"]
	var rsp, conv = chatgpt.Ask(nickname, message)
	var data = make(map[string]string)
	var msg = ""
	if rsp != nil {
		data["response"] = *rsp
		data["conversationId"] = conv
		data["text"] = *rsp
		data["ParentMessageId"] = ""
	} else {
		msg = "no available server"
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "error": msg, "response": data})
	//component.Response(c, 0, msg, data)
}

func restart(c *gin.Context) {
	var params map[string]string
	component.ParsePostMap(c, &params)
	if params == nil {
		return
	}
	var host = params["host"]
	if _, err := url.ParseRequestURI(host); err != nil {
		component.Response(c, 1, fmt.Sprintf("%s is not a host", host))
		return
	}
	var count = 0
	if host != "" {
		count = chatgpt.RestartServer(host)
	}
	component.Response(c, 0, fmt.Sprintf("restart %d server", count))
}

func newServer(c *gin.Context) {
	var params map[string]string
	component.ParsePostMap(c, &params)
	if params == nil {
		return
	}
	var host = params["host"]
	if _, err := url.ParseRequestURI(host); err != nil {
		component.Response(c, 1, fmt.Sprintf("%s is not a host", host))
		return
	}
	var msg = ""
	if chatgpt.AddServer(host) {
		msg = fmt.Sprintf("%s seems running.", host)
	}
	component.Response(c, 0, msg)
}

func serverList(c *gin.Context) {
	var arr = make([]interface{}, 0)
	var online = 0
	for _, v := range chatgpt.ServerMap {
		arr = append(arr, v.Info())
		if v.Status {
			online++
		}
	}
	var data = make(map[string]interface{})
	data["servers"] = arr
	data["count"] = len(arr)
	data["online"] = online
	data["offline"] = len(arr) - online
	component.Response(c, 0, "", data)
}
