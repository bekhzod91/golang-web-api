package config

import (
	"crypto/tls"
	"github.com/myproject/api/pkg/env"
	"github.com/spf13/cast"
	"log/slog"
	"os"
	"strings"
	"time"
)

type Config struct {
	Environment string

	// API config
	APIHost        string
	APIPort        int
	AllowedOrigins []string

	// Database config
	PostgresHost       string
	PostgresPort       int
	PostgresUser       string
	PostgresPassword   string
	PostgresDatabase   string
	PostgresSchema     string
	PostgresTLSEnabled bool
	PostgresTLSConfig  *tls.Config

	MigrationPath string

	// Cache config
	RedisHost       string
	RedisPort       int
	RedisUser       string
	RedisPassword   string
	RedisDatabase   int
	RedisTLSEnabled bool
	RedisTLSConfig  *tls.Config

	// Pulsar
	PulsarHost      string
	PulsarPort      int
	PulsarAdminHost string
	PulsarAdminPort int

	// context timeout in seconds
	CtxTimeout time.Duration
	LogLevel   slog.Level
	LogJSON    bool

	// AWS
	AWSBaseEndpoint    string
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSBucket          string

	// Google Maps
	GoogleMapsKey string

	// API Key
	APIKey string

	// Pagination
	PaginationLimit int
}

func NewConfig() Config {
	env.LoadEnv()

	config := Config{}
	config.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))
	config.APIHost = cast.ToString(getOrReturnDefault("APP_HOST", "0.0.0.0"))
	config.APIPort = cast.ToInt(getOrReturnDefault("APP_PORT", "8000"))
	config.AllowedOrigins = []string{"http://localhost"}

	// PostgresSQL
	config.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	config.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	config.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "postgres"))
	config.PostgresSchema = cast.ToString(getOrReturnDefault("POSTGRES_SCHEMA", "public"))
	config.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	config.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "postgres"))
	config.PostgresTLSEnabled = cast.ToBool(getOrReturnDefault("POSTGRES_TLS_ENABLED", false))
	config.PostgresTLSConfig = nil

	if config.PostgresTLSEnabled {
		config.PostgresTLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	config.MigrationPath = cast.ToString(getOrReturnDefault("MIGRATION_PATH", ""))

	// Redis
	config.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "localhost"))
	config.RedisPort = cast.ToInt(getOrReturnDefault("REDIS_PORT", 6379))
	config.RedisUser = cast.ToString(getOrReturnDefault("REDIS_USER", "default"))
	config.RedisPassword = cast.ToString(getOrReturnDefault("REDIS_PASSWORD", "postgres"))
	config.RedisDatabase = cast.ToInt(getOrReturnDefault("REDIS_NAME", 0))
	config.RedisTLSEnabled = cast.ToBool(getOrReturnDefault("REDIS_TLS_ENABLED", false))
	config.RedisTLSConfig = nil

	if config.RedisTLSEnabled {
		config.RedisTLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	// Pulsar
	config.PulsarHost = cast.ToString(getOrReturnDefault("PULSAR_HOST", "localhost"))
	config.PulsarPort = cast.ToInt(getOrReturnDefault("PULSAR_PORT", 6650))
	config.PulsarAdminHost = cast.ToString(getOrReturnDefault("PULSAR_ADMIN_HOST", "localhost"))
	config.PulsarAdminPort = cast.ToInt(getOrReturnDefault("PULSAR_ADMIN_PORT", 8080))

	config.CtxTimeout = 3 * time.Second
	config.LogLevel = getLogLevelFromEnv("LOG_LEVEL", slog.LevelDebug)
	config.LogJSON = cast.ToBool(getOrReturnDefault("LOG_JSON", false))

	// AWS S3
	config.AWSBaseEndpoint = cast.ToString(getOrReturnDefault("AWS_BASE_ENDPOINT", "https://dev-s3.myproject.zip24.com"))
	config.AWSRegion = cast.ToString(getOrReturnDefault("AWS_REGION", "us-east-1"))
	config.AWSAccessKeyID = cast.ToString(getOrReturnDefault("AWS_ACCESS_KEY_ID", "24vz49RGoAqopR84fTJe"))
	config.AWSSecretAccessKey = cast.ToString(getOrReturnDefault("AWS_SECRET_ACCESS_KEY", "QxVlmAdGN5n4xt0YoM8bccFNSaAWZTOh5Yj0Vk4r"))
	config.AWSBucket = cast.ToString(getOrReturnDefault("AWS_BUCKET", "myproject"))

	config.PaginationLimit = cast.ToInt(getOrReturnDefault("PAGINATION_LIMIT", 10))

	return config
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}

func getLogLevelFromEnv(key string, defaultValue slog.Level) slog.Level {
	_, exists := os.LookupEnv(key)
	if exists {
		switch strings.ToUpper(os.Getenv(key)) {
		case "DEBUG":
			return slog.LevelDebug
		case "INFO":
			return slog.LevelInfo
		case "WARN":
			return slog.LevelWarn
		case "ERROR":
			return slog.LevelError
		default:
			return defaultValue
		}
	}

	return defaultValue
}
