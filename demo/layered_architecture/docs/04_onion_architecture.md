# オニオンアーキテクチャ (Onion Architecture)

## 概要

オニオンアーキテクチャは、2008年にJeffrey Palermoによって提唱されました。
アプリケーションを**同心円状の層**で構成し、**ドメインモデルを中心**に据えるアーキテクチャです。

ヘキサゴナルの思想を継承しつつ、**層の構造をより明確に定義**しました。

---

## 1. 基本構造

### 1.1 概念図

```
        ┌─────────────────────────────────────────────────┐
        │               Infrastructure                    │
        │         (UI, DB, External Services)             │
        │                                                 │
        │    ┌───────────────────────────────────────┐    │
        │    │         Application Services          │    │
        │    │           (ユースケース)               │    │
        │    │                                       │    │
        │    │    ┌───────────────────────────┐      │    │
        │    │    │      Domain Services      │      │    │
        │    │    │    (ドメインサービス)       │      │    │
        │    │    │                           │      │    │
        │    │    │    ┌───────────────┐      │      │    │
        │    │    │    │ Domain Model │      │      │    │
        │    │    │    │ (Entity/VO)  │      │      │    │
        │    │    │    └───────────────┘      │      │    │
        │    │    │                           │      │    │
        │    │    └───────────────────────────┘      │    │
        │    │                                       │    │
        │    └───────────────────────────────────────┘    │
        │                                                 │
        └─────────────────────────────────────────────────┘

        依存の方向: 外側 → 内側 （常に中心に向かう）
```

### 1.2 各層の説明

```
┌─────────────────────────────────────────────────────────┐
│  Layer 1: Domain Model（最内層）                         │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ・Entity（エンティティ）: 識別子を持つオブジェクト        │
│  ・Value Object（値オブジェクト）: 不変の値               │
│  ・ビジネスルールそのもの                                │
│                                                         │
│  特徴:                                                  │
│  ・一切の外部依存なし                                    │
│  ・純粋なビジネスロジック                                │
│  ・最も変更されにくい                                    │
│                                                         │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│  Layer 2: Domain Services                               │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ・複数のEntityをまたぐ操作                              │
│  ・Entityに属さないドメインロジック                       │
│  ・Repositoryインターフェースの定義                       │
│                                                         │
│  例:                                                    │
│  ・TransferService（口座間送金）                         │
│  ・PricingService（価格計算）                            │
│                                                         │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│  Layer 3: Application Services                          │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ・ユースケースの実装                                    │
│  ・ドメイン層を調整（オーケストレーション）               │
│  ・トランザクション境界の管理                            │
│                                                         │
│  例:                                                    │
│  ・CreateUserUseCase                                    │
│  ・PlaceOrderUseCase                                    │
│                                                         │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│  Layer 4: Infrastructure（最外層）                       │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ・UIコンポーネント（Web, CLI, GUI）                     │
│  ・データベースアクセス実装                               │
│  ・外部サービス連携                                      │
│  ・フレームワーク依存コード                              │
│                                                         │
│  特徴:                                                  │
│  ・技術的な詳細                                          │
│  ・交換可能                                              │
│  ・内側の層に依存する                                    │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

## 2. 依存関係のルール

### 2.1 基本ルール

```
┌─────────────────────────────────────────────────────────┐
│                   絶対的なルール                          │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  外側の層は内側の層に依存できる                            │
│  内側の層は外側の層に依存してはならない                    │
│                                                         │
│  Infrastructure → Application → Domain → Domain Model   │
│       ✅              ✅           ✅         ✅        │
│                                                         │
│  Domain Model → Domain → Application → Infrastructure   │
│       ❌          ❌           ❌            ❌          │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 2.2 インターフェースの配置

