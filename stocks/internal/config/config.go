package config

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Ensure StockServiceConfig implements Config interface.
var _ Config = (*StockServiceConfig)(nil)

// Config interface provides methods to a retrieve a configuration values for stocks service.
type Config interface {
	Address() string
	GRPCAddress() string
	MetricsAddress() string
	SrvConfig() ServerConfig
	DbConfig() PostgresConfig
	GetKafkaBrokers() string
}

type StockServiceConfig struct {
	Server           ServerConfig
	Postgres         PostgresConfig
	ExternalServices ExternalServicesConfig
	Kafka            KafkaServiceConfig
}

type (
	// ServerConfig holds server configurations for stock service.
	ServerConfig struct {
		ServiceName  string        `env:"SERVICE_NAME,required"`
		HTTPPort     string        `env:"HTTP_PORT,required"`
		GRPCPORT     string        `env:"GRPC_PORT,required"`
		MetricsPort  string        `env:"METRICS_PORT,required"`
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
		ObservalityConfig struct {
			LogStashHost string `env:"LOGSTASH_HOST,required"`
		}
	}
	// KafkaServiceConfig holds needed configurations for stock service event producer.
	KafkaServiceConfig struct {
		Brokers string `env:"KAFKA_BROKERS,required"`
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
func NewStockServiceConfig() (*StockServiceConfig, error) {
	stockServiceConfig := &StockServiceConfig{}
	if err := env.Parse(stockServiceConfig); err != nil {
		return nil, fmt.Errorf("stockServiceConfig.Parse: %w", err)
	}

	return stockServiceConfig, nil
}

// Address returns the server address in format host:port.
func (c *StockServiceConfig) Address() string {
	return net.JoinHostPort("", c.Server.HTTPPort)
}

func (c *StockServiceConfig) GRPCAddress() string {
	return net.JoinHostPort("", c.Server.GRPCPORT)
}

func (c *StockServiceConfig) MetricsAddress() string {
	return net.JoinHostPort("", c.Server.MetricsPort)
}

func (c *StockServiceConfig) SrvConfig() ServerConfig {
	return c.Server
}

// DbConfig returns the postgresql configuration.
func (c *StockServiceConfig) DbConfig() PostgresConfig {
	return c.Postgres
}

func (c *StockServiceConfig) GetKafkaBrokers() string {
	if c.Kafka.Brokers == "" {
		return "kafka1:29091,kafka2:29092"
	}

	return c.Kafka.Brokers
}

// GenerateDSN returns a psql url.
func (p *PostgresConfig) GenerateDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, p.DBName,
	)
}
