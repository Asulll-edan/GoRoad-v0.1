package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppEnv  string `mapstructure:"APP_ENV"`
	APIPort string `mapstructure:"API_PORT"`
	WSPort  string `mapstructure:"WS_PORT"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`

	RedisURL  string `mapstructure:"REDIS_URL"`
	NATSURL   string `mapstructure:"NATS_URL"`
	MinioEndpoint string `mapstructure:"MINIO_ENDPOINT"`

	JWTSecret     string `mapstructure:"JWT_SECRET"`
	JWTExpiryHour int    `mapstructure:"JWT_EXPIRY_HOUR"`

	GeminiAPIKey string `mapstructure:"GEMINI_API_KEY"`

	LiveKitAPIKey    string `mapstructure:"LIVEKIT_API_KEY"`
	LiveKitAPISecret string `mapstructure:"LIVEKIT_API_SECRET"`
	LiveKitHost      string `mapstructure:"LIVEKIT_HOST"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		AppEnv:  getenv("APP_ENV", "development"),
		APIPort: getenv("API_PORT", "8080"),
		WSPort:  getenv("WS_PORT", "8081"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),

		RedisURL:      os.Getenv("REDIS_URL"),
		NATSURL:       os.Getenv("NATS_URL"),
		MinioEndpoint: os.Getenv("MINIO_ENDPOINT"),

		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiryHour: getint("JWT_EXPIRY_HOUR", 72),

		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),

		LiveKitAPIKey:    os.Getenv("LIVEKIT_API_KEY"),
		LiveKitAPISecret: os.Getenv("LIVEKIT_API_SECRET"),
		LiveKitHost:      os.Getenv("LIVEKIT_HOST"),
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getint(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func (c *Config) DBDSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=disable TimeZone=Asia/Jakarta"
}

func (c *Config) NATSUrl() string {
	if !strings.Contains(c.NATSURL, "://") {
		return "nats://" + c.NATSURL
	}
	return c.NATSURL
}