```
┌─────────────────────────────────────────────────────────┐
│       オニオンアーキテクチャの重要なポイント               │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  インターフェースは「使う側」で定義する                    │
│                                                         │
│  例: UserRepository                                     │
│                                                         │
│  ・Domain/Application層が「使う」                        │
│  ・Infrastructure層が「実装する」                        │
│                                                         │
│  → インターフェースはDomain層で定義                      │
│                                                         │
│       Domain層                Infrastructure層           │
│  ┌─────────────────┐      ┌─────────────────┐           │
│  │ UserRepository  │ ←─── │ PostgresUserRepo│           │
│  │  (interface)    │      │  (実装)         │           │
│  └─────────────────┘      └─────────────────┘           │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

## 3. Goでの実装例

### 3.1 ディレクトリ構造

```
onion/
├── main.go
├── domain/                           # Domain Model + Domain Services
│   ├── model/
│   │   ├── user.go                   # Entity
│   │   └── email.go                  # Value Object
│   ├── repository/
│   │   └── user_repository.go        # Repository Interface
│   └── service/
│       └── user_domain_service.go    # Domain Service
├── application/                      # Application Services
│   └── usecase/
│       ├── create_user.go
│       └── get_user.go
└── infrastructure/                   # Infrastructure
    ├── persistence/
    │   └── postgres_user_repository.go
    ├── web/
    │   └── handler/
    │       └── user_handler.go
    └── config/
        └── database.go
```

### 3.2 コード例

```go
// ========================================
// domain/model/email.go - Value Object
// ========================================
package model

import (
    "errors"
    "regexp"
)

// Email は値オブジェクト（不変）
type Email struct {
    value string
}

func NewEmail(value string) (Email, error) {
    if !isValidEmail(value) {
        return Email{}, errors.New("invalid email format")
    }
    return Email{value: value}, nil
}

func (e Email) Value() string {
    return e.value
}

func isValidEmail(email string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(email)
}

// ========================================
// domain/model/user.go - Entity
// ========================================
package model

import (
    "errors"
    "github.com/google/uuid"
)

// User はエンティティ（識別子を持つ）
type User struct {
    id    string
    name  string
    email Email
}

func NewUser(name string, email Email) (*User, error) {
    if name == "" {
        return nil, errors.New("name is required")
    }

    return &User{
        id:    uuid.New().String(),
        name:  name,
        email: email,
    }, nil
}

// Reconstruct はDBからの復元用
func ReconstructUser(id, name string, email Email) *User {
    return &User{
        id:    id,
        name:  name,
        email: email,
    }
}

func (u *User) ID() string    { return u.id }
func (u *User) Name() string  { return u.name }
func (u *User) Email() Email  { return u.email }

func (u *User) ChangeName(newName string) error {
    if newName == "" {
        return errors.New("name cannot be empty")
    }
    u.name = newName
    return nil
}

// ========================================
// domain/repository/user_repository.go - Repository Interface
// ========================================
package repository

import "app/domain/model"

// UserRepository はDomain層で定義されるインターフェース
// Infrastructure層で実装される
type UserRepository interface {
    Save(user *model.User) error
    FindByID(id string) (*model.User, error)
    FindByEmail(email model.Email) (*model.User, error)
    Exists(email model.Email) (bool, error)
}

// ========================================
// domain/service/user_domain_service.go - Domain Service
// ========================================
package service

import (
    "errors"
    "app/domain/model"
    "app/domain/repository"
)

// UserDomainService は複数のEntityをまたぐドメインロジック
type UserDomainService struct {
    userRepo repository.UserRepository
}

func NewUserDomainService(userRepo repository.UserRepository) *UserDomainService {
    return &UserDomainService{userRepo: userRepo}
}

// IsDuplicateEmail はメールアドレスの重複チェック
// Entity単体では判断できないビジネスルール
func (s *UserDomainService) IsDuplicateEmail(email model.Email) (bool, error) {
    return s.userRepo.Exists(email)
}

// ========================================
// application/usecase/create_user.go - Application Service
// ========================================
package usecase

import (
    "errors"
    "app/domain/model"
    "app/domain/repository"
    "app/domain/service"
)

type CreateUserInput struct {
    Name  string
    Email string
}

type CreateUserOutput struct {
    ID    string
    Name  string
    Email string
}

type CreateUserUseCase struct {
    userRepo          repository.UserRepository
    userDomainService *service.UserDomainService
}

func NewCreateUserUseCase(
    userRepo repository.UserRepository,
    userDomainService *service.UserDomainService,
) *CreateUserUseCase {
    return &CreateUserUseCase{
        userRepo:          userRepo,
        userDomainService: userDomainService,
    }
}

