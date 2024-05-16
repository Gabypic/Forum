package application

type Config struct {
	DatabaseURL string
}

func LoadConfig() Config {
	return Config{
		DatabaseURL: "forum.db",
	}
}
