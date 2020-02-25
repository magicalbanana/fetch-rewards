package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/magicalbanana/fetch-rewards/handler"
	"go.uber.org/zap"
)

func main() {
	lgr, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err.Error())
	}

	rts, err := loadRoutes(lgr)
	if err != nil {
		lgr.Fatal("failed to load routes", zap.String("error", err.Error()))
	}

	appAddress := os.Getenv("APP_ADDRESS")
	if appAddress == "" {
		appAddress = "0.0.0.0:5656"
	}

	lgr.Info("starting web server", zap.String("address", appAddress))
	srv := &http.Server{Addr: appAddress, Handler: rts}
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// sigterm signal sent from kubernetes
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			lgr.Error("HTTP Server Fatal Error, shutting down server", zap.String("error", err.Error()))
		}
		close(idleConnsClosed)
	}()

	lgr.Info("starting web server", zap.String("address", appAddress))
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// close DB connections
		lgr.Error("error while serving http server", zap.String("error", err.Error()))
	}

	<-idleConnsClosed
}

func loadRoutes(lgr *zap.Logger) (http.Handler, error) {
	router := mux.NewRouter()
	router.Path("/compare-versions").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := handler.VersionCompare{}
		base := handler.Base{}
		base.H = h.Compare
		base.Logger = lgr
		base.ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	return router, nil
}
