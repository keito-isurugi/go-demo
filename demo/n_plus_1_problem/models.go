package main

import (
	"fmt"
	"math/rand"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User ユーザーモデル
type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:100"`
	Posts []Post `gorm:"foreignKey:UserID"`
}

// Post 投稿モデル
type Post struct {
	ID       uint      `gorm:"primaryKey"`
	UserID   uint      `gorm:"index"`
	Title    string    `gorm:"size:200"`
	Content  string    `gorm:"size:1000"`
	Comments []Comment `gorm:"foreignKey:PostID"`
}

// Comment コメントモデル
type Comment struct {
	ID      uint   `gorm:"primaryKey"`
	PostID  uint   `gorm:"index"`
	Content string `gorm:"size:500"`
}

// InitDB データベースを初期化
func InitDB(showSQL bool) (*gorm.DB, error) {
	logLevel := logger.Silent
	if showSQL {
		logLevel = logger.Info
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// テーブルを作成
	if err := db.AutoMigrate(&User{}, &Post{}, &Comment{}); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	return db, nil
}

// SeedData サンプルデータを生成
func SeedData(db *gorm.DB, userCount, postsPerUser, commentsPerPost int) error {
	// 既存データをクリア
	db.Exec("DELETE FROM comments")
	db.Exec("DELETE FROM posts")
	db.Exec("DELETE FROM users")

	fmt.Printf("サンプルデータを生成中... (ユーザー: %d, 投稿/ユーザー: %d, コメント/投稿: %d)\n",
		userCount, postsPerUser, commentsPerPost)

	for i := 1; i <= userCount; i++ {
		user := User{
			Name: fmt.Sprintf("User_%d", i),
		}

		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// 各ユーザーに投稿を作成
		for j := 1; j <= postsPerUser; j++ {
			post := Post{
				UserID:  user.ID,
				Title:   fmt.Sprintf("Post %d by User %d", j, i),
				Content: fmt.Sprintf("This is the content of post %d by user %d", j, i),
			}

			if err := db.Create(&post).Error; err != nil {
				return fmt.Errorf("failed to create post: %w", err)
			}

			// 各投稿にコメントを作成
			for k := 1; k <= commentsPerPost; k++ {
				comment := Comment{
					PostID:  post.ID,
					Content: fmt.Sprintf("Comment %d on post %d", k, post.ID),
				}

				if err := db.Create(&comment).Error; err != nil {
					return fmt.Errorf("failed to create comment: %w", err)
				}
			}
		}
	}

	var totalUsers, totalPosts, totalComments int64
	db.Model(&User{}).Count(&totalUsers)
	db.Model(&Post{}).Count(&totalPosts)
	db.Model(&Comment{}).Count(&totalComments)

	fmt.Printf("データ生成完了: ユーザー %d件, 投稿 %d件, コメント %d件\n\n",
		totalUsers, totalPosts, totalComments)

	return nil
}

// generateRandomString ランダムな文字列を生成
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
