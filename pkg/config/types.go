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

type SASL struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type KafkaTopic struct {
	Name           string `yaml:"name"`
	Group          string `yaml:"group"`
	Partitions     int    `yaml:"partitions"`
	Replication    int    `yaml:"replication"`
	RetentionHours int    `yaml:"retention_hours"`
}

type Kafka struct {
	Brokers []string              `yaml:"brokers"`
	SASL    SASL                  `yaml:"sasl"`
	Topics  map[string]KafkaTopic `yaml:"topics"`
}

type JWT struct {
	Secret string `yaml:"secret"`
}

type AliyunEmail struct {
	AccessKey    string `yaml:"access-key"`
	AccessSecret string `yaml:"access-secret"`
	AccountName  string `yaml:"account-name"`
	Region       string `yaml:"region"`
	APIURL       string `yaml:"api-url"`
}

type Email struct {
	Aliyun AliyunEmail `yaml:"aliyun"`
}

type QiniuOSS struct {
	AccessKey  string `yaml:"access-key"`
	SecretKey  string `yaml:"secret-key"`
	Bucket     string `yaml:"bucket"`
	Domain     string `yaml:"domain"`
	UploadHost string `yaml:"upload-host"` // 七牛云上传地址，如 https://up-z2.qiniup.com
}

type OSS struct {
	Qiniu QiniuOSS `yaml:"qiniu"`
}

type OpenAI struct {
	APIKey string `yaml:"api_key"`
	APIURL string `yaml:"api_url"` // Embedding API 地址，不填则用默认 OpenAI
	Model  string `yaml:"model"`   // Embedding 模型，不填则用默认 text-embedding-3-small
}

type Qdrant struct {
	Addr string `yaml:"addr"`
}

type Service struct {
	Name          string `yaml:"name"`
	Addr          string `yaml:"addr"`
	WebSocketAddr string `yaml:"websocket_addr"`
}

type Config struct {
	Database Database           `yaml:"database"`
	Redis    Redis              `yaml:"redis"`
	Kafka    Kafka              `yaml:"kafka"`
	JWT      JWT                `yaml:"jwt"`
	Email    Email              `yaml:"email"`
	OSS      OSS                `yaml:"oss"`
	Services map[string]Service `yaml:"services"`
	Service  Service            `yaml:"service"`
	OpenAI   OpenAI             `yaml:"openai"`
	Qdrant   Qdrant             `yaml:"qdrant"`
}
