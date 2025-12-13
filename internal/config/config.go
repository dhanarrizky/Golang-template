package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string `mapstructure:"APP_ENV"`

	AppPort   string `mapstructure:"APP_PORT"`
	AppDomain string `mapstructure:"APP_DOMAIN"`
	AppName   string `mapstructure:"APP_NAME"`
	AppDebug  bool   `mapstructure:"APP_DEBUG"`

	DatabaseURL             string `mapstructure:"DATABASE_URL"`
	DatabaseMaxIdleConns    int    `mapstructure:"DATABASE_MAX_IDLE_CONNS"`
	DatabaseMaxOpenConns    int    `mapstructure:"DATABASE_MAX_OPEN_CONNS"`
	DatabaseConnMaxLifetime string `mapstructure:"DATABASE_CONN_MAX_LIFETIME"`

	RedisEnable   bool   `mapstructure:"REDIS_ENABLE"`
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
	RedisTLS      bool   `mapstructure:"REDIS_TLS"`

	JWTSecret    string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn string `mapstructure:"JWT_EXPIRES_IN"`
	JWTIssuer    string `mapstructure:"JWT_ISSUER"`
	JWTAudience  string `mapstructure:"JWT_AUDIENCE"`

	PasswordBcryptCost int `mapstructure:"PASSWORD_BCRYPT_COST"`

	CORSAllowedOrigins []string

	RateLimiterEnable bool `mapstructure:"RATE_LIMITER_ENABLE"`
	RateLimiterRPM    int  `mapstructure:"RATE_LIMITER_RPM"`

	CacheTTL string `mapstructure:"CACHE_TTL"`

	LogLevel string `mapstructure:"LOG_LEVEL"`

	SwaggerEnable   bool   `mapstructure:"SWAGGER_ENABLE"`
	SwaggerHost     string `mapstructure:"SWAGGER_HOST"`
	SwaggerBasePath string `mapstructure:"SWAGGER_BASE_PATH"`

	// Computed
	ParsedCacheTTL time.Duration
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	// Default values
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("DATABASE_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DATABASE_MAX_OPEN_CONNS", 50)
	viper.SetDefault("DATABASE_CONN_MAX_LIFETIME", "30m")
	viper.SetDefault("REDIS_ENABLE", true)
	viper.SetDefault("RATE_LIMITER_ENABLE", true)
	viper.SetDefault("CACHE_TTL", "10m")

	_ = viper.ReadInConfig() // optional file

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Parse CORS
	if origins := viper.GetString("CORS_ALLOWED_ORIGINS"); origins != "" {
		raw := strings.Split(origins, ",")
		for _, o := range raw {
			o = strings.TrimSpace(o)
			if o != "" {
				cfg.CORSAllowedOrigins = append(cfg.CORSAllowedOrigins, o)
			}
		}
	}

	// Parse duration safely
	ttl, err := time.ParseDuration(cfg.CacheTTL)
	if err != nil {
		cfg.ParsedCacheTTL = 10 * time.Minute
	} else {
		cfg.ParsedCacheTTL = ttl
	}

	return &cfg, nil
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "dev" || c.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "prod" || c.Environment == "production"
}
