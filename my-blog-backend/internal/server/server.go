// internal/server/server.go
package server

import (
	"log"

	"my-blog-backend/internal/app"

	"go.uber.org/zap"
)

// Run 启动服务器
// 主程序入口
func Run() {
	app, err := app.NewApplication()
	if err != nil {
		log.Fatal("Failed to create application", zap.Error(err))
	}

	if err := app.Run(); err != nil {
		log.Fatal("Application run failed", zap.Error(err))
	}
}
