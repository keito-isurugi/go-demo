package bank

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// DirtyReadDemoHandler Dirty Readを再現するデモ
func (h *Handler) DirtyReadDemoHandler(w http.ResponseWriter, r *http.Request) {
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

	action := r.URL.Query().Get("action")
	if action == "" {
		http.Error(w, "action parameter is required (read or update)", http.StatusBadRequest)
		return
	}

	switch action {
	case "update":
		h.handleDirtyReadUpdate(w, uint(accountID))
	case "read":
		h.handleDirtyReadRead(w, uint(accountID))
	default:
		http.Error(w, "Invalid action. Use 'read' or 'update'", http.StatusBadRequest)
	}
}

func (h *Handler) handleDirtyReadUpdate(w http.ResponseWriter, accountID uint) {
	tx := h.DB.Begin()
	tx.Exec("SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED")

	var account Account
	if err := tx.First(&account, accountID).Error; err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("Account not found: %v", err), http.StatusNotFound)
		return
	}

	account.Balance += 10000
	tx.Save(&account)

	time.Sleep(5 * time.Second)
	tx.Rollback()

	jsonResponse(w, map[string]interface{}{
		"success":    true,
		"message":    "Transaction rolled back (Dirty Read scenario)",
		"account_id": accountID,
	})
}

func (h *Handler) handleDirtyReadRead(w http.ResponseWriter, accountID uint) {
	tx := h.DB.Begin()
	tx.Exec("SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED")

	var account Account
	if err := tx.First(&account, accountID).Error; err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("Account not found: %v", err), http.StatusNotFound)
		return
	}

	time.Sleep(2 * time.Second)

	var accountAfter Account
	tx.First(&accountAfter, accountID)
	tx.Commit()

	jsonResponse(w, map[string]interface{}{
		"success":                    true,
		"message":                    "Dirty Read detected",
		"account_id":                 accountID,
		"balance_during_transaction": account.Balance,
		"balance_after_rollback":     accountAfter.Balance,
		"note":                       "If balance_during_transaction != balance_after_rollback, Dirty Read occurred",
	})
}

// PhantomReadDemoHandler Phantom Readを再現するデモ
func (h *Handler) PhantomReadDemoHandler(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	if action == "" {
		http.Error(w, "action parameter is required (read or insert)", http.StatusBadRequest)
		return
	}

	switch action {
	case "read":
		h.handlePhantomReadRead(w, r)
	case "insert":
		h.handlePhantomReadInsert(w)
	default:
		http.Error(w, "Invalid action. Use 'read' or 'insert'", http.StatusBadRequest)
	}
}

func (h *Handler) handlePhantomReadRead(w http.ResponseWriter, r *http.Request) {
	minBalanceStr := r.URL.Query().Get("min_balance")
	if minBalanceStr == "" {
		minBalanceStr = "50000"
	}

	minBalance, err := strconv.ParseInt(minBalanceStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid min_balance", http.StatusBadRequest)
		return
	}

	tx := h.DB.Begin()
	tx.Exec("SET TRANSACTION ISOLATION LEVEL READ COMMITTED")

	var accounts1 []Account
	tx.Where("balance >= ?", minBalance).Find(&accounts1)

	log.Printf("First read: Found %d accounts with balance >= %d", len(accounts1), minBalance)
	time.Sleep(5 * time.Second)

	var accounts2 []Account
	tx.Where("balance >= ?", minBalance).Find(&accounts2)
	tx.Commit()

	jsonResponse(w, map[string]interface{}{
		"success":               true,
		"message":               "Phantom Read check",
		"min_balance":           minBalance,
		"first_read_count":      len(accounts1),
		"second_read_count":     len(accounts2),
		"first_read_accounts":   accounts1,
		"second_read_accounts":  accounts2,
		"phantom_read_occurred": len(accounts2) > len(accounts1),
		"note":                  "If second_read_count > first_read_count, Phantom Read occurred",
	})
}

