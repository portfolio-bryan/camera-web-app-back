package config

import (
	"log"
	"os"
	"path"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort              string
	Environment             string
	POSTGRES_HOST           string
	POSTGRES_PORT           string
	POSTGRES_USERNAME       string
	POSTGRES_PASSWORD       string
	POSTGRES_DATABASE       string
	POSTGRES_SSLMODE        string
	POSTGRES_MAX_IDLE_CONNS int
	POSTGRES_MAX_OPEN_CONNS int
	Otel                    Otel
}

type Otel struct {
	ServiceName          string
	ExporterOtlpEndpoint string
	ExporterOtlpProtocol string
	ExporterOtlpHeaders  string
}

var config *Config

func InitConfig() error {
	environment, ok := os.LookupEnv("ENVIRONMENT")

	if !ok {
		return godotenv.Load()
	}

	if environment == "local" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		p := path.Join(cwd, "../../.env")
		return godotenv.Load(os.ExpandEnv(p))
	}

	if environment == "test" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		p := path.Join(cwd, "../../.test.env")
		if err := godotenv.Load(os.ExpandEnv(p)); err != nil {
			log.Println(".test.env file does not exist")
		}
	}

	return nil

}

func GetConfig() *Config {
	if config != nil {
		return config
	}

	_, err := strconv.Atoi(os.Getenv("SERVER_PORT"))

	if err != nil {
		log.Panic("serverPort was not be processed")
	}

	serverPort := os.Getenv("SERVER_PORT")

	environment := os.Getenv("ENVIRONMENT")
	postgresHost := os.Getenv("POSTGRES_HOST")
	_, err = strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		log.Panic("postgresPort was not be processed")
	}

	postgresPort := os.Getenv("POSTGRES_PORT")

	postgresUsername := os.Getenv("POSTGRES_USERNAME")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDatabase := os.Getenv("POSTGRES_DATABASE")
	postgresSSLMode := os.Getenv("POSTGRES_SSLMODE")
	postgresMaxIdleConns, err := strconv.Atoi(os.Getenv("POSTGRES_MAX_IDLE_CONNS"))

	if err != nil {
		log.Panic("postgresMaxIdleConns was not be processed")
	}
	postgresMaxOpenConns, err := strconv.Atoi(os.Getenv("POSTGRES_MAX_OPEN_CONNS"))

	if err != nil {
		log.Panic("postgresMaxOpenConns was not be processed")
	}

	otelServiceName := os.Getenv("OTEL_SERVICE_NAME")

	otelExporterOtlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")

	otelExporterOtlpProtocol := os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")

	otelExporterOtlpHeaders := os.Getenv("OTEL_EXPORTER_OTLP_HEADERS")

	config = &Config{
		ServerPort:              serverPort,
		Environment:             environment,
		POSTGRES_HOST:           postgresHost,
		POSTGRES_PORT:           postgresPort,
		POSTGRES_USERNAME:       postgresUsername,
		POSTGRES_PASSWORD:       postgresPassword,
		POSTGRES_DATABASE:       postgresDatabase,
		POSTGRES_SSLMODE:        postgresSSLMode,
		POSTGRES_MAX_IDLE_CONNS: postgresMaxIdleConns,
		POSTGRES_MAX_OPEN_CONNS: postgresMaxOpenConns,
		Otel: Otel{
			ServiceName:          otelServiceName,
			ExporterOtlpEndpoint: otelExporterOtlpEndpoint,
			ExporterOtlpProtocol: otelExporterOtlpProtocol,
			ExporterOtlpHeaders:  otelExporterOtlpHeaders,
		},
	}
	return config
}

// Provisional logic, and it is only useful in tests
func (c *Config) ChangePostgresPort(port string) {
	c.POSTGRES_PORT = port
}
