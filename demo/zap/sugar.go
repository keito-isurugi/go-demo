package main

import (
	"time"

	"go.uber.org/zap"
)

func suberDemo() {
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
}
