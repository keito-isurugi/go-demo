# レイヤードアーキテクチャ (Layered Architecture)

## 概要

レイヤードアーキテクチャは、アプリケーションを**水平方向の層（レイヤー）**に分割するアーキテクチャパターンです。
最も歴史が長く、多くの開発者にとって直感的に理解しやすい構造です。

---

## 1. 基本構造

### 1.1 概念図

```
┌─────────────────────────────────────────────────────────┐
│                  Presentation Layer                     │
│                 （プレゼンテーション層）                  │
│                                                         │
│  責務: UI表示、ユーザー入力の受付、リクエスト/レスポンス   │
│  例: HTTP Handler, CLI, GUI                             │
├─────────────────────────────────────────────────────────┤
│                         ↓↑                              │
├─────────────────────────────────────────────────────────┤
│                    Business Layer                       │
│                    （ビジネス層）                         │
│                                                         │
│  責務: ビジネスロジック、業務ルール、計算、判断            │
│  例: UserService, OrderService                          │
├─────────────────────────────────────────────────────────┤
│                         ↓↑                              │
├─────────────────────────────────────────────────────────┤
│                   Persistence Layer                     │
│                    （永続化層）                          │
│                                                         │
│  責務: データの保存・取得、CRUD操作                       │
│  例: UserRepository, OrderRepository                    │
├─────────────────────────────────────────────────────────┤
│                         ↓↑                              │
├─────────────────────────────────────────────────────────┤
│                    Database Layer                       │
│                  （データベース層）                       │
│                                                         │
│  責務: DB接続、トランザクション管理                       │
│  例: PostgreSQL, MySQL, MongoDB                         │
└─────────────────────────────────────────────────────────┘
```

### 1.2 変形パターン

実際のプロジェクトでは、層の数や名前が異なることがあります：

```
【3層アーキテクチャ】          【4層アーキテクチャ】

  Presentation                  Presentation
       ↓                             ↓
    Business                     Business
       ↓                             ↓
      Data                      Persistence
                                     ↓
                                  Database
```

---

## 2. 依存関係のルール

### 2.1 基本ルール

```
┌─────────────────────────────────────────────────────────┐
│                      ルール                              │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. 上位層は下位層に依存できる                            │
│  2. 下位層は上位層に依存してはならない                    │
│  3. 層をスキップして依存することは避ける（Strict Layering）│
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 2.2 図解

```
     Presentation
          │
          │ ✅ 依存OK
          ▼
       Business
          │
          │ ✅ 依存OK
          ▼
     Persistence
          │
          │ ✅ 依存OK
          ▼
       Database


     Presentation
          │
          │ ❌ 層スキップは避ける
          └──────────────────→ Persistence


       Business
          │
          │ ❌ 上位への依存は禁止
          ▼
     Presentation
```

### 2.3 Strict vs Relaxed Layering

```
【Strict Layering（厳格）】

  各層は直下の層のみに依存

  Presentation → Business → Persistence → Database
       │              │            │
       ×              ×            ×
       │              │            │
       └──────────────┴────────────┴──→ 直接アクセス禁止


【Relaxed Layering（緩和）】

  下位の層であればどこでもアクセス可能

  Presentation ─→ Business
       │              │
       └──────────────┴──→ Persistence → Database

  ※ パフォーマンス上の理由で緩和することがある
```

---

## 3. Goでの実装例

### 3.1 ディレクトリ構造

```
layered/
├── main.go
├── presentation/
│   └── handler/
│       └── user_handler.go
├── business/
│   └── service/
│       └── user_service.go
├── persistence/
│   └── repository/
│       └── user_repository.go
└── database/
    └── postgres.go
```

### 3.2 各層のコード例

```go
// ========================================
// database/postgres.go - Database Layer
// ========================================
package database

import (
    "database/sql"
    _ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
    var err error
    DB, err = sql.Open("postgres", "...")
    return err
}

// ========================================
// persistence/repository/user_repository.go - Persistence Layer
// ========================================
package repository

import (
    "app/database"
)

type User struct {
    ID    int
    Name  string
    Email string
}

type UserRepository struct{}

func (r *UserRepository) FindByID(id int) (*User, error) {
    row := database.DB.QueryRow(
        "SELECT id, name, email FROM users WHERE id = $1", id,
    )

    var user User
    err := row.Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) Save(user *User) error {
    _, err := database.DB.Exec(
        "INSERT INTO users (name, email) VALUES ($1, $2)",
        user.Name, user.Email,
    )
    return err
}

// ========================================
// business/service/user_service.go - Business Layer
// ========================================
package service

import (
    "errors"
    "app/persistence/repository"
)

type UserService struct {
    repo *repository.UserRepository  // 具象クラスに直接依存
}

func NewUserService() *UserService {
    return &UserService{
        repo: &repository.UserRepository{},
    }
}

func (s *UserService) GetUser(id int) (*repository.User, error) {
    return s.repo.FindByID(id)
}

func (s *UserService) CreateUser(name, email string) error {
    // ビジネスルール: メールは必須
    if email == "" {
        return errors.New("email is required")
    }

    user := &repository.User{
        Name:  name,
        Email: email,
    }
    return s.repo.Save(user)
}

// ========================================
// presentation/handler/user_handler.go - Presentation Layer
// ========================================
package handler

import (
    "encoding/json"
    "net/http"
    "strconv"
    "app/business/service"
)

type UserHandler struct {
    service *service.UserService
}

