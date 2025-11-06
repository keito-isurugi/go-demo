package bank

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// Handler は銀行振込APIのハンドラー
type Handler struct {
	DB *gorm.DB
}

// InitAccountsHandler テスト用の口座を初期化
func (h *Handler) InitAccountsHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.DB.AutoMigrate(&Account{}, &Transaction{}); err != nil {
		http.Error(w, fmt.Sprintf("Failed to migrate: %v", err), http.StatusInternalServerError)
		return
	}

	h.DB.Exec("TRUNCATE TABLE accounts CASCADE")
	h.DB.Exec("TRUNCATE TABLE transactions CASCADE")

	accounts := []Account{
		{AccountNo: "1001", Balance: 100000, OwnerName: "山田太郎"},
		{AccountNo: "1002", Balance: 50000, OwnerName: "佐藤花子"},
		{AccountNo: "1003", Balance: 200000, OwnerName: "鈴木一郎"},
	}

	for _, account := range accounts {
		if err := h.DB.Create(&account).Error; err != nil {
			http.Error(w, fmt.Sprintf("Failed to create account: %v", err), http.StatusInternalServerError)
			return
		}
	}

	jsonResponse(w, map[string]interface{}{
		"success":  true,
		"message":  "Accounts initialized successfully",
		"accounts": accounts,
	})
}

// NormalTransferHandler 通常の振込処理
func (h *Handler) NormalTransferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var fromAccount Account
	if err := tx.First(&fromAccount, req.FromAccountID).Error; err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("From account not found: %v", err), http.StatusNotFound)
		return
	}

	if fromAccount.Balance < req.Amount {
		tx.Rollback()
		jsonResponse(w, TransferResponse{
			Success: false,
			Message: "Insufficient balance",
		})
		return
	}

	var toAccount Account
	if err := tx.First(&toAccount, req.ToAccountID).Error; err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("To account not found: %v", err), http.StatusNotFound)
		return
	}

	fromAccount.Balance -= req.Amount
	toAccount.Balance += req.Amount

	if err := tx.Save(&fromAccount).Error; err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("Failed to update from account: %v", err), http.StatusInternalServerError)
		return
	}

	if err := tx.Save(&toAccount).Error; err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("Failed to update to account: %v", err), http.StatusInternalServerError)
		return
	}

	transaction, err := createTransaction(tx, req.FromAccountID, req.ToAccountID, req.Amount)
	if err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("Failed to create transaction: %v", err), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		http.Error(w, fmt.Sprintf("Failed to commit transaction: %v", err), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, TransferResponse{
		Success:       true,
		Message:       "Transfer completed successfully",
		TransactionID: transaction.ID,
		FromBalance:   fromAccount.Balance,
		ToBalance:     toAccount.Balance,
	})
}

// GetAccountHandler 口座情報を取得
func (h *Handler) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	accountIDStr := r.URL.Query().Get("account_id")
	if accountIDStr == "" {
		http.Error(w, "account_id parameter is required", http.StatusBadRequest)
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid account_id", http.StatusBadRequest)
		return
	}

	var account Account
	if err := h.DB.First(&account, accountID).Error; err != nil {
		http.Error(w, fmt.Sprintf("Account not found: %v", err), http.StatusNotFound)
		return
	}

	jsonResponse(w, account)
}

// ListAccountsHandler 全口座一覧を取得
func (h *Handler) ListAccountsHandler(w http.ResponseWriter, r *http.Request) {
	var accounts []Account
	if err := h.DB.Find(&accounts).Error; err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch accounts: %v", err), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, accounts)
}

// DeadlockAvoidanceHandler デッドロック回避策1: ロック順序の統一
func (h *Handler) DeadlockAvoidanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	fromAccount, toAccount, transaction, err := executeTransferWithLockOrder(
		ctx, h.DB, req.FromAccountID, req.ToAccountID, req.Amount,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, TransferResponse{
		Success:       true,
		Message:       "Transfer completed successfully (deadlock avoided by lock ordering)",
		TransactionID: transaction.ID,
		FromBalance:   fromAccount.Balance,
		ToBalance:     toAccount.Balance,
	})
}

// DeadlockTimeoutHandler デッドロック回避策2: タイムアウト設定
func (h *Handler) DeadlockTimeoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	maxRetries := 3
	initialDelay := 100 * time.Millisecond
	timeout := 5 * time.Second

	var fromAccount, toAccount *Account
	var transaction *Transaction
	var lastErr error

	err := retryWithBackoff(maxRetries, initialDelay, func(attempt int) error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		var err error
		fromAccount, toAccount, transaction, err = executeTransferWithLockOrder(
			ctx, h.DB, req.FromAccountID, req.ToAccountID, req.Amount,
		)

		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				log.Printf("Transaction timeout, retrying... (attempt %d/%d)", attempt+1, maxRetries)
				lastErr = err
				return err
			}
			lastErr = err
			return err
		}

		return nil
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Transaction failed: %v", lastErr), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, TransferResponse{
		Success:       true,
		Message:       "Transfer completed successfully (with retry mechanism)",
		TransactionID: transaction.ID,
		FromBalance:   fromAccount.Balance,
		ToBalance:     toAccount.Balance,
	})
}
