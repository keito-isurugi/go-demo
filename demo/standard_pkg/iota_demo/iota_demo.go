package iotademo

const (
	a = iota // iotaは0から始まる連続した整数を生成する
	b = iota // 1
	c = iota // 2
)

// 状態や種別の識別子
type UserStatus int
// ユーザーの状態を定数で管理
const (
    StatusActive UserStatus = iota
    StatusSuspended
    StatusDeactivated
)

// フラグのビット値管理（ビットマスク）
// ビット演算で状態のON/OFFを管理
const (
    FlagRead = 1 << iota // 1 << 0 → 0001
    FlagWrite            // 1 << 1 → 0010
    FlagDelete           // 1 << 2 → 0100
)

func hasPermission(flags, permission int) bool {
    return flags&permission != 0
}

// エラーコード番号の整理
type ErrorCode int

const (
    ErrNotFound ErrorCode = iota + 1000 // 1000
    ErrInvalidRequest                   // 1001
    ErrPermissionDenied                 // 1002
)