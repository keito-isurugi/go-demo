package bank

import "time"

// Account 口座モデル
type Account struct {
	ID        uint      `gorm:"primaryKey"`
	AccountNo string    `gorm:"uniqueIndex;size:20;not null"` // 口座番号
	Balance   int64     `gorm:"not null;default:0"`           // 残高（単位: 円）
	OwnerName string    `gorm:"size:100;not null"`            // 口座名義人
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Transaction 取引履歴モデル
type Transaction struct {
	ID              uint      `gorm:"primaryKey"`
	FromAccountID   uint      `gorm:"index;not null"`                     // 送金元口座ID
	ToAccountID     uint      `gorm:"index;not null"`                     // 送金先口座ID
	Amount          int64     `gorm:"not null"`                           // 送金額
	Status          string    `gorm:"size:20;not null;default:'pending'"` // pending, completed, failed
	TransactionType string    `gorm:"size:20;not null"`                   // transfer, deposit, withdrawal
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

// TransferRequest 振込リクエスト
type TransferRequest struct {
	FromAccountID uint  `json:"from_account_id"`
	ToAccountID   uint  `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferResponse 振込レスポンス
type TransferResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	TransactionID uint   `json:"transaction_id,omitempty"`
	FromBalance   int64  `json:"from_balance,omitempty"`
	ToBalance     int64  `json:"to_balance,omitempty"`
}

// 定数定義
const (
	TransactionStatusPending   = "pending"
	TransactionStatusCompleted = "completed"
	TransactionStatusFailed    = "failed"

	TransactionTypeTransfer   = "transfer"
	TransactionTypeDeposit    = "deposit"
	TransactionTypeWithdrawal = "withdrawal"
)
