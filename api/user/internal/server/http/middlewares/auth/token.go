package auth

import "github.com/gin-gonic/gin"

func Token() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO：统一校验token有效性
		c.Next()
	}
}
