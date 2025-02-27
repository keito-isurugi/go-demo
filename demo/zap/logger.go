package main

import (
	"time"

	"go.uber.org/zap"
)

func loggerDemo() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	url := "http://example.com"
	logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
		zap.String("msg", "loggerを使用"),
	)
}
