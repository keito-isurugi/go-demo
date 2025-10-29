## 概要

- キャッシュアルゴリズムの比較と実装
- LRU、LFU、Write-Through、Write-Backの4つの主要アルゴリズムを実装
- それぞれの特徴、メリット・デメリット、適切なユースケースを理解する

### 特徴

- **LRU (Least Recently Used)**: 最も最近使われていないデータを削除
- **LFU (Least Frequently Used)**: 最も使用頻度が低いデータを削除
- **Write-Through**: 書き込み時にキャッシュとDBを同期
- **Write-Back (Write-Behind)**: 書き込みをキャッシュに行い、非同期でDBに反映
- **スレッドセーフ**: すべての実装でmutexを使用した並行処理対応

## キャッシュとは

キャッシュは頻繁にアクセスされるデータを高速なストレージ（メモリ）に一時保存し、データ取得を高速化する技術。

### キャッシュの役割

1. **パフォーマンス向上**: データベースやディスクアクセスを削減
2. **レスポンス時間短縮**: メモリからの高速読み込み
3. **システム負荷軽減**: バックエンドへのリクエスト数を削減
4. **スケーラビリティ**: より多くのリクエストを処理可能

### キャッシュの課題

1. **容量制限**: メモリは有限なので削除戦略が必要
2. **データ一貫性**: キャッシュとDBの同期
3. **キャッシュ無効化**: いつデータを削除・更新するか

## 1. LRU (Least Recently Used) Cache

### 概要

最も最近使われていない（アクセスされていない）データを削除するアルゴリズム。

### アルゴリズムの特徴

- **時間的局所性**: 最近アクセスされたデータは近い将来も参照される可能性が高い
- **実装**: 双方向連結リスト + ハッシュマップ
- **計算量**: Get/Put ともに O(1)

### データ構造

```go
type LRUCache struct {
    capacity int                          // 最大容量
    cache    map[string]*list.Element     // キー -> リスト要素
    list     *list.List                   // 使用順序を管理（最新が先頭）
    mu       sync.RWMutex                 // スレッドセーフ
}

type lruEntry struct {
    key   string
    value interface{}
}
```

### 動作の流れ

#### Get操作

1. キーがキャッシュに存在するか確認
2. 存在する場合、そのエントリをリストの先頭に移動（最新としてマーク）
3. 値を返す

#### Put操作

1. キーが既に存在する場合、値を更新してリストの先頭に移動
2. 新規キーの場合、リストの先頭に追加
3. 容量超過の場合、リストの末尾（最も古い）を削除

### 実装のポイント

```go
// Get はキーに対応する値を取得（アクセス時に最新に移動）
func (c *LRUCache) Get(key string) (interface{}, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if elem, found := c.cache[key]; found {
        // アクセスされたので最新に移動
        c.list.MoveToFront(elem)
        return elem.Value.(*lruEntry).value, true
    }
    return nil, false
}

// Put はキーと値をキャッシュに保存
func (c *LRUCache) Put(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if elem, found := c.cache[key]; found {
        // 既存のエントリを更新
        c.list.MoveToFront(elem)
        elem.Value.(*lruEntry).value = value
        return
    }

    // 新規エントリを追加
    entry := &lruEntry{key: key, value: value}
    elem := c.list.PushFront(entry)
    c.cache[key] = elem

    // 容量超過なら最も古いものを削除
    if c.list.Len() > c.capacity {
        oldest := c.list.Back()
        if oldest != nil {
            c.list.Remove(oldest)
            delete(c.cache, oldest.Value.(*lruEntry).key)
        }
    }
}
```

### デモの動作例

```
容量3のLRUキャッシュを作成

データを追加:
  Put(a, Apple) -> Keys: [a]
  Put(b, Banana) -> Keys: [b a]
  Put(c, Cherry) -> Keys: [c b a]

データアクセス:
  Get(a) = Apple -> Keys: [a c b]  # aが最新に移動

容量超過（最も古いbが削除される）:
  Put(d, Date) -> Keys: [d a c]    # bが削除される

削除されたキーへのアクセス:
  Get(b) = not found

最終的なキャッシュの状態: [d a c] (サイズ: 3)
```

### メリット

- **実装がシンプル**: 標準ライブラリの連結リストを使用
- **O(1)の性能**: 高速な読み書き
- **直感的**: 人間の記憶モデルに近い
- **適用範囲が広い**: 多くのケースで有効

### デメリット

- **使用パターンに依存**: 頻度よりも時間を重視
- **スキャンに弱い**: 大量の一時アクセスでキャッシュが汚染される

### 適切なユースケース

