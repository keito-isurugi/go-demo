package main

import (
	"fmt"
	"math/rand"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Product 商品テーブル（インデックスなし用）
type Product struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:200"`
	CategoryID  int       // インデックスなし
	Price       int       // インデックスなし
	Stock       int       // インデックスなし
	Description string    `gorm:"size:500"`
	CreatedAt   time.Time // インデックスなし
}

// ProductIndexed 商品テーブル（インデックスあり用）
type ProductIndexed struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:200;index:idx_product_name"`                  // 単一カラムインデックス
	CategoryID  int       `gorm:"index:idx_product_category"`                       // 単一カラムインデックス
	Price       int       `gorm:"index:idx_product_price"`                          // 単一カラムインデックス
	Stock       int       `gorm:"index:idx_product_stock"`                          // 単一カラムインデックス
	Description string    `gorm:"size:500"`
	CreatedAt   time.Time `gorm:"index:idx_product_created_at"`                     // 単一カラムインデックス
}

// Order 注文テーブル（インデックスなし用）
type Order struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     int       // インデックスなし
	TotalPrice int
	Status     string    `gorm:"size:50"` // インデックスなし
	OrderDate  time.Time // インデックスなし
}

// OrderIndexed 注文テーブル（インデックスあり用）
type OrderIndexed struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     int       `gorm:"index:idx_order_user_id"`                        // 単一カラムインデックス
	TotalPrice int
	Status     string    `gorm:"size:50;index:idx_order_status"`                 // 単一カラムインデックス
	OrderDate  time.Time `gorm:"index:idx_order_date"`                           // 単一カラムインデックス
}

// Employee 従業員テーブル（複合インデックス用）
type Employee struct {
	ID         uint   `gorm:"primaryKey"`
	FirstName  string `gorm:"size:100"`
	LastName   string `gorm:"size:100"`
	Department string `gorm:"size:100"`
	Position   string `gorm:"size:100"`
	Salary     int
}

// EmployeeIndexed 従業員テーブル（複合インデックスあり）
type EmployeeIndexed struct {
	ID         uint   `gorm:"primaryKey"`
	FirstName  string `gorm:"size:100;index:idx_employee_name,priority:1"`            // 複合インデックス (FirstName, LastName)
	LastName   string `gorm:"size:100;index:idx_employee_name,priority:2"`            // 複合インデックス (FirstName, LastName)
	Department string `gorm:"size:100;index:idx_employee_dept_pos,priority:1"`        // 複合インデックス (Department, Position)
	Position   string `gorm:"size:100;index:idx_employee_dept_pos,priority:2"`        // 複合インデックス (Department, Position)
	Salary     int    `gorm:"index:idx_employee_salary"`                              // 単一カラムインデックス
}

// InitDB データベースを初期化
func InitDB(dbName string, showSQL bool) (*gorm.DB, error) {
	logLevel := logger.Silent
	if showSQL {
		logLevel = logger.Info
	}

	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return db, nil
}

// SetupNoIndexDB インデックスなしのDBをセットアップ
func SetupNoIndexDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&Product{}, &Order{}, &Employee{}); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}
	return nil
}

// SetupIndexedDB インデックスありのDBをセットアップ
func SetupIndexedDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&ProductIndexed{}, &OrderIndexed{}, &EmployeeIndexed{}); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}
	return nil
}

// SeedProducts 商品データを生成（インデックスなし）
func SeedProducts(db *gorm.DB, count int) error {
	db.Exec("DELETE FROM products")

	fmt.Printf("商品データを生成中... (%d件)\n", count)

	categories := []string{"Electronics", "Clothing", "Food", "Books", "Sports"}

	for i := 1; i <= count; i++ {
		product := Product{
			Name:        fmt.Sprintf("Product_%d", i),
			CategoryID:  rand.Intn(len(categories)),
			Price:       rand.Intn(10000) + 100,
			Stock:       rand.Intn(1000),
			Description: fmt.Sprintf("Description for product %d", i),
			CreatedAt:   time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour),
		}

		if err := db.Create(&product).Error; err != nil {
			return fmt.Errorf("failed to create product: %w", err)
		}
	}

	var total int64
	db.Model(&Product{}).Count(&total)
	fmt.Printf("商品データ生成完了: %d件\n", total)

	return nil
}

