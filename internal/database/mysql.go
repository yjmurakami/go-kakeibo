package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxLifetime  int    `yaml:"maxLifetime"`
}

func ValidateMySQLConfig(cnf MySQLConfig) error {
	if cnf.Host == "" {
		return fmt.Errorf("mysql.host is not specified")
	}

	if cnf.Port == "" {
		return fmt.Errorf("mysql.port is not specified")
	}

	if cnf.Username == "" {
		return fmt.Errorf("mysql.username is not specified")
	}

	if cnf.Password == "" {
		return fmt.Errorf("mysql.password is not specified")
	}

	if cnf.MaxOpenConns == 0 {
		return fmt.Errorf("mysql.maxOpenConns is not specified")
	}

	if cnf.MaxIdleConns == 0 {
		return fmt.Errorf("mysql.maxIdleConns is not specified")
	}

	if cnf.MaxLifetime == 0 {
		return fmt.Errorf("mysql.maxLifetime is not specified")
	}

	return nil
}

func OpenMySQL(conf MySQLConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true", conf.Username, conf.Password, conf.Host, conf.Port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(conf.MaxLifetime) * time.Minute)

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}