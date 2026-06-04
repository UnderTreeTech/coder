package cors

import (
	"time"

	"github.com/UnderTreeTech/waterdrop/pkg/server/http/middlewares"
	"github.com/gin-gonic/gin"
)

// Cors 跨域处理
func Cors() gin.HandlerFunc {
	cfg := middlewares.CORSConfig{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Length", "Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Cache-Control", "Content-Language", "Content-Type"},
	}
	return middlewares.NewCORS(cfg)
}