- **Webページキャッシュ**: 最近閲覧したページ
- **ブラウザキャッシュ**: 最近アクセスしたリソース
- **データベースクエリキャッシュ**: 最近実行されたクエリ結果
- **OSページキャッシュ**: メモリページの管理

## 2. LFU (Least Frequently Used) Cache

### 概要

最も使用頻度が低いデータを削除するアルゴリズム。

### アルゴリズムの特徴

- **頻度ベース**: 何回アクセスされたかを記録
- **実装**: ハッシュマップ + 頻度ごとの連結リスト
- **計算量**: Get/Put ともに O(1)

### データ構造

```go
type LFUCache struct {
    capacity  int                           // 最大容量
    cache     map[string]*lfuEntry          // キー -> エントリ
    freqMap   map[int]*list.List            // 頻度 -> リスト
    minFreq   int                           // 現在の最小頻度
    mu        sync.RWMutex                  // スレッドセーフ
}

type lfuEntry struct {
    key   string
    value interface{}
    freq  int                               // 使用頻度
    elem  *list.Element                     // リスト内の位置
}
```

### 動作の流れ

#### Get操作

1. キーがキャッシュに存在するか確認
2. 存在する場合、使用頻度を1増加
3. 頻度リストを更新（古い頻度リストから削除、新しい頻度リストに追加）
4. 値を返す

#### Put操作

1. キーが既に存在する場合、値を更新して頻度を増加
2. 新規キーの場合
   - 容量超過なら最小頻度のリストから最も古いものを削除
   - 新しいエントリを頻度1のリストに追加
3. 最小頻度を更新

### 実装のポイント

```go
// incrementFreq はエントリの頻度を増加
func (c *LFUCache) incrementFreq(entry *lfuEntry) {
    freq := entry.freq

    // 現在の頻度リストから削除
    c.freqMap[freq].Remove(entry.elem)
    if c.freqMap[freq].Len() == 0 {
        delete(c.freqMap, freq)
        if c.minFreq == freq {
            c.minFreq++
        }
    }

    // 新しい頻度リストに追加
    entry.freq++
    if c.freqMap[entry.freq] == nil {
        c.freqMap[entry.freq] = list.New()
    }
    entry.elem = c.freqMap[entry.freq].PushFront(entry)
}

// evict は最も頻度の低い要素を削除
func (c *LFUCache) evict() {
    if freqList := c.freqMap[c.minFreq]; freqList != nil {
        oldest := freqList.Back()
        if oldest != nil {
            entry := oldest.Value.(*lfuEntry)
            freqList.Remove(oldest)
            delete(c.cache, entry.key)
        }
    }
}
```

### デモの動作例

```
容量3のLFUキャッシュを作成

データを追加:
  Put(a, Apple) -> Size: 1, Freq(a): 1
  Put(b, Banana) -> Size: 2, Freq(b): 1
  Put(c, Cherry) -> Size: 3, Freq(c): 1

データアクセス（頻度を増加）:
  Get(a) -> Freq(a): 2
  Get(a) -> Freq(a): 3
  Get(b) -> Freq(b): 2

容量超過（最も頻度の低いcが削除される）:
  現在の頻度 - a: 3, b: 2, c: 1
  Put(d, Date) -> Size: 3

削除されたキーへのアクセス:
  Get(c) = not found

最終的な頻度 - a: 3, b: 2, d: 1
```

### メリット

- **長期的な人気度**: よく使われるデータを保持
- **スキャン耐性**: 一時的な大量アクセスに強い
- **O(1)の性能**: 高速な読み書き

### デメリット

- **実装が複雑**: 頻度管理が必要
- **古いエントリの滞留**: 過去に頻繁だったが今は使われないデータが残る
- **コールドスタート**: 新しいエントリが削除されやすい

### 適切なユースケース

- **動画配信**: 人気コンテンツの優先
- **CDN**: アクセス頻度の高いコンテンツ
- **データベースバッファプール**: 頻繁にアクセスされるページ
- **DNSキャッシュ**: よく解決されるドメイン

## 3. Write-Through Cache

### 概要

書き込み時にキャッシュとデータベースを同時に更新するアルゴリズム。

### アルゴリズムの特徴

- **同期書き込み**: キャッシュとDBを同時に更新
- **高い一貫性**: キャッシュとDBが常に同期
- **書き込み性能**: DBの書き込み待ちが発生

### データ構造

```go
type WriteThroughCache struct {
    cache map[string]interface{}     // キャッシュ
    db    map[string]interface{}     // 模擬DB
    mu    sync.RWMutex               // スレッドセーフ
}
```

### 動作の流れ

#### Read操作

