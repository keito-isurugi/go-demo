# DDD（ドメイン駆動設計）学習課題

## 概要

この課題では、実際の業務で活用できるDDDの知識を身につけることを目的としています。
ECサイトの「注文管理システム」を題材に、DDDの各概念を段階的に学習していきます。

### 学習の範囲について（重要）
- **課題1-8**: 「注文コンテキスト」という**単一の境界づけられたコンテキスト内**での戦術的設計を学びます
  - この段階では、OrderとStockは同じコンテキスト内に存在します
  - トランザクション整合性が必要な概念は同じコンテキスト内に配置します
- **課題9**: 複数のコンテキスト間の関係（戦略的設計）を学びます
  - 注文コンテキストと在庫管理コンテキストなど、異なる関心事を持つコンテキストの関係を学びます

この段階的アプローチにより、まず基礎を固めてから全体像を学べます。

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

### 成果物
`docs/ubiquitous_language.md` に以下を記載：
- 用語一覧（日本語・英語・定義）
- 用語間の関係性

### ヒント
- **名詞**: Customer（顧客）, Order（注文）, Cart（カート）, Product（商品）, Stock（在庫）, Payment（支払）, Shipment（出荷）
- **動詞**: Confirm（確定する）, Reserve（引き当てる）, Pay（支払う）, Ship（出荷する）
- **状態**: Draft（下書き）, Confirmed（確定済み）, Paid（支払済み）, Shipped（出荷済み）
- **重要**: 技術用語（Entity, Aggregateなど）ではなく、**ビジネスの言葉**で記述する

---

## 課題2: 値オブジェクト（Value Object）の実装

### 目的
不変で同一性を持たないオブジェクトを実装し、ドメインの概念をコードで表現する。

### DDDにおける値オブジェクトの重要性
エリック・エヴァンスは、値オブジェクトを「**概念的な同一性を持たず、オブジェクトの属性のみが重要なもの**」と定義しています。
値オブジェクトを適切に使うことで：
- **ドメインの概念を明示的に表現**できる（intやstringではなくMoneyやQuantity）
- **不変性により予期せぬ副作用を防ぐ**
- **バリデーションロジックを一箇所に集約**できる
- **ドメインロジックを値オブジェクト自身に持たせる**ことで表現力が高まる

### 課題内容

以下の値オブジェクトを実装してください：

1. **Money（金額）**
   - 金額と通貨を持つ
   - 加算・減算のメソッドを持つ
   - 異なる通貨同士の計算は許可しない
   - **重要**: 比較演算（IsGreaterThan, IsLessThan）も実装する

2. **Quantity（数量）**
   - 0以上の整数のみ許可
   - 加算・減算のメソッドを持つ（**値オブジェクト同士**の演算として実装）
   - **重要**: 乗算メソッド（単価×数量の計算に使用）も実装する

3. **Email（メールアドレス）**
   - 形式バリデーション
   - 正規化（小文字化など）

4. **ProductID（商品ID）**
   - 空文字列を許可しない
   - 識別子として使用される値オブジェクト

### 実装場所
```
domain/
  money.go
  money_test.go
  quantity.go
  quantity_test.go
  email.go
  email_test.go
  product_id.go
  product_id_test.go
```

### 実装のヒント（Money）
```go
package domain

// Money 値オブジェクト
// - 金額(amount)と通貨(currency)をフィールドに持つ
// - フィールドは非公開にして不変性を保つ

// Currency 通貨を表す型
// - JPY, USD などの定数を定義

// NewMoney コンストラクタ
// - 負の金額はエラーを返す

// Add 加算メソッド
// - 異なる通貨同士の計算はエラーを返す
// - 新しいMoneyを返す（不変性を保つ）

// Subtract 減算メソッド
// - 異なる通貨同士の計算はエラーを返す
// - 結果が負になる場合はエラーを返す

// IsGreaterThan 比較演算
// - 異なる通貨同士の比較はエラーを返す

// Equals 等価性の比較
// - 金額と通貨が両方一致すれば等しい
```

### 実装のヒント（Quantity）
```go
package domain

// Quantity 数量を表す値オブジェクト
// - value int をフィールドに持つ（非公開）

// NewQuantity コンストラクタ
// - 0未満の場合はエラーを返す

// Add 加算メソッド
// - **重要**: 引数として別のQuantityを受け取る（値オブジェクト同士の演算）
// - 新しいQuantityを返す（不変性を保つ）

// Subtract 減算メソッド
// - **重要**: 引数として別のQuantityを受け取る（値オブジェクト同士の演算）
// - 結果が負になる場合はエラーを返す

// Multiply 乗算メソッド
// - 単価×数量の計算に使用
// - 新しいQuantityを返す（不変性を保つ）

// IsZero ゼロかどうかを判定

// Value 値を取得するゲッター
```

### 実装のヒント（Email）
```go
package domain

// Email メールアドレスを表す値オブジェクト
// - value string をフィールドに持つ（正規化済み）

// NewEmail コンストラクタ
// - 形式バリデーション（@を含むか、など）
// - 正規化（小文字に変換、前後の空白を除去）

// String 文字列として取得

// Domain ドメイン部分を取得（@以降）

// Equals 等価性の比較
```

---

## 課題3: エンティティ（Entity）の実装

### 目的
一意の識別子を持ち、ライフサイクルを通じて同一性を保つオブジェクトを実装する。

### DDDにおけるエンティティの重要性
エリック・エヴァンスは、エンティティを「**連続性と同一性によって定義されるオブジェクト**」と定義しています。
エンティティの特徴：
- **一意の識別子**を持つ（属性が変わっても同じエンティティとして認識される）
- **可変である**が、**不変条件（Invariant）を常に維持**する
- **ビジネスルールをカプセル化**し、自身の整合性を守る
- **状態遷移のルール**をエンティティ内に持つ

### 課題内容

以下のエンティティを実装してください：

1. **Order（注文）**
   - OrderID（識別子）を持つ
   - 注文明細（OrderLine）のコレクションを持つ
   - 注文ステータスを持つ（Draft, Confirmed, Paid, Shipped, Delivered, Cancelled）
   - 合計金額を計算できる
   - 注文明細を追加・削除できる
   - **重要**: 不正な状態遷移を防ぐメソッドを実装（例: Paid状態からDraftには戻れない）

2. **OrderLine（注文明細）**
   - 商品ID、商品名、単価、数量を持つ
   - 小計を計算できる
   - **重要**: OrderLineはエンティティではなく**値オブジェクト**として扱う選択肢も検討する

### 実装場所
```
domain/
  order.go
  order_test.go
  order_line.go
  order_line_test.go
  order_status.go  # ステータスの型定義と遷移ルール
```

### ポイント
- **エンティティは可変だが、不変条件（Invariant）を守る**
  - 例: 注文明細が0件の状態では注文を確定できない
  - 例: 支払済みの注文はキャンセルできない（または特別な処理が必要）
