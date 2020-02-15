package config

type SQLiteConfig struct {
	FilePath string
}

type AppConfig struct {
	ListeningPort int
	MaxBodySize   int64
}

type Config struct {
	APP *AppConfig
	DB  *SQLiteConfig
}

func GetConfig() *Config {
	return &Config{
		APP: &AppConfig{
			ListeningPort: 8080,
			MaxBodySize:   1 * 1024 * 1024,
		},
		DB: &SQLiteConfig{
			FilePath: "_data/database.sqlite",
		},
	}
}
