package server

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

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
	IdleTimeout   string `yaml:"idleTimeout"`
	ReadTimeout   string `yaml:"readTimeout"`
	WriteTimeout  string `yaml:"writeTimeout"`
	JwtKey        string `yaml:"jwtKey"`
	JwtExpiration string `yaml:"jwtExpiration"`
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
		return nil, fmt.Errorf("%w : api.port is invalid", ErrInvalidConfig)
	}

	_, err = time.ParseDuration(cnf.Api.IdleTimeout)
	if err != nil {
		return nil, fmt.Errorf("%w : api.idleTimeout is invalid", ErrInvalidConfig)
	}

	_, err = time.ParseDuration(cnf.Api.ReadTimeout)
	if err != nil {
		return nil, fmt.Errorf("%w : api.readTimeout is invalid", ErrInvalidConfig)
	}

	_, err = time.ParseDuration(cnf.Api.WriteTimeout)
	if err != nil {
		return nil, fmt.Errorf("%w : api.writeTimeout is invalid", ErrInvalidConfig)
	}

	if cnf.Api.JwtKey == "" {
		return nil, fmt.Errorf("%w : api.jwtKey is invalid", ErrInvalidConfig)
	}

	_, err = time.ParseDuration(cnf.Api.JwtExpiration)
	if err != nil {
		return nil, fmt.Errorf("%w : api.jwtExpiration is invalid", ErrInvalidConfig)
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
	fmt.Fprintf(&b, "\tIdleTimeout\t: %s\n", c.Api.IdleTimeout)
	fmt.Fprintf(&b, "\tReadTimeout\t: %s\n", c.Api.ReadTimeout)
	fmt.Fprintf(&b, "\tWriteTimeout\t: %s\n", c.Api.WriteTimeout)
	fmt.Fprintf(&b, "\tJwtKey\t: %s\n", c.Api.JwtKey)
	fmt.Fprintf(&b, "\tJwtExpiration\t: %s\n", c.Api.JwtExpiration)

	// MySQL
	fmt.Fprintf(&b, "MySQL\n")
	fmt.Fprintf(&b, "\tHost\t: %s\n", c.MySQL.Host)
	fmt.Fprintf(&b, "\tPort\t: %s\n", c.MySQL.Port)
	fmt.Fprintf(&b, "\tUsername\t: %s\n", c.MySQL.Username)
	fmt.Fprintf(&b, "\tPassword\t: %s\n", c.MySQL.Password)
	fmt.Fprintf(&b, "\tMaxOpenConns\t: %d\n", c.MySQL.MaxOpenConns)
	fmt.Fprintf(&b, "\tMaxIdleConns\t: %d\n", c.MySQL.MaxIdleConns)
	fmt.Fprintf(&b, "\tMaxIdleTime \t: %s\n", c.MySQL.MaxIdleTime)

	return b.String()
}