// SeedProductsIndexed 商品データを生成（インデックスあり）
func SeedProductsIndexed(db *gorm.DB, count int) error {
	db.Exec("DELETE FROM product_indexeds")

	fmt.Printf("商品データ（インデックスあり）を生成中... (%d件)\n", count)

	categories := []string{"Electronics", "Clothing", "Food", "Books", "Sports"}

	for i := 1; i <= count; i++ {
		product := ProductIndexed{
			Name:        fmt.Sprintf("Product_%d", i),
			CategoryID:  rand.Intn(len(categories)),
			Price:       rand.Intn(10000) + 100,
			Stock:       rand.Intn(1000),
			Description: fmt.Sprintf("Description for product %d", i),
			CreatedAt:   time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour),
		}

		if err := db.Create(&product).Error; err != nil {
			return fmt.Errorf("failed to create product: %w", err)
		}
	}

	var total int64
	db.Model(&ProductIndexed{}).Count(&total)
	fmt.Printf("商品データ（インデックスあり）生成完了: %d件\n", total)

	return nil
}

// SeedOrders 注文データを生成（インデックスなし）
func SeedOrders(db *gorm.DB, count int) error {
	db.Exec("DELETE FROM orders")

	fmt.Printf("注文データを生成中... (%d件)\n", count)

	statuses := []string{"pending", "processing", "shipped", "delivered", "cancelled"}

	for i := 1; i <= count; i++ {
		order := Order{
			UserID:     rand.Intn(1000) + 1,
			TotalPrice: rand.Intn(100000) + 1000,
			Status:     statuses[rand.Intn(len(statuses))],
			OrderDate:  time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour),
		}

		if err := db.Create(&order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}
	}

	var total int64
	db.Model(&Order{}).Count(&total)
	fmt.Printf("注文データ生成完了: %d件\n", total)

	return nil
}

// SeedOrdersIndexed 注文データを生成（インデックスあり）
func SeedOrdersIndexed(db *gorm.DB, count int) error {
	db.Exec("DELETE FROM order_indexeds")

	fmt.Printf("注文データ（インデックスあり）を生成中... (%d件)\n", count)

	statuses := []string{"pending", "processing", "shipped", "delivered", "cancelled"}

	for i := 1; i <= count; i++ {
		order := OrderIndexed{
			UserID:     rand.Intn(1000) + 1,
			TotalPrice: rand.Intn(100000) + 1000,
			Status:     statuses[rand.Intn(len(statuses))],
			OrderDate:  time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour),
		}

		if err := db.Create(&order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}
	}

	var total int64
	db.Model(&OrderIndexed{}).Count(&total)
	fmt.Printf("注文データ（インデックスあり）生成完了: %d件\n", total)

	return nil
}

// SeedEmployees 従業員データを生成（インデックスなし）
func SeedEmployees(db *gorm.DB, count int) error {
	db.Exec("DELETE FROM employees")

	fmt.Printf("従業員データを生成中... (%d件)\n", count)

	departments := []string{"Engineering", "Sales", "Marketing", "HR", "Finance"}
	positions := []string{"Junior", "Senior", "Lead", "Manager", "Director"}
	firstNames := []string{"John", "Jane", "Mike", "Emily", "David", "Sarah", "Tom", "Lisa"}
	lastNames := []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis"}

	for i := 1; i <= count; i++ {
		employee := Employee{
			FirstName:  firstNames[rand.Intn(len(firstNames))],
			LastName:   lastNames[rand.Intn(len(lastNames))],
			Department: departments[rand.Intn(len(departments))],
			Position:   positions[rand.Intn(len(positions))],
			Salary:     rand.Intn(150000) + 50000,
		}

		if err := db.Create(&employee).Error; err != nil {
			return fmt.Errorf("failed to create employee: %w", err)
		}
	}

	var total int64
	db.Model(&Employee{}).Count(&total)
	fmt.Printf("従業員データ生成完了: %d件\n", total)

	return nil
}

// SeedEmployeesIndexed 従業員データを生成（インデックスあり）
func SeedEmployeesIndexed(db *gorm.DB, count int) error {
	db.Exec("DELETE FROM employee_indexeds")

	fmt.Printf("従業員データ（インデックスあり）を生成中... (%d件)\n", count)

	departments := []string{"Engineering", "Sales", "Marketing", "HR", "Finance"}
	positions := []string{"Junior", "Senior", "Lead", "Manager", "Director"}
	firstNames := []string{"John", "Jane", "Mike", "Emily", "David", "Sarah", "Tom", "Lisa"}
	lastNames := []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis"}

	for i := 1; i <= count; i++ {
		employee := EmployeeIndexed{
			FirstName:  firstNames[rand.Intn(len(firstNames))],
			LastName:   lastNames[rand.Intn(len(lastNames))],
			Department: departments[rand.Intn(len(departments))],
			Position:   positions[rand.Intn(len(positions))],
			Salary:     rand.Intn(150000) + 50000,
		}

		if err := db.Create(&employee).Error; err != nil {
			return fmt.Errorf("failed to create employee: %w", err)
		}
	}

	var total int64
	db.Model(&EmployeeIndexed{}).Count(&total)
	fmt.Printf("従業員データ（インデックスあり）生成完了: %d件\n", total)

	return nil
}
