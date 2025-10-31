package main

import (
	"fmt"

	"gorm.io/gorm"
)

// ============================================
// N+1問題のあるコード
// ============================================

// FetchUsersWithPostsBad N+1問題のある実装
// 1回のクエリでユーザーを取得し、各ユーザーごとに1回ずつ投稿を取得する
// ユーザーがN人いる場合、1 + N回のクエリが発行される
func FetchUsersWithPostsBad(db *gorm.DB) ([]User, error) {
	var users []User

	// 1回目: 全ユーザーを取得
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	// N回: 各ユーザーの投稿を個別に取得（N+1問題！）
	for i := range users {
		if err := db.Where("user_id = ?", users[i].ID).Find(&users[i].Posts).Error; err != nil {
			return nil, err
		}
	}

	return users, nil
}

// FetchPostsWithCommentsBad N+1問題のある実装（投稿とコメント）
// 投稿がN件ある場合、1 + N回のクエリが発行される
func FetchPostsWithCommentsBad(db *gorm.DB) ([]Post, error) {
	var posts []Post

	// 1回目: 全投稿を取得
	if err := db.Find(&posts).Error; err != nil {
		return nil, err
	}

	// N回: 各投稿のコメントを個別に取得（N+1問題！）
	for i := range posts {
		if err := db.Where("post_id = ?", posts[i].ID).Find(&posts[i].Comments).Error; err != nil {
			return nil, err
		}
	}

	return posts, nil
}

// FetchUsersWithPostsAndCommentsBad 多段階のN+1問題
// ユーザー → 投稿 → コメント と3階層で取得
// ユーザーがU人、投稿がP件ある場合、1 + U + P回のクエリが発行される
func FetchUsersWithPostsAndCommentsBad(db *gorm.DB) ([]User, error) {
	var users []User

	// 1回目: 全ユーザーを取得
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	// U回: 各ユーザーの投稿を個別に取得
	for i := range users {
		if err := db.Where("user_id = ?", users[i].ID).Find(&users[i].Posts).Error; err != nil {
			return nil, err
		}

		// P回: 各投稿のコメントを個別に取得
		for j := range users[i].Posts {
			if err := db.Where("post_id = ?", users[i].Posts[j].ID).Find(&users[i].Posts[j].Comments).Error; err != nil {
				return nil, err
			}
		}
	}

	return users, nil
}

// PrintUsersSummary ユーザー情報のサマリーを表示
func PrintUsersSummary(users []User) {
	fmt.Printf("取得したユーザー数: %d\n", len(users))
	for _, user := range users {
		fmt.Printf("  - %s: %d件の投稿\n", user.Name, len(user.Posts))
		for _, post := range user.Posts {
			if len(post.Comments) > 0 {
				fmt.Printf("    - %s: %d件のコメント\n", post.Title, len(post.Comments))
			}
		}
	}
}

// PrintPostsSummary 投稿情報のサマリーを表示
func PrintPostsSummary(posts []Post) {
	fmt.Printf("取得した投稿数: %d\n", len(posts))
	for i, post := range posts {
		if i < 3 { // 最初の3件のみ表示
			fmt.Printf("  - %s: %d件のコメント\n", post.Title, len(post.Comments))
		}
	}
	if len(posts) > 3 {
		fmt.Printf("  ... 他 %d件\n", len(posts)-3)
	}
}
