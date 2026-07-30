package sentry

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func InitSentry(dsn, environment string, tracesSampleRate float64) {
	if dsn == "" {
		log.Println("sentry: DSN not set, skipping initialization")
		return
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      environment,
		TracesSampleRate: tracesSampleRate,
		ProfilesSampleRate: 0.05,
		EnableTracing:    true,
		Release:          "go-road-backend@1.0.0",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %v", err)
	}
}

func GinMiddleware() gin.HandlerFunc {
	return sentrygin.New(sentrygin.Options{
		Repanic: true,
	})
}

func CaptureException(err error) {
	sentry.CaptureException(err)
}

func Flush() {
	sentry.Flush(2 * time.Second)
}
