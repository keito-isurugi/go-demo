# ヘキサゴナルアーキテクチャ (Hexagonal Architecture)

別名: **Ports and Adapters パターン**

## 概要

ヘキサゴナルアーキテクチャは、2005年にAlistair Cockburnによって提唱されました。
アプリケーションを**外部から隔離**し、**ポート（インターフェース）とアダプター（実装）**を通じて外部とやり取りするパターンです。

オニオン/クリーンアーキテクチャの**直接の先祖**であり、「Port/Adapter」という重要な概念を導入しました。

---

## 1. 基本構造

### 1.1 なぜ「六角形」か

```
六角形に特別な意味はない。

Alistair Cockburn曰く:
「四角形だと4面しか描けず、ポートの数が限定される印象を与える。
 六角形なら複数のポートを描きやすい」

つまり、六角形は「複数のポート」を視覚的に表現するための選択。
```

### 1.2 概念図

```
                    ┌─────────────────┐
                    │   Web UI        │
                    │   (Adapter)     │
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │   HTTP Port     │
                    │   (Input)       │
           ┌────────┴─────────────────┴────────┐
           │                                   │
    ┌──────┴──────┐                     ┌──────┴──────┐
    │ CLI Adapter │                     │ REST API    │
    │  (Input)    │                     │  (Input)    │
    └──────┬──────┘                     └──────┬──────┘
           │                                   │
           │      ┌───────────────────┐        │
           │      │                   │        │
           └──────►   Application     ◄────────┘
                  │      Core         │
           ┌──────►   (ビジネス       ◄────────┐
           │      │    ロジック)      │        │
           │      │                   │        │
           │      └───────────────────┘        │
           │                                   │
    ┌──────┴──────┐                     ┌──────┴──────┐
    │ DB Port     │                     │ External    │
    │  (Output)   │                     │ API Port    │
    └──────┬──────┘                     └──────┬──────┘
           │                                   │
    ┌──────┴──────┐                     ┌──────┴──────┐
    │ PostgreSQL  │                     │ Stripe      │
    │  Adapter    │                     │  Adapter    │
    └─────────────┘                     └─────────────┘
```

### 1.3 シンプル化した図

```
           外部世界
               │
        ┌──────▼──────┐
        │   Adapter   │  ← 外部をポートに変換
        │  (実装)     │
        └──────┬──────┘
               │
        ┌──────▼──────┐
        │    Port     │  ← インターフェース（契約）
        │ (Interface) │
        └──────┬──────┘
               │
        ┌──────▼──────┐
        │ Application │  ← ビジネスロジック
        │    Core     │     Portのみに依存
        └─────────────┘
```

---

## 2. Ports と Adapters

### 2.1 Port（ポート）とは

```
┌─────────────────────────────────────────────────────────┐
│  Port = インターフェース（契約）                          │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  アプリケーションコアが外部と通信するための「窓口」        │
│                                                         │
│  ・技術詳細を知らない                                    │
│  ・「何ができるか」だけを定義                             │
│  ・「どうやるか」は知らない                               │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

#### Input Port（駆動側ポート）

```go
// アプリケーションを「使う」側のインターフェース
// 外部 → アプリケーション

type UserService interface {
    CreateUser(name, email string) (*User, error)
    GetUser(id string) (*User, error)
}
```

#### Output Port（被駆動側ポート）

```go
// アプリケーションが「使う」側のインターフェース
// アプリケーション → 外部

type UserRepository interface {
    Save(user *User) error
    FindByID(id string) (*User, error)
}

type EmailSender interface {
    Send(to, subject, body string) error
}
```

### 2.2 Adapter（アダプター）とは

```
┌─────────────────────────────────────────────────────────┐
│  Adapter = Portの実装                                   │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  外部の技術詳細とPortを橋渡しする                         │
│                                                         │
│  ・具体的な技術を知っている                               │
│  ・Portインターフェースを実装する                         │
│  ・交換可能（差し替え可能）                               │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

#### Input Adapter（Primary Adapter）

```go
// HTTPハンドラー（Web UI用のAdapter）
type HTTPUserHandler struct {
    service UserService  // Input Portを使う
}

func (h *HTTPUserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // HTTPリクエストをサービス呼び出しに変換
    var req CreateUserRequest
    json.NewDecoder(r.Body).Decode(&req)

    user, err := h.service.CreateUser(req.Name, req.Email)
    // ...
}

// CLIコマンド（CLI用のAdapter）
type CLIUserCommand struct {
    service UserService  // 同じInput Portを使う
}

func (c *CLIUserCommand) Run(args []string) {
    user, err := c.service.CreateUser(args[0], args[1])
    // ...
}
```

#### Output Adapter（Secondary Adapter）

```go
// PostgreSQL用のAdapter
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) Save(user *User) error {
    _, err := r.db.Exec(
        "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)",
        user.ID, user.Name, user.Email,
    )
    return err
}

// MongoDB用のAdapter（差し替え可能）
type MongoUserRepository struct {
    collection *mongo.Collection
}

func (r *MongoUserRepository) Save(user *User) error {
    _, err := r.collection.InsertOne(context.Background(), user)
    return err
}
```

