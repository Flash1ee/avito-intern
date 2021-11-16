package configs

type RepositoryConnections struct {
	DataBaseUrl string `toml:"database_url"`
}

type Config struct {
	BindAddr         string                `toml:"bind_addr"`
	LogLevel         string                `toml:"log_level"`
	LogAddr          string                `toml:"log_path"`
	CurrencyAPI      string                `toml:"currency_api"`
	ServerRepository RepositoryConnections `toml:"server"`
}

func NewConfig() *Config {
	return &Config{}
}
