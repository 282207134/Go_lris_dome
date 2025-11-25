package config

import (
	"os"
	"strconv"
)

// Config 应用程序配置结构体
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	JWT      JWTConfig      `json:"jwt"`
	Log      LogConfig      `json:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string `json:"port"`
	Mode         string `json:"mode"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSL      string `json:"ssl"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret         string `json:"secret"`
	ExpirationTime int    `json:"expiration_time"`
	Issuer         string `json:"issuer"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	Output string `json:"output"`
}

var appConfig *Config

// GetConfig 获取应用程序配置
func GetConfig() *Config {
	if appConfig == nil {
		appConfig = loadConfig()
	}
	return appConfig
}

// loadConfig 加载配置
func loadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Mode:         getEnv("GIN_MODE", "debug"),
			ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 30),
			WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 30),
		},
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "sqlite"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			Database: getEnv("DB_NAME", "iris_sample.db"),
			Username: getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			SSL:      getEnv("DB_SSL", "disable"),
		},
		JWT: JWTConfig{
			Secret:         getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
			ExpirationTime: getEnvAsInt("JWT_EXPIRATION_TIME", 24*60*60), // 24小时
			Issuer:         getEnv("JWT_ISSUER", "iris-sample-project"),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
			Output: getEnv("LOG_OUTPUT", "stdout"),
		},
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为 int 类型
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool 获取环境变量并转换为 bool 类型
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}