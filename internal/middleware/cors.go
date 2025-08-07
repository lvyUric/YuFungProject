package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware CORS跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	config := cors.Config{
		AllowAllOrigins: true, // 临时允许所有来源
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"Accept",
			"X-Requested-With",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
			"Access-Control-Allow-Methods",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Authorization",
			"Access-Control-Allow-Origin",
		},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // 12小时
	}

	return cors.New(config)
}
