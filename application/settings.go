package application

// The DatabaseURL field stores the URL of the database used by the application.
type Config struct {
	DatabaseURL string
}

// Initializes the DatabaseURL field with the value "forum.db".
func LoadConfig() Config {
	return Config{
		DatabaseURL: "forum.db",
	}
}
