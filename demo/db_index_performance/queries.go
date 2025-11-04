package main

import "gorm.io/gorm"

// ============================================
// 商品テーブルのクエリ
// ============================================

// FindProductsByCategory カテゴリで商品を検索（インデックスなし）
func FindProductsByCategory(db *gorm.DB, categoryID int) ([]Product, error) {
	var products []Product
	err := db.Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}

// FindProductsByCategoryIndexed カテゴリで商品を検索（インデックスあり）
func FindProductsByCategoryIndexed(db *gorm.DB, categoryID int) ([]ProductIndexed, error) {
	var products []ProductIndexed
	err := db.Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}

// FindProductsByPriceRange 価格範囲で商品を検索（インデックスなし）
func FindProductsByPriceRange(db *gorm.DB, minPrice, maxPrice int) ([]Product, error) {
	var products []Product
	err := db.Where("price BETWEEN ? AND ?", minPrice, maxPrice).Find(&products).Error
	return products, err
}

// FindProductsByPriceRangeIndexed 価格範囲で商品を検索（インデックスあり）
func FindProductsByPriceRangeIndexed(db *gorm.DB, minPrice, maxPrice int) ([]ProductIndexed, error) {
	var products []ProductIndexed
	err := db.Where("price BETWEEN ? AND ?", minPrice, maxPrice).Find(&products).Error
	return products, err
}

// FindProductsByName 商品名で検索（LIKE）（インデックスなし）
func FindProductsByName(db *gorm.DB, keyword string) ([]Product, error) {
	var products []Product
	err := db.Where("name LIKE ?", "%"+keyword+"%").Find(&products).Error
	return products, err
}

// FindProductsByNameIndexed 商品名で検索（LIKE）（インデックスあり）
func FindProductsByNameIndexed(db *gorm.DB, keyword string) ([]ProductIndexed, error) {
	var products []ProductIndexed
	err := db.Where("name LIKE ?", "%"+keyword+"%").Find(&products).Error
	return products, err
}

// FindLowStockProducts 在庫が少ない商品を検索（インデックスなし）
func FindLowStockProducts(db *gorm.DB, threshold int) ([]Product, error) {
	var products []Product
	err := db.Where("stock < ?", threshold).Find(&products).Error
	return products, err
}

// FindLowStockProductsIndexed 在庫が少ない商品を検索（インデックスあり）
func FindLowStockProductsIndexed(db *gorm.DB, threshold int) ([]ProductIndexed, error) {
	var products []ProductIndexed
	err := db.Where("stock < ?", threshold).Find(&products).Error
	return products, err
}

// ============================================
// 注文テーブルのクエリ
// ============================================

