package main

import (
	"encoding/json"

	"go.uber.org/zap"
)

func basicConfiguration() {
	// For some users, the presets offered by the NewProduction, NewDevelopment,
	// and NewExample constructors won't be appropriate. For most of those
	// users, the bundled Config struct offers the right balance of flexibility
	// and convenience. (For more complex needs, see the AdvancedConfiguration
	// example.)
	//
	// See the documentation for Config and zapcore.EncoderConfig for all the
	// available options.
	rawJSON := []byte(`{
		  "level": "debug",
		  "encoding": "json",
		  "outputPaths": ["stdout", "/tmp/logs"],
		  "errorOutputPaths": ["stderr"],
		  "initialFields": {"foo": "bar"},
		  "encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		  }
		}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger := zap.Must(cfg.Build())
	defer logger.Sync()

	logger.Info("logger construction succeeded")
	logger.Info("zap BasicConfiguration")
}
