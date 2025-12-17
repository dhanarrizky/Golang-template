package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// =========================
// Main Application Config
// =========================
type Config struct {
	// =========================
	// Application Environment
	// =========================
	Environment string `mapstructure:"APP_ENV"`

	AppPort   string `mapstructure:"APP_PORT"`
	AppDomain string `mapstructure:"APP_DOMAIN"`
	AppName   string `mapstructure:"APP_NAME"`
	AppDebug  bool   `mapstructure:"APP_DEBUG"`

	// =========================
	// Database Configuration
	// =========================
	DatabaseURL             string `mapstructure:"DATABASE_URL"`
	DatabaseMaxIdleConns    int    `mapstructure:"DATABASE_MAX_IDLE_CONNS"`
	DatabaseMaxOpenConns    int    `mapstructure:"DATABASE_MAX_OPEN_CONNS"`
	DatabaseConnMaxLifetime string `mapstructure:"DATABASE_CONN_MAX_LIFETIME"`

	// =========================
	// Redis Configuration
	// =========================
	RedisEnable   bool   `mapstructure:"REDIS_ENABLE"`
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     int    `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
	RedisTLS      bool   `mapstructure:"REDIS_TLS"`

	// =========================
	// Authentication - JWT
	// =========================
	JWTSecret           string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn        string `mapstructure:"JWT_EXPIRES_IN"`
	JWTRefreshExpiresIn string `mapstructure:"JWT_REFRESH_EXPIRES_IN"`
	JWTIssuer           string `mapstructure:"JWT_ISSUER"`
	JWTAudience         string `mapstructure:"JWT_AUDIENCE"`

	// =========================
	// Security - Password (Argon2id)
	// =========================
	Password PasswordConfig

	// =========================
	// Security - Password (Argon2id)
	// =========================
	SecretToken string `mapstructure:"EMAIL_VERIFICATION_TOKEN_SECRET"`

	// =========================
	// Security - ID (HMAC-Signed Identifier)
	// =========================
	SecretKey string `mapstructure:"SECRET_KEY"`

	// =========================
	// CORS Configuration
	// =========================
	CORSAllowedOrigins []string

	// =========================
	// Rate Limiter
	// =========================
	RateLimiterEnable bool `mapstructure:"RATE_LIMITER_ENABLE"`
	RateLimiterRPM    int  `mapstructure:"RATE_LIMITER_RPM"`

	// =========================
	// Cache
	// =========================
	CacheTTL       string `mapstructure:"CACHE_TTL"`
	ParsedCacheTTL time.Duration

	// =========================
	// Logging
	// =========================
	LogLevel string `mapstructure:"LOG_LEVEL"`

	// =========================
	// Swagger / API Docs
	// =========================
	SwaggerEnable   bool   `mapstructure:"SWAGGER_ENABLE"`
	SwaggerHost     string `mapstructure:"SWAGGER_HOST"`
	SwaggerBasePath string `mapstructure:"SWAGGER_BASE_PATH"`

	// =========================
	// SMTP Email Configuration
	// =========================
	SMTPEngine      string `mapstructure:"SMTP_ENGINE"` // contoh: "smtp" atau "ses" di masa depan
	SMTPHost        string `mapstructure:"SMTP_HOST"`
	SMTPPort        string `mapstructure:"SMTP_PORT"` // biasanya "587" atau "465"
	SMTPUsername    string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword    string `mapstructure:"SMTP_PASSWORD"`
	SMTPFromAddress string `mapstructure:"SMTP_FROM_ADDRESS"` // email pengirim
	SMTPFromName    string `mapstructure:"SMTP_FROM_NAME"`    // nama pengirim (opsional)
}

// =========================
// Argon2id Password Config
// =========================
type PasswordConfig struct {
	Memory      uint32 `mapstructure:"PASSWORD_ARGON2_MEMORY"` // KB
	Iterations  uint32 `mapstructure:"PASSWORD_ARGON2_ITERATIONS"`
	Parallelism uint8  `mapstructure:"PASSWORD_ARGON2_PARALLELISM"`
	SaltLength  uint32 `mapstructure:"PASSWORD_ARGON2_SALT_LENGTH"`
	KeyLength   uint32 `mapstructure:"PASSWORD_ARGON2_KEY_LENGTH"`

	// Pepper management
	Peppers              map[int]string
	CurrentPepperVersion int `mapstructure:"PASSWORD_PEPPER_VERSION"`
}

// =========================
// Load Config
// =========================
func LoadConfig() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// =========================
	// Defaults
	// =========================
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")

	viper.SetDefault("DATABASE_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DATABASE_MAX_OPEN_CONNS", 50)
	viper.SetDefault("DATABASE_CONN_MAX_LIFETIME", "30m")

	viper.SetDefault("REDIS_ENABLE", true)
	viper.SetDefault("RATE_LIMITER_ENABLE", true)

	viper.SetDefault("CACHE_TTL", "10m")

	viper.SetDefault("SECRET_KEY", "secret-key-default")

	// Argon2id defaults (recommended)
	viper.SetDefault("PASSWORD_ARGON2_MEMORY", 64*1024) // 64 MB
	viper.SetDefault("PASSWORD_ARGON2_ITERATIONS", 4)
	viper.SetDefault("PASSWORD_ARGON2_PARALLELISM", 1)
	viper.SetDefault("PASSWORD_ARGON2_SALT_LENGTH", 16)
	viper.SetDefault("PASSWORD_ARGON2_KEY_LENGTH", 32)
	viper.SetDefault("PASSWORD_PEPPER_VERSION", 1)

	// secret token
	viper.SetDefault("EMAIL_VERIFICATION_TOKEN_SECRET", "super-long-random-secret")

	// =========================
	// Defaults untuk SMTP
	// =========================
	viper.SetDefault("SMTP_ENGINE", "smtp")
	viper.SetDefault("SMTP_HOST", "smtp.gmail.com")
	viper.SetDefault("SMTP_PORT", "587")
	viper.SetDefault("SMTP_FROM_ADDRESS", "no-reply@yourapp.com")
	viper.SetDefault("SMTP_FROM_NAME", "Your App")

	_ = viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// =========================
	// Parse CORS
	// =========================
	if origins := viper.GetString("CORS_ALLOWED_ORIGINS"); origins != "" {
		raw := strings.Split(origins, ",")
		for _, o := range raw {
			o = strings.TrimSpace(o)
			if o != "" {
				cfg.CORSAllowedOrigins = append(cfg.CORSAllowedOrigins, o)
			}
		}
	}

	// =========================
	// Parse Cache TTL
	// =========================
	ttl, err := time.ParseDuration(cfg.CacheTTL)
	if err != nil {
		cfg.ParsedCacheTTL = 10 * time.Minute
	} else {
		cfg.ParsedCacheTTL = ttl
	}

	// =========================
	// Load Peppers
	// =========================
	cfg.Password.Peppers = map[int]string{}
	if p := viper.GetString("PASSWORD_PEPPER_V1"); p != "" {
		cfg.Password.Peppers[1] = p
	}
	if p := viper.GetString("PASSWORD_PEPPER_V2"); p != "" {
		cfg.Password.Peppers[2] = p
	}

	return &cfg, nil
}

// =========================
// Helpers
// =========================
func (c *Config) IsDevelopment() bool {
	return c.Environment == "dev" || c.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "prod" || c.Environment == "production"
}
