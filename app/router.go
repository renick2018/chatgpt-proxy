package app

import (
	"chatgpt-proxy/midware"
	"github.com/gin-gonic/gin"
)

type Option func(engine *gin.Engine)

var options []Option

func Include(opts ...Option) {
	options = append(options, opts...)
}

// Init 初始化
func Init() *gin.Engine {
	r := gin.New()
	r.Use(midware.Cors())
	r.Use(midware.Auth())
	for _, opt := range options {
		opt(r)
	}
	return r
}

func Routers(e *gin.Engine) {
	var v = e.Group("/chatgpt")
	{
		v.POST("/ask", ask)
		v.POST("/restart", restart)
		v.POST("/add_server", newServer)
		v.POST("/servers", serverList)
	}
	e.POST("/ask", ask)
}
