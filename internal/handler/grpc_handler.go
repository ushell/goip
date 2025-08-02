package handler

import (
	"context"
	"time"

	pb "github.com/ushell/goip/api/proto"
	"github.com/ushell/goip/internal/ipquery"
	"github.com/ushell/goip/internal/service"
	"github.com/ushell/goip/pkg/logger"
)

// GRPCServer gRPC服务器
type GRPCServer struct {
	pb.UnimplementedIPQueryServiceServer
	service *service.IPService
	logger  *logger.Logger
}

// NewGRPCServer 创建新的gRPC服务器
func NewGRPCServer(service *service.IPService, logger *logger.Logger) *GRPCServer {
	return &GRPCServer{
		service: service,
		logger:  logger,
	}
}

// QueryIP 查询单个IP地址信息
func (s *GRPCServer) QueryIP(ctx context.Context, req *pb.QueryIPRequest) (*pb.QueryIPResponse, error) {
	s.logger.WithField("ip", req.Ip).Debug("收到gRPC查询IP请求")

	info, err := s.service.QueryIP(req.Ip)
	if err != nil {
		s.logger.WithError(err).WithField("ip", req.Ip).Error("查询IP失败")
		return &pb.QueryIPResponse{
			Info: &pb.IPInfo{
				Ip:           req.Ip,
				IsValid:      false,
				ErrorMessage: err.Error(),
			},
			Timestamp: time.Now().Unix(),
		}, nil
	}

	return &pb.QueryIPResponse{
		Info:      convertToProtoIPInfo(info),
		Timestamp: time.Now().Unix(),
	}, nil
}

// BatchQueryIP 批量查询IP地址信息
func (s *GRPCServer) BatchQueryIP(ctx context.Context, req *pb.BatchQueryIPRequest) (*pb.BatchQueryIPResponse, error) {
	s.logger.WithField("count", len(req.Ips)).Debug("收到gRPC批量查询IP请求")

	infos, err := s.service.BatchQueryIP(req.Ips)
	if err != nil {
		s.logger.WithError(err).Error("批量查询IP失败")
		return &pb.BatchQueryIPResponse{
			Infos:     []*pb.IPInfo{},
			Timestamp: time.Now().Unix(),
		}, nil
	}

	protoInfos := make([]*pb.IPInfo, 0, len(infos))
	for _, info := range infos {
		protoInfos = append(protoInfos, convertToProtoIPInfo(info))
	}

	return &pb.BatchQueryIPResponse{
		Infos:     protoInfos,
		Timestamp: time.Now().Unix(),
	}, nil
}

// GetServiceStatus 获取服务状态
func (s *GRPCServer) GetServiceStatus(ctx context.Context, req *pb.GetServiceStatusRequest) (*pb.GetServiceStatusResponse, error) {
	status := s.service.GetServiceStatus()

	return &pb.GetServiceStatusResponse{
		Status:     status["status"].(string),
		Version:    status["version"].(string),
		Uptime:     int64(status["uptime"].(float64)),
		QueryCount: status["query_count"].(int64),
	}, nil
}

// convertToProtoIPInfo 转换为protobuf IPInfo
func convertToProtoIPInfo(info *ipquery.IPInfo) *pb.IPInfo {
	return &pb.IPInfo{
		Ip:           info.IP,
		Country:      info.Country,
		CountryCode:  info.CountryCode,
		Region:       info.Region,
		City:         info.City,
		District:     info.District,
		Isp:          info.ISP,
		Latitude:     float32(info.Latitude),
		Longitude:    float32(info.Longitude),
		Timezone:     info.Timezone,
		PostalCode:   info.PostalCode,
		IsValid:      info.IsValid,
		ErrorMessage: info.ErrorMessage,
	}
}
