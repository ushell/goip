package service

import (
	"sync/atomic"
	"time"

	"github.com/ushell/goip/internal/config"
	"github.com/ushell/goip/internal/ipquery"
	"github.com/ushell/goip/pkg/errors"
	"github.com/ushell/goip/pkg/logger"
)

// IPService IP查询服务
type IPService struct {
	provider   ipquery.QueryProvider
	cache      *ipquery.MemoryCache
	config     *config.Config
	logger     *logger.Logger
	queryCount int64
	startTime  time.Time
}

// NewIPService 创建新的IP服务
func NewIPService(config *config.Config, logger *logger.Logger) (*IPService, error) {
	var provider ipquery.QueryProvider
	provider, err := ipquery.NewIP2RegionProvider(config.IPDatabase.Path)
	if err != nil {
		return nil, errors.NewWithError(errors.ErrCodeInternalError, "初始化IP查询提供者失败", err)
	}

	var cache *ipquery.MemoryCache
	if config.Cache.Enabled {
		cache = ipquery.NewMemoryCache(config.Cache.TTL)
	}

	return &IPService{
		provider:  provider,
		cache:     cache,
		config:    config,
		logger:    logger,
		startTime: time.Now(),
	}, nil
}

// QueryIP 查询单个IP地址信息
func (s *IPService) QueryIP(ip string) (*ipquery.IPInfo, error) {
	atomic.AddInt64(&s.queryCount, 1)

	// 验证IP地址
	if !ipquery.ValidateIP(ip) {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "无效的IP地址格式")
	}

	// 检查缓存
	if s.cache != nil {
		if cached, found := s.cache.Get(ip); found {
			s.logger.WithField("ip", ip).Debug("从缓存获取IP信息")
			return cached, nil
		}
	}

	// 查询IP信息
	info, err := s.provider.Query(ip)
	if err != nil {
		s.logger.WithError(err).WithField("ip", ip).Error("查询IP信息失败")
		return nil, errors.NewWithError(errors.ErrCodeInternalError, "查询IP信息失败", err)
	}

	// 缓存结果
	if s.cache != nil && info.IsValid {
		s.cache.Set(ip, info)
	}

	s.logger.WithField("ip", ip).WithField("country", info.Country).Info("查询IP信息成功")
	return info, nil
}

// BatchQueryIP 批量查询IP地址信息
func (s *IPService) BatchQueryIP(ips []string) ([]*ipquery.IPInfo, error) {
	if len(ips) == 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "IP列表不能为空")
	}

	if len(ips) > 100 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "单次查询IP数量不能超过100个")
	}

	atomic.AddInt64(&s.queryCount, int64(len(ips)))

	results := make([]*ipquery.IPInfo, 0, len(ips))

	for _, ip := range ips {
		info, err := s.QueryIP(ip)
		if err != nil {
			info = &ipquery.IPInfo{
				IP:           ip,
				IsValid:      false,
				ErrorMessage: errors.GetMessage(err),
			}
		}
		results = append(results, info)
	}

	s.logger.WithField("count", len(ips)).Info("批量查询IP信息成功")
	return results, nil
}

// GetServiceStatus 获取服务状态
func (s *IPService) GetServiceStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":      "running",
		"version":     "1.0.0",
		"uptime":      time.Since(s.startTime).Seconds(),
		"query_count": atomic.LoadInt64(&s.queryCount),
		"cache_size":  s.getCacheSize(),
	}
}

// getCacheSize 获取缓存大小
func (s *IPService) getCacheSize() int {
	if s.cache != nil {
		return s.cache.Size()
	}
	return 0
}

// Close 关闭服务
func (s *IPService) Close() error {
	if s.provider != nil {
		return s.provider.Close()
	}
	return nil
}
