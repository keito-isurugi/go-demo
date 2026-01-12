# コンテキストマップ

## 1. コンテキストマップ（図）

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                              ECサイト全体                                     │
│                                                                              │
│   ┌─────────────────────────────────────────┐                               │
│   │        注文コンテキスト                   │                               │
│   │       (Order Context)                   │                               │
│   │                                         │                               │
│   │  ┌─────────────┐  ┌─────────────┐       │                               │
│   │  │   Order     │  │   Stock     │       │                               │
│   │  │  (集約)     │  │  (集約)     │       │                               │
│   │  └─────────────┘  └─────────────┘       │                               │
│   │         │                │              │                               │
│   │         └───────┬────────┘              │                               │
│   │                 │                       │                               │
│   │    StockAllocationService               │                               │
│   │                                         │                               │
│   └─────────────────────────────────────────┘                               │
│          │                    │                                              │
│          │ OrderConfirmed     │ StockDepleted                               │
│          │ PaymentCompleted   │                                              │
│          ▼                    ▼                                              │
│   ┌──────────────┐     ┌──────────────────┐                                 │
│   │    支払      │     │    在庫管理       │                                 │
│   │  コンテキスト │     │   コンテキスト    │                                 │
│   │  (Payment)   │     │   (Inventory     │                                 │
│   │              │     │   Management)    │                                 │
│   │  ┌────────┐  │     │  ┌───────────┐   │                                 │
│   │  │Payment │  │     │  │ Inventory │   │                                 │
│   │  └────────┘  │     │  └───────────┘   │                                 │
│   └──────────────┘     └──────────────────┘                                 │
│          │                                                                   │
│          │ PaymentCompleted                                                  │
│          ▼                                                                   │
│   ┌──────────────┐                                                          │
│   │    配送      │                                                          │
│   │  コンテキスト │                                                          │
│   │  (Shipping)  │                                                          │
│   │              │                                                          │
│   │  ┌────────┐  │                                                          │
│   │  │Shipment│  │                                                          │
│   │  └────────┘  │                                                          │
│   └──────────────┘                                                          │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

## 2. 各コンテキストの責務とユビキタス言語

### 注文コンテキスト（Order Context）

| 用語 | 英語 | 定義 |
|------|------|------|
| 注文 | Order | 顧客が購入する商品のまとまり |
| 注文明細 | OrderLine | 注文内の1つの商品と数量 |
| 在庫 | Stock | 引当可能な商品の数量（available/reserved） |
| 引当 | Reserve | 注文のために在庫を確保すること |
| 確定 | Confirm | 注文を確定し、在庫を引き当てること |

**責務:**
- 注文の作成・管理
- 注文の状態遷移（Draft → Confirmed → Paid → Shipped → Delivered）
- 在庫の引当・解放

### 支払コンテキスト（Payment Context）

| 用語 | 英語 | 定義 |
|------|------|------|
| 支払 | Payment | 注文に対する決済処理 |
| 取引ID | TransactionID | 決済サービスが発行する一意の識別子 |
| 返金 | Refund | 支払い済みの金額を顧客に返すこと |

**責務:**
- 外部決済サービスとの連携
- 支払い状態の管理
- 返金処理

### 在庫管理コンテキスト（Inventory Management Context）

| 用語 | 英語 | 定義 |
|------|------|------|
| 在庫 | Inventory | 管理対象の商品（ロケーション、入荷日、賞味期限などを含む） |
| 入荷 | Receiving | 商品が倉庫に届くこと |
| 補充 | Replenishment | 在庫が少なくなったときに発注すること |
| 棚卸し | Stocktaking | 実際の在庫数を確認・調整すること |

**責務:**
- 入荷・補充の管理
- 在庫数の調整
- ロケーション管理
- 棚卸し

### 配送コンテキスト（Shipping Context）

| 用語 | 英語 | 定義 |
|------|------|------|
| 出荷 | Shipment | 商品を倉庫から発送すること |
| 配送 | Delivery | 商品を顧客に届けること |
| 追跡番号 | TrackingNumber | 配送業者が発行する追跡用の識別子 |

**責務:**
- 出荷準備
- 配送業者との連携
- 配送状況の追跡

## 3. コンテキスト間の連携方法

### 注文 → 支払（Pub/Sub）

```
パターン: 公開ホストサービス（Open Host Service）+ 公表された言語（Published Language）
```

**理由:**
- 注文コンテキストは支払の詳細を知る必要がない
- OrderConfirmedイベントを発行するだけで、支払コンテキストが処理を開始
- 疎結合を維持しつつ、イベント形式は共通フォーマットで定義

**連携フロー:**
```
1. Order.Confirm() → OrderConfirmedイベント生成
2. EventPublisher がイベントを発行
3. 支払コンテキストがイベントを購読し、決済処理を開始
```

### 注文 → 在庫管理（Pub/Sub）

```
パターン: 公開ホストサービス（Open Host Service）+ 公表された言語（Published Language）
```

**理由:**
- 在庫が枯渇したことを在庫管理コンテキストに通知
- 注文コンテキストは発注処理の詳細を知らない
- 在庫管理コンテキストが独自のルールで発注を判断

**連携フロー:**
```
1. Stock.Reserve() → 在庫が閾値以下になるとStockDepletedイベント生成
2. EventPublisher がイベントを発行
3. 在庫管理コンテキストがイベントを購読し、発注処理を開始
```