1. キャッシュを確認（キャッシュヒット）
2. キャッシュにない場合、DBから取得（キャッシュミス）
3. 値を返す

#### Write操作

1. **DBに書き込み**（時間がかかる）
2. **キャッシュに書き込み**
3. 両方成功したら完了

### 実装のポイント

```go
// Write はデータを書き込む（キャッシュとDBの両方に書き込み）
func (c *WriteThroughCache) Write(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()

    // DBに書き込み（時間がかかる処理をシミュレート）
    time.Sleep(10 * time.Millisecond)
    c.db[key] = value
    fmt.Printf("  [DB Write] key=%s\n", key)

    // キャッシュに書き込み
    c.cache[key] = value
    fmt.Printf("  [Cache Write] key=%s\n", key)
}

// Read はデータを読み込む（キャッシュ優先）
func (c *WriteThroughCache) Read(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    // キャッシュから取得
    if value, found := c.cache[key]; found {
        fmt.Printf("  [Cache Hit] key=%s\n", key)
        return value, true
    }

    // DBから取得
    if value, found := c.db[key]; found {
        fmt.Printf("  [Cache Miss] key=%s, loading from DB\n", key)
        return value, true
    }

    return nil, false
}
```

### デモの動作例

```
Write-Throughキャッシュを作成
（書き込み時にキャッシュとDBを同期）

データを書き込み:
  [DB Write] key=user:1
  [Cache Write] key=user:1
  所要時間: 10ms (キャッシュ+DB書き込み)

データを読み込み（キャッシュヒット）:
  [Cache Hit] key=user:1
  Read(user:1) = Alice
  所要時間: 0ms (キャッシュから高速取得)

複数のデータを書き込み:
  [DB Write] key=user:2
  [Cache Write] key=user:2
  [DB Write] key=user:3
  [Cache Write] key=user:3
  所要時間: 20ms (2回の書き込み)

最終的な状態 - Cache: 3, DB: 3

特徴:
  ✓ 書き込み時にキャッシュとDBを同時更新
  ✓ データの一貫性が高い
  ✗ 書き込みが遅い（DBの書き込み待ち）
```

### メリット

- **データの一貫性**: キャッシュとDBが常に同期
- **読み込みが高速**: キャッシュヒット時は即座に返す
- **シンプルな実装**: 書き込みロジックが明確
- **データロスなし**: 障害時もDBにデータが保存済み

### デメリット

- **書き込みが遅い**: DBの書き込み完了を待つ
- **書き込みスループット低下**: 同期処理がボトルネック
- **DBへの負荷**: すべての書き込みがDBに到達

### 適切なユースケース

- **金融システム**: データの一貫性が最重要
- **在庫管理**: 正確な在庫数が必要
- **予約システム**: ダブルブッキング防止
- **ユーザー認証**: セキュリティ情報の整合性

## 4. Write-Back (Write-Behind) Cache

### 概要

書き込みをキャッシュに行い、非同期でデータベースに反映するアルゴリズム。

### アルゴリズムの特徴

- **非同期書き込み**: キャッシュに書き込み後、バックグラウンドでDBに反映
- **高速な書き込み**: DB待ちがない
- **データロストリスク**: クラッシュ時にダーティデータが失われる可能性

### データ構造

```go
type WriteBackCache struct {
    cache      map[string]interface{}   // キャッシュ
    db         map[string]interface{}   // 模擬DB
    dirtyKeys  map[string]bool          // 変更されたキー
    mu         sync.RWMutex             // スレッドセーフ
    flushChan  chan string              // フラッシュキュー
    stopChan   chan struct{}            // 停止シグナル
}
```

### 動作の流れ

#### Read操作

1. キャッシュを確認（キャッシュヒット）
2. キャッシュにない場合、DBから取得（キャッシュミス）
3. 値を返す

#### Write操作

1. **キャッシュに書き込み**（高速）
2. ダーティフラグを設定
3. フラッシュキューに追加
4. すぐに完了を返す

#### バックグラウンドフラッシュ

1. 定期的にダーティなキーをチェック
2. DBに非同期で書き込み
3. ダーティフラグをクリア

### 実装のポイント

