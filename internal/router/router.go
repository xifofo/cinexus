package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cinexus/internal/controller"
	"cinexus/internal/middleware"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// 全局中间件
	r.Use(middleware.Cors())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// 创建控制器
	userController := controller.NewUserController()

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 无需认证的路由
		v1.POST("/auth/login", userController.Login)
		v1.POST("/auth/register", userController.Register)

		// 需要认证的路由
		auth := v1.Group("")
		auth.Use(middleware.JWT())
		{
			// 用户相关
			auth.GET("/user/info", userController.GetUserInfo)
			auth.PUT("/user/info", userController.UpdateUserInfo)
			auth.PUT("/user/password", userController.UpdatePassword)

			// 其他API路由...
		}
	}

	// 404处理
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "接口不存在",
		})
	})
}
