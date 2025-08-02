package ipquery

import (
	"fmt"
	"strings"

	ip2region "github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

// IP2RegionProvider 基于ip2region.xdb的真实IP查询提供者
type IP2RegionProvider struct {
	db          *ip2region.Searcher
	initialized bool
}

// NewIP2RegionProvider 创建新的基于ip2region.xdb的查询提供者
func NewIP2RegionProvider(dbPath string) (*IP2RegionProvider, error) {
	db, err := ip2region.NewWithFileOnly(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load ip2region database: %w", err)
	}

	return &IP2RegionProvider{
		db:          db,
		initialized: true,
	}, nil
}

// Query 查询单个IP地址信息
func (p *IP2RegionProvider) Query(ip string) (*IPInfo, error) {
	if !p.initialized {
		return nil, fmt.Errorf("provider not initialized")
	}

	if !ValidateIP(ip) {
		return &IPInfo{
			IP:           ip,
			IsValid:      false,
			ErrorMessage: "无效的IP地址格式",
		}, nil
	}

	if IsPrivateIP(ip) {
		return &IPInfo{
			IP:          ip,
			IsValid:     true,
			Country:     "局域网",
			CountryCode: "LAN",
			Region:      "局域网",
			City:        "局域网",
			ISP:         "局域网",
			Latitude:    0,
			Longitude:   0,
			Timezone:    "UTC",
			PostalCode:  "000000",
		}, nil
	}

	// 使用ip2region查询真实数据
	info, err := p.db.SearchByStr(ip)
	if err != nil {
		return &IPInfo{
			IP:           ip,
			IsValid:      false,
			ErrorMessage: fmt.Sprintf("IP查询失败: %v", err),
		}, nil
	}

	// 解析ip2region返回的数据格式
	// 格式: 国家|区域|省份|城市|ISP
	parts := strings.Split(info, "|")

	// 处理空值
	for i := range parts {
		if parts[i] == "0" {
			parts[i] = ""
		}
	}

	// 确保有足够的字段
	for len(parts) < 5 {
		parts = append(parts, "")
	}

	country := parts[0]
	region := parts[2]
	city := parts[3]
	isp := parts[4]

	// 设置国家代码
	countryCode := getCountryCode(country)

	return &IPInfo{
		IP:          ip,
		Country:     country,
		CountryCode: countryCode,
		Region:      region,
		City:        city,
		District:    "", // ip2region不提供区县信息
		ISP:         isp,
		Latitude:    0,  // ip2region不提供经纬度
		Longitude:   0,  // ip2region不提供经纬度
		Timezone:    "", // ip2region不提供时区
		PostalCode:  "", // ip2region不提供邮政编码
		IsValid:     true,
	}, nil
}

// BatchQuery 批量查询IP地址信息
func (p *IP2RegionProvider) BatchQuery(ips []string) ([]*IPInfo, error) {
	if !p.initialized {
		return nil, fmt.Errorf("provider not initialized")
	}

	results := make([]*IPInfo, 0, len(ips))
	for _, ip := range ips {
		info, err := p.Query(ip)
		if err != nil {
			info = &IPInfo{
				IP:           ip,
				IsValid:      false,
				ErrorMessage: err.Error(),
			}
		}
		results = append(results, info)
	}

	return results, nil
}

// Close 关闭提供者，释放资源
func (p *IP2RegionProvider) Close() error {
	if p.db != nil {
		p.db.Close()
	}
	p.initialized = false
	return nil
}

// getCountryCode 根据国家名称获取国家代码
func getCountryCode(country string) string {
	switch country {
	case "中国":
		return "CN"
	case "美国":
		return "US"
	case "日本":
		return "JP"
	case "韩国":
		return "KR"
	case "德国":
		return "DE"
	case "英国":
		return "GB"
	case "法国":
		return "FR"
	case "加拿大":
		return "CA"
	default:
		return "未知" // 未知国家
	}
}
