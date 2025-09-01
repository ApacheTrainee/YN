package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 解决跨域问题
func Core(c *gin.Context) {
	// 允许所有来源
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// 允许的方法
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	// 允许的头部
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// 是否允许携带凭证（如 Cookies）
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	// 处理预检请求
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent) // 返回 204
		return
	}

	c.Next()
}
