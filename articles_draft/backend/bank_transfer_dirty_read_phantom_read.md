## 概要

- Goで実装した銀行振込API
- データベースのトランザクション分離レベルによるDirty Read・Phantom Readの再現
- PostgreSQLを使用した実践的なトランザクション制御のデモ

### 特徴

- **銀行振込API**: 口座間の送金処理を実装
- **Dirty Read再現**: READ UNCOMMITTED分離レベルで未コミットデータの読み取りを再現
- **Phantom Read再現**: READ COMMITTED分離レベルでファントムレコードの出現を再現
- **トランザクション制御**: GORMとPostgreSQLを使った適切なトランザクション管理

## トランザクション分離レベルとは

### 基本概念

データベースのトランザクション分離レベルは、複数のトランザクションが同時実行される際のデータの整合性を保つための設定です。ANSI/ISO SQL標準では4つの分離レベルが定義されています。

**分離レベルの種類**:

1. **READ UNCOMMITTED**: 未コミットのデータも読み取れる（最も緩い）
2. **READ COMMITTED**: コミット済みのデータのみ読み取れる（デフォルト）
3. **REPEATABLE READ**: 同じトランザクション内で同じ読み取り結果を保証
4. **SERIALIZABLE**: 完全に分離され、直列化可能（最も厳しい）

### 分離レベルと発生する問題

| 分離レベル | Dirty Read | Non-Repeatable Read | Phantom Read |
|-----------|------------|---------------------|--------------|
| READ UNCOMMITTED | 発生する | 発生する | 発生する |
| READ COMMITTED | 発生しない | 発生する | 発生する |
| REPEATABLE READ | 発生しない | 発生しない | 発生する |
| SERIALIZABLE | 発生しない | 発生しない | 発生しない |

## Dirty Read（ダーティリード）とは

### 定義

Dirty Readは、**未コミットのデータを読み取ってしまう問題**です。トランザクションAがデータを更新している最中（まだコミットしていない）に、トランザクションBがそのデータを読み取ってしまうと、後でトランザクションAがロールバックした場合、トランザクションBが読み取ったデータは存在しないことになります。

### 発生シナリオ

```
時刻 | トランザクションA              | トランザクションB
-----|------------------------------|------------------
T1   | BEGIN                        |
T2   | UPDATE account SET balance = | 
     | 110000 WHERE id = 1          |
T3   |                              | BEGIN
T4   |                              | SELECT balance FROM account
     |                              | WHERE id = 1
     |                              | → 110000を読み取る（Dirty Read!）
T5   | ROLLBACK                     |
T6   |                              | COMMIT
```

この場合、トランザクションBは実際には存在しない残高110000を読み取ってしまっています。

### サンプルコード

```go
// DirtyReadDemoHandler Dirty Readを再現するデモ
func (h *BankTransferHandler) DirtyReadDemoHandler(w http.ResponseWriter, r *http.Request) {
	accountIDStr := r.URL.Query().Get("account_id")
	action := r.URL.Query().Get("action")

	if action == "update" {
		// トランザクション1: 残高を更新するがロールバックする
		tx1 := h.DB.Begin()
		
		// 分離レベルをREAD UNCOMMITTEDに設定（Dirty Readを許可）
		tx1.Exec("SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED")
		
		var account Account
		tx1.First(&account, accountID)
		
		// 残高を更新（まだコミットしていない）
		account.Balance += 10000
		tx1.Save(&account)

		// 5秒待機（この間に別のトランザクションがDirty Readを実行できる）
		time.Sleep(5 * time.Second)

		// ロールバック（更新を取り消す）
		tx1.Rollback()
	}

	if action == "read" {
		// トランザクション2: コミット前の未コミットデータを読み取る
		tx2 := h.DB.Begin()
		tx2.Exec("SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED")
		
		var account Account
		tx2.First(&account, accountID)

		// 少し待機（トランザクション1の更新を読み取る）
		time.Sleep(2 * time.Second)

		// 再度読み取り（トランザクション1がロールバックした後の状態）
		var accountAfter Account
		tx2.First(&accountAfter, accountID)

		tx2.Commit()
	}
}
```

### 実行方法

1. **口座を初期化**:
```bash
curl http://localhost:8080/api/bank/init
```

2. **別のターミナルで更新処理を開始**（ロールバックする）:
```bash
curl "http://localhost:8080/api/bank/dirty-read?account_id=1&action=update"
```

3. **別のターミナルで読み取り処理を実行**（未コミットデータを読み取る）:
```bash
curl "http://localhost:8080/api/bank/dirty-read?account_id=1&action=read"
```

結果として、`balance_during_transaction`と`balance_after_rollback`が異なる値になることで、Dirty Readが発生したことが確認できます。

**注意**: PostgreSQLでは`READ UNCOMMITTED`は実際には`READ COMMITTED`として扱われるため、真のDirty Readは発生しません。このコードは概念を示すためのものであり、実際のDirty Readを再現するにはMySQLなど別のデータベースが必要です。

