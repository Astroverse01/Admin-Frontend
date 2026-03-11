package config

import (
	"fmt"
	"os"
)

type Config struct {
	MongoURI   string
	JWTSecret  string
	Port       string
	AdminID    string
	EmailUser  string
	EmailPass  string
	SecretKey  string // for phoneNo decryption (AES)
	IV         string // hex-encoded 16-byte IV
	Salt       string
	Iterations int
	Keylen     int
	S3Bucket   string
	AWSRegion  string
}

func Load() *Config {
	iterations := 1
	if v := os.Getenv("ITERATIONS"); v != "" {
		if i, err := parseInt(v); err == nil {
			iterations = i
		}
	}
	keylen := 32
	if v := os.Getenv("KEYLEN"); v != "" {
		if i, err := parseInt(v); err == nil {
			keylen = i
		}
	}
	s3Bucket := os.Getenv("AWS_S3_BUCKET")
	if s3Bucket == "" {
		s3Bucket = getEnv("S3_BUCKET_NAME", "")
	}

	return &Config{
		MongoURI:   getEnv("MONGO_URI", "mongodb+srv://sapini8865:UFodVdiCQLsLELWI@cluster0.nrhp6.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
		Port:       getEnv("PORT", "8080"),
		AdminID:    getEnv("ADMIN_ID", ""),
		EmailUser:  getEnv("EMAIL_USER", ""),
		EmailPass:  getEnv("EMAIL_PASS", ""),
		SecretKey:  getEnv("SECRET_KEY", "cndsjvbjdsbvdfvgdkvbhhisuwrfj4947328ryu2ydb98yby8"),
		IV:         getEnv("IV", "d1ced75f02690f3d83a7e4e72c84a1ca"),
		Salt:       getEnv("SALT", "eb7ddfeab5349c80ab0290e03577f04f"),
		Iterations: iterations,
		Keylen:     keylen,
		S3Bucket:   s3Bucket,
		AWSRegion:  getEnv("AWS_REGION", ""),
	}
}

func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
