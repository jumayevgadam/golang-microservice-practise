package config

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Ensure StockServiceConfig implements Config interface.
var _ Config = (*CartServiceConfig)(nil)

// Config interface provides methods to a retrieve a configuration values for stocks service.
type Config interface {
	Address() string
	SrvConfig() ServerConfig
	DbConfig() PostgresConfig
	StockServiceURL() string
}

type CartServiceConfig struct {
	Server           ServerConfig
	Postgres         PostgresConfig
	ExternalServices ExternalServicesConfig
}

type (
	// ServerConfig holds server configurations for stock service.
	ServerConfig struct {
		HTTPPort     string        `env:"HTTP_PORT,required"`
		ReadTimeOut  time.Duration `env:"READ_TIMEOUT,required"`
		WriteTimeOut time.Duration `env:"WRITE_TIMEOUT,required"`
	}
	// PostgresConfig holds postgresql configurations for stock service.
	PostgresConfig struct {
		Host     string `env:"DB_HOST,required"`
		Port     string `env:"DB_PORT,required"`
		User     string `env:"DB_USER,required"`
		Password string `env:"DB_PASSWORD,required"`
		DBName   string `env:"DB_NAME,required"`
	}
	// ExternalServicesConfig holds ExternalServices configurations which need in stock service.
	ExternalServicesConfig struct {
		StockServiceURL string `env:"STOCK_SERVICE_URL,required"`
	}
)

// LoadEnv load environment variables.
func LoadEnv(path string) error {
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}

	return nil
}

// NewStockServiceConfig returns a new StockServiceConfig.
func NewCartServiceConfig() (*CartServiceConfig, error) {
	cartServiceConfig := &CartServiceConfig{}
	if err := env.Parse(cartServiceConfig); err != nil {
		return nil, fmt.Errorf("stockServiceConfig.Parse: %w", err)
	}

	return cartServiceConfig, nil
}

// Address returns the server address in format host:port.
func (c *CartServiceConfig) Address() string {
	return net.JoinHostPort("", c.Server.HTTPPort)
}

func (c *CartServiceConfig) SrvConfig() ServerConfig {
	return c.Server
}

// DbConfig returns the postgresql configuration.
func (c *CartServiceConfig) DbConfig() PostgresConfig {
	return c.Postgres
}

func (c *CartServiceConfig) StockServiceURL() string {
	return c.ExternalServices.StockServiceURL
}

// GenerateDSN returns a psql url.
func (p *PostgresConfig) GenerateDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, p.DBName,
	)
}
