package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/ushell/goip/api/proto"
	"github.com/ushell/goip/internal/config"
	"github.com/ushell/goip/internal/handler"
	"github.com/ushell/goip/internal/service"
	"github.com/ushell/goip/pkg/logger"
	"google.golang.org/grpc"
)

func main() {
	// 加载配置
	cfg, err := config.Load("./configs")
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	log := logger.New(cfg.Logging.Level, cfg.Logging.Format, cfg.Logging.Output)
	log.Info("启动IP查询服务...")

	// 创建IP服务
	ipService, err := service.NewIPService(cfg, log)
	if err != nil {
		log.WithError(err).Fatal("创建IP服务失败")
	}
	defer ipService.Close()

	// 创建HTTP服务器
	httpHandler := handler.NewHTTPHandler(ipService, log)
	router := gin.Default()
	httpHandler.SetupRoutes(router)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.HTTP.ReadTimeout,
		WriteTimeout: cfg.Server.HTTP.WriteTimeout,
		IdleTimeout:  cfg.Server.HTTP.IdleTimeout,
	}

	// 启动HTTP服务器
	go func() {
		log.WithField("addr", httpServer.Addr).Info("启动HTTP服务器")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("HTTP服务器启动失败")
		}
	}()

	// 启动gRPC服务器（暂时注释掉，等待proto文件生成）
	grpcServer := grpc.NewServer()
	pb.RegisterIPQueryServiceServer(grpcServer, handler.NewGRPCServer(ipService, log))

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.GRPC.Host, cfg.Server.GRPC.Port))
	if err != nil {
		log.WithError(err).Fatal("gRPC监听失败")
	}

	go func() {
		log.WithField("addr", lis.Addr()).Info("启动gRPC服务器")
		if err := grpcServer.Serve(lis); err != nil {
			log.WithError(err).Fatal("gRPC服务器启动失败")
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("正在关闭服务...")

	// 优雅关闭HTTP服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.WithError(err).Error("HTTP服务器关闭失败")
	}

	// 关闭gRPC服务器（暂时注释掉）
	grpcServer.GracefulStop()

	log.Info("服务已关闭")
}
