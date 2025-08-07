package configs

import (
	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	Log      LogConfig      `yaml:"log"`
	Upload   UploadConfig   `yaml:"upload"`
	Security SecurityConfig `yaml:"security"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MongoDB MongoDBConfig `yaml:"mongodb"`
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	URI         string `yaml:"uri"`
	Database    string `yaml:"database"`
	Timeout     string `yaml:"timeout"`
	MaxPoolSize int    `yaml:"max_pool_size"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret           string `yaml:"secret"`
	ExpiresIn        string `yaml:"expires_in"`
	RefreshExpiresIn string `yaml:"refresh_expires_in"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`       // 日志级别
	Format     string `yaml:"format"`      // 日志格式
	Output     string `yaml:"output"`      // 输出方式
	FilePath   string `yaml:"file_path"`   // 日志文件路径
	MaxSize    int    `yaml:"max_size"`    // 单个文件最大大小(MB)
	MaxAge     int    `yaml:"max_age"`     // 文件保留天数
	MaxBackups int    `yaml:"max_backups"` // 最大备份文件数
	Compress   bool   `yaml:"compress"`    // 是否压缩旧文件
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	MaxSize      string   `yaml:"max_size"`
	AllowedTypes []string `yaml:"allowed_types"`
	Path         string   `yaml:"path"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	PasswordMinLength int    `yaml:"password_min_length"`
	MaxLoginAttempts  int    `yaml:"max_login_attempts"`
	LockoutDuration   string `yaml:"lockout_duration"`
}

var AppConfig *Config

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath(".")

	// 设置环境变量前缀
	viper.SetEnvPrefix("INSURANCE")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// 修复JWT配置 - 手动从viper获取，确保不会丢失
	config.JWT.ExpiresIn = viper.GetString("jwt.expires_in")
	config.JWT.RefreshExpiresIn = viper.GetString("jwt.refresh_expires_in")

	return &config, nil
}
