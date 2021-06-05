package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yjmurakami/go-kakeibo/cmd/api/handler"
	"github.com/yjmurakami/go-kakeibo/cmd/api/service"
	"github.com/yjmurakami/go-kakeibo/internal/clock"
	"github.com/yjmurakami/go-kakeibo/internal/database"
	"github.com/yjmurakami/go-kakeibo/internal/repository"
)

func Start() error {
	cnf, err := readConfig("./config.yml")
	if err != nil {
		return err
	}

	db, err := database.OpenMySQL(cnf.MySQL)
	if err != nil {
		return err
	}
	defer db.Close()

	jwtExpiration, err := time.ParseDuration(cnf.Api.JwtExpiration)
	if err != nil {
		return err
	}

	timer := clock.SystemClock{}
	jwt := handler.NewJWT(timer, cnf.Api.JwtKey, jwtExpiration)
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	hc := handlerConfig{
		logger: logger,
		clock:  timer,
		db:     db,
		jwt:    jwt,
		config: cnf,
		repos:  repository.NewRepositories(),
	}

	mdl := handler.NewMiddlewareHandler(
		service.NewMiddlewareService(
			hc.db,
			hc.repos,
		),
		hc.jwt,
		hc.clock,
	)
	mux := newRouter(nil, hc) // TODO middleware
	mux = mdl.RecoverPanic(mux)

	idleTimeout, err := time.ParseDuration(cnf.Api.IdleTimeout)
	if err != nil {
		return err
	}

	readTimeout, err := time.ParseDuration(cnf.Api.ReadTimeout)
	if err != nil {
		return err
	}

	writeTimeout, err := time.ParseDuration(cnf.Api.WriteTimeout)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cnf.Api.Port),
		ErrorLog:     log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile), // TODO
		Handler:      mux,
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	log.Printf("server is running\n%s", cnf)
	return srv.ListenAndServe()
}