- **状態遷移のルールをエンティティ内でカプセル化する**
  - setter経由ではなく、ビジネスの意味を持つメソッド（Confirm(), Pay(), Cancel()）で状態を変更
- **ファクトリメソッド**を使ってエンティティを生成し、生成時の不変条件を保証する

### 状態遷移図とテストケース

注文の状態遷移を図で表現し、許可される遷移と禁止される遷移を明確にしてください：

```
状態遷移図:

Draft → Confirmed → Paid → Shipped → Delivered
  ↓         ↓         ↓
Cancelled Cancelled   ✗
                   (キャンセル不可)
```

**許可される遷移:**
- Draft → Confirmed: 注文確定
- Draft → Cancelled: 確定前のキャンセル
- Confirmed → Paid: 支払い完了
- Confirmed → Cancelled: 支払い前のキャンセル
- Paid → Shipped: 出荷
- Shipped → Delivered: 配送完了

**禁止される遷移（エラーになる）:**
- Confirmed → Draft: 確定後は下書きに戻せない
- Paid → Draft: 支払い後は下書きに戻せない
- Paid → Cancelled: 支払い後のキャンセルは不可（または特別な返金処理が必要）
- Shipped → Cancelled: 出荷後はキャンセル不可
- Delivered → Cancelled: 配送完了後はキャンセル不可

---

## 課題4: 集約（Aggregate）の設計

### 目的
トランザクション整合性の境界を定義し、集約ルートを通じたアクセスを実装する。

### DDDにおける集約の重要性
エリック・エヴァンスは、集約を「**関連するオブジェクトの集まりで、データ変更の単位として扱われるもの**」と定義しています。
集約の設計原則：
- **トランザクション整合性の境界**を定義する（集約内は強い整合性、集約間は結果整合性）
- **集約ルート（Aggregate Root）**を通じてのみ、集約内のエンティティにアクセスする
- **集約はできるだけ小さく保つ**（パフォーマンスと並行性のため）
- **不変条件（Invariant）**を常に維持する
- **集約IDで他の集約を参照する**（オブジェクト参照ではなく）

### 課題内容

以下の2つの集約を設計・実装してください：

#### 4-1. 注文集約（Order Aggregate）

1. **集約の境界を決める**
   - Order（集約ルート）
   - OrderLine（集約内エンティティ、または値オブジェクト）
   - **重要**: CustomerやProductは別の集約なので、IDで参照する（オブジェクト参照しない）

2. **集約ルート経由のアクセス**
   - OrderLineへの操作はすべてOrder経由で行う
   - 外部からOrderLineを直接変更できないようにする
   - **実装例**:
     ```go
     // ❌ 悪い例: OrderLineを直接操作
     orderLine := order.Lines[0]
     orderLine.ChangeQuantity(newQty)

     // ✅ 良い例: Order経由で操作
     order.ChangeOrderLineQuantity(productID, newQty)
     ```

3. **不変条件の保護**

   集約の最も重要な責務は、**不変条件（Invariant）を常に維持すること**です。以下の不変条件を実装してください：

   - 注文の合計金額が上限を超えないこと（例: 100万円）
   - 注文明細が0件の注文は確定できないこと
   - **同じ商品の重複注文明細を許可しない**（同じProductIDのOrderLineは1つのみ）
     ```go
     func (o *Order) AddLine(line OrderLine) error {
         if o.status != Draft {
             return errors.New("can only add lines to draft orders")
         }

         // 重複チェック（重要な不変条件）
         for _, existingLine := range o.lines {
             if existingLine.ProductID() == line.ProductID() {
                 return errors.New("product already exists in order")
             }
         }

         // 合計金額の上限チェック
         newTotal, err := o.calculateTotalWithNewLine(line)
         if err != nil {
             return err
         }

         // Money型を使った比較
         maxAmount, _ := NewMoney(1000000, JPY)
         if newTotal.IsGreaterThan(maxAmount) {
             return errors.New("order total exceeds maximum allowed amount")
         }

         o.lines = append(o.lines, line)
         return nil
     }
     ```
   - Draft状態でのみ明細の追加・削除を許可する

   **不変条件を守るテストケース:**
   ```go
   func TestOrder_AddLine_DuplicateProduct_Error(t *testing.T) {
       // 同じProductIDの明細を2回追加しようとするとエラーになることを確認
   }

   func TestOrder_AddLine_ExceedMaxAmount_Error(t *testing.T) {
       // 合計金額が100万円を超える明細を追加しようとするとエラーになることを確認
   }

   func TestOrder_AddLine_NotDraftStatus_Error(t *testing.T) {
       // Confirmed状態の注文に明細を追加しようとするとエラーになることを確認
   }
   ```

   **重要**: 不変条件は集約ルート（Order）のメソッド内で必ずチェックし、違反する操作は拒否してください。

#### 4-2. 在庫集約（Stock Aggregate）

**重要**: この課題を始める前に、Stockエンティティを実装してください。課題5（ドメインサービス）で使用します。

1. **Stockエンティティの実装**
   - ProductID（商品ID）を識別子として持つ
   - available（利用可能在庫数）とreserved（引当済み在庫数）を管理
   - 在庫の引当・解放のメソッドを持つ

2. **不変条件の保護**
   - 利用可能在庫数は0未満にならない
   - 引当済み在庫数は0未満にならない
   - 引当時は利用可能在庫から引当済みに移動

### 実装場所
```
domain/
  order/
    order.go
    order_test.go
  stock.go
  stock_test.go
```

※ 集約は機能ドメイン単位でディレクトリを切る（技術的分類ではなく）

### Stockエンティティ実装のヒント
```go
package domain

// Stock 在庫エンティティ（Stock集約のルート）
type Stock struct {
    productID ProductID
    available Quantity  // 利用可能在庫数
    reserved  Quantity  // 引当済み在庫数
}

// NewStock コンストラクタ
func NewStock(productID ProductID, available Quantity) (*Stock, error) {
    if productID == "" {
        return nil, errors.New("product ID cannot be empty")
    }
    reserved, _ := NewQuantity(0)
    return &Stock{
        productID: productID,
        available: available,
        reserved:  reserved,
    }, nil
}

// Reserve 在庫を引き当てる
// - 利用可能在庫から引当済みに移動
// - 在庫不足の場合はエラーを返す
func (s *Stock) Reserve(quantity Quantity) error {
    if !s.CanReserve(quantity) {
        return errors.New("insufficient stock")
    }

    // 利用可能在庫を減らし、引当済みを増やす
    s.available, _ = s.available.Subtract(quantity)
    s.reserved, _ = s.reserved.Add(quantity)
    return nil
}

// Release 引当を解放する
// - 引当済みから利用可能在庫に戻す
func (s *Stock) Release(quantity Quantity) error {
    // 引当済み在庫が十分あるかチェック
    if s.reserved.Value() < quantity.Value() {
        return errors.New("cannot release more than reserved")
    }

    s.reserved, _ = s.reserved.Subtract(quantity)
    s.available, _ = s.available.Add(quantity)
    return nil
}

// CanReserve 引当可能かチェック
func (s Stock) CanReserve(quantity Quantity) bool {
    return s.available.Value() >= quantity.Value()
}

// ProductID ゲッター
func (s Stock) ProductID() ProductID {
    return s.productID
}

// Available ゲッター
func (s Stock) Available() Quantity {
    return s.available
}

// Reserved ゲッター
func (s Stock) Reserved() Quantity {
    return s.reserved
}
```

