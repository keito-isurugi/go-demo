# クリーンアーキテクチャ (Clean Architecture)

## 概要

クリーンアーキテクチャは、2012年にRobert C. Martin（Uncle Bob）によって提唱されました。
ヘキサゴナル、オニオン、その他のアーキテクチャを**統合・一般化**し、より**明確な原則と層定義**を提供します。

---

## 1. 基本構造

### 1.1 概念図（有名な同心円図）

```
        ┌────────────────────────────────────────────────────────┐
        │                Frameworks & Drivers                    │
        │              (Web, UI, DB, Devices)                    │
        │                                                        │
        │    ┌────────────────────────────────────────────────┐  │
        │    │             Interface Adapters                 │  │
        │    │      (Controllers, Gateways, Presenters)       │  │
        │    │                                                │  │
        │    │    ┌────────────────────────────────────────┐  │  │
        │    │    │           Application                  │  │  │
        │    │    │           Business Rules               │  │  │
        │    │    │            (Use Cases)                 │  │  │
        │    │    │                                        │  │  │
        │    │    │    ┌────────────────────────────────┐  │  │  │
        │    │    │    │         Enterprise             │  │  │  │
        │    │    │    │       Business Rules           │  │  │  │
        │    │    │    │         (Entities)             │  │  │  │
        │    │    │    └────────────────────────────────┘  │  │  │
        │    │    │                                        │  │  │
        │    │    └────────────────────────────────────────┘  │  │
        │    │                                                │  │
        │    └────────────────────────────────────────────────┘  │
        │                                                        │
        └────────────────────────────────────────────────────────┘

        依存の方向: 外側 → 内側 （The Dependency Rule）
```

### 1.2 4つの層

```
┌─────────────────────────────────────────────────────────────┐
│  Layer 1: Entities（エンティティ）                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ・Enterprise Business Rules（企業全体のビジネスルール）     │
│  ・複数のアプリケーションで共有可能                          │
│  ・最も汎用的で、変更されにくい                              │
│                                                             │
│  例:                                                        │
│  ・User（ユーザーのビジネスルール）                          │
│  ・Order（注文のビジネスルール）                             │
│  ・Money（金額の計算ルール）                                 │
│                                                             │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  Layer 2: Use Cases（ユースケース）                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ・Application Business Rules（アプリ固有のビジネスルール）   │
│  ・システムの振る舞いを定義                                  │
│  ・Entitiesへの・からのデータフローを調整                    │
│                                                             │
│  例:                                                        │
│  ・CreateUserUseCase（ユーザー作成）                         │
│  ・PlaceOrderUseCase（注文実行）                             │
│  ・TransferMoneyUseCase（送金）                              │
│                                                             │
│  特徴:                                                      │
│  ・Input Port（入力境界）とOutput Port（出力境界）を定義    │
│  ・Interactor（相互作用子）がロジックを実装                 │
│                                                             │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  Layer 3: Interface Adapters（インターフェースアダプター）    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ・外部と内部のデータ形式を変換                              │
│  ・Use Casesが必要とする形式 ↔ 外部の形式                   │
│                                                             │
│  構成要素:                                                  │
│  ・Controller: 入力を受けてUse Caseを呼ぶ                   │
│  ・Presenter: Use Caseの出力を表示形式に変換                │
│  ・Gateway: 外部サービスとの通信を抽象化                     │
│  ・Repository実装: DBアクセスの実装                         │
│                                                             │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  Layer 4: Frameworks & Drivers（フレームワークとドライバ）    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ・最も外側の層                                              │
│  ・フレームワーク、ライブラリ、DBなど                        │
│  ・「糊」のようなコード                                      │
│                                                             │
│  例:                                                        │
│  ・Webサーバー（Echo, Gin, net/http）                       │
│  ・データベースドライバ（PostgreSQL, MySQL）                 │
│  ・外部API クライアント                                      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 2. The Dependency Rule（依存ルール）

### 2.1 絶対的なルール

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│   "ソースコードの依存は、内側に向かってのみ指すことができる"  │
│                                                             │
│   "Nothing in an inner circle can know anything at all      │
│    about something in an outer circle."                     │
│                                                             │
│   内側の円にあるものは、外側の円にあるものについて           │
│   一切知ってはならない                                       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 具体例

```
┌─────────────────────────────────────────────────────────────┐
│  ✅ OK（外→内への依存）                                     │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Controller → UseCase                                       │
│  Gateway → UseCase                                          │
│  UseCase → Entity                                           │
│  Presenter → UseCase                                        │
│                                                             │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  ❌ NG（内→外への依存）                                     │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Entity → Repository実装                                    │
│  UseCase → Controller                                       │
│  UseCase → Database                                         │
│  Entity → Framework                                         │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. Input/Output Ports パターン

