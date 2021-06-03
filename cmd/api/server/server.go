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

	timer := clock.SystemClock{}
	jwt := handler.NewJWT(timer, cnf.Api.JwtKey, cnf.Api.JwtExpiration)
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	hc := handlerConfig{
		logger:    logger,
		clock:     timer,
		db:        db,
		jwt:       jwt,
		config:    cnf,
		container: newContainer(),
	}

	mdl := handler.NewMiddlewareHandler(
		service.NewMiddlewareService(
			hc.db,
			hc.container.userRepository,
		),
		hc.jwt,
		hc.clock,
	)
	mux := newRouter(mdl, hc)
	mux = mdl.RecoverPanic(mux)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cnf.Api.Port),
		ErrorLog:     log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile), // TODO
		Handler:      mux,
		IdleTimeout:  time.Duration(cnf.Api.IdleTimeout) * time.Second,
		ReadTimeout:  time.Duration(cnf.Api.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cnf.Api.WriteTimeout) * time.Second,
	}

	log.Printf("server is running\n%s", cnf)
	return srv.ListenAndServe()
}