### Stockエンティティのテストケース例

```go
// stock_test.go
package domain_test

import (
    "testing"
    "yourproject/domain"
)

func TestStock_Reserve_Success(t *testing.T) {
    // 利用可能在庫100、10個を引き当て
    available, _ := domain.NewQuantity(100)
    stock, _ := domain.NewStock("PROD-001", available)
    reserveQty, _ := domain.NewQuantity(10)

    err := stock.Reserve(reserveQty)

    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if stock.Available().Value() != 90 {
        t.Errorf("expected available 90, got %d", stock.Available().Value())
    }
    if stock.Reserved().Value() != 10 {
        t.Errorf("expected reserved 10, got %d", stock.Reserved().Value())
    }
}

func TestStock_Reserve_InsufficientStock_Error(t *testing.T) {
    // 利用可能在庫10、20個を引き当てようとする
    available, _ := domain.NewQuantity(10)
    stock, _ := domain.NewStock("PROD-001", available)
    reserveQty, _ := domain.NewQuantity(20)

    err := stock.Reserve(reserveQty)

    if err == nil {
        t.Error("expected error for insufficient stock")
    }
}

func TestStock_Release_Success(t *testing.T) {
    // 引当済み10個を解放
    available, _ := domain.NewQuantity(90)
    stock, _ := domain.NewStock("PROD-001", available)
    reserveQty, _ := domain.NewQuantity(10)
    stock.Reserve(reserveQty)

    releaseQty, _ := domain.NewQuantity(10)
    err := stock.Release(releaseQty)

    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if stock.Available().Value() != 100 {
        t.Errorf("expected available 100, got %d", stock.Available().Value())
    }
    if stock.Reserved().Value() != 0 {
        t.Errorf("expected reserved 0, got %d", stock.Reserved().Value())
    }
}

func TestStock_CanReserve(t *testing.T) {
    available, _ := domain.NewQuantity(100)
    stock, _ := domain.NewStock("PROD-001", available)

    qty1, _ := domain.NewQuantity(50)
    if !stock.CanReserve(qty1) {
        t.Error("expected to be able to reserve 50")
    }

    qty2, _ := domain.NewQuantity(150)
    if stock.CanReserve(qty2) {
        t.Error("expected not to be able to reserve 150")
    }
}
```

**必ず実装してテストしてください**: 課題5（ドメインサービス）でStockエンティティを使用します。

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

┌─────────────────────────────────────┐
│           Stock Aggregate           │
│  ┌───────────────────────────────┐  │
│  │    Stock (Aggregate Root)     │  │
│  │  - ProductID                  │  │
│  │  - Available                  │  │
│  │  - Reserved                   │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘
```

---

## 課題5: ドメインサービスの実装

### 目的
単一のエンティティに属さないドメインロジックを実装する。

### DDDにおけるドメインサービスの重要性
エリック・エヴァンスは、ドメインサービスを「**エンティティや値オブジェクトに自然に属さないドメインロジック**」を実装する場所と定義しています。
ドメインサービスの特徴：
- **ステートレス**（状態を持たない）
- **複数の集約を協調させる**役割を担う
- **ドメインロジックを実装する**（アプリケーション層の処理ではない）
- **エンティティのメソッドを呼び出す**ことでドメインロジックを再利用
- **過度な使用は避ける**（ドメインロジックは可能な限りエンティティや値オブジェクトに配置）

### 前提
この課題を始める前に、課題4-2で**Stockエンティティ**を実装済みであることを確認してください。
StockAllocationServiceは、Stockエンティティの`Reserve()`, `CanReserve()`, `Release()`メソッドを使用します。

### 課題内容

以下のドメインサービスを実装してください：

1. **在庫引当サービス（StockAllocationService）**
   - **複数の集約（OrderとStock）をまたがる**処理を担当
   - 注文に対して在庫を引き当てる
   - 在庫不足の場合は部分引当または引当失敗を返す
   - **重要**: Stockエンティティのメソッド（Reserve, CanReserveなど）を呼び出して引当処理を行う
   - **ドメインサービス自身は状態を変更しない**（エンティティに委譲する）
   - **ドメインサービスはビジネスルールを判断**する（どの在庫から引き当てるか、部分引当を許可するかなど）

2. **割引計算サービス（DiscountService）**（オプション課題）
   - 顧客ランク、購入金額、クーポンに基づいて割引額を計算
   - 複数の割引ルールを組み合わせる
   - **重要**: 複数のドメインオブジェクト（Order, Customer, DiscountPolicy）を協調させる
   - **割引計算のビジネスルール**をカプセル化する

### 実装場所
```
domain/
  stock_allocation_service.go
  stock_allocation_service_test.go
  discount_service.go         # オプション
  discount_service_test.go    # オプション
```

### StockAllocationService実装のヒント

```go
package domain

// StockAllocationService 在庫引当サービス
// 複数の集約（OrderとStock）をまたがるロジックを担当
type StockAllocationService struct {
    // ステートレス（状態を持たない）
}

// AllocationResult 引当結果
type AllocationResult struct {
    Success        bool
    AllocatedItems []AllocatedItem
    FailedItems    []ProductID
    Message        string
}

// Allocate 注文に対して在庫を引き当てる
// 引数:
//   - order: 引当対象の注文
//   - stocks: 商品ごとの在庫エンティティ（map[ProductID]*Stock）
// 戻り値:
//   - AllocationResult: 引当結果
//   - error: エラー
func (s *StockAllocationService) Allocate(
    order Order,
    stocks map[ProductID]*Stock,
) (AllocationResult, error) {
    // 1. 注文の各明細行をループ
    // 2. 各商品について:
    //    a. Stockエンティティが存在するか確認
    //    b. stock.CanReserve(quantity) で引当可能かチェック
    //    c. 引当可能なら stock.Reserve(quantity) で引当実行
    //    d. 引当不可なら失敗としてマーク
    //    e. 結果をAllocationResultに追加
    // 3. すべての商品が引当成功した場合のみSuccess=true
    // 4. 結果を返す

    // 注意: エラーが発生した場合の補償トランザクション（引当済みを戻す）は
    // アプリケーション層のトランザクション管理で対応
}

