package main

import (
	"time"

	"go.uber.org/zap"
)

func main() {
	loggerSugar, _ := zap.NewProduction()
	defer loggerSugar.Sync() // flushes buffer, if any

	url := "http://example.com"
	sugar := loggerSugar.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
	sugar.Infof("sugarを使用")


	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
		zap.String("msg", "loggerを使用"),
)

}
