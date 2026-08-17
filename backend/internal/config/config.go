package config

import (
	"os"
	"strconv"
)

// Config 集中解析全部环境变量配置。
type Config struct {
	AppName    string
	AppEnv     string
	ServerPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	JWTSecret        string
	JWTAccessExpire  int
	JWTRefreshExpire int

	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string
	MinIOUseSSL    bool

	CORSOrigins string
}

// Load 从环境变量加载配置，未设置时使用默认值，保证任意目录下可启动。
func Load() *Config {
	return &Config{
		AppName:    getEnv("APP_NAME", "flowershop"),
		AppEnv:     getEnv("APP_ENV", "development"),
		ServerPort: getEnv("SERVER_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "flowershop_user"),
		DBPassword: getEnv("DB_PASSWORD", "flowershop_pwd"),
		DBName:     getEnv("DB_NAME", "flowershop_db"),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),

		JWTSecret:        getEnv("JWT_SECRET", "change_me_to_a_long_random_string"),
		JWTAccessExpire:  getEnvInt("JWT_ACCESS_EXPIRE", 86400),
		JWTRefreshExpire: getEnvInt("JWT_REFRESH_EXPIRE", 604800),

		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey: getEnv("MINIO_ROOT_USER", "minioadmin"),
		MinIOSecretKey: getEnv("MINIO_ROOT_PASSWORD", "minioadmin"),
		MinIOBucket:    getEnv("MINIO_BUCKET", "flowershop"),
		MinIOUseSSL:    getEnvBool("MINIO_USE_SSL", false),

		CORSOrigins: getEnv("CORS_ORIGINS", "*"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getEnvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}
