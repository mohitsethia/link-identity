package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/link-identity/app/application"
	appconfig "github.com/link-identity/app/config"
	httpHandler "github.com/link-identity/app/http"
	"github.com/link-identity/app/infrastructure"
	"github.com/link-identity/app/infrastructure/repository"
	"github.com/link-identity/app/infrastructure/sql"
	"github.com/link-identity/app/utils"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

var (
	// zap logger instance
	logEntry *logrus.Entry
)

func init() {
	logger := infrastructure.NewLogger(os.Stdout, "info", "test")
	hostname, _ := os.Hostname()
	logEntry = logger.WithFields(logrus.Fields{
		"env":          "test",
		"pod_id":       hostname,
		"program":      "test-app",
		"channel":      "http",
		"request_path": "",
		"remote_addr":  "",
		"status_code":  "",
	})
}

func main() {
	// setup database connection
	db := sql.NewDBConnection()

	repo := repository.NewContactRepository(db)

	service := application.NewService(repo)
	handler := httpHandler.NewLinkIdentityHandler(service)

	// setup the http server
	router := SetupRouters(handler)

	// service address will be changed as port in next PR.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", appconfig.Values.Server.Port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      router}
	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, Cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer Cancel()
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()
		// Trigger graceful shutdown
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()
	// Run the server
	logEntry.Info("Starting application at port: " + appconfig.Values.Server.Port)
	errServer := srv.ListenAndServe()
	if errServer != nil && errServer != http.ErrServerClosed {
		log.Fatal(errServer)
	}
	// Wait for server context to be stopped
	<-serverCtx.Done()
	logEntry.Info("Application stopped gracefully!")
}

// SetupRouters ...
func SetupRouters(handler *httpHandler.LinkIdentityHandler) *chi.Mux {
	// Base route initialize.
	router := chi.NewRouter()
	router.Use(infrastructure.NewLoggerMiddleware(logEntry).Wrap)

	//Health check registration
	router.Get("/health/check", GetHealthCheck)
	router.Get("/", GetHealthCheck)

	// Register Contact get handler
	{
		router.Post("/identify", handler.Identify)
	}
	return router
}

// GetHealthCheck ...
func GetHealthCheck(w http.ResponseWriter, _ *http.Request) {
	res := utils.ResponseDTO{
		StatusCode: http.StatusOK,
		Data:       "success",
	}
	utils.ResponseJSON(w, http.StatusOK, res)
}