## Phantom Read（ファントムリード）とは

### 定義

Phantom Readは、**同じクエリを2回実行した際に、1回目には存在しなかったレコードが2回目に現れる問題**です。トランザクションAが条件に一致するレコードを読み取っている間に、トランザクションBがその条件に一致する新しいレコードを挿入し、トランザクションAが再度同じクエリを実行すると、新しいレコードが「幻のように」現れます。

### 発生シナリオ

```
時刻 | トランザクションA                    | トランザクションB
-----|--------------------------------------|------------------
T1   | BEGIN                               |
T2   | SELECT * FROM accounts              |
     | WHERE balance >= 50000              |
     | → 2件取得                           |
T3   |                                      | BEGIN
T4   |                                      | INSERT INTO accounts 
     |                                      | (account_no, balance, owner_name)
     |                                      | VALUES ('2001', 150000, 'New User')
T5   |                                      | COMMIT
T6   | SELECT * FROM accounts              |
     | WHERE balance >= 50000              |
     | → 3件取得（Phantom Read!）          |
T7   | COMMIT                              |
```

この場合、トランザクションAは同じクエリを2回実行したにもかかわらず、2回目に新しいレコードが現れています。

### サンプルコード

```go
// PhantomReadDemoHandler Phantom Readを再現するデモ
func (h *BankTransferHandler) PhantomReadDemoHandler(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")

	if action == "read" {
		// トランザクション1: 条件に一致するレコードを読み取る
		tx1 := h.DB.Begin()
		
		// 分離レベルをREAD COMMITTEDに設定（Phantom Readが発生する可能性がある）
		tx1.Exec("SET TRANSACTION ISOLATION LEVEL READ COMMITTED")
		
		minBalance := 50000
		
		// 最初の読み取り
		var accounts1 []Account
		tx1.Where("balance >= ?", minBalance).Find(&accounts1)

		// 5秒待機（この間に別のトランザクションが新しいレコードを挿入できる）
		time.Sleep(5 * time.Second)

		// 再度読み取り（Phantom Readが発生する可能性がある）
		var accounts2 []Account
		tx1.Where("balance >= ?", minBalance).Find(&accounts2)

		tx1.Commit()
	}

	if action == "insert" {
		// トランザクション2: 条件に一致する新しいレコードを挿入
		tx2 := h.DB.Begin()
		tx2.Exec("SET TRANSACTION ISOLATION LEVEL READ COMMITTED")
		
		// 新しい口座を作成（残高が高い）
		newAccount := Account{
			AccountNo: fmt.Sprintf("PHANTOM-%d", time.Now().Unix()),
			Balance:   150000,
			OwnerName: "Phantom User",
		}

		tx2.Create(&newAccount)
		tx2.Commit()
	}
}
```

### 実行方法

1. **口座を初期化**:
```bash
curl http://localhost:8080/api/bank/init
```

2. **別のターミナルで読み取り処理を開始**:
```bash
curl "http://localhost:8080/api/bank/phantom-read?action=read&min_balance=50000"
```

3. **別のターミナルで新しい口座を挿入**:
```bash
curl "http://localhost:8080/api/bank/phantom-read?action=insert"
```

結果として、`first_read_count`と`second_read_count`が異なる値になることで、Phantom Readが発生したことが確認できます。

## 銀行振込APIの実装

### データモデル

```go
// Account 口座モデル
type Account struct {
	ID        uint      `gorm:"primaryKey"`
	AccountNo string    `gorm:"uniqueIndex;size:20;not null"` // 口座番号
	Balance   int64     `gorm:"not null;default:0"`           // 残高（単位: 円）
	OwnerName string    `gorm:"size:100;not null"`            // 口座名義人
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Transaction 取引履歴モデル
type Transaction struct {
	ID              uint      `gorm:"primaryKey"`
	FromAccountID   uint      `gorm:"index;not null"` // 送金元口座ID
	ToAccountID     uint      `gorm:"index;not null"` // 送金先口座ID
	Amount          int64     `gorm:"not null"`       // 送金額
	Status          string    `gorm:"size:20;not null;default:'pending'"` // pending, completed, failed
	TransactionType string    `gorm:"size:20;not null"`                   // transfer, deposit, withdrawal
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}
```

### 通常の振込処理

```go
// NormalTransferHandler 通常の振込処理（トランザクション分離レベル: READ COMMITTED）
func (h *BankTransferHandler) NormalTransferHandler(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest
	json.NewDecoder(r.Body).Decode(&req)

	// トランザクション開始（デフォルトの分離レベル: READ COMMITTED）
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 送金元口座を取得
	var fromAccount Account
	tx.First(&fromAccount, req.FromAccountID)

	// 残高チェック
	if fromAccount.Balance < req.Amount {
		tx.Rollback()
		// エラーレスポンスを返す
		return
	}

	// 送金先口座を取得
	var toAccount Account
	tx.First(&toAccount, req.ToAccountID)

	// 残高を更新
	fromAccount.Balance -= req.Amount
	toAccount.Balance += req.Amount
	tx.Save(&fromAccount)
	tx.Save(&toAccount)

	// 取引履歴を作成
	transaction := Transaction{
		FromAccountID:   req.FromAccountID,
		ToAccountID:     req.ToAccountID,
		Amount:          req.Amount,
		Status:          "completed",
		TransactionType: "transfer",
	}
	tx.Create(&transaction)

	// トランザクションコミット
	tx.Commit()
}
```