---

## 3. Goでの実装例

### 3.1 ディレクトリ構造

```
hexagonal/
├── main.go
├── core/                          # Application Core
│   ├── domain/
│   │   └── user.go                # エンティティ
│   ├── port/
│   │   ├── input.go               # Input Ports
│   │   └── output.go              # Output Ports
│   └── service/
│       └── user_service.go        # ビジネスロジック
└── adapter/
    ├── input/                     # Input Adapters
    │   ├── http/
    │   │   └── user_handler.go
    │   └── cli/
    │       └── user_command.go
    └── output/                    # Output Adapters
        ├── persistence/
        │   ├── postgres_user_repo.go
        │   └── memory_user_repo.go
        └── notification/
            └── smtp_email_sender.go
```

### 3.2 コード例

```go
// ========================================
// core/domain/user.go - ドメインモデル
// ========================================
package domain

import "errors"

type User struct {
    ID    string
    Name  string
    Email string
}

func NewUser(name, email string) (*User, error) {
    if name == "" {
        return nil, errors.New("name is required")
    }
    if email == "" {
        return nil, errors.New("email is required")
    }
    return &User{
        ID:    generateID(),
        Name:  name,
        Email: email,
    }, nil
}

// ========================================
// core/port/input.go - Input Ports
// ========================================
package port

import "app/core/domain"

// アプリケーションを使うためのインターフェース
type UserService interface {
    CreateUser(name, email string) (*domain.User, error)
    GetUser(id string) (*domain.User, error)
}

// ========================================
// core/port/output.go - Output Ports
// ========================================
package port

import "app/core/domain"

// アプリケーションが使うインターフェース
type UserRepository interface {
    Save(user *domain.User) error
    FindByID(id string) (*domain.User, error)
}

type EmailSender interface {
    Send(to, subject, body string) error
}

// ========================================
// core/service/user_service.go - ビジネスロジック
// ========================================
package service

import (
    "app/core/domain"
    "app/core/port"
)

// UserServiceImpl は port.UserService を実装
type UserServiceImpl struct {
    repo   port.UserRepository  // Output Portに依存
    mailer port.EmailSender     // Output Portに依存
}

func NewUserService(repo port.UserRepository, mailer port.EmailSender) *UserServiceImpl {
    return &UserServiceImpl{
        repo:   repo,
        mailer: mailer,
    }
}

func (s *UserServiceImpl) CreateUser(name, email string) (*domain.User, error) {
    // ドメインロジック
    user, err := domain.NewUser(name, email)
    if err != nil {
        return nil, err
    }

    // Output Port経由でDB保存
    if err := s.repo.Save(user); err != nil {
        return nil, err
    }

    // Output Port経由でメール送信
    s.mailer.Send(email, "Welcome!", "ようこそ！")

    return user, nil
}

func (s *UserServiceImpl) GetUser(id string) (*domain.User, error) {
    return s.repo.FindByID(id)
}

// ========================================
// adapter/input/http/user_handler.go - Input Adapter
// ========================================
package http

import (
    "encoding/json"
    "net/http"
    "app/core/port"
)

type UserHandler struct {
    service port.UserService  // Input Portに依存
}

func NewUserHandler(service port.UserService) *UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    user, err := h.service.CreateUser(req.Name, req.Email)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

// ========================================
// adapter/output/persistence/postgres_user_repo.go - Output Adapter
// ========================================
package persistence

import (
    "database/sql"
    "app/core/domain"
)

type PostgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(user *domain.User) error {
    _, err := r.db.Exec(
        "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)",
        user.ID, user.Name, user.Email,
    )
    return err
}

func (r *PostgresUserRepository) FindByID(id string) (*domain.User, error) {
    row := r.db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id)
    var user domain.User
    err := row.Scan(&user.ID, &user.Name, &user.Email)
    return &user, err
}

// ========================================
// adapter/output/persistence/memory_user_repo.go - テスト用Adapter
// ========================================
package persistence

import (
    "errors"
    "app/core/domain"
)

type InMemoryUserRepository struct {
    users map[string]*domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{
        users: make(map[string]*domain.User),
    }
}

func (r *InMemoryUserRepository) Save(user *domain.User) error {
    r.users[user.ID] = user
    return nil
}

func (r *InMemoryUserRepository) FindByID(id string) (*domain.User, error) {
    if user, ok := r.users[id]; ok {
        return user, nil
    }
    return nil, errors.New("user not found")
}

// ========================================
// main.go - 組み立て
// ========================================
package main

import (
    "database/sql"
    "net/http"

    "app/core/service"
    httpAdapter "app/adapter/input/http"
    "app/adapter/output/persistence"
    "app/adapter/output/notification"
)

func main() {
    // Output Adaptersの作成
    db, _ := sql.Open("postgres", "...")
    userRepo := persistence.NewPostgresUserRepository(db)
    mailer := notification.NewSMTPEmailSender("smtp.example.com")

    // Application Coreの作成
    userService := service.NewUserService(userRepo, mailer)

    // Input Adaptersの作成
    userHandler := httpAdapter.NewUserHandler(userService)

    // ルーティング
    http.HandleFunc("/users", userHandler.CreateUser)
    http.ListenAndServe(":8080", nil)
}
```

