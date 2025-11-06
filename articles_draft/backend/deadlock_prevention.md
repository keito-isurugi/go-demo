## 概要

- Goで実装した銀行振込APIにおけるデッドロックの再現と回避策
- PostgreSQLを使用した実践的なデッドロック対策のデモ
- ロック順序の統一、タイムアウト設定、リトライロジックなどの実装

### 特徴

- **デッドロック再現**: 2つのトランザクションが相互にロックを待つ状況を再現
- **回避策1: ロック順序の統一**: 常に同じ順序でロックを取得することでデッドロックを防止
- **回避策2: タイムアウト設定**: タイムアウトとリトライロジックでデッドロックを検出・回復
- **実践的な実装**: 銀行振込APIを例にした実際のユースケース

## デッドロックとは

### 基本概念

デッドロック（Deadlock）は、2つ以上のトランザクションが互いに相手が保持しているロックを待つ状態になり、どちらも進行できなくなる問題です。

### デッドロックの発生条件

デッドロックが発生するには、以下の4つの条件がすべて満たされる必要があります（Coffman条件）：

1. **相互排他（Mutual Exclusion）**: リソースが同時に複数のプロセスで使用できない
2. **保持と待機（Hold and Wait）**: プロセスがリソースを保持しながら、他のリソースを待つ
3. **非先取り（No Preemption）**: リソースを強制的に奪うことができない
4. **循環待機（Circular Wait）**: プロセス間で循環的な待機関係が発生

### デッドロックの例

銀行振込の例でデッドロックを説明します：

```
時刻 | トランザクションA              | トランザクションB
-----|------------------------------|------------------
T1   | BEGIN                        | BEGIN
T2   | SELECT * FROM accounts       |
     | WHERE id = 1 FOR UPDATE      |
     | → 口座1をロック               |
T3   |                              | SELECT * FROM accounts
     |                              | WHERE id = 2 FOR UPDATE
     |                              | → 口座2をロック
T4   | SELECT * FROM accounts       |
     | WHERE id = 2 FOR UPDATE      |
     | → 口座2をロックしようとする   |
     | → トランザクションBが待機     |
T5   |                              | SELECT * FROM accounts
     |                              | WHERE id = 1 FOR UPDATE
     |                              | → 口座1をロックしようとする
     |                              | → トランザクションAが待機
T6   | デッドロック発生！           | デッドロック発生！
```

この場合、トランザクションAは口座2を待ち、トランザクションBは口座1を待つため、どちらも進行できません。

## デッドロックを再現するサンプルコード

### 実装

```go
// DeadlockDemoHandler デッドロックを再現するデモ
func (h *BankTransferHandler) DeadlockDemoHandler(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	account1Str := r.URL.Query().Get("account1")
	account2Str := r.URL.Query().Get("account2")

	if action == "tx1" {
		// トランザクション1: 口座1をロック → 口座2をロックしようとする
		tx1 := h.DB.Begin()
		
		// 口座1をロック
		var account1 Account
		tx1.Raw("SELECT * FROM accounts WHERE id = ? FOR UPDATE", account1ID).Scan(&account1)

		// 少し待機（トランザクション2が口座2をロックする時間を与える）
		time.Sleep(2 * time.Second)

		// 口座2をロックしようとする（ここでデッドロックが発生する可能性がある）
		var account2 Account
		if err := tx1.Raw("SELECT * FROM accounts WHERE id = ? FOR UPDATE", account2ID).Scan(&account2).Error; err != nil {
			// デッドロックが検出された
			tx1.Rollback()
			// エラーレスポンスを返す
		}
	}

	if action == "tx2" {
		// トランザクション2: 口座2をロック → 口座1をロックしようとする
		tx2 := h.DB.Begin()
		
		// 口座2をロック
		var account2 Account
		tx2.Raw("SELECT * FROM accounts WHERE id = ? FOR UPDATE", account2ID).Scan(&account2)

		// 少し待機（トランザクション1が口座1をロックする時間を与える）
		time.Sleep(2 * time.Second)

		// 口座1をロックしようとする（ここでデッドロックが発生する可能性がある）
		var account1 Account
		if err := tx2.Raw("SELECT * FROM accounts WHERE id = ? FOR UPDATE", account1ID).Scan(&account1).Error; err != nil {
			// デッドロックが検出された
			tx2.Rollback()
			// エラーレスポンスを返す
		}
	}
}
```

### 実行方法

1. **口座を初期化**:
```bash
curl http://localhost:8080/api/bank/init
```

2. **2つのターミナルで同時に実行**:
```bash
# ターミナル1: トランザクション1（口座1 → 口座2）
curl "http://localhost:8080/api/bank/deadlock?action=tx1&account1=1&account2=2"

# ターミナル2: トランザクション2（口座2 → 口座1）
curl "http://localhost:8080/api/bank/deadlock?action=tx2&account1=1&account2=2"
```

