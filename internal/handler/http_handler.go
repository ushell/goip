package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ushell/goip/internal/service"
	"github.com/ushell/goip/pkg/errors"
	"github.com/ushell/goip/pkg/logger"
)

// HTTPHandler HTTP处理器
type HTTPHandler struct {
	service *service.IPService
	logger  *logger.Logger
}

// NewHTTPHandler 创建新的HTTP处理器
func NewHTTPHandler(service *service.IPService, logger *logger.Logger) *HTTPHandler {
	return &HTTPHandler{
		service: service,
		logger:  logger,
	}
}

// QueryIP 查询单个IP地址信息
func (h *HTTPHandler) QueryIP(c *gin.Context) {
	ip := c.Param("ip")
	ip = strings.TrimSpace(ip)

	if ip == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidRequest,
			"message": "IP地址不能为空",
		})
		return
	}

	info, err := h.service.QueryIP(ip)
	if err != nil {
		h.logger.WithError(err).WithField("ip", ip).Error("查询IP失败")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.GetCode(err),
			"message": errors.GetMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": info,
	})
}

// BatchQueryIP 批量查询IP地址信息
func (h *HTTPHandler) BatchQueryIP(c *gin.Context) {
	var req struct {
		IPs []string `json:"ips" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrCodeInvalidRequest,
			"message": "请求格式错误",
		})
		return
	}

	infos, err := h.service.BatchQueryIP(req.IPs)
	if err != nil {
		h.logger.WithError(err).Error("批量查询IP失败")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.GetCode(err),
			"message": errors.GetMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": infos,
	})
}

// HealthCheck 健康检查
func (h *HTTPHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Unix(),
	})
}

// GetServiceStatus 获取服务状态
func (h *HTTPHandler) GetServiceStatus(c *gin.Context) {
	status := h.service.GetServiceStatus()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": status,
	})
}

// GetClientIP 获取客户端IP
func (h *HTTPHandler) GetClientIP(c *gin.Context) {
	clientIP := c.ClientIP()

	info, err := h.service.QueryIP(clientIP)
	if err != nil {
		h.logger.WithError(err).WithField("ip", clientIP).Error("查询客户端IP失败")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.GetCode(err),
			"message": errors.GetMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": info,
	})
}

// SetupRoutes 设置路由
func (h *HTTPHandler) SetupRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		// IP查询
		v1.GET("/ip/:ip", h.QueryIP)
		v1.POST("/ip/batch", h.BatchQueryIP)

		// 客户端IP查询
		v1.GET("/ip/client", h.GetClientIP)

		// 服务状态
		v1.GET("/health", h.HealthCheck)
		v1.GET("/status", h.GetServiceStatus)
	}
}
