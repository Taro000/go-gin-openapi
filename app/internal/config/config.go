package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port     string
	DB       DBConfig
	Firebase FirebaseConfig
	Auth     AuthConfig
}

type DBConfig struct {
	DSNEnv   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func (c DBConfig) DSN() string {
	if c.DSNEnv != "" {
		return c.DSNEnv
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Asia%%2FTokyo",
		c.User, c.Password, c.Host, c.Port, c.Name,
	)
}

type FirebaseConfig struct {
	APIKey string

	// Admin SDK credentials
	ProjectID            string
	CredentialsFile      string
	ServiceAccountJSON   string
	AuthEmulatorHostport string
}

type AuthConfig struct {
	Bypass bool
}

func Load() Config {
	return Config{
		Port: env("PORT", "8080"),
		DB: DBConfig{
			DSNEnv:   os.Getenv("DB_DSN"),
			Host:     env("DB_HOST", "db"),
			Port:     env("DB_PORT", "3306"),
			User:     env("DB_USER", "root"),
			Password: env("DB_PASSWORD", "root"),
			Name:     env("DB_NAME", "go-gin-webapi"),
		},
		Firebase: FirebaseConfig{
			APIKey:               os.Getenv("FIREBASE_API_KEY"),
			ProjectID:            os.Getenv("FIREBASE_PROJECT_ID"),
			CredentialsFile:      os.Getenv("FIREBASE_CREDENTIALS_FILE"),
			ServiceAccountJSON:   os.Getenv("FIREBASE_SERVICE_ACCOUNT_JSON"),
			AuthEmulatorHostport: os.Getenv("FIREBASE_AUTH_EMULATOR_HOST"),
		},
		Auth: AuthConfig{
			Bypass: envBool("AUTH_BYPASS", false),
		},
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return def
}