```go
// Write はデータを書き込む（キャッシュにのみ書き込み、後でDBに反映）
func (c *WriteBackCache) Write(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()

    // キャッシュに書き込み（高速）
    c.cache[key] = value
    c.dirtyKeys[key] = true
    fmt.Printf("  [Cache Write] key=%s (dirty)\n", key)

    // フラッシュキューに追加
    select {
    case c.flushChan <- key:
    default:
        // キューが満杯の場合は同期的にフラッシュ
        c.flushKey(key)
    }
}

// flushWorker はバックグラウンドでダーティなキーをDBに書き込む
func (c *WriteBackCache) flushWorker(interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            c.flushAll()
        case key := <-c.flushChan:
            // 少し遅延させてバッチ処理
            time.Sleep(50 * time.Millisecond)
            c.mu.Lock()
            c.flushKey(key)
            c.mu.Unlock()
        case <-c.stopChan:
            return
        }
    }
}

// flushKey は特定のキーをDBに書き込む
func (c *WriteBackCache) flushKey(key string) {
    if !c.dirtyKeys[key] {
        return
    }

    if value, found := c.cache[key]; found {
        // DBに書き込み（時間がかかる処理をシミュレート）
        time.Sleep(10 * time.Millisecond)
        c.db[key] = value
        delete(c.dirtyKeys, key)
        fmt.Printf("  [DB Flush] key=%s\n", key)
    }
}
```

### デモの動作例

```
Write-Backキャッシュを作成
（書き込みをキャッシュに行い、非同期でDBに反映）

データを書き込み:
  [Cache Write] key=user:1 (dirty)
  所要時間: 0ms (キャッシュのみ、高速)

現在の状態 - Cache: 1, DB: 0, Dirty: 1

データを読み込み（キャッシュヒット）:
  [Cache Hit] key=user:1
  Read(user:1) = Alice

複数のデータを書き込み:
  [Cache Write] key=user:2 (dirty)
  [Cache Write] key=user:3 (dirty)
  所要時間: 0ms (2回の書き込み、すべて高速)

現在の状態 - Cache: 3, DB: 0, Dirty: 3

バックグラウンドでDBに書き込み中...
  [DB Flush] key=user:1
  [DB Flush] key=user:2
  [DB Flush] key=user:3

最終的な状態 - Cache: 3, DB: 3, Dirty: 0

特徴:
  ✓ 書き込みが高速（キャッシュのみ）
  ✓ DBへの書き込みを非同期で実行
  ✗ データロストのリスク（クラッシュ時）
  ✗ データの一貫性に遅延がある
```

### メリット

- **高速な書き込み**: DBの書き込み待ちがない
- **書き込みスループット向上**: バッチ処理で効率化
- **DB負荷軽減**: 複数の書き込みをまとめて実行
- **書き込み集約**: 同じキーへの複数書き込みを1回に

### デメリット

- **データロストリスク**: クラッシュ時にダーティデータが失われる
- **一貫性の遅延**: キャッシュとDBの不一致期間が存在
- **実装が複雑**: バックグラウンドワーカーの管理
- **メモリ使用量増加**: ダーティフラグの管理

### 適切なユースケース

- **ログ書き込み**: 多少のデータロストは許容可能
- **カウンター**: アクセスカウント、いいね数など
- **分析データ**: リアルタイム性は不要
- **セッション管理**: 短時間の一時データ

## アルゴリズムの比較

### 削除戦略の比較

| アルゴリズム | 削除基準 | 計算量 | 実装難易度 |
|----------|--------|--------|----------|
| **LRU** | 最も古くアクセス | O(1) | 中 |
| **LFU** | 最も使用頻度が低い | O(1) | 高 |

### 書き込み戦略の比較

| アルゴリズム | 書き込み速度 | データ一貫性 | データロストリスク | 実装難易度 |
|----------|----------|-----------|--------------|----------|
| **Write-Through** | 遅い | 高い | なし | 低 |
| **Write-Back** | 高速 | 遅延あり | あり | 高 |

### 性能比較

```
読み込み性能（キャッシュヒット時）:
  すべて O(1) で高速

書き込み性能:
  Write-Through: ~10ms (DB書き込み待ち)
  Write-Back:    ~0ms  (キャッシュのみ)

削除戦略の効果:
  LRU: 時間的局所性が高い場合に有効
  LFU: アクセスパターンが一定の場合に有効
```

## 実装のベストプラクティス

### 1. スレッドセーフ

```go
// 読み取りはRLock、書き込みはLockを使用
func (c *LRUCache) Get(key string) (interface{}, bool) {
    c.mu.Lock()  // Getでも書き込み（リストの移動）があるのでLock
    defer c.mu.Unlock()
    // ...
}
```

### 2. エラーハンドリング

```go
// キャッシュミス時の処理
if value, found := cache.Get(key); found {
    return value, nil
}

// DBから取得してキャッシュに保存
value, err := fetchFromDB(key)
if err != nil {
    return nil, err
}
cache.Put(key, value)
return value, nil
```

### 3. TTL（有効期限）の追加

