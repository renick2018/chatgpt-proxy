package component

import (
	"chatgpt-proxy/lib/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func ParsePostMap(c *gin.Context, pointer interface{}) {
	if err := c.ShouldBindBodyWith(pointer, binding.JSON); err != nil {
		logger.Warning(fmt.Sprintf("parse requesr params err: %+v", err))
		c.JSON(http.StatusOK, gin.H{"err": "params error"})
	}
}

func Response(c *gin.Context, code int, message string, data ...interface{}) {
	if len(data) != 1 {
		c.JSON(http.StatusOK, gin.H{"code": code, "error": message})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": code, "error": message, "data": data[0]})
	}
}
