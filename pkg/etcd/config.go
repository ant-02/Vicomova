package etcd

// EtcdConfig holds etcd connection configuration
type EtcdConfig struct {
	Endpoints []string `mapstructure:"endpoints"`
	Username  string   `mapstructure:"username"`
	Password  string   `mapstructure:"password"`
	Namespace string   `mapstructure:"namespace"` // 如 "vicomova"
}
