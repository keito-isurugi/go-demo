# DDD（ドメイン駆動設計）学習課題

## 概要

この課題では、実際の業務で活用できるDDDの知識を身につけることを目的としています。
ECサイトの「注文管理システム」を題材に、DDDの各概念を段階的に学習していきます。

---

## 課題1: ユビキタス言語の定義

### 目的
ドメインエキスパートと開発者が共通認識を持つための「ユビキタス言語」を定義する。

### 課題内容

以下のビジネス要件を読み、ユビキタス言語（用語集）を作成してください。

**ビジネス要件:**
> 顧客は商品をカートに入れ、注文を確定できます。
> 注文が確定すると在庫が引き当てられ、支払いが処理されます。
> 支払いが完了すると注文は「支払済み」になり、出荷準備が始まります。
> 在庫が不足している場合は「入荷待ち」となり、入荷次第出荷されます。

### 成果物
`docs/ubiquitous_language.md` に以下を記載：
- 用語一覧（日本語・英語・定義）
- 用語間の関係性

### ヒント
- Customer, Order, Cart, Product, Stock, Payment, Shipment などが候補
- 「注文確定」「在庫引当」などの動詞も重要

---

## 課題2: 値オブジェクト（Value Object）の実装

### 目的
不変で同一性を持たないオブジェクトを実装し、ドメインの概念をコードで表現する。

### 課題内容

以下の値オブジェクトを実装してください：

1. **Money（金額）**
   - 金額と通貨を持つ
   - 加算・減算のメソッドを持つ
   - 異なる通貨同士の計算は許可しない

2. **Quantity（数量）**
   - 0以上の整数のみ許可
   - 加算・減算のメソッドを持つ

3. **Email（メールアドレス）**
   - 形式バリデーション
   - 正規化（小文字化など）

### 実装場所
```
domain/
  money.go
  money_test.go
  quantity.go
  quantity_test.go
  email.go
  email_test.go
```

### 実装例（Money）
```go
package domain

import "errors"

type Currency string

const (
    JPY Currency = "JPY"
    USD Currency = "USD"
)

type Money struct {
    amount   int
    currency Currency
}

func NewMoney(amount int, currency Currency) (Money, error) {
    if amount < 0 {
        return Money{}, errors.New("amount must be non-negative")
    }
    return Money{amount: amount, currency: currency}, nil
}

func (m Money) Add(other Money) (Money, error) {
    if m.currency != other.currency {
        return Money{}, errors.New("currency mismatch")
    }
    return NewMoney(m.amount+other.amount, m.currency)
}

func (m Money) Amount() int {
    return m.amount
}

func (m Money) Currency() Currency {
    return m.currency
}

func (m Money) Equals(other Money) bool {
    return m.amount == other.amount && m.currency == other.currency
}
```

---

## 課題3: エンティティ（Entity）の実装

### 目的
一意の識別子を持ち、ライフサイクルを通じて同一性を保つオブジェクトを実装する。

### 課題内容

以下のエンティティを実装してください：

1. **Order（注文）**
   - OrderID（識別子）を持つ
   - 注文明細（OrderLine）のコレクションを持つ
   - 注文ステータスを持つ（Draft, Confirmed, Paid, Shipped, Delivered, Cancelled）
   - 合計金額を計算できる
   - 注文明細を追加・削除できる

2. **OrderLine（注文明細）**
   - 商品ID、商品名、単価、数量を持つ
   - 小計を計算できる

### 実装場所
```
domain/
  order.go
  order_test.go
  order_line.go
```

### ポイント
- エンティティは可変だが、不変条件（Invariant）を守る
- 状態遷移のルールをエンティティ内でカプセル化する

---

## 課題4: 集約（Aggregate）の設計

### 目的
トランザクション整合性の境界を定義し、集約ルートを通じたアクセスを実装する。

### 課題内容

「注文」集約を設計・実装してください：

1. **集約の境界を決める**
   - Order（集約ルート）
   - OrderLine（集約内エンティティ）
   - どこまでを1つの集約とするか検討

2. **集約ルート経由のアクセス**
   - OrderLineへの操作はすべてOrder経由で行う
   - 外部からOrderLineを直接変更できないようにする

3. **不変条件の保護**
   - 注文の合計金額が上限を超えないこと
   - 注文明細が0件の注文は確定できないこと

### 実装場所
```
domain/
  order/
    order.go
    order_test.go
```

※ 集約は機能ドメイン単位でディレクトリを切る（技術的分類ではなく）

### 設計図
```
┌─────────────────────────────────────┐
│           Order Aggregate           │
│  ┌───────────────────────────────┐  │
│  │    Order (Aggregate Root)     │  │
│  │  - OrderID                    │  │
│  │  - CustomerID                 │  │
│  │  - Status                     │  │
│  │  - OrderLines []              │  │
│  │  - TotalAmount                │  │
│  └───────────────────────────────┘  │
│              │                      │
│              ▼                      │
│  ┌───────────────────────────────┐  │
│  │        OrderLine              │  │
│  │  - ProductID                  │  │
│  │  - ProductName                │  │
│  │  - UnitPrice                  │  │
│  │  - Quantity                   │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘
```

