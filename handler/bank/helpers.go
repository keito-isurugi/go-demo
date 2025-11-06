package bank

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// lockAccount 口座をロックして取得
func lockAccount(tx *gorm.DB, accountID uint) (*Account, error) {
	var account Account
	if err := tx.Raw("SELECT * FROM accounts WHERE id = ? FOR UPDATE", accountID).Scan(&account).Error; err != nil {
		return nil, fmt.Errorf("failed to lock account %d: %w", accountID, err)
	}
	return &account, nil
}

// getAccountIDsInOrder ロック順序を統一するため、IDを小さい順にソート
func getAccountIDsInOrder(id1, id2 uint) (firstID, secondID uint) {
	if id1 <= id2 {
		return id1, id2
	}
	return id2, id1
}

// transferFunds 残高を更新（ロック済みの口座を使用）
func transferFunds(fromAccount, toAccount *Account, amount int64) error {
	if fromAccount.Balance < amount {
		return fmt.Errorf("insufficient balance: %d < %d", fromAccount.Balance, amount)
	}

	fromAccount.Balance -= amount
	toAccount.Balance += amount
	return nil
}

// createTransaction 取引履歴を作成
func createTransaction(tx *gorm.DB, fromAccountID, toAccountID uint, amount int64) (*Transaction, error) {
	transaction := Transaction{
		FromAccountID:   fromAccountID,
		ToAccountID:     toAccountID,
		Amount:          amount,
		Status:          TransactionStatusCompleted,
		TransactionType: TransactionTypeTransfer,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return &transaction, nil
}

// executeTransferWithLockOrder ロック順序を統一した振込処理
func executeTransferWithLockOrder(
	ctx context.Context,
	db *gorm.DB,
	fromAccountID, toAccountID uint,
	amount int64,
) (*Account, *Account, *Transaction, error) {
	// ロック順序の統一
	firstID, secondID := getAccountIDsInOrder(fromAccountID, toAccountID)

	tx := db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1つ目の口座をロック
	firstAccount, err := lockAccount(tx, firstID)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	// 2つ目の口座をロック
	secondAccount, err := lockAccount(tx, secondID)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	// 送金元・送金先を特定
	var fromAccount, toAccount *Account
	if fromAccountID == firstID {
		fromAccount = firstAccount
		toAccount = secondAccount
	} else {
		fromAccount = secondAccount
		toAccount = firstAccount
	}

	// 残高チェックと更新
	if err := transferFunds(fromAccount, toAccount, amount); err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	// 残高を保存
	if err := tx.Save(fromAccount).Error; err != nil {
		tx.Rollback()
		return nil, nil, nil, fmt.Errorf("failed to update from account: %w", err)
	}

	if err := tx.Save(toAccount).Error; err != nil {
		tx.Rollback()
		return nil, nil, nil, fmt.Errorf("failed to update to account: %w", err)
	}

	// 取引履歴を作成
	transaction, err := createTransaction(tx, fromAccountID, toAccountID, amount)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	// コミット
	if err := tx.Commit().Error; err != nil {
		return nil, nil, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return fromAccount, toAccount, transaction, nil
}

// retryWithBackoff 指数バックオフでリトライ
func retryWithBackoff(
	maxRetries int,
	initialDelay time.Duration,
	fn func(attempt int) error,
) error {
	delay := initialDelay
	for attempt := 0; attempt < maxRetries; attempt++ {
		if err := fn(attempt); err != nil {
			if attempt == maxRetries-1 {
				return err
			}
			time.Sleep(delay)
			delay *= 2 // 指数バックオフ
			continue
		}
		return nil
	}
	return fmt.Errorf("max retries exceeded")
}

// jsonResponse JSONレスポンスを返すヘルパー
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// レスポンス書き込みエラーはログに記録（既にレスポンスを書き始めているため、エラーを返せない）
		// 実際のプロダクションでは適切なロギングライブラリを使用
		_ = err
	}
}
