package ipquery

import (
	"net"
	"strings"
)

// IPInfo IP信息结构体
type IPInfo struct {
	IP           string  `json:"ip"`
	Country      string  `json:"country"`
	CountryCode  string  `json:"country_code"`
	Region       string  `json:"region"`
	City         string  `json:"city"`
	District     string  `json:"district"`
	ISP          string  `json:"isp"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Timezone     string  `json:"timezone"`
	PostalCode   string  `json:"postal_code"`
	IsValid      bool    `json:"is_valid"`
	ErrorMessage string  `json:"error_message,omitempty"`
}

// QueryProvider IP查询提供者接口
type QueryProvider interface {
	Query(ip string) (*IPInfo, error)
	BatchQuery(ips []string) ([]*IPInfo, error)
	Close() error
}

// ValidateIP 验证IP地址
func ValidateIP(ip string) bool {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return false
	}

	// 检查IPv4
	if net.ParseIP(ip) != nil {
		return true
	}

	// 检查IPv6
	if strings.Contains(ip, ":") {
		_, _, err := net.ParseCIDR(ip + "/128")
		return err == nil
	}

	return false
}

// IsPrivateIP 检查是否为私有IP
func IsPrivateIP(ip string) bool {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return false
	}

	// IPv4私有地址范围
	privateIPv4 := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"169.254.0.0/16",
	}

	// IPv6私有地址范围
	privateIPv6 := []string{
		"fc00::/7",
		"fe80::/10",
		"::1/128",
	}

	allPrivate := append(privateIPv4, privateIPv6...)

	for _, cidr := range allPrivate {
		_, subnet, _ := net.ParseCIDR(cidr)
		if subnet.Contains(ipAddr) {
			return true
		}
	}

	return false
}

// IsPublicIP 检查是否为公网IP
func IsPublicIP(ip string) bool {
	return ValidateIP(ip) && !IsPrivateIP(ip)
}