クリーンアーキテクチャの特徴的な概念です。

### 3.1 概念図

```
                    Controller
                        │
                        │ calls
                        ▼
              ┌─────────────────┐
              │   Input Port    │  ← インターフェース
              │   (Boundary)    │
              └────────┬────────┘
                       │
                       │ implements
                       ▼
              ┌─────────────────┐
              │    Interactor   │  ← Use Case の実装
              │   (Use Case)    │
              └────────┬────────┘
                       │
                       │ uses
                       ▼
              ┌─────────────────┐
              │  Output Port    │  ← インターフェース
              │   (Boundary)    │
              └────────┬────────┘
                       │
                       │ implements
                       ▼
              ┌─────────────────┐
              │    Presenter    │  or  Repository
              │    (Gateway)    │
              └─────────────────┘
```

### 3.2 各コンポーネントの役割

```
┌─────────────────────────────────────────────────────────────┐
│  Input Port（入力境界）                                     │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ・Use Caseを呼び出すためのインターフェース                  │
│  ・外部（Controller等）が使う                                │
│  ・Use Case層で定義                                         │
│                                                             │
│  type CreateUserInputPort interface {                       │
│      Execute(input CreateUserInput) (CreateUserOutput, error)│
│  }                                                          │
│                                                             │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  Interactor（相互作用子）                                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ・Use Caseの具体的な実装                                    │
│  ・Input Portを実装する                                     │
│  ・Output Portを使ってデータを取得・保存                    │
│                                                             │
│  type CreateUserInteractor struct {                         │
│      userRepo UserRepository  // Output Port                │
│  }                                                          │
│                                                             │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  Output Port（出力境界）                                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ・Use Caseが外部リソースを使うためのインターフェース        │
│  ・Repository、Gateway、Presenterなど                       │
│  ・Use Case層で定義、Adapter層で実装                        │
│                                                             │
│  type UserRepository interface {                            │
│      Save(user *entity.User) error                          │
│      FindByID(id string) (*entity.User, error)              │
│  }                                                          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 4. Goでの実装例

### 4.1 ディレクトリ構造

```
clean/
├── main.go
├── domain/                              # Entities
│   └── entity/
│       └── user.go
├── usecase/                             # Use Cases
│   ├── port/
│   │   ├── input/
│   │   │   └── user_input.go            # Input Ports
│   │   └── output/
│   │       └── user_repository.go       # Output Ports
│   ├── interactor/
│   │   └── user_interactor.go           # Interactors
│   └── dto/
│       └── user_dto.go                  # Input/Output DTOs
├── adapter/                             # Interface Adapters
│   ├── controller/
│   │   └── user_controller.go
│   ├── gateway/
│   │   └── user_repository_impl.go
│   └── presenter/
│       └── user_presenter.go
└── infrastructure/                      # Frameworks & Drivers
    ├── web/
    │   └── router.go
    └── database/
        └── postgres.go
```

### 4.2 コード例

```go
// ========================================
// domain/entity/user.go - Entity
// ========================================
package entity

import (
    "errors"
    "time"
)

// User は Enterprise Business Rules を表す
type User struct {
    ID        string
    Name      string
    Email     string
    CreatedAt time.Time
}

