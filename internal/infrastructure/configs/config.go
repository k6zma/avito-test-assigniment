package configs

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/k6zma/avito-test-assigniment/pkg/validators"
)

const (
	configSubName     = "config"
	configsFileType   = "yaml"
	configsFolderPath = "configs/"
)

type Config struct {
	ServerCfg   *Server           `mapstructure:"server"   validate:"required"`
	PostgresCfg *PostgresDatabase `mapstructure:"postgres" validate:"required"`
	LoggingCfg  *Logging          `mapstructure:"logging"  validate:"required"`
}

type Server struct {
	AppName         string        `mapstructure:"app_name"         validate:"required,printascii"`
	Host            string        `mapstructure:"host"             validate:"required,hostname|ip"`
	Port            int           `mapstructure:"port"             validate:"required,min=1,max=65535"`
	Concurrency     int           `mapstructure:"concurrency"      validate:"required,min=1,max=65535"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"     validate:"required,gt=0"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"    validate:"required,gt=0"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"     validate:"required,gt=0"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" validate:"required,gt=0"`
	Prefork         bool          `mapstructure:"prefork"`
}

type PostgresDatabase struct {
	Host              string        `mapstructure:"host"                validate:"required,hostname|ip"`
	Port              int           `mapstructure:"port"                validate:"required,min=1,max=65535"`
	Name              string        `mapstructure:"name"                validate:"required,printascii"`
	User              string        `mapstructure:"user"                validate:"required,printascii"`
	Password          string        `mapstructure:"password"            validate:"required"`
	MaxConnections    int           `mapstructure:"max_connections"     validate:"required,min=1,max=2000"`
	MinConnections    int           `mapstructure:"min_connections"     validate:"required,min=1,max=2000,ltfield=MaxConnections"`
	ConnectionTimeout time.Duration `mapstructure:"connection_timeout"  validate:"required,gt=0"`
	MaxConnLifetime   time.Duration `mapstructure:"max_conn_lifetime"   validate:"required,gt=0"`
	MaxConnIdleTime   time.Duration `mapstructure:"max_conn_idle_time"  validate:"required,gt=0"`
	HealthCheckPeriod time.Duration `mapstructure:"health_check_period" validate:"required,gt=0"`
}

type Logging struct {
	Level  string `mapstructure:"level"  validate:"required,oneof=debug info warn error"`
	Format string `mapstructure:"format" validate:"required,oneof=text json"`
}

func LoadConfig(environment string) (*Config, error) {
	viper.SetConfigName(environment + "-" + configSubName)
	viper.SetConfigType(configsFileType)
	viper.AddConfigPath(configsFolderPath)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	bindEnv()

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed read config: %w", err)
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed parse config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed validate config: %w", err)
	}

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("server.app_name", "Peerly ")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.concurrency", 128)
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.write_timeout", "15s")
	viper.SetDefault("server.idle_timeout", "15s")
	viper.SetDefault("server.shutdown_timeout", "30s")
	viper.SetDefault("server.prefork", false)

	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.name", "hub_dev")
	viper.SetDefault("postgres.user", "dev_user")
	viper.SetDefault("postgres.password", "dev_password")
	viper.SetDefault("postgres.max_connections", 20)
	viper.SetDefault("postgres.min_connections", 5)
	viper.SetDefault("postgres.connection_timeout", "10s")
	viper.SetDefault("postgres.max_conn_lifetime", "30m")
	viper.SetDefault("postgres.max_conn_idle_time", "5m")
	viper.SetDefault("postgres.health_check_period", "30s")

	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
}

func bindEnv() {
	pairs := map[string]string{
		"server.app_name":         "SERVER_APP_NAME",
		"server.host":             "SERVER_HOST",
		"server.port":             "SERVER_PORT",
		"server.concurrency":      "SERVER_CONCURRENCY",
		"server.read_timeout":     "SERVER_READ_TIMEOUT",
		"server.write_timeout":    "SERVER_WRITE_TIMEOUT",
		"server.shutdown_timeout": "SERVER_SHUTDOWN_TIMEOUT",
		"server.idle_timeout":     "SERVER_IDLE_TIMEOUT",
		"server.prefork":          "SERVER_PREFORK",

		"postgres.host":                "POSTGRES_HOST",
		"postgres.port":                "POSTGRES_PORT",
		"postgres.name":                "POSTGRES_DB_NAME",
		"postgres.user":                "POSTGRES_USER",
		"postgres.password":            "POSTGRES_PASSWORD",
		"postgres.max_connections":     "POSTGRES_MAX_CONNECTIONS",
		"postgres.min_connections":     "POSTGRES_MIN_CONNECTIONS",
		"postgres.connection_timeout":  "POSTGRES_CONNECTION_TIMEOUT",
		"postgres.max_conn_lifetime":   "POSTGRES_MAX_CONN_LIFETIME",
		"postgres.max_conn_idle_time":  "POSTGRES_MAX_CONN_IDLE_TIME",
		"postgres.health_check_period": "POSTGRES_HEALTH_CHECK_PERIOD",

		"logging.level":  "LOGGING_LEVEL",
		"logging.format": "LOGGING_FORMAT",
	}

	for k, env := range pairs {
		err := viper.BindEnv(k, env)
		if err != nil {
			panic(err)
		}
	}
}

func (c *Config) Validate() error {
	return validators.Validator.Struct(c)
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"Server Config: %+v, Postgres Config: %+v, Logging Config: %+v",
		c.ServerCfg,
		c.PostgresCfg,
		c.LoggingCfg,
	)
}
