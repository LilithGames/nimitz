package utils

import (
	"github.com/getsentry/sentry-go"
	"log"
	"solarland/infra/annunciation/nimitz/pkg/config"
	"time"
)

func SentrySetup() {
	cfg := config.GetCfg()
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         cfg.GetString("sentry.dsn"),
		Environment: cfg.GetString("sentry.environment"),
		Release:     "",
		Debug:       true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)
}
