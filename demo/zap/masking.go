package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func execNoMasking() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	password := "secret123"

	logger.Info("user login",
		zap.String("username", "alice"),
		zap.String("password", password), // これだと機密情報が丸見えに
	)
}

func maskSecret(_ string) string {
	return "****"
}

func execMasking() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	username := "alice"
	password := "secret123"

	logger.Info("user login",
		zap.String("username", username),
		zap.String("password", maskSecret(password)), // ここでマスク
	)
}

// マスキングしたい項目を保持する
var sensitiveFields = []string{"password", "token", "apikey"}

// マスキング対象かどうか
func isSensitive(key string) bool {
	for _, k := range sensitiveFields {
		if k == key {
			return true
		}
	}
	return false
}

// カスタムCore: 指定されたキーならマスクする
type maskingCore struct {
	zapcore.Core
}

func (c maskingCore) With(fields []zapcore.Field) zapcore.Core {
	newFields := make([]zapcore.Field, 0, len(fields))
	for _, f := range fields {
		if isSensitive(f.Key) && f.Type == zapcore.StringType {
			newFields = append(newFields, zap.String(f.Key, maskSecret(f.String)))
		} else {
			newFields = append(newFields, f)
		}
	}
	return maskingCore{c.Core.With(newFields)}
}

func (c maskingCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return c.Core.Check(ent, ce)
}

func (c maskingCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	newFields := make([]zapcore.Field, 0, len(fields))
	for _, f := range fields {
		if isSensitive(f.Key) && f.Type == zapcore.StringType {
			newFields = append(newFields, zap.String(f.Key, maskSecret(f.String)))
		} else {
			newFields = append(newFields, f)
		}
	}
	return c.Core.Write(ent, newFields)
}

func newLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)
	maskedCore := maskingCore{core}
	return zap.New(maskedCore)
}

func execMasking2() {
	logger := newLogger()
	defer logger.Sync()

	// ログ出力テスト
	logger.Info("user login",
		zap.String("username", "taro"),
		zap.String("password", "mySecret123"),
		zap.String("token", "abc123token"),
		zap.String("message", "ログイン成功"),
	)
}