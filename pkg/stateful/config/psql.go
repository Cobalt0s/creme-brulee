package config

import (
	"fmt"
	"net/url"
)

type PsqlConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
	SslMode      string
	SslRootCert  string
}

func NewPsqlConfig() (*PsqlConfig, error) {

	envPostgresHost, err := GetEnv("POSTGRES_HOST")
	if err != nil {
		return nil, err
	}
	envPostgresPort, err := GetEnv("POSTGRES_PORT")
	if err != nil {
		return nil, err
	}
	envPostgresUser, err := GetEnv("POSTGRES_USER")
	if err != nil {
		return nil, err
	}
	envPostgresPassword, err := GetEnv("POSTGRES_PASSWORD")
	if err != nil {
		return nil, err
	}
	envPostgresDatabaseName, err := GetEnv("POSTGRES_DATABASE_NAME")
	if err != nil {
		return nil, err
	}
	envPostgresSslMode, err := GetEnv("POSTGRES_SSL_MODE")
	if err != nil {
		return nil, err
	}
	var envPostgresSslRootCertificate string
	if envPostgresSslMode != "disable" {
		envPostgresSslRootCertificate, err = GetEnv("POSTGRES_SSL_ROOT_CERT")
		if err != nil {
			return nil, err
		}
	}

	return &PsqlConfig{
		Host:         envPostgresHost,
		Port:         envPostgresPort,
		User:         envPostgresUser,
		Password:     envPostgresPassword,
		DatabaseName: envPostgresDatabaseName,
		SslMode:      envPostgresSslMode,
		SslRootCert:  envPostgresSslRootCertificate,
	}, nil
}

func (c *PsqlConfig) GetDataSourcePSQL() *url.URL {
	postgresURL := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     fmt.Sprintf("%v:%v", c.Host, c.Port),
		Path:     fmt.Sprintf("/%v", c.DatabaseName),
		RawQuery: "sslmode=disable",
	}

	if c.SslMode != "disable" {
		postgresURL.RawQuery = fmt.Sprintf("sslmode=%v&sslrootcert=%v", c.SslMode, c.SslRootCert)
	}

	return postgresURL
}
