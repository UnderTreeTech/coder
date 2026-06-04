package http

import (
	"github.com/gin-gonic/gin"
)

func registerAPI(engine *gin.Engine) {
	login := engine.Group("/api")
	{
		login.GET("/app/user", getUserInfo) // 获取用户信息
	}
}
