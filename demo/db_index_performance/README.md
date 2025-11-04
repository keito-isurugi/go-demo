# データベースインデックスのパフォーマンス比較デモ

このプログラムは、データベースのインデックスの有無によるパフォーマンスの違いを実際に比較するデモです。

## インデックスとは

インデックスは、データベースのテーブルに対して作成される検索用の索引構造です。書籍の索引のように、特定のデータを高速に見つけるために使用されます。

**インデックスなし:**
```
全レコードをスキャン（フルテーブルスキャン）
→ データ量に比例して処理時間が増加
```

**インデックスあり:**
```
インデックスを使って効率的に検索
→ データ量が多くても高速にアクセス可能
```

## 実行方法

```bash
cd demo/db_index_performance
go run .
```

## ファイル構成

- [main.go](main.go) - メインプログラム（ベンチマーク実行）
- [models.go](models.go) - データベースモデルとサンプルデータ生成
- [queries.go](queries.go) - 各種検索クエリ
- [benchmark.go](benchmark.go) - ベンチマーク計測ユーティリティ

## データモデル

### 1. 商品テーブル（10,000件）
- **単一カラムインデックス**: category_id, price, stock, name, created_at

### 2. 注文テーブル（10,000件）
- **単一カラムインデックス**: user_id, status, order_date

### 3. 従業員テーブル（5,000件）
- **複合インデックス**: (first_name, last_name), (department, position)
- **単一カラムインデックス**: salary

## ベンチマーク内容

### 1. カテゴリで商品を検索
```sql
WHERE category_id = ?
```
- 単一カラムインデックスの効果を測定

### 2. 価格範囲で商品を検索
```sql
WHERE price BETWEEN ? AND ?
```
- 範囲検索でのインデックスの効果を測定

### 3. ユーザーIDで注文を検索
```sql
WHERE user_id = ?
```
- 外部キーへのインデックスの重要性を測定

### 4. ステータスで注文を検索
```sql
WHERE status = ?
```
- 文字列カラムへのインデックスの効果を測定

### 5. 名前で従業員を検索（複合インデックス）
```sql
WHERE first_name = ? AND last_name = ?
```
- 複合インデックスの効果を測定

### 6. 部署と役職で従業員を検索（複合インデックス）
```sql
WHERE department = ? AND position = ?
```
- 複数条件での複合インデックスの効果を測定

### 7. 価格でソート
```sql
ORDER BY price DESC
```
- ソート処理でのインデックスの効果を測定

### 8. 日付でソート
```sql
ORDER BY order_date DESC
```
- 日付カラムのソートでのインデックスの効果を測定

## 実行結果の例

### 顕著な改善が見られるケース

**ユーザーID検索:**
```
インデックスなし: 439µs
インデックスあり: 93µs
改善率: 4.72倍速くなりました
```

**名前検索（複合インデックス）:**
```
インデックスなし: 2.363ms
インデックスあり: 595µs
改善率: 3.97倍速くなりました
```

**日付ソート:**
```
インデックスなし: 874µs
インデックスあり: 261µs
改善率: 3.34倍速くなりました
```

## インデックスの種類

### 1. 単一カラムインデックス

1つのカラムにインデックスを作成

```go
type Product struct {
    CategoryID int `gorm:"index:idx_product_category"`
    Price      int `gorm:"index:idx_product_price"`
}
```

**用途:**
- WHERE句での単一条件検索
- ORDER BY句でのソート
- 外部キー

### 2. 複合インデックス

複数のカラムを組み合わせたインデックス

```go
type Employee struct {
    FirstName  string `gorm:"index:idx_employee_name,priority:1"`
    LastName   string `gorm:"index:idx_employee_name,priority:2"`
}
```

**用途:**
- 複数条件での検索
- カラムの順序が重要（左端のカラムから使用される）

**注意点:**
```sql
-- 複合インデックス: (first_name, last_name)

-- ✓ インデックスが使われる
WHERE first_name = 'John'
WHERE first_name = 'John' AND last_name = 'Smith'

-- ✗ インデックスが使われない
WHERE last_name = 'Smith'
```

