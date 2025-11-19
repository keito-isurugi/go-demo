package algorithm

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// ============================================================
// LRU (Least Recently Used) Cache
// ============================================================

// LRUCache は最も最近使われていない要素を削除するキャッシュ
type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
	mu       sync.RWMutex
}

type lruEntry struct {
	key   string
	value interface{}
}

// NewLRUCache はLRUキャッシュを作成
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

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

// Size は現在のキャッシュサイズを返す
func (c *LRUCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.list.Len()
}

// Keys は現在のキーを新しい順に返す
func (c *LRUCache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, c.list.Len())
	for elem := c.list.Front(); elem != nil; elem = elem.Next() {
		keys = append(keys, elem.Value.(*lruEntry).key)
	}
	return keys
}

// ============================================================
// LFU (Least Frequently Used) Cache
// ============================================================

// LFUCache は最も使用頻度の低い要素を削除するキャッシュ
type LFUCache struct {
	capacity  int
	cache     map[string]*lfuEntry
	freqMap   map[int]*list.List // 頻度ごとのリスト
	minFreq   int
	mu        sync.RWMutex
}

type lfuEntry struct {
	key   string
	value interface{}
	freq  int
	elem  *list.Element
}

// NewLFUCache はLFUキャッシュを作成
func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		capacity: capacity,
		cache:    make(map[string]*lfuEntry),
		freqMap:  make(map[int]*list.List),
		minFreq:  0,
	}
}

// Get はキーに対応する値を取得（アクセス頻度を増加）
func (c *LFUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, found := c.cache[key]
	if !found {
		return nil, false
	}

	// 頻度を増加
	c.incrementFreq(entry)
	return entry.value, true
}

// Put はキーと値をキャッシュに保存
func (c *LFUCache) Put(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.capacity == 0 {
		return
	}

	if entry, found := c.cache[key]; found {
		// 既存のエントリを更新
		entry.value = value
		c.incrementFreq(entry)
		return
	}

	// 容量超過なら最も頻度の低いものを削除
	if len(c.cache) >= c.capacity {
		c.evict()
	}

	// 新規エントリを追加
	entry := &lfuEntry{
		key:   key,
		value: value,
		freq:  1,
	}

	if c.freqMap[1] == nil {
		c.freqMap[1] = list.New()
	}
	entry.elem = c.freqMap[1].PushFront(entry)
	c.cache[key] = entry
	c.minFreq = 1
}

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

// Size は現在のキャッシュサイズを返す
func (c *LFUCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

// GetFrequency はキーの使用頻度を返す
func (c *LFUCache) GetFrequency(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if entry, found := c.cache[key]; found {
		return entry.freq
	}
	return 0
}

// ============================================================
// Write-Through Cache
// ============================================================

// WriteThroughCache は書き込み時にキャッシュとDBを同期するキャッシュ
type WriteThroughCache struct {
	cache map[string]interface{}
	db    map[string]interface{} // 模擬DB
	mu    sync.RWMutex
}

// NewWriteThroughCache はWrite-Throughキャッシュを作成
func NewWriteThroughCache() *WriteThroughCache {
	return &WriteThroughCache{
		cache: make(map[string]interface{}),
		db:    make(map[string]interface{}),
	}
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

// GetCacheSize はキャッシュのサイズを返す
func (c *WriteThroughCache) GetCacheSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

// GetDBSize はDBのサイズを返す
func (c *WriteThroughCache) GetDBSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.db)
}

// ============================================================
// Write-Back (Write-Behind) Cache
// ============================================================

// WriteBackCache は書き込みをキャッシュに行い、非同期でDBに反映するキャッシュ
type WriteBackCache struct {
	cache      map[string]interface{}
	db         map[string]interface{} // 模擬DB
	dirtyKeys  map[string]bool        // 変更されたキー
	mu         sync.RWMutex
	flushChan  chan string
	stopChan   chan struct{}
}

// NewWriteBackCache はWrite-Backキャッシュを作成
func NewWriteBackCache(flushInterval time.Duration) *WriteBackCache {
	c := &WriteBackCache{
		cache:     make(map[string]interface{}),
		db:        make(map[string]interface{}),
		dirtyKeys: make(map[string]bool),
		flushChan: make(chan string, 100),
		stopChan:  make(chan struct{}),
	}

	// バックグラウンドでフラッシュ処理
	go c.flushWorker(flushInterval)

	return c
}

// Read はデータを読み込む（キャッシュ優先）
func (c *WriteBackCache) Read(key string) (interface{}, bool) {
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

// flushKey は特定のキーをDBに書き込む（ロックは呼び出し元が取得済み）
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

// flushAll はすべてのダーティなキーをDBに書き込む
func (c *WriteBackCache) flushAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.dirtyKeys {
		c.flushKey(key)
	}
}

