package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config 全局配置结构体
type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Logging     LoggingConfig     `mapstructure:"logging"`
	IPDatabase  IPDatabaseConfig  `mapstructure:"ip_database"`
	Cache       CacheConfig       `mapstructure:"cache"`
	Metrics     MetricsConfig     `mapstructure:"metrics"`
	HealthCheck HealthCheckConfig `mapstructure:"health_check"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	HTTP HTTPConfig `mapstructure:"http"`
	GRPC GRPCConfig `mapstructure:"grpc"`
}

// HTTPConfig HTTP服务器配置
type HTTPConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// GRPCConfig gRPC服务器配置
type GRPCConfig struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	MaxConnectionIdle time.Duration `mapstructure:"max_connection_idle"`
	MaxConnectionAge  time.Duration `mapstructure:"max_connection_age"`
	Timeout           time.Duration `mapstructure:"timeout"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// IPDatabaseConfig IP数据库配置
type IPDatabaseConfig struct {
	Type           string        `mapstructure:"type"`
	Path           string        `mapstructure:"path"`
	CacheSize      int           `mapstructure:"cache_size"`
	AutoReload     bool          `mapstructure:"auto_reload"`
	ReloadInterval time.Duration `mapstructure:"reload_interval"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enabled bool          `mapstructure:"enabled"`
	Type    string        `mapstructure:"type"`
	TTL     time.Duration `mapstructure:"ttl"`
	MaxSize int           `mapstructure:"max_size"`
}

// MetricsConfig 监控配置
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
	Port    int    `mapstructure:"port"`
}

// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
	Enabled bool          `mapstructure:"enabled"`
	Path    string        `mapstructure:"path"`
	Timeout time.Duration `mapstructure:"timeout"`
}

// Load 加载配置
func Load(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	// 设置默认值
	viper.SetDefault("server.http.host", "0.0.0.0")
	viper.SetDefault("server.http.port", 8080)
	viper.SetDefault("server.grpc.host", "0.0.0.0")
	viper.SetDefault("server.grpc.port", 50051)
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("ip_database.type", "local")
	viper.SetDefault("cache.enabled", true)
	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("health_check.enabled", true)

	// 读取环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