### APIエンドポイント

| エンドポイント | メソッド | 説明 |
|--------------|---------|------|
| `/api/bank/init` | GET | テスト用の口座を初期化 |
| `/api/bank/transfer` | POST | 通常の振込処理 |
| `/api/bank/account` | GET | 口座情報を取得（`account_id`パラメータ必要） |
| `/api/bank/accounts` | GET | 全口座一覧を取得 |
| `/api/bank/dirty-read` | GET | Dirty Readデモ（`action`と`account_id`パラメータ必要） |
| `/api/bank/phantom-read` | GET | Phantom Readデモ（`action`パラメータ必要） |

### 使用例

#### 1. 口座の初期化

```bash
curl http://localhost:8080/api/bank/init
```

レスポンス:
```json
{
  "success": true,
  "message": "Accounts initialized successfully",
  "accounts": [
    {
      "ID": 1,
      "AccountNo": "1001",
      "Balance": 100000,
      "OwnerName": "山田太郎"
    },
    {
      "ID": 2,
      "AccountNo": "1002",
      "Balance": 50000,
      "OwnerName": "佐藤花子"
    },
    {
      "ID": 3,
      "AccountNo": "1003",
      "Balance": 200000,
      "OwnerName": "鈴木一郎"
    }
  ]
}
```

#### 2. 振込処理

```bash
curl -X POST http://localhost:8080/api/bank/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": 10000
  }'
```

レスポンス:
```json
{
  "success": true,
  "message": "Transfer completed successfully",
  "transaction_id": 1,
  "from_balance": 90000,
  "to_balance": 60000
}
```

#### 3. 口座情報の取得

```bash
curl "http://localhost:8080/api/bank/account?account_id=1"
```

## 分離レベルの選択指針

### 各分離レベルの使い分け

1. **READ UNCOMMITTED**
   - 使用場面: 統計情報の取得など、多少の不整合が許容される場合
   - 注意点: Dirty Readが発生するため、金融取引などでは使用しない

2. **READ COMMITTED**（PostgreSQLのデフォルト）
   - 使用場面: 一般的なWebアプリケーション
   - 注意点: Phantom Readが発生する可能性がある

3. **REPEATABLE READ**
   - 使用場面: 同じトランザクション内で複数回読み取る必要がある場合
   - 注意点: デッドロックのリスクが増加する可能性がある

4. **SERIALIZABLE**
   - 使用場面: 金融取引など、完全な整合性が要求される場合
   - 注意点: パフォーマンスが低下する可能性がある

### PostgreSQLでの設定方法

```go
// トランザクション開始時に分離レベルを設定
tx := db.Begin()
tx.Exec("SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED")
// または
tx.Exec("SET TRANSACTION ISOLATION LEVEL READ COMMITTED")
tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")
tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
```

**PostgreSQLの注意点**:
- PostgreSQLでは`READ UNCOMMITTED`は実際には`READ COMMITTED`として扱われます
- そのため、PostgreSQLでは真のDirty Readは発生しません
- Phantom Readは`READ COMMITTED`レベルで発生する可能性があります
- `REPEATABLE READ`では、PostgreSQLの実装によりPhantom Readも防がれます（ANSI標準では防がれないが、PostgreSQLの実装では防がれる）

## まとめ

### 重要なポイント

1. **Dirty Read**: 未コミットのデータを読み取ってしまう問題。READ UNCOMMITTEDで発生する。
2. **Phantom Read**: 同じクエリを2回実行した際に、新しいレコードが現れる問題。READ COMMITTEDでも発生する。
3. **分離レベルの選択**: アプリケーションの要件に応じて適切な分離レベルを選択する。
4. **トランザクション管理**: GORMの`Begin()`と`Commit()`/`Rollback()`を適切に使用する。

### 実践的なアドバイス

- 金融取引など、整合性が重要な場合は`SERIALIZABLE`を使用
- 一般的なWebアプリケーションでは`READ COMMITTED`（デフォルト）で十分
- パフォーマンスと整合性のバランスを考慮して選択
- デッドロックを避けるため、トランザクションの実行時間を短く保つ

### 参考資料

- [PostgreSQL Documentation - Transaction Isolation](https://www.postgresql.org/docs/current/transaction-iso.html)
- [ANSI/ISO SQL Standard - Isolation Levels](https://en.wikipedia.org/wiki/Isolation_(database_systems))
- [GORM Documentation - Transactions](https://gorm.io/docs/transactions.html)