// Flush は明示的にすべてのダーティなキーをDBに書き込む
func (c *WriteBackCache) Flush() {
	c.flushAll()
}

// Stop はバックグラウンドワーカーを停止
func (c *WriteBackCache) Stop() {
	close(c.stopChan)
	c.flushAll()
}

// GetCacheSize はキャッシュのサイズを返す
func (c *WriteBackCache) GetCacheSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

// GetDBSize はDBのサイズを返す
func (c *WriteBackCache) GetDBSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.db)
}

// GetDirtyCount はダーティなキーの数を返す
func (c *WriteBackCache) GetDirtyCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.dirtyKeys)
}

// ============================================================
// デモ関数
// ============================================================

// DemoCacheAlgorithms はすべてのキャッシュアルゴリズムのデモを実行
func DemoCacheAlgorithms() {
	fmt.Println("==============================================")
	fmt.Println("キャッシュアルゴリズムの比較")
	fmt.Println("==============================================")

	demoLRU()
	fmt.Println()

	demoLFU()
	fmt.Println()

	demoWriteThrough()
	fmt.Println()

	demoWriteBack()
}

// demoLRU はLRUキャッシュのデモ
func demoLRU() {
	fmt.Println("1. LRU (Least Recently Used) Cache")
	fmt.Println("----------------------------------------------")

	cache := NewLRUCache(3)

	fmt.Println("容量3のLRUキャッシュを作成")
	fmt.Println()

	// データ追加
	fmt.Println("データを追加:")
	cache.Put("a", "Apple")
	fmt.Printf("  Put(a, Apple) -> Keys: %v\n", cache.Keys())

	cache.Put("b", "Banana")
	fmt.Printf("  Put(b, Banana) -> Keys: %v\n", cache.Keys())

	cache.Put("c", "Cherry")
	fmt.Printf("  Put(c, Cherry) -> Keys: %v\n", cache.Keys())

	fmt.Println()

	// データアクセス（aが最新に移動）
	fmt.Println("データアクセス:")
	if value, found := cache.Get("a"); found {
		fmt.Printf("  Get(a) = %v -> Keys: %v\n", value, cache.Keys())
	}

	fmt.Println()

	// 容量超過（最も古いbが削除される）
	fmt.Println("容量超過（最も古いbが削除される）:")
	cache.Put("d", "Date")
	fmt.Printf("  Put(d, Date) -> Keys: %v\n", cache.Keys())

	fmt.Println()

	// 削除されたキーへのアクセス
	fmt.Println("削除されたキーへのアクセス:")
	if _, found := cache.Get("b"); !found {
		fmt.Println("  Get(b) = not found")
	}

	fmt.Println()
	fmt.Printf("最終的なキャッシュの状態: %v (サイズ: %d)\n", cache.Keys(), cache.Size())
}

// demoLFU はLFUキャッシュのデモ
func demoLFU() {
	fmt.Println("2. LFU (Least Frequently Used) Cache")
	fmt.Println("----------------------------------------------")

	cache := NewLFUCache(3)

	fmt.Println("容量3のLFUキャッシュを作成")
	fmt.Println()

	// データ追加
	fmt.Println("データを追加:")
	cache.Put("a", "Apple")
	fmt.Printf("  Put(a, Apple) -> Size: %d, Freq(a): %d\n", cache.Size(), cache.GetFrequency("a"))

	cache.Put("b", "Banana")
	fmt.Printf("  Put(b, Banana) -> Size: %d, Freq(b): %d\n", cache.Size(), cache.GetFrequency("b"))

	cache.Put("c", "Cherry")
	fmt.Printf("  Put(c, Cherry) -> Size: %d, Freq(c): %d\n", cache.Size(), cache.GetFrequency("c"))

	fmt.Println()

	// データアクセス（頻度を増加）
	fmt.Println("データアクセス（頻度を増加）:")
	cache.Get("a")
	fmt.Printf("  Get(a) -> Freq(a): %d\n", cache.GetFrequency("a"))

	cache.Get("a")
	fmt.Printf("  Get(a) -> Freq(a): %d\n", cache.GetFrequency("a"))

	cache.Get("b")
	fmt.Printf("  Get(b) -> Freq(b): %d\n", cache.GetFrequency("b"))

	fmt.Println()

	// 容量超過（最も頻度の低いcが削除される）
	fmt.Println("容量超過（最も頻度の低いcが削除される）:")
	fmt.Printf("  現在の頻度 - a: %d, b: %d, c: %d\n",
		cache.GetFrequency("a"),
		cache.GetFrequency("b"),
		cache.GetFrequency("c"))

	cache.Put("d", "Date")
	fmt.Printf("  Put(d, Date) -> Size: %d\n", cache.Size())

	fmt.Println()

	// 削除されたキーへのアクセス
	fmt.Println("削除されたキーへのアクセス:")
	if _, found := cache.Get("c"); !found {
		fmt.Println("  Get(c) = not found")
	}

	fmt.Println()
	fmt.Printf("最終的な頻度 - a: %d, b: %d, d: %d\n",
		cache.GetFrequency("a"),
		cache.GetFrequency("b"),
		cache.GetFrequency("d"))
}

