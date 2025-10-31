package main

import "gorm.io/gorm"

// ============================================
// N+1問題を解決したコード
// ============================================

// FetchUsersWithPostsGood Preloadを使った最適化
// 2回のクエリで全データを取得（JOINまたはIN句を使用）
// - 1回目: 全ユーザーを取得
// - 2回目: 全投稿を一括取得（WHERE user_id IN (...)）
func FetchUsersWithPostsGood(db *gorm.DB) ([]User, error) {
	var users []User

	// Preloadで関連データを事前読み込み
	if err := db.Preload("Posts").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// FetchPostsWithCommentsGood Preloadを使った最適化
// 2回のクエリで全データを取得
func FetchPostsWithCommentsGood(db *gorm.DB) ([]Post, error) {
	var posts []Post

	// Preloadで関連データを事前読み込み
	if err := db.Preload("Comments").Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

// FetchUsersWithPostsAndCommentsGood 多段階Preloadを使った最適化
// 3回のクエリで全データを取得
// - 1回目: 全ユーザーを取得
// - 2回目: 全投稿を一括取得
// - 3回目: 全コメントを一括取得
func FetchUsersWithPostsAndCommentsGood(db *gorm.DB) ([]User, error) {
	var users []User

	// ネストしたPreloadで多段階の関連データを事前読み込み
	if err := db.Preload("Posts.Comments").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// FetchUsersWithPostsJoin JOINを使った最適化（別のアプローチ）
// 1回のクエリで取得（LEFT JOINを使用）
// データの重複が発生するため、データ量によっては逆に遅くなる可能性がある
func FetchUsersWithPostsJoin(db *gorm.DB) ([]User, error) {
	var users []User

	// JOINを使った取得
	if err := db.Joins("LEFT JOIN posts ON posts.user_id = users.id").
		Preload("Posts").
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// FetchUsersWithConditionGood 条件付きPreload
// 特定の条件に合う関連データのみを取得
func FetchUsersWithConditionGood(db *gorm.DB, minPostID uint) ([]User, error) {
	var users []User

	// 条件を指定してPreload
	if err := db.Preload("Posts", "id >= ?", minPostID).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
