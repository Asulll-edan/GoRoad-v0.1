package config

import (
	"github.com/spf13/viper"
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
	viper.SetDefault("API_PORT", "8080")
	viper.SetDefault("WS_PORT", "8081")
	viper.SetDefault("JWT_EXPIRY_HOUR", 72)
	viper.SetDefault("APP_ENV", "development")

	viper.AutomaticEnv()

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
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