func NewUser(id, name, email string) (*User, error) {
    if name == "" {
        return nil, errors.New("name is required")
    }
    if email == "" {
        return nil, errors.New("email is required")
    }

    return &User{
        ID:        id,
        Name:      name,
        Email:     email,
        CreatedAt: time.Now(),
    }, nil
}

// ビジネスルール: 名前を変更できる
func (u *User) ChangeName(name string) error {
    if name == "" {
        return errors.New("name cannot be empty")
    }
    u.Name = name
    return nil
}

// ========================================
// usecase/dto/user_dto.go - DTOs
// ========================================
package dto

// Input DTOs（リクエスト）
type CreateUserInput struct {
    Name  string
    Email string
}

type GetUserInput struct {
    ID string
}

// Output DTOs（レスポンス）
type UserOutput struct {
    ID        string
    Name      string
    Email     string
    CreatedAt string
}

// ========================================
// usecase/port/input/user_input.go - Input Ports
// ========================================
package input

import "app/usecase/dto"

// CreateUserInputPort は CreateUser ユースケースの入力境界
type CreateUserInputPort interface {
    Execute(input dto.CreateUserInput) (*dto.UserOutput, error)
}

// GetUserInputPort は GetUser ユースケースの入力境界
type GetUserInputPort interface {
    Execute(input dto.GetUserInput) (*dto.UserOutput, error)
}

// ========================================
// usecase/port/output/user_repository.go - Output Ports
// ========================================
package output

import "app/domain/entity"

// UserRepository は Output Port（永続化）
type UserRepository interface {
    Save(user *entity.User) error
    FindByID(id string) (*entity.User, error)
    FindByEmail(email string) (*entity.User, error)
    NextID() string
}

// UserPresenter は Output Port（表示）
type UserPresenter interface {
    Output(user *entity.User) *dto.UserOutput
}

// ========================================
// usecase/interactor/user_interactor.go - Interactors
// ========================================
package interactor

import (
    "errors"
    "app/domain/entity"
    "app/usecase/dto"
    "app/usecase/port/output"
)

// CreateUserInteractor は CreateUser ユースケースの実装
type CreateUserInteractor struct {
    userRepo  output.UserRepository
    presenter output.UserPresenter
}

func NewCreateUserInteractor(
    userRepo output.UserRepository,
    presenter output.UserPresenter,
) *CreateUserInteractor {
    return &CreateUserInteractor{
        userRepo:  userRepo,
        presenter: presenter,
    }
}

// Execute は input.CreateUserInputPort を実装
func (i *CreateUserInteractor) Execute(input dto.CreateUserInput) (*dto.UserOutput, error) {
    // 1. 重複チェック
    existing, _ := i.userRepo.FindByEmail(input.Email)
    if existing != nil {
        return nil, errors.New("email already exists")
    }

    // 2. Entity生成
    id := i.userRepo.NextID()
    user, err := entity.NewUser(id, input.Name, input.Email)
    if err != nil {
        return nil, err
    }

    // 3. 永続化
    if err := i.userRepo.Save(user); err != nil {
        return nil, err
    }

    // 4. 出力（Presenterを通して変換）
    return i.presenter.Output(user), nil
}

// GetUserInteractor は GetUser ユースケースの実装
type GetUserInteractor struct {
    userRepo  output.UserRepository
    presenter output.UserPresenter
}

func NewGetUserInteractor(
    userRepo output.UserRepository,
    presenter output.UserPresenter,
) *GetUserInteractor {
    return &GetUserInteractor{
        userRepo:  userRepo,
        presenter: presenter,
    }
}

func (i *GetUserInteractor) Execute(input dto.GetUserInput) (*dto.UserOutput, error) {
    user, err := i.userRepo.FindByID(input.ID)
    if err != nil {
        return nil, err
    }
    return i.presenter.Output(user), nil
}

// ========================================
// adapter/presenter/user_presenter.go - Presenter
// ========================================
package presenter

import (
    "app/domain/entity"
    "app/usecase/dto"
)

type UserPresenter struct{}

func NewUserPresenter() *UserPresenter {
    return &UserPresenter{}
}

// Output は output.UserPresenter を実装
func (p *UserPresenter) Output(user *entity.User) *dto.UserOutput {
    return &dto.UserOutput{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
    }
}