// FindOrdersByUserID ユーザーIDで注文を検索（インデックスなし）
func FindOrdersByUserID(db *gorm.DB, userID int) ([]Order, error) {
	var orders []Order
	err := db.Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

// FindOrdersByUserIDIndexed ユーザーIDで注文を検索（インデックスあり）
func FindOrdersByUserIDIndexed(db *gorm.DB, userID int) ([]OrderIndexed, error) {
	var orders []OrderIndexed
	err := db.Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

// FindOrdersByStatus ステータスで注文を検索（インデックスなし）
func FindOrdersByStatus(db *gorm.DB, status string) ([]Order, error) {
	var orders []Order
	err := db.Where("status = ?", status).Find(&orders).Error
	return orders, err
}

// FindOrdersByStatusIndexed ステータスで注文を検索（インデックスあり）
func FindOrdersByStatusIndexed(db *gorm.DB, status string) ([]OrderIndexed, error) {
	var orders []OrderIndexed
	err := db.Where("status = ?", status).Find(&orders).Error
	return orders, err
}

// FindRecentOrders 最近の注文を検索（インデックスなし）
func FindRecentOrders(db *gorm.DB, days int) ([]Order, error) {
	var orders []Order
	err := db.Where("order_date >= datetime('now', '-' || ? || ' days')", days).Find(&orders).Error
	return orders, err
}

// FindRecentOrdersIndexed 最近の注文を検索（インデックスあり）
func FindRecentOrdersIndexed(db *gorm.DB, days int) ([]OrderIndexed, error) {
	var orders []OrderIndexed
	err := db.Where("order_date >= datetime('now', '-' || ? || ' days')", days).Find(&orders).Error
	return orders, err
}

// ============================================
// 従業員テーブルのクエリ（複合インデックス）
// ============================================

// FindEmployeesByName 名前で従業員を検索（インデックスなし）
func FindEmployeesByName(db *gorm.DB, firstName, lastName string) ([]Employee, error) {
	var employees []Employee
	err := db.Where("first_name = ? AND last_name = ?", firstName, lastName).Find(&employees).Error
	return employees, err
}

// FindEmployeesByNameIndexed 名前で従業員を検索（複合インデックスあり）
func FindEmployeesByNameIndexed(db *gorm.DB, firstName, lastName string) ([]EmployeeIndexed, error) {
	var employees []EmployeeIndexed
	err := db.Where("first_name = ? AND last_name = ?", firstName, lastName).Find(&employees).Error
	return employees, err
}

// FindEmployeesByDepartmentAndPosition 部署と役職で従業員を検索（インデックスなし）
func FindEmployeesByDepartmentAndPosition(db *gorm.DB, department, position string) ([]Employee, error) {
	var employees []Employee
	err := db.Where("department = ? AND position = ?", department, position).Find(&employees).Error
	return employees, err
}

// FindEmployeesByDepartmentAndPositionIndexed 部署と役職で従業員を検索（複合インデックスあり）
func FindEmployeesByDepartmentAndPositionIndexed(db *gorm.DB, department, position string) ([]EmployeeIndexed, error) {
	var employees []EmployeeIndexed
	err := db.Where("department = ? AND position = ?", department, position).Find(&employees).Error
	return employees, err
}

// FindEmployeesBySalaryRange 給与範囲で従業員を検索（インデックスなし）
func FindEmployeesBySalaryRange(db *gorm.DB, minSalary, maxSalary int) ([]Employee, error) {
	var employees []Employee
	err := db.Where("salary BETWEEN ? AND ?", minSalary, maxSalary).Find(&employees).Error
	return employees, err
}

// FindEmployeesBySalaryRangeIndexed 給与範囲で従業員を検索（インデックスあり）
func FindEmployeesBySalaryRangeIndexed(db *gorm.DB, minSalary, maxSalary int) ([]EmployeeIndexed, error) {
	var employees []EmployeeIndexed
	err := db.Where("salary BETWEEN ? AND ?", minSalary, maxSalary).Find(&employees).Error
	return employees, err
}

// ============================================
// ソート処理のクエリ
// ============================================

// FindProductsSortedByPrice 価格でソート（インデックスなし）
func FindProductsSortedByPrice(db *gorm.DB, limit int) ([]Product, error) {
	var products []Product
	err := db.Order("price DESC").Limit(limit).Find(&products).Error
	return products, err
}

// FindProductsSortedByPriceIndexed 価格でソート（インデックスあり）
func FindProductsSortedByPriceIndexed(db *gorm.DB, limit int) ([]ProductIndexed, error) {
	var products []ProductIndexed
	err := db.Order("price DESC").Limit(limit).Find(&products).Error
	return products, err
}

// FindOrdersSortedByDate 日付でソート（インデックスなし）
func FindOrdersSortedByDate(db *gorm.DB, limit int) ([]Order, error) {
	var orders []Order
	err := db.Order("order_date DESC").Limit(limit).Find(&orders).Error
	return orders, err
}

// FindOrdersSortedByDateIndexed 日付でソート（インデックスあり）
func FindOrdersSortedByDateIndexed(db *gorm.DB, limit int) ([]OrderIndexed, error) {
	var orders []OrderIndexed
	err := db.Order("order_date DESC").Limit(limit).Find(&orders).Error
	return orders, err
}
