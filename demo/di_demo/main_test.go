package main

import (
	"testing"
	"time"
)

// 現在時刻を提供するインターフェースのモック
type MockTimeProvider struct {
	FixedTime time.Time
}

func (m *MockTimeProvider) Now() time.Time {
	return m.FixedTime
}

func TestIsExpired(t *testing.T) {
	deadline := time.Date(2025, 5, 9, 12, 0, 0, 0, time.Local)
	tests := []struct {
		name     string
		now      time.Time
		expected bool
	}{
		{"期限前", time.Date(2025, 5, 9, 11, 0, 0, 0, time.Local), false},
		{"期限後", time.Date(2025, 5, 9, 13, 0, 0, 0, time.Local), true},
	}

	for _, tt := range tests {
		// モックを使って現在時刻を制御
		// これにより、テストの実行時刻に依存しない
		mock := &MockTimeProvider{FixedTime: tt.now}
		checker := NewExpiryChecker(mock)

		got := checker.IsExpired(deadline)
		if got != tt.expected {
			t.Errorf("[%s] expected %v, got %v", tt.name, tt.expected, got)
		}
	}
}