---

## 課題5: ドメインサービスの実装

### 目的
単一のエンティティに属さないドメインロジックを実装する。

### 課題内容

以下のドメインサービスを実装してください：

1. **在庫引当サービス（StockAllocationService）**
   - 注文に対して在庫を引き当てる
   - 在庫不足の場合は部分引当または引当失敗を返す

2. **割引計算サービス（DiscountService）**
   - 顧客ランク、購入金額、クーポンに基づいて割引を計算
   - 複数の割引ルールを組み合わせる

### 実装場所
```
domain/
  stock_allocation_service.go
  stock_allocation_service_test.go
  discount_service.go
  discount_service_test.go
```

### ポイント
- ドメインサービスは状態を持たない
- 複数の集約をまたがるロジックに適している

---

## 課題6: リポジトリパターンの実装

### 目的
ドメインオブジェクトの永続化を抽象化し、インフラストラクチャ層への依存を分離する。

### 課題内容

1. **リポジトリインターフェースの定義（ドメイン層）**

```go
// domain/repository/order_repository.go
package repository

import "context"

type OrderRepository interface {
    Save(ctx context.Context, order *Order) error
    FindByID(ctx context.Context, id OrderID) (*Order, error)
    FindByCustomerID(ctx context.Context, customerID CustomerID) ([]*Order, error)
    Delete(ctx context.Context, id OrderID) error
}
```

2. **リポジトリ実装（インフラストラクチャ層）**

```go
// infrastructure/persistence/order_repository_impl.go
```

### 実装場所
```
domain/
  order_repository.go

infrastructure/
  persistence/
    order_repository.go
    order_repository_test.go
```

### ポイント
- インターフェースはドメイン層に置く
- 実装はインフラストラクチャ層に置く
- テストではインメモリ実装を使用

---

## 課題7: ドメインイベントの実装

### 目的
集約間の整合性を疎結合に保ちながら実現する。

### 課題内容

以下のドメインイベントを実装してください：

1. **OrderConfirmed（注文確定）**
   - 注文が確定されたときに発行
   - 在庫引当処理をトリガー

2. **PaymentCompleted（支払完了）**
   - 支払いが完了したときに発行
   - 出荷準備処理をトリガー

3. **StockDepleted（在庫枯渇）**
   - 在庫が閾値を下回ったときに発行
   - 発注処理をトリガー

### 実装場所
```
domain/
  event.go
  order_confirmed.go
  payment_completed.go
  stock_depleted.go

application/
  order_event_handler.go
```

### 実装例
```go
// domain/event.go
package domain

import "time"

type DomainEvent interface {
    OccurredAt() time.Time
    AggregateID() string
    EventType() string
}

// domain/order_confirmed.go
type OrderConfirmed struct {
    orderID    string
    customerID string
    totalAmount int
    occurredAt time.Time
}

func NewOrderConfirmed(orderID, customerID string, totalAmount int) OrderConfirmed {
    return OrderConfirmed{
        orderID:     orderID,
        customerID:  customerID,
        totalAmount: totalAmount,
        occurredAt:  time.Now(),
    }
}
```

---

## 課題8: アプリケーションサービスの実装

### 目的
ユースケースを実装し、ドメインオブジェクトを協調させる。

### 課題内容

以下のユースケースを実装してください：

1. **注文作成ユースケース**
   - 顧客IDと商品リストから注文を作成
   - カートの内容を注文に変換

2. **注文確定ユースケース**
   - 注文を確定状態に変更
   - 在庫を引き当てる
   - ドメインイベントを発行

3. **支払処理ユースケース**
   - 外部決済サービスと連携
   - 支払い結果を注文に反映

### 実装場所
```
application/
  create_order.go
  create_order_test.go
  confirm_order.go
  confirm_order_test.go
  process_payment.go
  process_payment_test.go
```

### 実装例
```go
// application/confirm_order.go
package application

type ConfirmOrderUseCase struct {
    orderRepo      domain.OrderRepository
    stockService   domain.StockAllocationService
    eventPublisher domain.Publisher
}

type ConfirmOrderInput struct {
    OrderID string
}

type ConfirmOrderOutput struct {
    Success bool
    Message string
}

func (uc *ConfirmOrderUseCase) Execute(ctx context.Context, input ConfirmOrderInput) (*ConfirmOrderOutput, error) {
    // 1. 注文を取得
    order, err := uc.orderRepo.FindByID(ctx, OrderID(input.OrderID))
    if err != nil {
        return nil, err
    }

    // 2. 在庫を引き当て
    err = uc.stockService.Allocate(ctx, order)
    if err != nil {
        return nil, err
    }

    // 3. 注文を確定
    err = order.Confirm()
    if err != nil {
        return nil, err
    }

    // 4. 永続化
    err = uc.orderRepo.Save(ctx, order)
    if err != nil {
        return nil, err
    }

    // 5. ドメインイベント発行
    uc.eventPublisher.Publish(domain.NewOrderConfirmed(
        order.ID().String(),
        order.CustomerID().String(),
        order.TotalAmount().Amount(),
    ))

    return &ConfirmOrderOutput{Success: true, Message: "Order confirmed"}, nil
}
```

