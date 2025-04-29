package main

import (
	"fmt"
	"time"
)

// 現在時刻を提供するインターフェース
type TimeProvider interface {
	Now() time.Time
}

type RealTimeProvider struct{}

func (tp *RealTimeProvider) Now() time.Time {
	return time.Now()
}

// ExpiryCheckerは期限切れをチェックする構造体
type ExpiryChecker struct {
	TimeProvider TimeProvider
}

func NewExpiryChecker(tp TimeProvider) *ExpiryChecker {
	return &ExpiryChecker{TimeProvider: tp}
}

func (ec *ExpiryChecker) IsExpired(deadline time.Time) bool {
	// time.Now()を直接呼び出すとテストの実行時刻に依存してしまう
	// そのため、TimeProviderインターフェースを通じて現在時刻を取得する
	// これにより、テスト時にモックを使って現在時刻を制御できる
	return ec.TimeProvider.Now().After(deadline)
}

func main() {
	checker := NewExpiryChecker(&RealTimeProvider{})

	// DBなどに保存されてる「期限」など
	deadline := time.Date(2025, 5, 9, 12, 0, 0, 0, time.Local)

	if checker.IsExpired(deadline) {
		fmt.Println("期限切れです")
	} else {
		fmt.Println("まだ有効です")
	}
}