func (uc *CreateUserUseCase) Execute(input CreateUserInput) (*CreateUserOutput, error) {
    // 1. 値オブジェクトの生成（バリデーション含む）
    email, err := model.NewEmail(input.Email)
    if err != nil {
        return nil, err
    }

    // 2. ドメインサービスで重複チェック
    isDuplicate, err := uc.userDomainService.IsDuplicateEmail(email)
    if err != nil {
        return nil, err
    }
    if isDuplicate {
        return nil, errors.New("email already exists")
    }

    // 3. エンティティの生成
    user, err := model.NewUser(input.Name, email)
    if err != nil {
        return nil, err
    }

    // 4. 永続化
    if err := uc.userRepo.Save(user); err != nil {
        return nil, err
    }

    // 5. 出力を返す
    return &CreateUserOutput{
        ID:    user.ID(),
        Name:  user.Name(),
        Email: user.Email().Value(),
    }, nil
}

// ========================================
// infrastructure/persistence/postgres_user_repository.go
// ========================================
package persistence

import (
    "database/sql"
    "app/domain/model"
    "app/domain/repository"
)

// PostgresUserRepository はrepository.UserRepositoryを実装
type PostgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(user *model.User) error {
    _, err := r.db.Exec(
        "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)",
        user.ID(), user.Name(), user.Email().Value(),
    )
    return err
}

func (r *PostgresUserRepository) FindByID(id string) (*model.User, error) {
    row := r.db.QueryRow(
        "SELECT id, name, email FROM users WHERE id = $1", id,
    )

    var userID, name, emailStr string
    if err := row.Scan(&userID, &name, &emailStr); err != nil {
        return nil, err
    }

    email, _ := model.NewEmail(emailStr)
    return model.ReconstructUser(userID, name, email), nil
}

func (r *PostgresUserRepository) FindByEmail(email model.Email) (*model.User, error) {
    row := r.db.QueryRow(
        "SELECT id, name, email FROM users WHERE email = $1",
        email.Value(),
    )

    var userID, name, emailStr string
    if err := row.Scan(&userID, &name, &emailStr); err != nil {
        return nil, err
    }

    e, _ := model.NewEmail(emailStr)
    return model.ReconstructUser(userID, name, e), nil
}

func (r *PostgresUserRepository) Exists(email model.Email) (bool, error) {
    var count int
    err := r.db.QueryRow(
        "SELECT COUNT(*) FROM users WHERE email = $1",
        email.Value(),
    ).Scan(&count)
    return count > 0, err
}

// ========================================
// infrastructure/web/handler/user_handler.go
// ========================================
package handler

import (
    "encoding/json"
    "net/http"
    "app/application/usecase"
)

type UserHandler struct {
    createUserUseCase *usecase.CreateUserUseCase
}

func NewUserHandler(createUserUseCase *usecase.CreateUserUseCase) *UserHandler {
    return &UserHandler{createUserUseCase: createUserUseCase}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    output, err := h.createUserUseCase.Execute(usecase.CreateUserInput{
        Name:  req.Name,
        Email: req.Email,
    })
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(output)
}

// ========================================
// main.go - 組み立て
// ========================================
package main

import (
    "database/sql"
    "net/http"

    "app/application/usecase"
    domainService "app/domain/service"
    "app/infrastructure/persistence"
    "app/infrastructure/web/handler"
)

func main() {
    // Infrastructure
    db, _ := sql.Open("postgres", "...")

    // Repository（Infrastructure層の実装）
    userRepo := persistence.NewPostgresUserRepository(db)

    // Domain Service
    userDomainService := domainService.NewUserDomainService(userRepo)

    // Application Service（UseCase）
    createUserUseCase := usecase.NewCreateUserUseCase(userRepo, userDomainService)

    // Handler（Infrastructure層）
    userHandler := handler.NewUserHandler(createUserUseCase)

    // Router
    http.HandleFunc("/users", userHandler.CreateUser)
    http.ListenAndServe(":8080", nil)
}
```

---

## 4. 依存関係の可視化

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│     Infrastructure                    Infrastructure        │
│    (UserHandler)                    (PostgresUserRepo)      │
│          │                                 │                │
│          │ depends on                      │ implements     │
│          ▼                                 ▼                │
│    ┌───────────────────────────────────────────────┐        │
│    │              Application Layer                │        │
│    │            (CreateUserUseCase)                │        │
│    │                     │                         │        │
│    │                     │ depends on              │        │
│    │                     ▼                         │        │
│    │    ┌───────────────────────────────────┐     │        │
│    │    │          Domain Layer             │     │        │
│    │    │  (UserDomainService, Repository)  │     │        │
│    │    │               │                   │     │        │
│    │    │               │ depends on        │     │        │
│    │    │               ▼                   │     │        │
│    │    │    ┌───────────────────────┐      │     │        │
│    │    │    │     Domain Model      │      │     │        │
│    │    │    │    (User, Email)      │      │     │        │
│    │    │    └───────────────────────┘      │     │        │
│    │    └───────────────────────────────────┘     │        │
│    └───────────────────────────────────────────────┘        │
│                                                             │
└─────────────────────────────────────────────────────────────┘

矢印の方向 = 依存の方向（すべて内側を向いている）
```

