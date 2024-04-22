package main

import (
	"context"
	"fmt"
	"github.com/link-identity/app/infrastructure/mysql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appconfig "github.com/link-identity/app/config"
	domain "github.com/link-identity/app/domain"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// commitHash should be set at compile time with current git hash
	commitHash string
	// tag should be set at compile time with current branch or tag
	tag string
	// zap logger instance
	logEntryZap *zap.Logger
)

func init() {
	hostname, _ := os.Hostname()

	level, err := zap.ParseAtomicLevel("debug")
	if err != nil {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	config := zap.NewProductionConfig()
	config.Level = level
	config.EncoderConfig = getEncoderConfig()
	logZap := zap.Must(config.Build())

	logEntryZap = logZap.With(
		zap.String("env", "dev"),
		zap.String("pod_id", hostname),
		zap.String("program", "link-identity"),
		zap.String("channel", "http"),
		zap.String("request_path", ""),
		zap.String("remote_addr", ""),
		zap.String("status_code", ""),
	)
}

func main() {
	// logger defers
	defer logEntryZap.Sync()

	// setup database connection
	db := mysql.NewDBConnection()

	// setup the http server
	router := SetupRouters(db)

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
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)
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
	errServer := srv.ListenAndServe()
	if errServer != nil && errServer != http.ErrServerClosed {
		log.Fatal(errServer)
	}
	// Wait for server context to be stopped
	<-serverCtx.Done()
	logEntryZap.Info("Application stopped gracefully!")
}

// SetupRouters ...
func SetupRouters(db *mysql.DbConn) *gin.Engine {
	// Base route initialize.
	router := gin.Default()
	//router := chi.NewRouter()

	//Health check registration
	router.GET("/health/check", GetHealthCheck())

	// Register customer get handler
	{
		//router.Group(func(r chi.Router) {
		//})
	}
	return router
}

func GetHealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := map[string]interface{}{
			"success": "true",
			"message": "connected",
		}
		domain.HttpResponse(http.StatusOK).Data(res).Send(c)
	}
}

// getEncoderConfig
func getEncoderConfig() zapcore.EncoderConfig {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	return encoderCfg
}
