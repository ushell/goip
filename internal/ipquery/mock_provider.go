package ipquery

import (
	"fmt"
	"math/rand"
	"time"
)

// MockProvider 模拟IP查询提供者
type MockProvider struct {
	initialized bool
}

// NewMockProvider 创建新的模拟提供者
func NewMockProvider() *MockProvider {
	return &MockProvider{
		initialized: true,
	}
}

// Query 查询单个IP地址信息
func (m *MockProvider) Query(ip string) (*IPInfo, error) {
	if !m.initialized {
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

	// 模拟真实数据
	rand.Seed(time.Now().UnixNano())
	countries := []string{"中国", "美国", "日本", "韩国", "德国", "英国", "法国", "加拿大"}
	countryCodes := []string{"CN", "US", "JP", "KR", "DE", "GB", "FR", "CA"}
	regions := []string{"北京", "上海", "广东", "江苏", "浙江", "四川", "湖北", "河南"}
	cities := []string{"北京", "上海", "广州", "深圳", "杭州", "南京", "成都", "武汉"}
	isps := []string{"中国电信", "中国联通", "中国移动", "中国铁通", "教育网"}

	idx := rand.Intn(len(countries))

	return &IPInfo{
		IP:          ip,
		Country:     countries[idx],
		CountryCode: countryCodes[idx],
		Region:      regions[rand.Intn(len(regions))],
		City:        cities[rand.Intn(len(cities))],
		District:    fmt.Sprintf("%s区", cities[rand.Intn(len(cities))]),
		ISP:         isps[rand.Intn(len(isps))],
		Latitude:    30.0 + rand.Float64()*20.0,
		Longitude:   100.0 + rand.Float64()*20.0,
		Timezone:    "Asia/Shanghai",
		PostalCode:  fmt.Sprintf("%06d", 100000+rand.Intn(900000)),
		IsValid:     true,
	}, nil
}

// BatchQuery 批量查询IP地址信息
func (m *MockProvider) BatchQuery(ips []string) ([]*IPInfo, error) {
	if !m.initialized {
		return nil, fmt.Errorf("provider not initialized")
	}

	results := make([]*IPInfo, 0, len(ips))
	for _, ip := range ips {
		info, err := m.Query(ip)
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

// Close 关闭提供者
func (m *MockProvider) Close() error {
	m.initialized = false
	return nil
}