// 実装例:
func (s *StockAllocationService) Allocate(
    order Order,
    stocks map[ProductID]*Stock,
) (AllocationResult, error) {
    result := AllocationResult{
        Success: true,
        AllocatedItems: []AllocatedItem{},
        FailedItems: []ProductID{},
    }

    for _, line := range order.Lines() {
        productID := line.ProductID()
        quantity := line.Quantity()

        stock, exists := stocks[productID]
        if !exists {
            result.Success = false
            result.FailedItems = append(result.FailedItems, productID)
            continue
        }

        if !stock.CanReserve(quantity) {
            result.Success = false
            result.FailedItems = append(result.FailedItems, productID)
            continue
        }

        if err := stock.Reserve(quantity); err != nil {
            result.Success = false
            result.FailedItems = append(result.FailedItems, productID)
            continue
        }

        result.AllocatedItems = append(result.AllocatedItems, AllocatedItem{
            ProductID: productID,
            Quantity:  quantity,
        })
    }

    if result.Success {
        result.Message = "すべての商品の在庫引当に成功しました"
    } else {
        result.Message = fmt.Sprintf("一部の商品の在庫が不足しています: %v", result.FailedItems)
    }

    return result, nil
}
```

### StockAllocationServiceのテストケース例

```go
// stock_allocation_service_test.go
package domain_test

import (
    "testing"
    "yourproject/domain"
)

func TestStockAllocationService_Allocate_AllSuccess(t *testing.T) {
    // テストデータの準備
    service := &domain.StockAllocationService{}

    // 注文を作成（2つの商品）
    order := createTestOrder(
        orderLine("PROD-001", 10),
        orderLine("PROD-002", 5),
    )

    // 在庫を作成（十分な在庫あり）
    stocks := map[domain.ProductID]*domain.Stock{
        "PROD-001": createTestStock("PROD-001", 100),
        "PROD-002": createTestStock("PROD-002", 50),
    }

    // 引当実行
    result, err := service.Allocate(order, stocks)

    // 検証
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if !result.Success {
        t.Error("expected success")
    }
    if len(result.FailedItems) != 0 {
        t.Errorf("expected no failed items, got %d", len(result.FailedItems))
    }
    // Stockの状態も確認
    if stocks["PROD-001"].Available().Value() != 90 {
        t.Errorf("expected available 90, got %d", stocks["PROD-001"].Available().Value())
    }
}

func TestStockAllocationService_Allocate_PartialFailure(t *testing.T) {
    service := &domain.StockAllocationService{}

    // 注文を作成
    order := createTestOrder(
        orderLine("PROD-001", 10),
        orderLine("PROD-002", 100), // 在庫不足
    )

    // 在庫を作成（PROD-002は在庫不足）
    stocks := map[domain.ProductID]*domain.Stock{
        "PROD-001": createTestStock("PROD-001", 100),
        "PROD-002": createTestStock("PROD-002", 50), // 不足
    }

    // 引当実行
    result, err := service.Allocate(order, stocks)

    // 検証
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if result.Success {
        t.Error("expected failure")
    }
    if len(result.FailedItems) != 1 {
        t.Errorf("expected 1 failed item, got %d", len(result.FailedItems))
    }
    if result.FailedItems[0] != "PROD-002" {
        t.Error("expected PROD-002 to fail")
    }
}
```

### ポイント
- **ドメインサービスは状態を持たない**（ステートレス）
- **複数の集約をまたがるロジックに適している**
- **エンティティのメソッドを呼び出す**ことでドメインロジックを再利用
- ドメインサービス自身がビジネスルールを実装するのではなく、**エンティティを協調させる役割**
- 在庫の状態変更は必ずStockエンティティのメソッド経由で行う（カプセル化の維持）
- **テストでは、Stockエンティティの状態変化も確認**する

### アンチパターン: 貧血ドメインモデル
**避けるべきパターン**: ドメインサービスにすべてのロジックを集約し、エンティティをデータホルダーにしてしまうこと
```go
// ❌ 悪い例: 貧血ドメインモデル
type Order struct {
    ID     OrderID
    Status OrderStatus  // ただのデータホルダー
}

type OrderService struct {}
func (s *OrderService) ConfirmOrder(order *Order) {
    order.Status = Confirmed  // ビジネスロジックがサービスに集中
}

// ✅ 良い例: リッチドメインモデル
type Order struct {
    id     OrderID
    status OrderStatus
}

func (o *Order) Confirm() error {
    if o.status != Draft {
        return errors.New("can only confirm draft orders")
    }
    o.status = Confirmed  // ビジネスロジックはエンティティ内に
    return nil
}
```

---

## 課題6: リポジトリパターンの実装

### 目的
ドメインオブジェクトの永続化を抽象化し、インフラストラクチャ層への依存を分離する。

### DDDにおけるリポジトリの重要性
エリック・エヴァンスは、リポジトリを「**集約をカプセル化し、永続化の詳細を隠蔽するメカニズム**」と定義しています。
リポジトリの原則：
- **インターフェースはドメイン層に配置**（依存性逆転の原則）
- **実装はインフラストラクチャ層に配置**
- **集約ルート単位でリポジトリを作成**（OrderRepository, StockRepository）
- **コレクション指向のインターフェース**を提供（DBではなくメモリ上のコレクションのように扱える）
- **検索条件はドメインの言葉で表現**（SQLの詳細を隠蔽）

### 課題内容

1. **リポジトリインターフェースの定義（ドメイン層）**

   **OrderRepository**
   - CRUD操作のメソッドを持つ（Save, FindByID, FindByCustomerID, Delete）
   - contextを第一引数に取る
   - **仕様パターン**を使った検索メソッドも検討（例: FindBySpecification）

   **StockRepository**
   - CRUD操作のメソッドを持つ（Save, FindByProductID, FindByProductIDs, Delete）
   - 複数の商品の在庫を一括取得するメソッドも定義
   - contextを第一引数に取る

2. **リポジトリ実装（インフラストラクチャ層）**
   - インターフェースを実装する構造体を作成
   - データベース接続を依存として受け取る
   - トランザクション管理を考慮
   - **重要**: ドメインオブジェクトとDBエンティティのマッピングを適切に行う

### 実装場所
```
domain/
  order_repository.go
  stock_repository.go

infrastructure/
  persistence/
    order_repository.go
    order_repository_test.go
    stock_repository.go
    stock_repository_test.go
```

### StockRepository実装のヒント
```go
// domain/stock_repository.go
package domain

import "context"

