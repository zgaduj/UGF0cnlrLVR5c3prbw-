package config

type SQLiteConfig struct {
	FilePath string
}

type AppConfig struct {
	ListeningPort int
}

type Config struct {
	APP *AppConfig
	DB  *SQLiteConfig
}

func GetConfig() *Config {
	return &Config{
		APP: &AppConfig{
			ListeningPort: 8080,
		},
		DB: &SQLiteConfig{
			FilePath: "_data/database.sqlite",
		},
	}
}