func NewUserHandler() *UserHandler {
    return &UserHandler{
        service: service.NewUserService(),
    }
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, _ := strconv.Atoi(idStr)

    user, err := h.service.GetUser(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(user)
}
```

### 3.3 依存関係の可視化

```
main.go
   │
   └──→ presentation/handler/user_handler.go
              │
              └──→ business/service/user_service.go
                        │
                        └──→ persistence/repository/user_repository.go
                                   │
                                   └──→ database/postgres.go
```

---

## 4. 長所と短所

### 4.1 長所

```
┌─────────────────────────────────────────────────────────┐
│  ✅ 長所                                                │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. 理解しやすい                                         │
│     水平方向の分割は直感的で、新人でも把握しやすい         │
│                                                         │
│  2. 実装が簡単                                           │
│     複雑なインターフェース設計が不要                      │
│                                                         │
│  3. 役割分担が明確                                       │
│     「UI担当」「ロジック担当」とチーム分割できる          │
│                                                         │
│  4. フレームワークとの相性が良い                          │
│     多くのWebフレームワークがこの構造を前提としている      │
│                                                         │
│  5. 学習コストが低い                                     │
│     すぐに開発を始められる                               │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 4.2 短所

```
┌─────────────────────────────────────────────────────────┐
│  ❌ 短所                                                │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. ビジネス層がインフラに依存                            │
│     Business → Persistence → Database                   │
│     DBを変えるとビジネス層も影響を受ける                  │
│                                                         │
│  2. テストが困難                                         │
│     ビジネスロジックのテストに実DBが必要になりがち        │
│                                                         │
│  3. DB中心設計になりやすい                               │
│     「まずテーブル設計」という発想になりがち              │
│                                                         │
│  4. 変更の影響範囲が大きい                               │
│     下位層の変更が上位層に波及する                        │
│                                                         │
│  5. ビジネスルールが分散しやすい                          │
│     複数のServiceにロジックが散らばる                    │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

## 5. 典型的な問題パターン

### 5.1 問題1: DBへの直接依存

```go
// ❌ 問題: UserServiceがPostgreSQLに直接依存
type UserService struct {
    db *sql.DB  // PostgreSQL固有
}

func (s *UserService) GetUser(id int) (*User, error) {
    // PostgreSQL固有のプレースホルダ $1
    row := s.db.QueryRow("SELECT * FROM users WHERE id = $1", id)
    // ...
}

// MySQLに変更したい場合:
// - SQLの書き換え（$1 → ?）
// - UserServiceのコード変更が必要
```

### 5.2 問題2: テストの困難さ

```go
// ❌ 問題: テストに実DBが必要
func TestUserService_GetUser(t *testing.T) {
    // テスト用DBのセットアップが必要
    setupTestDatabase()
    defer teardownTestDatabase()

    service := NewUserService()
    user, err := service.GetUser(1)

    // ...アサーション
}

// 問題点:
// - テストが遅い
// - 環境依存
// - CI/CDで動かすのが大変
```

### 5.3 問題3: 循環依存の誘発

```
// ❌ 問題: 層をまたいだ型の共有で依存が複雑化

persistence/
└── user.go  (User struct定義)
        ↑
business/
└── user_service.go  (Userを使用)
        ↑
presentation/
└── user_handler.go  (Userを使用)

// Userの変更が全層に影響
```

---

## 6. 改善パターン

### 6.1 インターフェースの導入（DIP適用）

```go
// ========================================
// business/repository.go - インターフェース定義
// ========================================
package business

// Business層でインターフェースを定義
type UserRepository interface {
    FindByID(id int) (*User, error)
    Save(user *User) error
}

type User struct {
    ID    int
    Name  string
    Email string
}

// ========================================
// business/service/user_service.go
// ========================================
package service

import "app/business"

type UserService struct {
    repo business.UserRepository  // インターフェースに依存
}

func NewUserService(repo business.UserRepository) *UserService {
    return &UserService{repo: repo}
}

// ========================================
// persistence/repository/user_repository.go
// ========================================
package repository

import "app/business"

// インターフェースを実装
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) FindByID(id int) (*business.User, error) {
    // 実装
}
```

**注意**: この改善を進めると、実質的にオニオン/クリーンアーキテクチャに近づきます。

---

## 7. いつ使うべきか

### 7.1 適しているケース

```
✅ レイヤードアーキテクチャが適している:

・小規模プロジェクト（数人、数ヶ月）
・シンプルなCRUDアプリケーション
・プロトタイプ、MVP
・チームがアーキテクチャに不慣れ
・外部依存が少ない（DBのみ等）
・変更頻度が低いシステム
```

### 7.2 適さないケース

```
❌ レイヤードアーキテクチャが適さない:

・大規模・長期プロジェクト
・複雑なビジネスロジック
・高いテストカバレッジが必要
・DBやフレームワークの変更可能性がある
・マイクロサービス化の可能性がある
```

---

## 8. 他のアーキテクチャとの違い

| 観点 | レイヤード | ヘキサゴナル/オニオン/クリーン |
|-----|-----------|------------------------------|
| 依存の方向 | 上→下 | 外→内 |
| DBの位置 | 最下層（基盤） | 最外層（詳細） |
| インターフェース | オプション | 必須 |
| テスタビリティ | 低〜中 | 高 |
| 複雑さ | 低 | 中〜高 |

---

## まとめ

レイヤードアーキテクチャは：

1. **シンプルで理解しやすい**最も基本的なアーキテクチャ
2. **小規模プロジェクトに最適**だが、スケールには限界がある
3. **問題点を認識**した上で使えば有効なツール
4. **進化の起点**として、必要に応じてオニオン/クリーンへ移行可能

---

次のドキュメント: [03_hexagonal_architecture.md](./03_hexagonal_architecture.md) - ヘキサゴナルアーキテクチャ
