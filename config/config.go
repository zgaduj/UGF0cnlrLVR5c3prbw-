package config

type DBConfig struct {
	//FilePath string
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASS"`
	Database string `env:"DB_BASE"`
}

type APIConfig struct {
	ListeningPort int   `env:"APP_ENDPOINT_PORT" envDefault:"8080"`
	MaxBodySize   int64 `env:"APP_HTTP_MAX_BODY_SIZE" envDefault:"1"`
}

type WorkerConfig struct {
	Timeout          int64 `env:"TIMEOUT_HTTP_REQUEST" envDefault:"5"`
	IntervalCheckURL int   `env:"INTERVAL_CHECK_URL" envDefault:"1"`
}