// ========================================
// adapter/gateway/user_repository_impl.go - Repository実装
// ========================================
package gateway

import (
    "database/sql"
    "github.com/google/uuid"
    "app/domain/entity"
)

type PostgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(user *entity.User) error {
    _, err := r.db.Exec(
        "INSERT INTO users (id, name, email, created_at) VALUES ($1, $2, $3, $4)",
        user.ID, user.Name, user.Email, user.CreatedAt,
    )
    return err
}

func (r *PostgresUserRepository) FindByID(id string) (*entity.User, error) {
    row := r.db.QueryRow(
        "SELECT id, name, email, created_at FROM users WHERE id = $1", id,
    )

    var user entity.User
    err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *PostgresUserRepository) FindByEmail(email string) (*entity.User, error) {
    row := r.db.QueryRow(
        "SELECT id, name, email, created_at FROM users WHERE email = $1", email,
    )

    var user entity.User
    err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *PostgresUserRepository) NextID() string {
    return uuid.New().String()
}

// ========================================
// adapter/controller/user_controller.go - Controller
// ========================================
package controller

import (
    "encoding/json"
    "net/http"
    "app/usecase/dto"
    "app/usecase/port/input"
)

type UserController struct {
    createUser input.CreateUserInputPort
    getUser    input.GetUserInputPort
}

func NewUserController(
    createUser input.CreateUserInputPort,
    getUser input.GetUserInputPort,
) *UserController {
    return &UserController{
        createUser: createUser,
        getUser:    getUser,
    }
}

func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    // DTOに変換してUse Caseを呼ぶ
    output, err := c.createUser.Execute(dto.CreateUserInput{
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

func (c *UserController) Get(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")

    output, err := c.getUser.Execute(dto.GetUserInput{ID: id})
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(output)
}

// ========================================
// main.go - 組み立て（Dependency Injection）
// ========================================
package main

import (
    "database/sql"
    "net/http"

    "app/adapter/controller"
    "app/adapter/gateway"
    "app/adapter/presenter"
    "app/usecase/interactor"
)

func main() {
    // Infrastructure
    db, _ := sql.Open("postgres", "...")

    // Adapters（Output Ports の実装）
    userRepo := gateway.NewPostgresUserRepository(db)
    userPresenter := presenter.NewUserPresenter()

    // Use Cases（Interactors）
    createUserInteractor := interactor.NewCreateUserInteractor(userRepo, userPresenter)
    getUserInteractor := interactor.NewGetUserInteractor(userRepo, userPresenter)

    // Adapters（Input を受ける）
    userController := controller.NewUserController(createUserInteractor, getUserInteractor)

    // Router
    http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodPost:
            userController.Create(w, r)
        case http.MethodGet:
            userController.Get(w, r)
        }
    })

    http.ListenAndServe(":8080", nil)
}
```

---

## 5. データの流れ

### 5.1 リクエスト〜レスポンスの流れ

```
HTTP Request
     │
     ▼
┌─────────────┐
│  Controller │  1. HTTPリクエストを受け取る
│  (Adapter)  │  2. Input DTOに変換
└──────┬──────┘
       │ Input DTO
       ▼
┌─────────────┐
│ Input Port  │  インターフェース
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Interactor  │  3. ビジネスロジック実行
│ (Use Case)  │  4. Output Port経由でDB操作
└──────┬──────┘
       │
       ▼
┌─────────────┐
│Output Port  │  インターフェース（Repository等）
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Gateway    │  5. 実際のDB操作
│ (Adapter)   │
└──────┬──────┘
       │ Entity
       ▼
┌─────────────┐
│ Interactor  │  6. Entity を受け取る
└──────┬──────┘
       │ Entity
       ▼
┌─────────────┐
│ Presenter   │  7. Output DTOに変換
│ (Adapter)   │
└──────┬──────┘
       │ Output DTO
       ▼
┌─────────────┐
│ Controller  │  8. HTTPレスポンスに変換
└──────┬──────┘
       │
       ▼
