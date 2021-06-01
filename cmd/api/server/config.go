package server

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/yjmurakami/go-kakeibo/internal/database"
	"gopkg.in/yaml.v3"
)

var ErrInvalidConfig = errors.New("config is invalid")

type config struct {
	Api       apiConfig            `yaml:"api"`
	MySQL     database.MySQLConfig `yaml:"mysql"`
	MySQLTest database.MySQLConfig `yaml:"mysqlTest"`
}

type apiConfig struct {
	Port          string `yaml:"port"`
	IdleTimeout   int    `yaml:"idleTimeout"`
	ReadTimeout   int    `yaml:"readTimeout"`
	WriteTimeout  int    `yaml:"writeTimeout"`
	JwtKey        string `yaml:"jwtKey"`
	JwtExpiration int    `yaml:"jwtExpiration"`
}

func readConfig(filepath string) (*config, error) {
	buf, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	cnf := &config{}
	err = yaml.Unmarshal(buf, cnf)
	if err != nil {
		return nil, err
	}

	if cnf.Api.Port == "" {
		return nil, fmt.Errorf("%w : api.port is not specified", ErrInvalidConfig)
	}

	if cnf.Api.IdleTimeout == 0 {
		return nil, fmt.Errorf("%w : api.idleTimeout is not specified", ErrInvalidConfig)
	}

	if cnf.Api.ReadTimeout == 0 {
		return nil, fmt.Errorf("%w : api.readTimeout is not specified", ErrInvalidConfig)
	}

	if cnf.Api.WriteTimeout == 0 {
		return nil, fmt.Errorf("%w : api.writeTimeout is not specified", ErrInvalidConfig)
	}

	if cnf.Api.JwtKey == "" {
		return nil, fmt.Errorf("%w : api.jwtKey is not specified", ErrInvalidConfig)
	}
	if cnf.Api.JwtExpiration == 0 {
		return nil, fmt.Errorf("%w : api.jwtExpiration is not specified", ErrInvalidConfig)
	}

	err = database.ValidateMySQLConfig(cnf.MySQL)
	if err != nil {
		return nil, fmt.Errorf("%w : %s", ErrInvalidConfig, err.Error())
	}

	return cnf, nil
}

func (c *config) String() string {
	b := strings.Builder{}

	// API
	fmt.Fprintf(&b, "API\n")
	fmt.Fprintf(&b, "\tPort\t: %s\n", c.Api.Port)
	fmt.Fprintf(&b, "\tIdleTimeout\t: %d\n", c.Api.IdleTimeout)
	fmt.Fprintf(&b, "\tReadTimeout\t: %d\n", c.Api.ReadTimeout)
	fmt.Fprintf(&b, "\tWriteTimeout\t: %d\n", c.Api.WriteTimeout)
	fmt.Fprintf(&b, "\tJwtKey\t: %s\n", c.Api.JwtKey)
	fmt.Fprintf(&b, "\tJwtExpiration\t: %d\n", c.Api.JwtExpiration)

	// MySQL
	fmt.Fprintf(&b, "MySQL\n")
	fmt.Fprintf(&b, "\tHost\t: %s\n", c.MySQL.Host)
	fmt.Fprintf(&b, "\tPort\t: %s\n", c.MySQL.Port)
	fmt.Fprintf(&b, "\tUsername\t: %s\n", c.MySQL.Username)
	fmt.Fprintf(&b, "\tPassword\t: %s\n", c.MySQL.Password)
	fmt.Fprintf(&b, "\tMaxOpenConns\t: %d\n", c.MySQL.MaxOpenConns)
	fmt.Fprintf(&b, "\tMaxIdleConns\t: %d\n", c.MySQL.MaxIdleConns)
	fmt.Fprintf(&b, "\tMaxLifetime \t: %d\n", c.MySQL.MaxLifetime)

	return b.String()
}
