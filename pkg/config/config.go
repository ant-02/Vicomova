package config

type Database struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	User         string `json:"user"`
	Password     string `json:"password"`
	DBName       string `json:"dbname"`
	MaxOpenConns int    `json:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
}

type Redis struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type Kafka struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
}

type JWT struct {
	Secret string `json:"secret"`
}

type Email struct {
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
	AccountName  string `json:"account_name"`
	Region       string `json:"region"`
	APIURL       string `json:"api_url"`
}

type Etcd struct {
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Namespace string   `json:"namespace"`
}

type Service struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type Config struct {
	Database Database          `json:"database"`
	Redis    Redis             `json:"redis"`
	Kafka    Kafka             `json:"kafka"`
	JWT      JWT               `json:"jwt"`
	Email    Email             `json:"email"`
	Etcd     Etcd              `json:"etcd"`
	Services map[string]Service `json:"services"`
	Service  Service
}