---

## 5. 長所と短所

### 5.1 長所

```
┌─────────────────────────────────────────────────────────┐
│  ✅ 長所                                                │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. ドメイン中心の設計                                   │
│     ビジネスロジックが最も保護される                      │
│     → ビジネスの変化に強い                               │
│                                                         │
│  2. テスタビリティ                                       │
│     Domain層は純粋なロジック                             │
│     → モックなしでテスト可能                             │
│                                                         │
│  3. 技術からの独立                                       │
│     DB、フレームワークの変更が容易                        │
│     → 長期運用に適している                               │
│                                                         │
│  4. DDDとの親和性                                        │
│     Entity、Value Object、Domain Serviceが自然に配置    │
│     → DDDの実践に最適                                   │
│                                                         │
│  5. 明確な層構造                                         │
│     ヘキサゴナルより層の役割が明確                        │
│     → チームでの共通理解が得やすい                       │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 5.2 短所

```
┌─────────────────────────────────────────────────────────┐
│  ❌ 短所                                                │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. 学習コスト                                           │
│     DDD、依存関係逆転の理解が必要                        │
│     → 初心者には難しい                                   │
│                                                         │
│  2. コード量の増加                                       │
│     インターフェース、DTO、ValueObjectが増える           │
│     → 小規模プロジェクトには過剰                        │
│                                                         │
│  3. Domain Services層の曖昧さ                           │
│     Application Servicesとの境界が不明確になることがある │
│     → チーム内で認識を合わせる必要                       │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

## 6. ヘキサゴナルとの違い

```
┌─────────────────────────────────────────────────────────┐
│              ヘキサゴナル vs オニオン                     │
├────────────────────────┬────────────────────────────────┤
│     ヘキサゴナル        │         オニオン               │
├────────────────────────┼────────────────────────────────┤
│                        │                                │
│  Port/Adapterという    │  同心円の層という               │
│  概念で説明            │  概念で説明                    │
│                        │                                │
│  Core内部の構造は      │  Domain Model, Domain Service,│
│  未定義                │  Application Serviceを明確化   │
│                        │                                │
│  入力/出力の対称性を   │  ドメインモデルの中心性を      │
│  強調                  │  強調                          │
│                        │                                │
│  DDDとの関連は         │  DDDとの親和性を               │
│  明示されていない      │  明示している                  │
│                        │                                │
└────────────────────────┴────────────────────────────────┘
```

**結論**: オニオンは、ヘキサゴナルの「Core」の内部構造を詳細化したもの

---

## 7. クリーンアーキテクチャとの違い

次のドキュメントで詳しく解説しますが、簡単に：

```
┌─────────────────────────────────────────────────────────┐
│                オニオン vs クリーン                      │
├────────────────────────┬────────────────────────────────┤
│       オニオン          │         クリーン               │
├────────────────────────┼────────────────────────────────┤
│                        │                                │
│  Domain Services層     │  Use Cases層                   │
│  Application Services層│  (より明確な責務)              │
│  という名前            │                                │
│                        │                                │
│  実装寄りの指針        │  より抽象的・理論的            │
│                        │                                │
│  Jeffrey Palermo       │  Robert C. Martin              │
│  (2008)                │  (2012)                        │
│                        │                                │
└────────────────────────┴────────────────────────────────┘
```

**本質的には非常に似ている**。クリーンはオニオンを含む複数のアーキテクチャを統合したもの。

---

## まとめ

オニオンアーキテクチャは：

1. **ドメインモデルを中心**に据えた同心円構造
2. **依存は常に内側へ**向かう
3. **DDDとの親和性**が高い
4. ヘキサゴナルの「Core」を**詳細化**したもの
5. **中〜大規模プロジェクト**に適している

---

次のドキュメント: [05_clean_architecture.md](./05_clean_architecture.md) - クリーンアーキテクチャ
