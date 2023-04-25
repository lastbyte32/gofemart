package config

import (
	"flag"

	"github.com/caarlos0/env/v7"
)

type config struct {
	AccrualHost string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	ApiHost     string `env:"RUN_ADDRESS"`
	DSN         string `env:"DATABASE_URI"`
	SigningKey  string `env:"SIGNING_KEY"`
}

func (c *config) GetSigningKey() string {
	return c.SigningKey
}

func (c *config) GetAccrual() string {
	return c.AccrualHost
}

func (c *config) GetApiHost() string {
	return c.ApiHost
}

func (c *config) GetDSN() string {
	return c.DSN
}

func (c *config) env() error {
	if err := env.Parse(c); err != nil {
		return err
	}
	return nil
}

func (c *config) flags() {

	flag.StringVar(&c.AccrualHost, "r", defaultAccrualHost, "accrual system")
	flag.StringVar(&c.ApiHost, "a", defaultApiHost, "api host")
	flag.StringVar(&c.SigningKey, "k", defaultSigningKey, "key sign")
	flag.StringVar(&c.DSN, "d", defaultDatabaseDSN, "dsn")
	flag.Parse()
}

func New() (*config, error) {
	c := &config{}
	c.flags()
	if err := c.env(); err != nil {
		return nil, err
	}
	return c, nil
}
