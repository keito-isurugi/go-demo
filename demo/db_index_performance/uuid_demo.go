package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// OrderUUID UUIDを主キーとして使うパターン
type OrderUUID struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)"` // UUID
	UserID     int       `gorm:"index:idx_order_uuid_user_id"`
	TotalPrice int
	Status     string    `gorm:"size:50"`
	OrderDate  time.Time
}

// OrderSequential 連番を主キーとして使うパターン（比較用）
type OrderSequential struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"` // 連番
	UserID     int       `gorm:"index:idx_order_seq_user_id"`
	TotalPrice int
	Status     string    `gorm:"size:50"`
	OrderDate  time.Time
}

func main() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("UUID vs 連番のインデックスパフォーマンス比較")
	fmt.Println(strings.Repeat("=", 80))

	// UUID版のDB
	dbUUID, _ := gorm.Open(sqlite.Open("uuid_test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	dbUUID.AutoMigrate(&OrderUUID{})
	dbUUID.Exec("DELETE FROM order_uuids")

	// 連番版のDB
	dbSeq, _ := gorm.Open(sqlite.Open("sequential_test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	dbSeq.AutoMigrate(&OrderSequential{})
	dbSeq.Exec("DELETE FROM order_sequentials")

	count := 10000

	// UUID版：データ挿入
	fmt.Printf("\n[UUID版] %d件のデータを挿入中...\n", count)
	startUUID := time.Now()
	for i := 0; i < count; i++ {
		order := OrderUUID{
			ID:         uuid.New().String(),
			UserID:     i % 1000,
			TotalPrice: 5000,
			Status:     "pending",
			OrderDate:  time.Now(),
		}
		dbUUID.Create(&order)
	}
	uuidInsertTime := time.Since(startUUID)
	fmt.Printf("UUID版 挿入時間: %v\n", uuidInsertTime)

	// 連番版：データ挿入
	fmt.Printf("\n[連番版] %d件のデータを挿入中...\n", count)
	startSeq := time.Now()
	for i := 0; i < count; i++ {
		order := OrderSequential{
			UserID:     i % 1000,
			TotalPrice: 5000,
			Status:     "pending",
			OrderDate:  time.Now(),
		}
		dbSeq.Create(&order)
	}
	seqInsertTime := time.Since(startSeq)
	fmt.Printf("連番版 挿入時間: %v\n", seqInsertTime)

	fmt.Printf("\n📊 挿入速度比較: 連番版は UUID版の %.2f倍速い\n",
		float64(uuidInsertTime)/float64(seqInsertTime))

	// 検索パフォーマンス比較
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("検索パフォーマンス比較")
	fmt.Println(strings.Repeat("=", 80))

	// UUID版で検索
	var ordersUUID []OrderUUID
	startUUID = time.Now()
	dbUUID.Where("user_id = ?", 500).Find(&ordersUUID)
	uuidSearchTime := time.Since(startUUID)
	fmt.Printf("UUID版 検索時間: %v (件数: %d)\n", uuidSearchTime, len(ordersUUID))

	// 連番版で検索
	var ordersSeq []OrderSequential
	startSeq = time.Now()
	dbSeq.Where("user_id = ?", 500).Find(&ordersSeq)
	seqSearchTime := time.Since(startSeq)
	fmt.Printf("連番版 検索時間: %v (件数: %d)\n", seqSearchTime, len(ordersSeq))

	if uuidSearchTime > seqSearchTime {
		fmt.Printf("\n📊 検索速度比較: 連番版は UUID版の %.2f倍速い\n",
			float64(uuidSearchTime)/float64(seqSearchTime))
	} else {
		fmt.Printf("\n📊 検索速度比較: ほぼ同等\n")
	}
}
