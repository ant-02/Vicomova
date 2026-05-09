package config

// ConfigChangeCallback is called when config changes
type ConfigChangeCallback func(*Config)

type Database struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Kafka struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

type JWT struct {
	Secret string `yaml:"secret"`
}

type Email struct {
	AccessKey    string `yaml:"access_key"`
	AccessSecret string `yaml:"access_secret"`
	AccountName  string `yaml:"account_name"`
	Region       string `yaml:"region"`
	APIURL       string `yaml:"api_url"`
}

type Service struct {
	Name string `yaml:"name"`
	Addr string `yaml:"addr"`
}

type Config struct {
	Database Database           `yaml:"database"`
	Redis    Redis              `yaml:"redis"`
	Kafka    Kafka              `yaml:"kafka"`
	JWT      JWT                `yaml:"jwt"`
	Email    Email              `yaml:"email"`
	Services map[string]Service `yaml:"services"`
	Service  Service
}
