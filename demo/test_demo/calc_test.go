package main

import "testing"

func TestAdd(t *testing.T) {
	got := Add(2, 3)
	want := 5

	if got != want {
		t.Errorf("Add(2, 3) = %d; want %d", got, want)
	}
}

func TestAddTableDriven(t *testing.T) {
	// テストケースの一覧（テーブル）を定義
	tests := []struct {
		a, b int // 入力値
		want int // 期待する戻り値
	}{
		{1, 2, 3},         // 1 + 2 = 3
		{0, 0, 0},         // 0 + 0 = 0
		{-1, 1, 0},        // -1 + 1 = 0
		{100, 200, 300},   // 100 + 200 = 300
		{10000, 1, 0},     // 10000 + 1 = 0 (テストで通らないif文)
	}

	// テストケースごとに繰り返し実行
	for _, tt := range tests {
		// Add関数を実行し、結果を取得
		got := Add(tt.a, tt.b)

		// 結果と期待値が一致しない場合はエラーを出力
		if got != tt.want {
			t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
		}
	}
}