### 支払 → 配送（Pub/Sub）

```
パターン: 公開ホストサービス（Open Host Service）+ 公表された言語（Published Language）
```

**理由:**
- 支払完了を配送コンテキストに通知
- 配送準備は支払が完了してから開始される
- 支払コンテキストは配送の詳細を知らない

**連携フロー:**
```
1. Order.Pay() → PaymentCompletedイベント生成
2. EventPublisher がイベントを発行
3. 配送コンテキストがイベントを購読し、出荷準備を開始
```

### 注文 ↔ 外部決済サービス（腐敗防止層）

```
パターン: 腐敗防止層（Anti-Corruption Layer）
```

**理由:**
- 外部決済サービスのモデル（StatusCode: 0/1など）がドメインに侵入するのを防ぐ
- ドメインの言葉（Success: true/false）に変換する
- 外部サービスの変更がドメインに影響しないようにする

## 4. 各コンテキストでの「在庫」の意味の違い

### 注文コンテキストの「Stock」

```go
type Stock struct {
    productID ProductID
    available Available  // 引当可能な数量
    reserved  Reserved   // 引当済みの数量
}
```

**関心事:**
- 注文に対して在庫を引き当てられるか？
- 引当後の残り在庫数は？

**操作:**
- Reserve（引当）
- Release（解放）
- CanReserve（引当可能か確認）

### 在庫管理コンテキストの「Inventory」

```go
type Inventory struct {
    productID    ProductID
    quantity     int           // 総在庫数
    location     Location      // 保管場所（棚番号など）
    receivedDate time.Time     // 入荷日
    expiryDate   *time.Time    // 賞味期限（あれば）
    lotNumber    string        // ロット番号
}
```

**関心事:**
- どこに保管されているか？
- いつ入荷したか？
- 賞味期限は？
- 補充が必要か？

**操作:**
- Receive（入荷登録）
- Adjust（数量調整）
- Transfer（ロケーション移動）
- CountStock（棚卸し）

### なぜ分けるのか？

| 観点 | 注文コンテキスト | 在庫管理コンテキスト |
|------|-----------------|---------------------|
| 目的 | 注文を処理する | 在庫を管理する |
| 関心事 | 引当可能か | どこにあるか、いつ届いたか |
| 操作 | Reserve/Release | Receive/Adjust/Transfer |
| 変更頻度 | 注文ごと | 入荷・棚卸しごと |

**同じ「在庫」という言葉でも、コンテキストによって意味が異なる**のがDDDの核心です。

## 5. 腐敗防止層（ACL）の実装例

### インターフェース定義（アプリケーション層）

```go
// application/payment_gateway.go

// PaymentResult ドメインの言葉で表現された支払い結果
type PaymentResult struct {
    Success       bool
    TransactionID string
    ErrorMessage  string
}

// PaymentGateway 外部決済サービスとの連携インターフェース
type PaymentGateway interface {
    Charge(ctx context.Context, orderID string, amount int, currency string) (*PaymentResult, error)
    Refund(ctx context.Context, transactionID string, amount int, currency string) (*PaymentResult, error)
}
```

### アダプター実装（インフラストラクチャ層）

```go
// infrastructure/payment/stripe_adapter.go

// StripeAdapter Stripe決済サービスのアダプター
type StripeAdapter struct {
    client *stripe.Client
}

func (a *StripeAdapter) Charge(ctx context.Context, orderID string, amount int, currency string) (*application.PaymentResult, error) {
    // 1. Stripeのモデルで決済を実行
    resp, err := a.client.Charge(&stripe.ChargeParams{
        Amount:   int64(amount),
        Currency: currency,
        Metadata: map[string]string{"order_id": orderID},
    })

    if err != nil {
        return nil, err
    }

    // 2. Stripeのレスポンスをドメインの言葉に変換（腐敗防止）
    //    外部サービスのモデル → ドメインのモデル
    return &application.PaymentResult{
        Success:       resp.Status == "succeeded",  // Stripeの"succeeded"を変換
        TransactionID: resp.ID,                     // "ch_xxx"のような形式
        ErrorMessage:  resp.FailureMessage,
    }, nil
}
```

### 変換のポイント

| 外部サービス（Stripe） | ドメイン |
|----------------------|---------|
| Status: "succeeded" | Success: true |
| Status: "failed" | Success: false |
| ID: "ch_xxx" | TransactionID: "ch_xxx" |
| FailureMessage | ErrorMessage |

**メリット:**
- Stripeの仕様変更があってもアダプター内で吸収
- ドメイン層は外部サービスの存在を知らない
- テスト時にモックに差し替えやすい

## まとめ

```
┌─────────────────────────────────────────────────────────────┐
│                     コンテキスト間の関係                      │
├─────────────────────────────────────────────────────────────┤
│  注文 ─── Pub/Sub ───→ 支払 ─── Pub/Sub ───→ 配送          │
│    │                                                        │
│    └──── Pub/Sub ───→ 在庫管理                              │
│                                                             │
│  注文 ←── ACL ───→ 外部決済サービス                          │
└─────────────────────────────────────────────────────────────┘

【パターン選択の理由】
- Pub/Sub: 疎結合を維持、非同期処理が可能
- ACL: 外部サービスのモデルからドメインを保護
```
