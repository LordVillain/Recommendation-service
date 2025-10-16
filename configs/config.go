package configs

import (
	"os"
	"time"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	Db DbConfig
	OAuth OAuthConfig
	// Kafka KafkaConfig
}

type DbConfig struct {
	Dsn string
}

type OAuthConfig struct {
	Secret string
	JWTTTL time.Duration
}

// type KafkaConfig struct {
// 	Brokers []string
// 	Topic   string
// 	GroupID string
// 	ClientID string
// }

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file, using default config", err.Error())
	}

	// brokersRaw := os.Getenv("KAFKA_BROKERS")
	// brokers := strings.Split(brokersRaw, ",")

	jwtTTLStr := os.Getenv("JWT_TTL")
	if jwtTTLStr == "" {
		jwtTTLStr = "15m"
	}
	jwtTTL, err := time.ParseDuration(jwtTTLStr)
	if err != nil {
		logger.Error("Invalid JWT_TTL, using default 1h", err.Error())
		jwtTTL = 15 * time.Minute
	}

	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		OAuth: OAuthConfig{
			Secret: os.Getenv("SECRET"),
			JWTTTL: jwtTTL,
		},
		// Kafka: KafkaConfig{
		// 	Brokers: brokers,
		// 	Topic:   os.Getenv("KAFKA_TOPIC"),
		// 	GroupID: os.Getenv("KAFKA_GROUP_ID"),
		// 	ClientID: os.Getenv("KAFKA_CLIENT_ID"),
		// },
	}
}
