package cmd

import (
	"cinexus/config"
	"cinexus/internal/database"
	"cinexus/internal/middleware"
	"cinexus/internal/model"
	"cinexus/internal/router"
	"cinexus/pkg/logger"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "cinexus",
	Short: "cinexus",
	Long:  `Film Fusion`,
	Run: func(cmd *cobra.Command, args []string) {
		godotenv.Load()

		// 设置运行模式
		if config.Conf.Server.RunMode == "debug" {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}

		// 初始化数据库
		if err := initDB(); err != nil {
			logger.Error("数据库初始化失败", zap.Error(err))
			return
		}

		// 创建gin引擎
		r := gin.New()
		r.Use(middleware.Logger(), middleware.Recovery())

		// 注册路由
		router.RegisterRoutes(r)

		// 创建HTTP服务器
		srv := &http.Server{
			Addr:    ":" + config.Conf.Server.Port,
			Handler: r,
		}

		// 启动服务器
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatal("监听失败", zap.Error(err))
			}
		}()

		// 等待中断信号优雅地关闭服务器
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Info("关闭服务器...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatal("服务器强制关闭", zap.Error(err))
		}

		logger.Info("服务器已退出")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println("Run Server Error: ", err)
		os.Exit(1)
	}
}

func initDB() error {
	// 初始化数据库连接
	if err := database.Init(); err != nil {
		return err
	}

	// 自动迁移数据库表结构
	err := database.DB.AutoMigrate(
		&model.User{},
		// 添加其他模型...
	)

	if err != nil {
		return err
	}

	logger.Info("数据库初始化成功")
	return nil
}