---

## 4. 依存関係の可視化

```
┌─────────────────────────────────────────────────────────────┐
│                       main.go                               │
│                     (組み立て)                               │
└──────────────────────────┬──────────────────────────────────┘
                           │ 生成・注入
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ UserHandler  │  │ UserService  │  │ PostgresRepo │
│(Input Adapter)│  │   (Core)     │  │(Output Adapter)│
└───────┬──────┘  └───────┬──────┘  └───────┬──────┘
        │                 │                 │
        │    uses         │    uses         │
        ▼                 ▼                 ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ UserService  │  │UserRepository│  │UserRepository│
│ (Input Port) │  │(Output Port) │  │(Output Port) │
│  interface   │  │  interface   │  │  interface   │
└──────────────┘  └──────────────┘  └──────────────┘


依存の方向:
  Adapter → Port ← Core

  ・AdapterはPortに依存
  ・CoreもPortに依存
  ・CoreはAdapterを知らない
```

---

## 5. 長所と短所

### 5.1 長所

```
┌─────────────────────────────────────────────────────────┐
│  ✅ 長所                                                │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. 技術からの独立                                       │
│     CoreはDB、Web、外部APIを知らない                     │
│     → 技術変更の影響を受けない                           │
│                                                         │
│  2. テスタビリティ                                       │
│     Adapterを差し替えてテスト可能                        │
│     → InMemoryRepositoryでDBなしテスト                  │
│                                                         │
│  3. 対称性                                               │
│     Input/Output両方に同じパターン適用                   │
│     → 一貫した設計                                       │
│                                                         │
│  4. プラグイン可能                                       │
│     新しいAdapter追加が容易                              │
│     → HTTP + CLI + gRPC を同時サポート                  │
│                                                         │
│  5. 明確な境界                                           │
│     Port = 契約が明示される                              │
│     → チーム間の協業がしやすい                           │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 5.2 短所

```
┌─────────────────────────────────────────────────────────┐
│  ❌ 短所                                                │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. 初期コストが高い                                     │
│     インターフェース定義が必要                           │
│     → 小規模プロジェクトには過剰                        │
│                                                         │
│  2. 間接化のオーバーヘッド                               │
│     すべてがインターフェース経由                         │
│     → コードの追跡が少し面倒                            │
│                                                         │
│  3. 層の定義が曖昧                                       │
│     「Core」の内部構造は未定義                           │
│     → オニオン/クリーンで補完される                      │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

## 6. レイヤードアーキテクチャとの違い

```
┌─────────────────────────────────────────────────────────┐
│                レイヤードアーキテクチャ                   │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Presentation → Business → Persistence → Database       │
│                                                         │
│  ・依存は上から下                                        │
│  ・Business層がPersistence層の具象に依存                │
│  ・DBは「基盤」という位置づけ                            │
│                                                         │
└─────────────────────────────────────────────────────────┘

                         ↓ 違い

┌─────────────────────────────────────────────────────────┐
│                ヘキサゴナルアーキテクチャ                 │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Adapter → Port ← Core → Port ← Adapter                │
│                                                         │
│  ・依存は外から内へ                                      │
│  ・CoreはPortインターフェースのみに依存                  │
│  ・DBは「外部詳細」という位置づけ                        │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 図解比較

```
【レイヤード】             【ヘキサゴナル】

    UI                     UI (Adapter)
     │                          │
     ▼                          ▼
  Business               Port (Interface)
     │                          │
     ▼                          ▼
    DB ←── 基盤              Core
                                │
                                ▼
                         Port (Interface)
                                │
                                ▼
                          DB (Adapter) ←── 外部詳細
```

---

## 7. オニオン/クリーンとの関係

```
ヘキサゴナル (2005)
      │
      │ 「Port/Adapter」という概念を導入
      │
      ├────────────────────────────────┐
      │                                │
      ▼                                ▼
オニオン (2008)                  クリーン (2012)
      │                                │
      │ 同心円モデルで                  │ Port/Adapterを
      │ 層を明確化                      │ Input/Outputに分類
      │                                │
      └────────────────────────────────┘
                     │
                     ▼
            本質的には同じ思想:
            「ビジネスロジックを外部から守る」
```

---

## まとめ

ヘキサゴナルアーキテクチャは：

1. **Port/Adapter**という重要な概念を導入した
2. アプリケーションを**外部技術から隔離**する
3. オニオン/クリーンの**直接の先祖**
4. 「依存は外から内へ」という原則の**起源**

---

次のドキュメント: [04_onion_architecture.md](./04_onion_architecture.md) - オニオンアーキテクチャ