結果として、どちらかのトランザクションがデッドロックエラーを検出し、ロールバックされます。PostgreSQLは自動的にデッドロックを検出し、一方のトランザクションをロールバックします。

## デッドロックの回避策

### 回避策1: ロック順序の統一

最も効果的なデッドロック回避策は、**常に同じ順序でロックを取得する**ことです。これにより、循環待機を防ぐことができます。

#### 実装

```go
// DeadlockAvoidanceHandler デッドロック回避策1: ロック順序の統一
func (h *BankTransferHandler) DeadlockAvoidanceHandler(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest
	json.NewDecoder(r.Body).Decode(&req)

	// ロック順序の統一: 常にIDの小さい順にロックを取得
	firstAccountID := req.FromAccountID
	secondAccountID := req.ToAccountID
	if firstAccountID > secondAccountID {
		firstAccountID, secondAccountID = secondAccountID, firstAccountID
	}

	tx := h.DB.Begin()

	// 1つ目の口座をロック（常にIDの小さい方）
	var firstAccount Account
	tx.Raw("SELECT * FROM accounts WHERE id = ? FOR UPDATE", firstAccountID).Scan(&firstAccount)

	// 2つ目の口座をロック（常にIDの大きい方）
	var secondAccount Account
	tx.Raw("SELECT * FROM accounts WHERE id = ? FOR UPDATE", secondAccountID).Scan(&secondAccount)

	// 残高を更新
	// ...

	tx.Commit()
}
```

#### なぜこれでデッドロックが防げるのか

ロック順序を統一することで、すべてのトランザクションが同じ順序でロックを取得します：

```
トランザクションA: 口座1（ID=1）→ 口座2（ID=2）
トランザクションB: 口座2（ID=2）→ 口座1（ID=1）
```

これを以下のように統一します：

```
トランザクションA: 口座1（ID=1）→ 口座2（ID=2）
トランザクションB: 口座1（ID=1）→ 口座2（ID=2）
```

この場合、トランザクションBは口座1をロックしようとしますが、トランザクションAが既にロックしているため待機します。トランザクションAが完了すると、トランザクションBが続行できます。循環待機が発生しないため、デッドロックは発生しません。

#### 使用例

```bash
curl -X POST http://localhost:8080/api/bank/transfer-safe \
  -H "Content-Type: application/json" \
  -d '{
    "from_account_id": 2,
    "to_account_id": 1,
    "amount": 10000
  }'
```

この場合、`from_account_id=2`と`to_account_id=1`ですが、内部では常にIDの小さい順（1 → 2）でロックを取得するため、デッドロックが発生しません。

### 回避策2: タイムアウト設定とリトライ

デッドロックが発生した場合にタイムアウトで検出し、自動的にリトライする方法です。

#### 実装

```go
// DeadlockTimeoutHandler デッドロック回避策2: タイムアウト設定
func (h *BankTransferHandler) DeadlockTimeoutHandler(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest
	json.NewDecoder(r.Body).Decode(&req)

	maxRetries := 3
	retryDelay := 100 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		// タイムアウト付きコンテキストを作成
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// トランザクション開始
		tx := h.DB.WithContext(ctx).Begin()
		
		// ロック順序の統一
		firstAccountID := req.FromAccountID
		secondAccountID := req.ToAccountID
		if firstAccountID > secondAccountID {
			firstAccountID, secondAccountID = secondAccountID, firstAccountID
		}

		// 1つ目の口座をロック
		var firstAccount Account
		if err := tx.Raw("SELECT * FROM accounts WHERE id = ? FOR UPDATE", firstAccountID).Scan(&firstAccount).Error; err != nil {
			tx.Rollback()
			if ctx.Err() == context.DeadlineExceeded {
				// タイムアウト: リトライ
				log.Printf("Transaction timeout, retrying... (attempt %d/%d)", attempt+1, maxRetries)
				time.Sleep(retryDelay)
				retryDelay *= 2 // 指数バックオフ
				continue
			}
			// その他のエラー
			return
		}

		// 2つ目の口座をロック
		// ...

		// 残高を更新
		// ...

		// コミット成功
		tx.Commit()
		return
	}

	// 最大リトライ回数に達した
	http.Error(w, "Transaction failed after maximum retries", http.StatusInternalServerError)
}
```

#### 指数バックオフ

リトライの間隔を徐々に増やすことで、システムへの負荷を軽減します：

