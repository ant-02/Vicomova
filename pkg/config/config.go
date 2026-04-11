package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig     `mapstructure:"server"`
	Database DatabaseConfig  `mapstructure:"database"`
	Redis    RedisConfig     `mapstructure:"redis"`
	Kafka    KafkaConfig     `mapstructure:"kafka"`
	JWT      JWTConfig       `mapstructure:"jwt"`
	Email    AliyunEmailConfig `mapstructure:"email"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

type AliyunEmailConfig struct {
	AccessKey    string `mapstructure:"access_key"`
	AccessSecret string `mapstructure:"access_secret"`
	AccountName  string `mapstructure:"account_name"`
	Region       string `mapstructure:"region"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host string `mapstructure:"host"`
	Port int `mapstructure:"port"`
	User string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName string `mapstructure:"dbname"`
	MaxOpenConns int `mapstructure:"max_open_conns"`
	MaxIdleConns int `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB int `mapstructure:"db"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Topic string `mapstructure:"topic"`
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.DBName)
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func Load(configPath string) (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	// 优先从指定路径加载
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// 从当前目录和上级目录查找
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("../config")
	}

	// 读取环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func LoadWithBase(basePath string) (*Config, error) {
	dir, _ := filepath.Split(basePath)
	if dir != "" {
		viper.AddConfigPath(dir)
	}
	return Load(basePath)
}

func DefaultAddr(serviceName string, port int) string {
	return fmt.Sprintf("%s:%d", serviceName, port)
}
