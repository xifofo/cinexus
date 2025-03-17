package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"cinexus/pkg/jwt"
	"cinexus/pkg/logger"
)

// JWT 中间件，用于验证JWT令牌
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未提供授权令牌",
			})
			c.Abort()
			return
		}

		// 检查令牌格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "授权令牌格式错误",
			})
			c.Abort()
			return
		}

		// 解析令牌
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			logger.Error("解析令牌失败", zap.Error(err))

			var msg string
			switch err {
			case jwt.ErrTokenExpired:
				msg = "令牌已过期"
			case jwt.ErrTokenNotValidYet:
				msg = "令牌尚未生效"
			case jwt.ErrTokenMalformed:
				msg = "令牌格式错误"
			default:
				msg = "无效的令牌"
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  msg,
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
