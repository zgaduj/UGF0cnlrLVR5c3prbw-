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

//type AppConfig struct {
//	ListeningPort int
//	MaxBodySize   int64
//}
//
//type WorkerConfig struct {
//	Timeout  int64
//	Interval int
//}
//
//type Config struct {
//	App    *AppConfig
//	DB     *DBConfig
//	Worker *WorkerConfig
//}

//func GetConfig() *Config {
//	return &Config{
//		App: &AppConfig{
//			ListeningPort: 8080,
//			MaxBodySize:   1 * 1024 * 1024,
//		},
//		Worker: &WorkerConfig{
//			Timeout:  5,
//			Interval: 2,
//		},
//		DB: &DBConfig{
//			//FilePath: "_data/database.sqlite",
//			Host:     "db",
//			Port:     3306,
//			User:     "user",
//			Password: "password",
//			Database: "db",
//		},
//	}
//}