HTTP Response
```

---

## 6. 長所と短所

### 6.1 長所

```
┌─────────────────────────────────────────────────────────────┐
│  ✅ 長所                                                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. 明確な層と責務                                           │
│     4層 + Port/Adapterで役割が明確                          │
│     → 設計の指針が具体的                                    │
│                                                             │
│  2. 高いテスタビリティ                                       │
│     すべての層がインターフェースで分離                        │
│     → 各層を独立してテスト可能                              │
│                                                             │
│  3. フレームワーク独立                                       │
│     Entityは一切のフレームワークを知らない                   │
│     → フレームワーク変更の影響を受けない                    │
│                                                             │
│  4. ビジネスルールの保護                                     │
│     最内層にあり、外部依存がない                             │
│     → 長期的に安定                                          │
│                                                             │
│  5. 統一された原則                                           │
│     複数のアーキテクチャを統合                               │
│     → 業界標準の共通言語                                    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 6.2 短所

```
┌─────────────────────────────────────────────────────────────┐
│  ❌ 短所                                                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. 複雑さとコード量                                         │
│     インターフェース、DTO、Portが多い                        │
│     → 小規模プロジェクトには過剰                            │
│                                                             │
│  2. 学習コスト                                               │
│     概念の理解に時間がかかる                                 │
│     → チーム全員の理解が必要                                │
│                                                             │
│  3. 初期開発コスト                                           │
│     構造を整えるまでに時間がかかる                           │
│     → MVPやプロトタイプには不向き                           │
│                                                             │
│  4. 過度な抽象化のリスク                                     │
│     厳密に適用しすぎると冗長になる                           │
│     → プラグマティックな判断が必要                          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 7. オニオンアーキテクチャとの違い

```
┌─────────────────────────────────────────────────────────────┐
│                 オニオン vs クリーン                         │
├──────────────────────────┬──────────────────────────────────┤
│        オニオン           │          クリーン                │
├──────────────────────────┼──────────────────────────────────┤
│                          │                                  │
│  Domain Model            │  Entities                        │
│  Domain Services         │  （統合）                         │
│  Application Services    │  Use Cases                       │
│  Infrastructure          │  Interface Adapters              │
│                          │  Frameworks & Drivers            │
│                          │                                  │
├──────────────────────────┼──────────────────────────────────┤
│                          │                                  │
│  層の名前が実装寄り      │  層の名前が抽象的・概念的        │
│                          │                                  │
│  Port/Adapterは          │  Input Port/Output Portを        │
│  暗黙的                  │  明示的に定義                    │
│                          │                                  │
│  Presenterは未定義       │  Presenterを明示的に配置         │
│                          │                                  │
└──────────────────────────┴──────────────────────────────────┘
```

**本質は同じ**: 依存は内側へ、ビジネスロジックを守る

---

## 8. いつ使うべきか

```
┌─────────────────────────────────────────────────────────────┐
│  ✅ クリーンアーキテクチャが適している                       │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ・大規模・長期プロジェクト                                  │
│  ・複雑なビジネスドメイン                                    │
│  ・高いテストカバレッジが必要                                │
│  ・複数チームでの開発                                        │
│  ・技術スタックの変更可能性がある                            │
│  ・DDDを実践したい                                           │
│                                                             │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  ❌ クリーンアーキテクチャが適さない                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ・小規模・短期プロジェクト                                  │
│  ・シンプルなCRUD                                            │
│  ・プロトタイプ、MVP                                         │
│  ・チームがアーキテクチャに不慣れ                            │
│  ・スピード優先の開発                                        │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## まとめ

クリーンアーキテクチャは：

1. **The Dependency Rule**が核心: 依存は常に内側へ
2. **4層構造**: Entities → Use Cases → Adapters → Frameworks
3. **Input/Output Port**でデータフローを明示
4. ヘキサゴナル、オニオンを**統合・一般化**したもの
5. **大規模・長期プロジェクト**に最適

---

次のドキュメント: [06_comparison.md](./06_comparison.md) - 4つのアーキテクチャの比較まとめ