```go
retryDelay := 100 * time.Millisecond
for attempt := 0; attempt < maxRetries; attempt++ {
	// ...
	if timeout {
		time.Sleep(retryDelay)
		retryDelay *= 2 // 100ms → 200ms → 400ms
		continue
	}
}
```

#### 使用例

```bash
curl -X POST http://localhost:8080/api/bank/transfer-timeout \
  -H "Content-Type: application/json" \
  -d '{
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": 10000
  }'
```

## PostgreSQLのデッドロック検出

PostgreSQLは自動的にデッドロックを検出し、一方のトランザクションをロールバックします。

### デッドロック検出の仕組み

1. **デッドロック検出プロセス**: PostgreSQLは定期的に（通常は1秒ごと）デッドロックを検出します
2. **犠牲者の選択**: デッドロックが検出されると、コストが最も低いトランザクションがロールバックされます
3. **エラーメッセージ**: ロールバックされたトランザクションには`deadlock detected`エラーが返されます

### エラーハンドリング

```go
if err := tx.Raw("SELECT * FROM accounts WHERE id = ? FOR UPDATE", accountID).Scan(&account).Error; err != nil {
	if strings.Contains(err.Error(), "deadlock detected") {
		// デッドロックが検出された
		tx.Rollback()
		// リトライまたはエラーレスポンスを返す
		return
	}
	// その他のエラー
	return
}
```

## その他のデッドロック回避策

### 1. トランザクションの実行時間を短くする

トランザクションが長く実行されるほど、デッドロックが発生する可能性が高くなります。

```go
// 悪い例: トランザクション内で長時間の処理
tx := db.Begin()
// データ取得
// 複雑な計算処理（時間がかかる）
// データ更新
tx.Commit()

// 良い例: トランザクション外で処理
// データ取得
// 複雑な計算処理（トランザクション外）
tx := db.Begin()
// データ更新
tx.Commit()
```

### 2. 必要なロックのみを取得する

不要なロックを取得しないことで、デッドロックのリスクを減らします。

```go
// 悪い例: すべてのレコードをロック
tx.Where("status = ?", "active").Find(&accounts) // ロック不要

// 良い例: 必要なレコードのみをロック
tx.Raw("SELECT * FROM accounts WHERE id = ? FOR UPDATE", accountID).Scan(&account)
```

### 3. インデックスの適切な使用

インデックスを使用することで、ロックの範囲を狭め、デッドロックのリスクを減らします。

```sql
-- インデックスがある場合: 特定の行のみをロック
CREATE INDEX idx_account_id ON accounts(id);
SELECT * FROM accounts WHERE id = 1 FOR UPDATE; -- 行ロックのみ

-- インデックスがない場合: テーブル全体をスキャンしてロック
SELECT * FROM accounts WHERE id = 1 FOR UPDATE; -- より広い範囲をロック
```

### 4. 分離レベルの調整

分離レベルを下げることで、ロックの競合を減らすことができます（ただし、整合性のリスクが増加します）。

```go
tx.Exec("SET TRANSACTION ISOLATION LEVEL READ COMMITTED")
```

## デバッグとモニタリング

### PostgreSQLのデッドロックログ

PostgreSQLのログでデッドロックを確認できます：

```sql
-- デッドロックの検出ログを確認
SELECT * FROM pg_stat_database WHERE datname = 'go_demo';
```

### アプリケーションレベルのモニタリング

```go
// デッドロックの発生をログに記録
if err := tx.Commit().Error; err != nil {
	if strings.Contains(err.Error(), "deadlock detected") {
		log.Printf("Deadlock detected: %v", err)
		// メトリクスを記録
		metrics.IncrementDeadlockCount()
	}
	return err
}
```

## まとめ

### 重要なポイント

1. **デッドロックの原因**: 複数のトランザクションが異なる順序でロックを取得することで発生
2. **回避策1: ロック順序の統一**: 最も効果的で推奨される方法
3. **回避策2: タイムアウトとリトライ**: デッドロックが発生した場合の回復手段
4. **PostgreSQLの自動検出**: PostgreSQLは自動的にデッドロックを検出し、一方をロールバック

### 実践的なアドバイス

- **ロック順序の統一を最優先**: すべてのトランザクションで同じ順序でロックを取得する
- **トランザクションの実行時間を短く**: トランザクション内での処理を最小限に
- **適切なエラーハンドリング**: デッドロックエラーを検出し、リトライまたは適切なエラーレスポンスを返す
- **モニタリング**: デッドロックの発生頻度を監視し、問題があれば対策を検討

### 参考資料

- [PostgreSQL Documentation - Deadlocks](https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS)
- [Database Deadlocks Explained](https://www.postgresql.org/docs/current/explicit-locking.html)
- [GORM Documentation - Transactions](https://gorm.io/docs/transactions.html)