// StockRepository 在庫リポジトリのインターフェース
type StockRepository interface {
    // Save 在庫を保存
    Save(ctx context.Context, stock *Stock) error

    // FindByProductID 商品IDで在庫を取得
    FindByProductID(ctx context.Context, productID ProductID) (*Stock, error)

    // FindByProductIDs 複数の商品IDで在庫を一括取得
    // 注文の引当処理で使用
    FindByProductIDs(ctx context.Context, productIDs []ProductID) (map[ProductID]*Stock, error)

    // Delete 在庫を削除
    Delete(ctx context.Context, productID ProductID) error
}
```

### ポイント
- **インターフェースはドメイン層に置く**（依存性逆転の原則）
  ```
  domain/           ← リポジトリインターフェース定義
    order_repository.go
  infrastructure/   ← リポジトリ実装
    persistence/
      order_repository_impl.go
  ```
- **実装はインフラストラクチャ層に置く**
- テストではインメモリ実装を使用
- **集約ごとにリポジトリを作成**（Order集約用、Stock集約用）
- リポジトリは集約ルート単位で定義する
- **集約全体を取得・永続化する**（集約内の部分的な操作は行わない）

### リポジトリとDAOの違い
- **DAO（Data Access Object）**: データベーステーブル単位でCRUD操作を提供（技術的な視点）
- **Repository**: 集約単位でドメインオブジェクトの永続化を抽象化（ドメインの視点）

```go
// ❌ DAOパターン（テーブル単位）
type OrderDAO interface {
    Insert(order *OrderTable) error
    Update(order *OrderTable) error
}

type OrderLineDAO interface {
    InsertBatch(lines []*OrderLineTable) error
}

// ✅ Repositoryパターン（集約単位）
type OrderRepository interface {
    Save(ctx context.Context, order *Order) error  // Order集約全体を保存
    FindByID(ctx context.Context, id OrderID) (*Order, error)  // Order集約全体を取得
}
```

---

## 課題7: ドメインイベントの実装

### 目的
集約間の整合性を疎結合に保ちながら実現する。

### DDDにおけるドメインイベントの重要性
エリック・エヴァンスは、ドメインイベントを「**ドメインエキスパートが関心を持つ、ドメイン内で起きた出来事**」と定義しています。
ドメインイベントの特徴：
- **過去形で命名**する（OrderConfirmed, PaymentCompleted）
- **不変（イミュータブル）**である
- **集約間の疎結合**を実現する
- **結果整合性**を実現するメカニズム
- **監査ログ**やイベントソーシングの基礎となる

### ドメインイベントの2つの用途
1. **同一境界づけられたコンテキスト内の整合性維持**（課題7で学習）
   - 例: 注文確定 → 在庫引当
2. **異なる境界づけられたコンテキスト間の連携**（課題9で学習）
   - 例: 注文コンテキスト → 配送コンテキスト

### 課題内容

以下のドメインイベントを実装してください：

1. **OrderConfirmed（注文確定）**
   - 注文が確定されたときに発行
   - 在庫引当処理をトリガー
   - **重要**: エンティティ（Order）内でイベントを生成し、保持する

2. **PaymentCompleted（支払完了）**
   - 支払いが完了したときに発行
   - 出荷準備処理をトリガー

3. **StockDepleted（在庫枯渇）**
   - 在庫が閾値を下回ったときに発行
   - 発注処理をトリガー
   - **重要**: ビジネスルールに基づいてイベントを発行（例: available < 10）

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

### 実装のヒント
```go
// domain/event.go
package domain

// DomainEvent インターフェース
// - OccurredAt() time.Time  // イベント発生日時
// - AggregateID() string    // 集約のID
// - EventType() string      // イベントの種類

// domain/order_confirmed.go
// OrderConfirmed 構造体
// - orderID, customerID, totalAmount, occurredAt をフィールドに持つ
// - NewOrderConfirmed コンストラクタで生成
// - DomainEvent インターフェースを実装

// エンティティ内でイベントを保持する例
type Order struct {
    id           OrderID
    status       OrderStatus
    domainEvents []DomainEvent  // イベントを保持
}

func (o *Order) Confirm() error {
    if o.status != Draft {
        return errors.New("can only confirm draft orders")
    }
    o.status = Confirmed

    // ドメインイベントを生成して保持
    event := NewOrderConfirmed(o.id, o.customerID, o.calculateTotal())
    o.domainEvents = append(o.domainEvents, event)

    return nil
}

// イベントを取得して発行（アプリケーション層で使用）
func (o *Order) DomainEvents() []DomainEvent {
    return o.domainEvents
}

func (o *Order) ClearDomainEvents() {
    o.domainEvents = nil
}
```

### ポイント
- **エンティティがドメインイベントを生成**する
- **イベントは不変**（すべてのフィールドを非公開にし、getterのみ提供）
- **アプリケーション層でイベントを発行**する（エンティティは生成のみ）
- **イベント駆動による疎結合**を実現

---

## 課題8: アプリケーションサービスの実装

### 目的
ユースケースを実装し、ドメインオブジェクトを協調させる。

### DDDにおけるアプリケーションサービスの重要性
エリック・エヴァンスの定義では、アプリケーション層は「**タスクを調整し、ドメインオブジェクトに作業を委譲する薄い層**」です。
アプリケーションサービスの責務：
- **ユースケースの実装**（ビジネスロジックではない）
- **トランザクション境界の管理**
- **ドメインオブジェクトの協調**
- **ドメインイベントの発行**
- **インフラストラクチャ層（リポジトリ等）との調整**

### アプリケーションサービスがやってはいけないこと
- **ビジネスロジックの実装**（それはドメイン層の責務）
- **状態の保持**（ステートレスであるべき）
- **ドメインオブジェクトの内部状態の直接操作**

### 課題内容

以下のユースケースを実装してください：

1. **注文作成ユースケース（CreateOrderUseCase）**
   - 顧客IDと商品リストから注文を作成
   - カートの内容を注文に変換
   - **重要**: ドメインの不変条件を保証（商品が存在するか、数量が有効か等）

2. **注文確定ユースケース（ConfirmOrderUseCase）**
   - 注文を確定状態に変更
   - 在庫を引き当てる
   - ドメインイベントを発行
   - **重要**: トランザクション管理（注文保存と在庫保存を同一トランザクションで）

3. **支払処理ユースケース（ProcessPaymentUseCase）**
   - 外部決済サービスと連携
   - 支払い結果を注文に反映
   - **重要**: 外部サービスとの連携は腐敗防止層を通す

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

### 実装のヒント
```go
// application/confirm_order.go
package application

import "context"

// ConfirmOrderUseCase ユースケース構造体
// 依存として以下を持つ:
type ConfirmOrderUseCase struct {
    orderRepo   domain.OrderRepository      // 注文の永続化
    stockRepo   domain.StockRepository      // 在庫の取得・永続化
    allocSvc    *domain.StockAllocationService  // 在庫引当サービス
    publisher   EventPublisher              // イベント発行
}

// ConfirmOrderInput 入力DTO
type ConfirmOrderInput struct {
    OrderID domain.OrderID
}

// ConfirmOrderOutput 出力DTO
type ConfirmOrderOutput struct {
    Success bool
    Message string
}