func (h *Handler) handlePhantomReadInsert(w http.ResponseWriter) {
	tx := h.DB.Begin()
	tx.Exec("SET TRANSACTION ISOLATION LEVEL READ COMMITTED")

	newAccount := Account{
		AccountNo: fmt.Sprintf("PHANTOM-%d", time.Now().Unix()),
		Balance:   150000,
		OwnerName: "Phantom User",
	}

	if err := tx.Create(&newAccount).Error; err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("Failed to create account: %v", err), http.StatusInternalServerError)
		return
	}

	time.Sleep(2 * time.Second)

	if err := tx.Commit().Error; err != nil {
		http.Error(w, fmt.Sprintf("Failed to commit: %v", err), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"message": "New account inserted (Phantom Read scenario)",
		"account": newAccount,
	})
}

// DeadlockDemoHandler デッドロックを再現するデモ
func (h *Handler) DeadlockDemoHandler(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	if action == "" {
		http.Error(w, "action parameter is required (tx1 or tx2)", http.StatusBadRequest)
		return
	}

	account1Str := r.URL.Query().Get("account1")
	account2Str := r.URL.Query().Get("account2")
	if account1Str == "" || account2Str == "" {
		http.Error(w, "account1 and account2 parameters are required", http.StatusBadRequest)
		return
	}

	account1ID, err := strconv.ParseUint(account1Str, 10, 32)
	if err != nil {
		http.Error(w, "Invalid account1", http.StatusBadRequest)
		return
	}

	account2ID, err := strconv.ParseUint(account2Str, 10, 32)
	if err != nil {
		http.Error(w, "Invalid account2", http.StatusBadRequest)
		return
	}

	switch action {
	case "tx1":
		h.handleDeadlockTX1(w, uint(account1ID), uint(account2ID))
	case "tx2":
		h.handleDeadlockTX2(w, uint(account1ID), uint(account2ID))
	default:
		http.Error(w, "Invalid action. Use 'tx1' or 'tx2'", http.StatusBadRequest)
	}
}

func (h *Handler) handleDeadlockTX1(w http.ResponseWriter, account1ID, account2ID uint) {
	tx := h.DB.Begin()

	log.Printf("TX1: Locking account %d", account1ID)
	account1, err := lockAccount(tx, account1ID)
	if err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("Failed to lock account1: %v", err), http.StatusInternalServerError)
		return
	}

	time.Sleep(2 * time.Second)

	log.Printf("TX1: Trying to lock account %d", account2ID)
	account2, err := lockAccount(tx, account2ID)
	if err != nil {
		tx.Rollback()
		jsonResponse(w, map[string]interface{}{
			"success": false,
			"message": "Deadlock detected in TX1",
			"error":   err.Error(),
		})
		return
	}

	account1.Balance -= 1000
	account2.Balance += 1000
	tx.Save(account1)
	tx.Save(account2)

	if err := tx.Commit().Error; err != nil {
		jsonResponse(w, map[string]interface{}{
			"success": false,
			"message": "Failed to commit TX1",
			"error":   err.Error(),
		})
		return
	}

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"message": "TX1 completed successfully",
	})
}

func (h *Handler) handleDeadlockTX2(w http.ResponseWriter, account1ID, account2ID uint) {
	tx := h.DB.Begin()

	log.Printf("TX2: Locking account %d", account2ID)
	account2, err := lockAccount(tx, account2ID)
	if err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("Failed to lock account2: %v", err), http.StatusInternalServerError)
		return
	}

	time.Sleep(2 * time.Second)

	log.Printf("TX2: Trying to lock account %d", account1ID)
	account1, err := lockAccount(tx, account1ID)
	if err != nil {
		tx.Rollback()
		jsonResponse(w, map[string]interface{}{
			"success": false,
			"message": "Deadlock detected in TX2",
			"error":   err.Error(),
		})
		return
	}

	account2.Balance -= 1000
	account1.Balance += 1000
	tx.Save(account2)
	tx.Save(account1)

	if err := tx.Commit().Error; err != nil {
		jsonResponse(w, map[string]interface{}{
			"success": false,
			"message": "Failed to commit TX2",
			"error":   err.Error(),
		})
		return
	}

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"message": "TX2 completed successfully",
	})
}
