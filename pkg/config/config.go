package config

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type LogConfig struct {
	Level    string
	Encoding string
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	// Load .env file if exists (ignore error if not found)
	_ = godotenv.Load()

	// Set default config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Enable environment variable override
	viper.AutomaticEnv()

	// Set default values
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Config file not found, using defaults and environment variables: %v", err)
	}

	var config Config

	// App config
	config.App = AppConfig{
		Name: viper.GetString("app.name"),
		Env:  viper.GetString("app.env"),
		Port: viper.GetString("app.port"),
	}

	// Database config
	config.Database = DatabaseConfig{
		Host:            viper.GetString("database.host"),
		Port:            viper.GetString("database.port"),
		User:            viper.GetString("database.user"),
		Password:        viper.GetString("database.password"),
		Name:            viper.GetString("database.name"),
		SSLMode:         viper.GetString("database.sslmode"),
		MaxOpenConns:    viper.GetInt("database.max_open_conns"),
		MaxIdleConns:    viper.GetInt("database.max_idle_conns"),
		ConnMaxLifetime: viper.GetDuration("database.conn_max_lifetime"),
	}

	// JWT config
	config.JWT = JWTConfig{
		Secret:     viper.GetString("jwt.secret"),
		Expiration: viper.GetDuration("jwt.expiration"),
	}

	// Log config
	config.Log = LogConfig{
		Level:    viper.GetString("log.level"),
		Encoding: viper.GetString("log.encoding"),
	}

	// Override with environment variables if present
	if appPort := viper.GetString("APP_PORT"); appPort != "" {
		config.App.Port = appPort
	}
	if dbHost := viper.GetString("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}
	if dbPort := viper.GetString("DB_PORT"); dbPort != "" {
		config.Database.Port = dbPort
	}
	if dbUser := viper.GetString("DB_USER"); dbUser != "" {
		config.Database.User = dbUser
	}
	if dbPassword := viper.GetString("DB_PASSWORD"); dbPassword != "" {
		config.Database.Password = dbPassword
	}
	if dbName := viper.GetString("DB_NAME"); dbName != "" {
		config.Database.Name = dbName
	}
	if jwtSecret := viper.GetString("JWT_SECRET"); jwtSecret != "" {
		config.JWT.Secret = jwtSecret
	}

	return &config, nil
}

func setDefaults() {
	// App defaults
	viper.SetDefault("app.name", "go-clean-boiler")
	viper.SetDefault("app.env", "development")
	viper.SetDefault("app.port", "8080")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.name", "go_clean_boiler")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 25)
	viper.SetDefault("database.conn_max_lifetime", 5*time.Minute)

	// JWT defaults
	viper.SetDefault("jwt.secret", "your-secret-key-change-this-in-production")
	viper.SetDefault("jwt.expiration", 24*time.Hour)

	// Log defaults
	viper.SetDefault("log.level", "debug")
	viper.SetDefault("log.encoding", "console")
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}
