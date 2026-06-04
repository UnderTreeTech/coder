package auth

import "github.com/gin-gonic/gin"

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO：统一校验接口访问权限
		c.Next()
	}
}