```go
type entry struct {
    key       string
    value     interface{}
    expiresAt time.Time
}

func (c *Cache) isExpired(entry *entry) bool {
    return time.Now().After(entry.expiresAt)
}
```

### 4. メトリクス収集

```go
type CacheMetrics struct {
    Hits   int64
    Misses int64
    Evictions int64
}

func (c *Cache) GetHitRate() float64 {
    total := c.metrics.Hits + c.metrics.Misses
    if total == 0 {
        return 0
    }
    return float64(c.metrics.Hits) / float64(total)
}
```

### 5. グレースフルシャットダウン

```go
func (c *WriteBackCache) Stop() {
    close(c.stopChan)  // ワーカーに停止を通知
    c.flushAll()       // すべてのダーティデータをフラッシュ
}
```

## 実際のプロダクトでの使用例

### 1. Redis

- **デフォルト**: LRU with approximate algorithm
- **オプション**: LFU, volatile-lru, allkeys-random など
- **設定**: `maxmemory-policy allkeys-lru`

### 2. Memcached

- **削除戦略**: LRU
- **書き込み**: Write-Through（明示的な削除が必要）

### 3. CPU Cache

- **L1/L2/L3キャッシュ**: 主にLRU系のアルゴリズム
- **Write-Back**: L1キャッシュで使用

### 4. OS Page Cache

- **Linux**: LRU with two lists (active/inactive)
- **Write-Back**: ページキャッシュからディスクへの書き込み

### 5. Webブラウザ

- **ブラウザキャッシュ**: LRU
- **Service Worker Cache**: プログラマブルなキャッシュ戦略

## 適切なアルゴリズムの選び方

### LRUを選ぶべきとき

- ✓ 時間的局所性が高い（最近のデータがよく使われる）
- ✓ 実装をシンプルに保ちたい
- ✓ 汎用的なキャッシュが必要
- 例: Webページキャッシュ、セッション管理

### LFUを選ぶべきとき

- ✓ 頻度が重要（人気コンテンツを保持）
- ✓ スキャン攻撃への耐性が必要
- ✓ アクセスパターンが比較的安定
- 例: CDN、動画配信、DNSキャッシュ

### Write-Throughを選ぶべきとき

- ✓ データの一貫性が最重要
- ✓ データロストが許容できない
- ✓ 読み込みが多く書き込みが少ない
- 例: 金融システム、在庫管理、認証

### Write-Backを選ぶべきとき

- ✓ 書き込み性能が最重要
- ✓ 多少のデータロストは許容可能
- ✓ 書き込みが非常に多い
- 例: ログ、カウンター、分析データ

## 複合戦略

実際のシステムでは複数のアルゴリズムを組み合わせることが多い：

### 1. LRU + TTL

```go
// 一定時間後に自動削除されるLRU
type LRUWithTTL struct {
    *LRUCache
    ttl time.Duration
}
```

### 2. Write-Back + Periodic Flush

```go
// 定期的にフラッシュするWrite-Back
// ダーティ期間を制限してデータロストリスクを軽減
```

### 3. 2-Level Cache (L1/L2)

```go
// L1: Write-Back (高速)
// L2: Write-Through (永続性)
```

### 4. Adaptive Replacement Cache (ARC)

- LRUとLFUを動的に切り替え
- ワークロードに自動適応

## まとめ

### 学んだこと

1. **LRU**: 時間ベースの削除戦略、実装がシンプル
2. **LFU**: 頻度ベースの削除戦略、人気コンテンツを保持
3. **Write-Through**: 同期書き込み、一貫性が高い
4. **Write-Back**: 非同期書き込み、書き込みが高速

### キャッシュ設計の原則

1. **ワークロードを理解**: アクセスパターンを分析
2. **トレードオフを考慮**: 性能 vs 一貫性 vs 複雑さ
3. **測定して最適化**: 推測せず実測する
4. **エラーハンドリング**: キャッシュ障害時の対応

### 実装のポイント

- **スレッドセーフ**: mutexで並行アクセスを保護
- **O(1)性能**: ハッシュマップと連結リストの組み合わせ
- **グレースフルシャットダウン**: ダーティデータのフラッシュ
- **メトリクス**: ヒット率やエビクション数を追跡

### 実際の開発での応用

キャッシュアルゴリズムはあらゆる層で使われている：

- **アプリケーション層**: Redis, Memcached
- **データベース層**: バッファプール、クエリキャッシュ
- **OS層**: ページキャッシュ、ファイルシステムキャッシュ
- **ハードウェア層**: CPUキャッシュ（L1/L2/L3）

適切なアルゴリズムを選択し、正しく実装することで、システムのパフォーマンスを劇的に向上させることができます。