## インデックスの効果

### 効果的なケース

1. **WHERE句での検索**
   - 特定の値での検索（`=`, `IN`）
   - 範囲検索（`BETWEEN`, `>`, `<`）

2. **ORDER BY句でのソート**
   - インデックスがあればソートを省略できる

3. **JOIN操作**
   - 結合キーにインデックスがあると高速化

4. **GROUP BY句での集約**
   - グループ化のキーにインデックスがあると効率的

### 効果的でないケース

1. **中間一致検索**
   ```sql
   -- ✗ インデックスが効かない
   WHERE name LIKE '%keyword%'

   -- ✓ インデックスが効く（前方一致）
   WHERE name LIKE 'keyword%'
   ```

2. **選択性が低いカラム**
   - データの種類が少ない（例: 性別、true/false）
   - 全体の大部分を占める値での検索

3. **小規模なテーブル**
   - 数百件程度ならフルスキャンの方が速い場合も

4. **関数を使った検索**
   ```sql
   -- ✗ インデックスが効かない
   WHERE YEAR(created_at) = 2024

   -- ✓ インデックスが効く
   WHERE created_at >= '2024-01-01' AND created_at < '2025-01-01'
   ```

## インデックスのデメリット

1. **ディスク容量の消費**
   - インデックスはテーブルとは別にデータを保持

2. **書き込み処理のオーバーヘッド**
   - INSERT/UPDATE/DELETE時にインデックスも更新が必要

3. **メンテナンスコスト**
   - インデックスの断片化が発生する可能性

## ベストプラクティス

### インデックスを作成すべきカラム

- WHERE句で頻繁に使用するカラム
- JOIN条件のカラム
- ORDER BY句で使用するカラム
- 外部キー
- 一意性制約が必要なカラム（UNIQUE INDEX）

### インデックスの設計

1. **選択性の高いカラムを選ぶ**
   - 値のバリエーションが多いカラム

2. **複合インデックスの順序**
   - よく使う条件を左に配置
   - 選択性の高いカラムを左に配置

3. **不要なインデックスは削除**
   - 使われていないインデックスは削除

4. **インデックスの数を適切に**
   - 多すぎると書き込み性能が低下

### 実行計画の確認

SQLiteでは `EXPLAIN QUERY PLAN` を使用：

```sql
EXPLAIN QUERY PLAN
SELECT * FROM products WHERE category_id = 2;
```

出力例：
```
-- インデックスなし
SCAN TABLE products

-- インデックスあり
SEARCH TABLE products USING INDEX idx_product_category (category_id=?)
```

## GORMでのインデックス定義

### 単一カラムインデックス

```go
type User struct {
    Email string `gorm:"index"`                    // デフォルト名
    Name  string `gorm:"index:idx_user_name"`      // カスタム名
}
```

### 複合インデックス

```go
type User struct {
    FirstName string `gorm:"index:idx_name,priority:1"`
    LastName  string `gorm:"index:idx_name,priority:2"`
}
```

### ユニークインデックス

```go
type User struct {
    Email string `gorm:"uniqueIndex"`
}
```

### 部分インデックス（条件付き）

```go
type User struct {
    DeletedAt *time.Time `gorm:"index:idx_deleted_at,where:deleted_at IS NULL"`
}
```

## まとめ

- **インデックスは検索・ソート処理を大幅に高速化**
- **適切に使わないと逆効果になる可能性もある**
- **WHERE句で頻繁に使用するカラムに作成**
- **複合インデックスはカラムの順序が重要**
- **書き込み処理のオーバーヘッドとのトレードオフ**
- **実行計画を確認して効果を測定**

## 関連リンク

- [GORM公式ドキュメント - Indexes](https://gorm.io/docs/indexes.html)
- [SQLite - Query Planning](https://www.sqlite.org/queryplanner.html)
