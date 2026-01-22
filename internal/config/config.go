package config

import (
	"os"
)

type Config struct {
	MongoURI  string
	JWTSecret string
	Port      string
	AdminID   string
	EmailUser string
	EmailPass string
}

func Load() *Config {
	return &Config{
		MongoURI:  getEnv("MONGO_URI", "mongodb+srv://sapini8865:UFodVdiCQLsLELWI@cluster0.nrhp6.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
		Port:      getEnv("PORT", "8080"),
		AdminID:   getEnv("ADMIN_ID", ""),
		EmailUser: getEnv("EMAIL_USER", ""),
		EmailPass: getEnv("EMAIL_PASS", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