// demoWriteThrough はWrite-Throughキャッシュのデモ
func demoWriteThrough() {
	fmt.Println("3. Write-Through Cache")
	fmt.Println("----------------------------------------------")

	cache := NewWriteThroughCache()

	fmt.Println("Write-Throughキャッシュを作成")
	fmt.Println("（書き込み時にキャッシュとDBを同期）")
	fmt.Println()

	// データ書き込み
	fmt.Println("データを書き込み:")
	start := time.Now()
	cache.Write("user:1", "Alice")
	duration := time.Since(start)
	fmt.Printf("  所要時間: %v (キャッシュ+DB��き込み)\n", duration)

	fmt.Println()

	// データ読み込み（キャッシュヒット）
	fmt.Println("データを読み込み（キャッシュヒット）:")
	start = time.Now()
	if value, found := cache.Read("user:1"); found {
		duration = time.Since(start)
		fmt.Printf("  Read(user:1) = %v\n", value)
		fmt.Printf("  所要時間: %v (キャッシュから高速取得)\n", duration)
	}

	fmt.Println()

	// 複数のデータ書き込み
	fmt.Println("複数のデータを書き込み:")
	start = time.Now()
	cache.Write("user:2", "Bob")
	cache.Write("user:3", "Charlie")
	duration = time.Since(start)
	fmt.Printf("  所要時間: %v (2回の書き込み)\n", duration)

	fmt.Println()
	fmt.Printf("最終的な状態 - Cache: %d, DB: %d\n",
		cache.GetCacheSize(), cache.GetDBSize())

	fmt.Println()
	fmt.Println("特徴:")
	fmt.Println("  ✓ 書き込み時にキャッシュとDBを同時更新")
	fmt.Println("  ✓ データの一貫性が高い")
	fmt.Println("  ✗ 書き込みが遅い（DBの書き込み待ち）")
}

// demoWriteBack はWrite-Backキャッシュのデモ
func demoWriteBack() {
	fmt.Println("4. Write-Back (Write-Behind) Cache")
	fmt.Println("----------------------------------------------")

	cache := NewWriteBackCache(2 * time.Second)
	defer cache.Stop()

	fmt.Println("Write-Backキャッシュを作成")
	fmt.Println("（書き込みをキャッシュに行い、非同期でDBに反映）")
	fmt.Println()

	// データ書き込み
	fmt.Println("データを書き込み:")
	start := time.Now()
	cache.Write("user:1", "Alice")
	duration := time.Since(start)
	fmt.Printf("  所要時間: %v (キャッシュのみ、高速)\n", duration)

	fmt.Println()
	fmt.Printf("現在の状態 - Cache: %d, DB: %d, Dirty: %d\n",
		cache.GetCacheSize(), cache.GetDBSize(), cache.GetDirtyCount())

	fmt.Println()

	// データ読み込み（キャッシュヒット）
	fmt.Println("データを読み込み（キャッシュヒット）:")
	if value, found := cache.Read("user:1"); found {
		fmt.Printf("  Read(user:1) = %v\n", value)
	}

	fmt.Println()

	// 複数のデータ書き込み
	fmt.Println("複数のデータを書き込み:")
	start = time.Now()
	cache.Write("user:2", "Bob")
	cache.Write("user:3", "Charlie")
	duration = time.Since(start)
	fmt.Printf("  所要時間: %v (2回の書き込み、すべて高速)\n", duration)

	fmt.Println()
	fmt.Printf("現在の状態 - Cache: %d, DB: %d, Dirty: %d\n",
		cache.GetCacheSize(), cache.GetDBSize(), cache.GetDirtyCount())

	fmt.Println()
	fmt.Println("バックグラウンドでDBに書き込み中...")
	time.Sleep(3 * time.Second)

	fmt.Println()
	fmt.Printf("最終的な状態 - Cache: %d, DB: %d, Dirty: %d\n",
		cache.GetCacheSize(), cache.GetDBSize(), cache.GetDirtyCount())

	fmt.Println()
	fmt.Println("特徴:")
	fmt.Println("  ✓ 書き込みが高速（キャッシュのみ）")
	fmt.Println("  ✓ DBへの書き込みを非同期で実行")
	fmt.Println("  ✗ データロストのリスク（クラッシュ時）")
	fmt.Println("  ✗ データの一貫性に遅延がある")
}