// Execute メソッドの処理フロー:
func (u *ConfirmOrderUseCase) Execute(ctx context.Context, input ConfirmOrderInput) (ConfirmOrderOutput, error) {
    // 1. 注文を取得
    order, err := u.orderRepo.FindByID(ctx, input.OrderID)

    // 2. 注文の商品IDリストを取得
    productIDs := extractProductIDs(order)

    // 3. 在庫を取得（StockRepository経由）
    stocks, err := u.stockRepo.FindByProductIDs(ctx, productIDs)

    // 4. 在庫を引き当て（StockAllocationService使用）
    result, err := u.allocSvc.Allocate(order, stocks)

    // 5. 注文を確定（ドメインのメソッド呼び出し）
    err = order.Confirm()

    // 6. 永続化（注文と在庫の両方）
    err = u.orderRepo.Save(ctx, order)
    for _, stock := range stocks {
        err = u.stockRepo.Save(ctx, stock)
    }

    // 7. ドメインイベント発行
    // 注意: order.Total()はMoney型を返すべき（値オブジェクトの一貫性）
    total, err := order.Total()
    if err != nil {
        return ConfirmOrderOutput{}, err
    }
    event := domain.NewOrderConfirmed(order.ID(), order.CustomerID(), total)
    u.publisher.Publish(ctx, event)

    return ConfirmOrderOutput{Success: true, Message: "注文を確定しました"}, nil
}

// 注意: Order.Total()の実装
// ❌ 悪い例: intを返す
// func (o Order) Total() int

// ✅ 良い例: Moneyを返す（値オブジェクトの一貫性を保つ）
// func (o Order) Total() (Money, error)
```

### 重要なポイント
- **アプリケーションサービスはトランザクション境界**を管理
- **複数のリポジトリ**（OrderRepositoryとStockRepository）を使用
- **ドメインサービス**（StockAllocationService）を呼び出してビジネスロジックを実行
- **エンティティのメソッド**を呼び出して状態を変更（直接フィールドを操作しない）
- **ドメインイベント**を発行して他のコンテキストに通知
- **DTOを使用**して入出力をドメインオブジェクトから分離

### アプリケーション層とドメイン層の責務分離
```go
// ❌ 悪い例: アプリケーション層にビジネスロジック
func (u *ConfirmOrderUseCase) Execute(ctx context.Context, input ConfirmOrderInput) error {
    order, _ := u.orderRepo.FindByID(ctx, input.OrderID)

    // ビジネスロジックがアプリケーション層に漏れている
    if order.Status != Draft {
        return errors.New("can only confirm draft orders")
    }
    order.Status = Confirmed  // 直接状態を変更
}

// ✅ 良い例: ビジネスロジックはドメイン層に
func (u *ConfirmOrderUseCase) Execute(ctx context.Context, input ConfirmOrderInput) error {
    order, _ := u.orderRepo.FindByID(ctx, input.OrderID)

    // ビジネスロジックはドメインのメソッドに委譲
    if err := order.Confirm(); err != nil {
        return err
    }

    // アプリケーション層は協調と永続化のみ
    u.orderRepo.Save(ctx, order)
}
```

---

## 課題9: 境界づけられたコンテキストの設計

### 目的
複数のコンテキスト間の関係を理解し、コンテキストマップを作成する。

### DDDにおける境界づけられたコンテキストの重要性
エリック・エヴァンスは、境界づけられたコンテキストを「**ユビキタス言語が適用される境界**」と定義しています。
これはDDDの戦略的設計の中核概念です。

境界づけられたコンテキストの原則：
- **コンテキスト内ではユビキタス言語が統一される**
- **コンテキスト間では言葉の意味が異なる可能性がある**
- **コンテキストは組織構造やチーム境界と対応することが多い**
- **技術的な境界（マイクロサービス等）とは必ずしも一致しない**
- **コンテキスト間の関係性を明示的に定義する**（コンテキストマップ）

### コンテキストマップのパターン
エリック・エヴァンスとヴォーン・ヴァーノンが定義した主要なパターン：
- **共有カーネル（Shared Kernel）**: 2つのチームが共有するドメインモデルの一部
- **顧客-供給者（Customer-Supplier）**: 上流チームと下流チームの関係
- **順応者（Conformist）**: 下流が上流のモデルをそのまま採用
- **腐敗防止層（Anti-Corruption Layer, ACL）**: 他のコンテキストのモデルから自コンテキストを保護
- **公開ホストサービス（Open Host Service）**: 統一されたAPIを提供
- **公表された言語（Published Language）**: 共通のデータ交換フォーマット

### 課題内容

ECサイトを以下のコンテキストに分割し、設計してください：

1. **注文コンテキスト（Order Context）**
   - 注文の作成・管理
   - **注文のための在庫引当**（StockAllocationService）
   - **注意**: 課題4-8で実装したOrderとStockは同じコンテキスト内に存在
   - 在庫は「注文における引当対象」として扱う

2. **在庫管理コンテキスト（Inventory Management Context）**
   - 在庫の入荷・補充管理
   - 在庫数の調整
   - 在庫の棚卸し
   - **注意**: 注文コンテキストのStockとは異なる視点（管理業務）

3. **支払コンテキスト（Payment Context）**
   - 決済処理
   - 支払い方法の管理

4. **配送コンテキスト（Shipping Context）**
   - 出荷・配送管理
   - 配送業者との連携

### 重要な設計判断

**Q: なぜOrderとStockは同じコンテキスト内なのか？**
- 注文における在庫引当は、注文プロセスの一部として密接に関連
- トランザクション境界が同じ（注文確定と在庫引当は同時に行われる必要がある）
- 課題4-8では**単一の境界づけられたコンテキスト内での戦術的設計**を学ぶ
- エリック・エヴァンスの定義では、「**強い整合性が必要な概念は同じコンテキスト内に配置**」する

**Q: 在庫管理コンテキストとの違いは？**
- 注文コンテキストの`Stock`: 「引当可能な在庫（available/reserved）」を表現
- 在庫管理コンテキストの`Inventory`: 「入荷・補充・棚卸し・ロケーション管理」を表現
- 同じ「在庫」という言葉でも、**コンテキストによって関心事が異なる**（ユビキタス言語の違い）
- これは境界づけられたコンテキストの核心概念: **コンテキストごとにモデルが異なる**

### 成果物

`docs/context_map.md` に以下を記載：

1. **コンテキストマップ**（図で表現）
2. **各コンテキストの責務とユビキタス言語**
3. **コンテキスト間の連携方法**
   - パターン名（共有カーネル、顧客-供給者、腐敗防止層など）を明記
   - なぜそのパターンを選択したかの理由
4. **各コンテキストでの「在庫」の意味の違い**を説明
5. **腐敗防止層の実装例**（外部システムとの連携時）

### コンテキストマップ例
```
┌──────────────────────────────────────────────────────────────────────┐
│                          ECサイト全体                                │
│                                                                      │
│   ┌────────────────────────────────┐                                │
│   │   注文コンテキスト             │  ← 課題1-8で実装する範囲       │
│   │  (Order Context)               │                                │
│   │ ┌────────┐ ┌──────┐           │                                │
│   │ │ Order  │ │Stock │           │  ← 同じコンテキスト内          │
│   │ └────────┘ └──────┘           │     (強い整合性)               │
│   │  StockAllocationService        │                                │
│   └────────────────────────────────┘                                │
│          │                  │                                       │
│          │ OrderConfirmed   │ StockDepleted                         │
│          │ (Pub/Sub)        │ (Pub/Sub)                             │
│          ▼                  ▼                                       │
│   ┌─────────────┐    ┌──────────────────┐                          │
│   │   支払      │    │  在庫管理        │  ← 課題9で設計する        │
│   │  Context    │    │  Context         │                          │
│   │             │    │ (Inventory       │                          │
│   │             │    │  Management)     │                          │
│   └─────────────┘    └──────────────────┘                          │
│          │                                                          │
│          │ PaymentCompleted (Pub/Sub)                               │
│          ▼                                                          │
│   ┌─────────────┐                                                  │
│   │   配送      │                                                  │
│   │  Context    │                                                  │
│   └─────────────┘                                                  │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘

