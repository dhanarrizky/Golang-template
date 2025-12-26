package config

import (
	"strings"

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
	// key - PublicIDCodec
	// =========================
	PublicIdAesKey string `mapstructure:"PUBLIC_ID_AES_KEY"`

	// =========================
	// key - PublicIDCodec
	// =========================
	TokenHmacSecret string `mapstructure:"TOKEN_HMAC_SECRET"`
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

	viper.SetDefault("RATE_LIMITER_ENABLE", true)

	viper.SetDefault("SECRET_KEY", "secret-key-default")

	// Argon2id defaults (recommended)
	viper.SetDefault("PASSWORD_ARGON2_MEMORY", 64*1024) // 64 MB
	viper.SetDefault("PASSWORD_ARGON2_ITERATIONS", 4)
	viper.SetDefault("PASSWORD_ARGON2_PARALLELISM", 1)
	viper.SetDefault("PASSWORD_ARGON2_SALT_LENGTH", 16)
	viper.SetDefault("PASSWORD_ARGON2_KEY_LENGTH", 32)
	viper.SetDefault("PASSWORD_PEPPER_VERSION", 1)

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
