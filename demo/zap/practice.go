package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// コンテキストキー
type ctxKey string

const requestIDKey ctxKey = "request_id"

// マスク対象のキー
var sensitiveKeys = map[string]bool{
	"password":    true,
	"api_key":     true,
	"credit_card": true,
}

// センシティブ情報をマスクする関数
func maskSensitiveData(key, value string) string {
	if _, exists := sensitiveKeys[key]; exists {
		if len(value) > 6 {
			return value[:3] + "****" + value[len(value)-3:] // 例: abc****xyz
		}
		return "****" // 短すぎる場合は完全マスク
	}
	return value
}

// カスタムエンコーダ
func customEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.LevelKey = "level"
	encoderConfig.MessageKey = "message"
	encoderConfig.CallerKey = "caller"
	encoderConfig.StacktraceKey = "stacktrace"
	return zapcore.NewJSONEncoder(encoderConfig)
}

// カスタムロガー
func initLogger() *zap.Logger {
	core := zapcore.NewCore(
		zapcore.RegisterHooks(
			customEncoder(),
			func(entry zapcore.Entry) error {
				// ここでログのフック処理を追加可能（例: 特定のレベル以上のログを別の出力に送る）
				return nil
			},
		),
		zapcore.AddSync(zapcore.Lock(zapcore.AddSync(zap.NewExample().Core()))),
		zap.InfoLevel,
	)

	return zap.New(core)
}

// リクエストIDを取得
func getRequestID(ctx context.Context) string {
	if val, ok := ctx.Value(requestIDKey).(string); ok {
		return val
	}
	return "unknown"
}

// 構造化ログを出力する関数
func logRequest(ctx context.Context, logger *zap.Logger, data map[string]string) {
	reqID := getRequestID(ctx)

	fields := []zap.Field{
		zap.String("request_id", reqID),
		zap.String("endpoint", "/example"),
		zap.String("method", "POST"),
		zap.Time("timestamp", time.Now()),
	}

	// データのマスキング
	for key, value := range data {
		maskedValue := maskSensitiveData(key, value)
		fields = append(fields, zap.String(key, maskedValue))
	}

	logger.Info("Request received", fields...)
}

// エラーログを出力する関数
func logError(ctx context.Context, logger *zap.Logger, err error) {
	reqID := getRequestID(ctx)

	logger.Error("Error occurred",
		zap.String("request_id", reqID),
		zap.String("error", err.Error()),
		zap.Stack("stacktrace"),
	)
}

func practice() {
	logger := initLogger()
	defer logger.Sync() // バッファをフラッシュ

	// コンテキストにリクエストIDを設定
	ctx := context.WithValue(context.Background(), requestIDKey, "123456")

	// リクエストデータ（センシティブ情報を含む）
	requestData := map[string]string{
		"username":    "john_doe",
		"password":    "supersecretpassword",
		"api_key":     "1234567890abcdef",
		"credit_card": "4111111111111111",
		"email":       "john.doe@example.com",
	}

	// ログ出力
	logRequest(ctx, logger, requestData)

	// エラーログ出力
	err := fmt.Errorf("database connection failed")
	logError(ctx, logger, err)
}