【コンテキスト間の関係パターン】
- 注文 → 支払: **Pub/Sub (結果整合性)**
  - OrderConfirmedイベント経由で支払処理を開始
- 注文 → 在庫管理: **Pub/Sub (結果整合性)**
  - StockDepletedイベント経由で発注処理を開始
- 支払 → 配送: **Pub/Sub (結果整合性)**
  - PaymentCompletedイベント経由で出荷準備を開始

【各コンテキストでの「在庫」の意味（ユビキタス言語の違い）】
- **注文コンテキスト**の`Stock`:
  - 意味: 「引当可能な在庫」
  - 属性: available（利用可能数）, reserved（引当済み数）
  - 関心事: 注文に対して在庫を引き当てられるか
- **在庫管理コンテキスト**の`Inventory`:
  - 意味: 「管理対象の在庫」
  - 属性: location（保管場所）, expiryDate（賞味期限）, receivedDate（入荷日）
  - 関心事: 入荷、補充、棚卸し、ロケーション管理
```

### 重要な学習ポイント
1. **課題1-8は単一の境界づけられたコンテキスト内**での戦術的設計を学ぶ
2. **課題9で初めて複数コンテキスト間の関係**（戦略的設計）を学ぶ
3. **同じ言葉でもコンテキストごとに意味が異なる**のは正常（ユビキタス言語の境界）
4. **コンテキスト間は結果整合性**で連携（ドメインイベント経由）

### 腐敗防止層（ACL）の実装例
```go
// 外部決済サービスのレスポンス（外部モデル）
type ExternalPaymentResponse struct {
    TransactionID string
    StatusCode    int  // 0: success, 1: failure
    ErrorMsg      string
}

// 腐敗防止層: 外部モデルをドメインモデルに変換
type PaymentGatewayAdapter struct {
    externalClient *ExternalPaymentClient
}

func (a *PaymentGatewayAdapter) ProcessPayment(payment *Payment) (*PaymentResult, error) {
    // 外部サービスの呼び出し
    extResp, err := a.externalClient.Charge(payment.Amount(), payment.CardToken())

    // 外部モデルをドメインモデルに変換（腐敗防止層の役割）
    if extResp.StatusCode == 0 {
        return &PaymentResult{
            Success:       true,
            TransactionID: domain.TransactionID(extResp.TransactionID),
        }, nil
    } else {
        return &PaymentResult{
            Success: false,
            Error:   errors.New(extResp.ErrorMsg),
        }, nil
    }
}
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
│   ├── order/              # 注文集約（集約ルート + 関連エンティティ）
│   │   ├── order.go
│   │   ├── order_line.go
│   │   └── order_test.go
│   ├── stock.go            # 在庫エンティティ（Stock集約）
│   ├── stock_test.go
│   ├── money.go            # 値オブジェクト
│   ├── quantity.go
│   ├── email.go
│   ├── event.go            # ドメインイベントインターフェース
│   ├── order_confirmed.go
│   ├── payment_completed.go
│   ├── stock_depleted.go
│   ├── order_repository.go      # リポジトリインターフェース
│   ├── stock_repository.go      # 在庫リポジトリインターフェース
│   ├── stock_allocation_service.go  # ドメインサービス
│   └── discount_service.go
├── application/
│   ├── create_order.go
│   ├── confirm_order.go        # StockRepositoryを使用
│   ├── process_payment.go
│   └── order_event_handler.go
├── infrastructure/
│   ├── persistence/
│   │   ├── order_repository.go     # リポジトリ実装
│   │   ├── order_repository_test.go
│   │   ├── stock_repository.go     # 在庫リポジトリ実装
│   │   └── stock_repository_test.go
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
  - ❌ 悪い例: `domain/valueobjects/`, `domain/entities/`, `domain/services/`
  - ✅ 良い例: `domain/order/`, `domain/stock.go`, `domain/money.go`
  - エリック・エヴァンスは「**ドメインの概念を表現する構造**」を推奨
- **機能ドメイン単位（order等）でディレクトリを切る**
  - Goでは大きな集約の場合のみサブディレクトリを作成
  - 小さな集約や値オブジェクトはdomainパッケージ直下に配置
- **各集約ごとにエンティティを作成**（Order集約、Stock集約）
- **各集約ごとにリポジトリインターフェースを定義**
- **ドメインサービスは複数の集約を協調させる役割**（StockAllocationServiceはOrderとStockを協調）
- **課題1-8で実装したコードは単一の境界づけられたコンテキスト**内に収まる

### エリック・エヴァンスのレイヤードアーキテクチャ
DDDでは以下の4層構造を推奨：

1. **プレゼンテーション層（Interfaces）**: ユーザーとの対話、リクエスト/レスポンスの変換
2. **アプリケーション層（Application）**: ユースケースの実装、ドメインオブジェクトの協調
3. **ドメイン層（Domain）**: ビジネスロジック、エンティティ、値オブジェクト、集約、ドメインサービス
4. **インフラストラクチャ層（Infrastructure）**: 永続化、外部サービス連携、技術的詳細

**依存性の方向**: 外側から内側へ（プレゼンテーション → アプリケーション → ドメイン ← インフラストラクチャ）

---

## 学習の進め方

### 推奨スケジュール
1. **課題1-2**: DDDの基礎概念を理解する（1-2日）
   - ユビキタス言語の定義とコミュニケーション
   - 値オブジェクトによるドメイン概念の表現
2. **課題3-4**: エンティティと集約の設計を学ぶ（2-3日）
   - ライフサイクルと同一性の管理
   - トランザクション境界の設計
3. **課題5-6**: ドメインサービスとリポジトリを実装する（2-3日）
   - リッチドメインモデル vs 貧血ドメインモデル
   - 永続化の抽象化
4. **課題7-8**: イベント駆動とユースケースを実装する（3-4日）
   - 疎結合な設計
   - アプリケーション層とドメイン層の分離