---

## 課題9: 境界づけられたコンテキストの設計

### 目的
複数のコンテキスト間の関係を理解し、コンテキストマップを作成する。

### 課題内容

ECサイトを以下のコンテキストに分割し、設計してください：

1. **注文コンテキスト（Order Context）**
   - 注文の作成・管理

2. **在庫コンテキスト（Inventory Context）**
   - 在庫の管理・引当

3. **支払コンテキスト（Payment Context）**
   - 決済処理

4. **配送コンテキスト（Shipping Context）**
   - 出荷・配送管理

### 成果物

1. **コンテキストマップ**（図で表現）
2. **各コンテキストの責務**
3. **コンテキスト間の連携方法**
   - 共有カーネル / 顧客-供給者 / 腐敗防止層 など

### コンテキストマップ例
```
┌─────────────────────────────────────────────────────────────┐
│                        ECサイト                             │
│                                                             │
│   ┌─────────────┐    Pub/Sub    ┌─────────────┐            │
│   │   注文      │──────────────▶│   在庫      │            │
│   │  Context    │               │  Context    │            │
│   └─────────────┘               └─────────────┘            │
│          │                                                  │
│          │ ACL                                              │
│          ▼                                                  │
│   ┌─────────────┐               ┌─────────────┐            │
│   │   支払      │◀─────────────│   配送      │            │
│   │  Context    │   Conformist  │  Context    │            │
│   └─────────────┘               └─────────────┘            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 課題10: 総合演習 - 完全なDDDアプリケーション

### 目的
これまでの学習内容を統合し、実践的なアプリケーションを構築する。

### 課題内容

「注文管理システム」を完成させてください：

**機能要件:**
1. 顧客が商品を選択して注文を作成できる
2. 注文を確定すると在庫が引き当てられる
3. 支払いを処理できる
4. 支払い完了後、出荷準備に移行する
5. 注文のキャンセルができる（条件付き）

**非機能要件:**
1. クリーンアーキテクチャに基づいた層分離
2. 単体テストのカバレッジ80%以上
3. ドメインイベントによる疎結合な設計

### 最終的なディレクトリ構成
```
demo/ddd/
├── cmd/
│   └── server/
│       └── main.go
├── domain/
│   ├── order/           # 注文集約（集約ルート + 関連エンティティ）
│   │   ├── order.go
│   │   ├── order_line.go
│   │   └── order_test.go
│   ├── money.go          # 値オブジェクト
│   ├── quantity.go
│   ├── email.go
│   ├── event.go          # ドメインイベントインターフェース
│   ├── order_confirmed.go
│   ├── order_repository.go  # リポジトリインターフェース
│   └── stock_allocation_service.go  # ドメインサービス
├── application/
│   ├── create_order.go
│   ├── confirm_order.go
│   ├── process_payment.go
│   └── order_event_handler.go
├── infrastructure/
│   ├── persistence/
│   │   └── order_repository.go  # リポジトリ実装
│   └── external/
├── interfaces/
│   └── api/
│       └── handler/
├── docs/
│   └── ubiquitous_language.md
└── go.mod
```

### ディレクトリ構成のポイント
- **技術的分類（valueobject, entity, service等）でディレクトリを切らない**
- **機能ドメイン単位（order等）でディレクトリを切る**
- domainパッケージはフラットに保ち、インポートパスをシンプルに
- 集約が大きくなった場合のみサブディレクトリを検討

---

## 学習の進め方

1. **課題1-2**: DDDの基礎概念を理解する（1-2日）
2. **課題3-4**: エンティティと集約の設計を学ぶ（2-3日）
3. **課題5-6**: ドメインサービスとリポジトリを実装する（2-3日）
4. **課題7-8**: イベント駆動とユースケースを実装する（3-4日）
5. **課題9-10**: 全体設計と統合（3-5日）

---

## 参考資料

### 書籍
- 「エリック・エヴァンスのドメイン駆動設計」
- 「実践ドメイン駆動設計」
- 「ドメイン駆動設計入門」

### オンラインリソース
- [DDD Reference](https://www.domainlanguage.com/ddd/reference/)
- [Go + DDDのサンプル実装](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example)

---

## チェックリスト

各課題完了時に確認してください：

- [ ] 課題1: ユビキタス言語を定義した
- [ ] 課題2: 値オブジェクトを実装した（テスト含む）
- [ ] 課題3: エンティティを実装した（テスト含む）
- [ ] 課題4: 集約を設計・実装した
- [ ] 課題5: ドメインサービスを実装した
- [ ] 課題6: リポジトリパターンを実装した
- [ ] 課題7: ドメインイベントを実装した
- [ ] 課題8: アプリケーションサービスを実装した
- [ ] 課題9: コンテキストマップを作成した
- [ ] 課題10: 総合演習を完了した
