package config

type DBConfig struct {
	//FilePath string
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type AppConfig struct {
	ListeningPort int
	MaxBodySize   int64
}

type Config struct {
	APP *AppConfig
	DB  *DBConfig
}

func GetConfig() *Config {
	return &Config{
		APP: &AppConfig{
			ListeningPort: 8080,
			MaxBodySize:   1 * 1024 * 1024,
		},
		DB: &DBConfig{
			//FilePath: "_data/database.sqlite",
			Host:     "db",
			Port:     3306,
			User:     "user",
			Password: "password",
			Database: "db",
		},
	}
}