5. **課題9-10**: 全体設計と統合（3-5日）
   - 戦略的設計
   - コンテキストマップの作成

### 学習のコツ
- **コードを書く前にドメインエキスパートとの対話をシミュレーション**する
  - エリック・エヴァンス: 「**ユビキタス言語がコードとドキュメントで一貫している**ことが重要」
- **ビジネスルールをコードで表現できているか常に確認**する
  - 貧血ドメインモデルに陥っていないかチェック
- **技術的な制約にとらわれず、まずドメインモデルを設計**する
  - データベーススキーマから設計を始めない
- **テストを書くことでドメインモデルの使いやすさを検証**する
  - テストが書きにくい = 設計に問題がある可能性
- **リファクタリングを恐れない**（より深い洞察を得たら設計を改善）
  - エリック・エヴァンス: 「**深いモデルは反復的なリファクタリングから生まれる**」
- **小さく始めて段階的に複雑さに対応**する
  - 最初から完璧なモデルを目指さない

---

## 参考資料

### 必読書籍
1. **「エリック・エヴァンスのドメイン駆動設計」（青本）**
   - DDDの原典、戦術的設計と戦略的設計の両方をカバー
   - この課題は青本の構成に沿って作成されています
2. **「実践ドメイン駆動設計」ヴォーン・ヴァーノン著（赤本）**
   - 実装例が豊富、集約の設計パターンが詳しい
3. **「ドメイン駆動設計入門」成瀬允宣著**
   - 日本語で読みやすい、初学者向け

### オンラインリソース
- [DDD Reference](https://www.domainlanguage.com/ddd/reference/) - エリック・エヴァンス公式
- [Domain-Driven Design Reference](https://www.domainlanguage.com/wp-content/uploads/2016/05/DDD_Reference_2015-03.pdf) - PDF版
- [Go + DDDのサンプル実装](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example)

### DDDコミュニティ
- [DDD Community](https://www.dddcommunity.org/)
- [Virtual DDD](https://virtualddd.com/) - オンラインミートアップ

---

## チェックリスト

各課題完了時に確認してください：

- [ ] 課題1: ユビキタス言語を定義した
  - [ ] ビジネスの言葉で記述している（技術用語を使っていない）
  - [ ] コードで使う用語と一致している
- [ ] 課題2: 値オブジェクト（Money, Quantity, Email, ProductID）を実装した（テスト含む）
  - [ ] すべて不変（イミュータブル）
  - [ ] 等価性で比較している
  - [ ] バリデーションロジックを含んでいる
  - [ ] 値オブジェクト同士の演算を実装している（Quantity.Add(other Quantity)など）
  - [ ] 計算結果も値オブジェクトの型で返している（Subtotal()やTotal()がMoneyを返す）
  - [ ] **注意**: Currency型のスペルが正しいか確認（Currensyではなく）
- [ ] 課題3: エンティティ（Order, OrderLine）を実装した（テスト含む）
  - [ ] 一意の識別子を持っている
  - [ ] ビジネスルールをカプセル化している
  - [ ] 状態遷移のルールを実装している
  - [ ] 状態遷移図を作成した
  - [ ] すべての状態遷移についてテストを書いた
- [ ] 課題4: 集約を設計・実装した
  - [ ] Order集約（Order + OrderLine）
  - [ ] Stock集約（Stockエンティティ）
  - [ ] 各集約の不変条件を保護している
  - [ ] 重複商品チェックの実装（同じProductIDの明細を2回追加できない）
  - [ ] 合計金額の上限チェックの実装
  - [ ] 集約ルート経由でのみアクセスしている
  - [ ] 集約間はIDで参照している
  - [ ] Stockエンティティのテストを書いた（Reserve, Release, CanReserve）
- [ ] 課題5: ドメインサービスを実装した
  - [ ] StockAllocationService（Stockエンティティのメソッドを使用）
  - [ ] DiscountService
  - [ ] ドメインサービスはステートレス
  - [ ] StockAllocationServiceのテストを書いた（成功ケース、部分失敗ケース）
  - [ ] Stockの状態変化もテストで確認している
- [ ] 課題6: リポジトリパターンを実装した
  - [ ] OrderRepositoryインターフェース（ドメイン層）
  - [ ] StockRepositoryインターフェース（ドメイン層）
  - [ ] リポジトリ実装（インフラストラクチャ層）
- [ ] 課題7: ドメインイベントを実装した
- [ ] 課題8: アプリケーションサービスを実装した
  - [ ] StockRepositoryを使用して在庫を取得
  - [ ] StockAllocationServiceを使用して引当処理
  - [ ] 複数のリポジトリを協調させている
  - [ ] Order.Total()がMoney型を返している（値オブジェクトの一貫性）
- [ ] 課題9: コンテキストマップを作成した
  - [ ] OrderとStockが同じコンテキスト内にある理由を理解
  - [ ] 各コンテキストでの「在庫」の意味の違いを説明できる
- [ ] 課題10: 総合演習を完了した

### DDDベストプラクティスの確認（エリック・エヴァンスの原則）

#### 戦術的設計パターン
- [ ] **値オブジェクト**: 不変で、等価性で比較している
- [ ] **エンティティ**: 一意の識別子を持ち、ライフサイクルを通じて同一性を保つ
- [ ] **エンティティは自身のビジネスルールをカプセル化**している（リッチドメインモデル）
- [ ] **集約**: トランザクション整合性の境界を明確に定義している
- [ ] **集約ルート経由でのみ集約内のエンティティにアクセス**している
- [ ] **集約はできるだけ小さく保つ**（パフォーマンスのため）
- [ ] **集約間はIDで参照**している（オブジェクト参照を避ける）
- [ ] **ドメインサービス**: ステートレスで、複数の集約を協調させる役割に徹している
- [ ] **リポジトリ**: インターフェースはドメイン層、実装はインフラ層に配置
- [ ] **ファクトリ**: 複雑なオブジェクト生成ロジックをカプセル化している

#### 戦略的設計パターン
- [ ] **ユビキタス言語**: コードとドキュメントで一貫した用語を使用
- [ ] **境界づけられたコンテキスト**: コンテキストの境界を明確に定義
- [ ] **コンテキストマップ**: コンテキスト間の関係性を明示的に文書化
- [ ] **腐敗防止層**: 外部システムから自コンテキストを保護

#### レイヤード アーキテクチャ
- [ ] **アプリケーションサービスがトランザクション境界を管理**している
- [ ] **ドメインロジックはドメイン層に集中**し、アプリケーション層には配置していない
- [ ] **依存性の方向**: 外側から内側へ（インフラ層はドメイン層のインターフェースに依存）

#### アンチパターンの回避
- [ ] **貧血ドメインモデルを避ける**: エンティティがデータホルダーになっていない
- [ ] **神オブジェクトを避ける**: 責務が適切に分散されている
- [ ] **技術的分類でディレクトリを切らない**: ドメイン単位で整理している


---
