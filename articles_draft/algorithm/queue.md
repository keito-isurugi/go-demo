## 概要

- Goでキュー（Queue）を実装
- キューはFIFO（First In First Out：先入れ先出し）のデータ構造。最初に追加した要素が最初に取り出される。

### 特徴

- **FIFO構造**: 最初に入れたものが最初に出る
- **順序保証**: 追加した順番で処理される
- **実用的**: タスクキュー、BFS、プリンタキューなど幅広く使われる
- **効率的**: Enqueue（追加）はO(1)、Dequeue（取り出し）は実装により異なる

## キューとは

### 基本的な考え方

キューは、データを待ち行列のように管理するデータ構造。レジの行列や銀行の窓口を想像すると理解しやすい。

**イメージ**:
```
行列（キュー）:

追加（Enqueue）:        取り出し（Dequeue）:
      ↓                        ↑
[1, 2, 3, 4, 5]
 ↑           ↑
 最初         最後
（先頭）     （末尾）
```

### なぜ重要か？

**プログラミングの基礎**:
- タスクキュー（順番に処理）
- 幅優先探索（BFS）
- メッセージキュー
- プリンタスプーラ

**実用例**:
```
1. プリンタの印刷待ち
   ジョブキュー: [job1] → [job2] → [job3]
   先に追加されたjob1から順番に印刷

2. タスクの実行キュー
   タスク: [task1] → [task2] → [task3]
   先に追加されたtask1から順番に実行

3. BFS（幅優先探索）
   探索キュー: [node1] → [node2] → [node3]
   levelごとに順番に探索
```

## サンプルコード

### 基本実装（汎用型）

```go
package main

import "fmt"

// Queue はFIFO（先入れ先出し）のデータ構造
type Queue struct {
	items []interface{}
}

// NewQueue は新しいキューを生成
func NewQueue() *Queue {
	return &Queue{
		items: []interface{}{},
	}
}

// Enqueue はキューに要素を追加（末尾に追加）
func (q *Queue) Enqueue(item interface{}) {
	q.items = append(q.items, item)
}

// Dequeue はキューから要素を取り出す（先頭から取り出し）
func (q *Queue) Dequeue() (interface{}, bool) {
	if q.IsEmpty() {
		return nil, false
	}

	item := q.items[0]
	q.items = q.items[1:]

	return item, true
}

// Peek はキューの先頭要素を確認（取り出さない）
func (q *Queue) Peek() (interface{}, bool) {
	if q.IsEmpty() {
		return nil, false
	}

	return q.items[0], true
}

// IsEmpty はキューが空かどうか
func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

// Size はキューのサイズを返す
func (q *Queue) Size() int {
	return len(q.items)
}

// Clear はキューをクリア
func (q *Queue) Clear() {
	q.items = []interface{}{}
}
```

**使用例**:
```go
queue := NewQueue()

// Enqueue操作
queue.Enqueue(10)
queue.Enqueue(20)
queue.Enqueue(30)
// キュー: [10, 20, 30]

// Peek操作（先頭要素の確認）
if front, ok := queue.Peek(); ok {
    fmt.Println(front) // 10
}

// Dequeue操作
if item, ok := queue.Dequeue(); ok {
    fmt.Println(item) // 10
}
// キュー: [20, 30]
```

### 整数専用キュー

型安全性を高めるため、整数専用のキューも実装できる。

```go
// IntQueue は整数専用のキュー
type IntQueue struct {
	items []int
}

// NewIntQueue は新しい整数キューを生成
func NewIntQueue() *IntQueue {
	return &IntQueue{
		items: []int{},
	}
}

// Enqueue はキューに要素を追加
func (q *IntQueue) Enqueue(item int) {
	q.items = append(q.items, item)
}

// Dequeue はキューから要素を取り出す
func (q *IntQueue) Dequeue() (int, bool) {
	if q.IsEmpty() {
		return 0, false
	}

	item := q.items[0]
	q.items = q.items[1:]

	return item, true
}

// Peek はキューの先頭要素を確認（取り出さない）
func (q *IntQueue) Peek() (int, bool) {
	if q.IsEmpty() {
		return 0, false
	}

	return q.items[0], true
}

// IsEmpty はキューが空かどうか
func (q *IntQueue) IsEmpty() bool {
	return len(q.items) == 0
}

// Size はキューのサイズを返す
func (q *IntQueue) Size() int {
	return len(q.items)
}
```

## 動作原理

### Enqueue操作の流れ

```
初期状態: キュー = []

Enqueue(10):
  ↓
[10]
 ↑
先頭

Enqueue(20):
  ↓
[10, 20]
 ↑   ↑
先頭 末尾

Enqueue(30):
  ↓
[10, 20, 30]
 ↑       ↑
先頭    末尾
```

### Dequeue操作の流れ

```
キュー: [10, 20, 30]
         ↑
        先頭

Dequeue() → 10を取り出す
  ↓
[20, 30]
 ↑
新しい先頭

Dequeue() → 20を取り出す
  ↓
[30]
 ↑
先頭

Dequeue() → 30を取り出す
  ↓
[]
空キュー
```

### Peek vs Dequeue

```
キュー: [10, 20, 30]

Peek() → 10を確認（取り出さない）
  キューは変わらない: [10, 20, 30]

Dequeue() → 10を取り出す
  キューが変わる: [20, 30]
```

## 実用例

### 1. プリンタキューのシミュレーション

複数の印刷ジョブを順番に処理する。

```go
type PrintJob struct {
	id   int
	name string
}

func simulatePrinterQueue() {
	printerQueue := NewQueue()

	// 印刷ジョブの追加
	jobs := []PrintJob{
		{1, "document1.pdf"},
		{2, "photo.jpg"},
		{3, "report.docx"},
		{4, "invoice.pdf"},
	}

	fmt.Println("印刷ジョブの追加:")
	for _, job := range jobs {
		printerQueue.Enqueue(job)
		fmt.Printf("ジョブ追加: ID=%d, ファイル=%s\n", job.id, job.name)
	}

	// 印刷処理
	fmt.Println("\n印刷処理の実行:")
	jobNumber := 1
	for !printerQueue.IsEmpty() {
		if job, ok := printerQueue.Dequeue(); ok {
			printJob := job.(PrintJob)
			fmt.Printf("%d. 印刷中: ID=%d, ファイル=%s （残り: %d件）\n",
				jobNumber, printJob.id, printJob.name, printerQueue.Size())
			jobNumber++
		}
	}
}
```

**動作例**:
```
印刷ジョブの追加:
ジョブ追加: ID=1, ファイル=document1.pdf
ジョブ追加: ID=2, ファイル=photo.jpg
ジョブ追加: ID=3, ファイル=report.docx
ジョブ追加: ID=4, ファイル=invoice.pdf

印刷処理の実行:
1. 印刷中: ID=1, ファイル=document1.pdf （残り: 3件）
2. 印刷中: ID=2, ファイル=photo.jpg （残り: 2件）
3. 印刷中: ID=3, ファイル=report.docx （残り: 1件）
4. 印刷中: ID=4, ファイル=invoice.pdf （残り: 0件）
```

### 2. タスクキューの処理

タスクを順番に実行する。

```go
type Task struct {
	id       int
	priority string
	name     string
}

func processTaskQueue() {
	taskQueue := NewQueue()

	// タスクの追加
	tasks := []Task{
		{1, "高", "データベースバックアップ"},
		{2, "中", "ログ分析"},
		{3, "高", "セキュリティスキャン"},
		{4, "低", "レポート生成"},
	}

	fmt.Println("タスクの追加:")
	for _, task := range tasks {
		taskQueue.Enqueue(task)
		fmt.Printf("タスク追加: [%s] %s\n", task.priority, task.name)
	}

	// タスクの実行
	fmt.Println("\nタスクの実行:")
	for !taskQueue.IsEmpty() {
		if task, ok := taskQueue.Dequeue(); ok {
			t := task.(Task)
			fmt.Printf("実行中: [%s] %s （残りタスク: %d）\n",
				t.priority, t.name, taskQueue.Size())
		}
	}
}
```

**特徴**:
- 追加した順番で処理される
- 優先度を考慮しない（単純なFIFO）
- 優先度キューとの違いを理解する

### 3. BFS（幅優先探索）

グラフやツリーをレベルごとに探索する。

```go
func demoBFS() {
	// グラフの隣接リスト表現
	// 0 → 1, 2
	// 1 → 3, 4
	// 2 → 5
	graph := map[int][]int{
		0: {1, 2},
		1: {3, 4},
		2: {5},
		3: {},
		4: {},
		5: {},
	}

	visited := make(map[int]bool)
	queue := NewIntQueue()

	start := 0
	queue.Enqueue(start)
	visited[start] = true

	fmt.Printf("開始ノード: %d\n", start)
	fmt.Println("\n探索順序:")

	level := 0
	for !queue.IsEmpty() {
		levelSize := queue.Size()
		fmt.Printf("レベル %d: ", level)

		for i := 0; i < levelSize; i++ {
			node, _ := queue.Dequeue()
			fmt.Printf("%d ", node)

			// 隣接ノードをキューに追加
			for _, neighbor := range graph[node] {
				if !visited[neighbor] {
					queue.Enqueue(neighbor)
					visited[neighbor] = true
				}
			}
		}
		fmt.Println()
		level++
	}
}
```

**出力例**:
```
開始ノード: 0

探索順序:
レベル 0: 0
レベル 1: 1 2
レベル 2: 3 4 5
```

**BFSの特徴**:
- キューを使って実装
- 近い順に探索（levelごと）
- 最短経路の発見に有効

### 4. ホットポテトゲーム

順番にポテトを渡し、最後に残った人が勝者。

```go
func demoHotPotato() {
	players := []string{"Alice", "Bob", "Charlie", "David", "Eve"}
	queue := NewQueue()

	// プレイヤーをキューに追加
	for _, player := range players {
		queue.Enqueue(player)
	}

	passes := 7 // ポテトを渡す回数

	for queue.Size() > 1 {
		// ポテトを渡す
		for i := 0; i < passes; i++ {
			if player, ok := queue.Dequeue(); ok {
				queue.Enqueue(player)
			}
		}

		// 現在ポテトを持っている人が脱落
		if eliminated, ok := queue.Dequeue(); ok {
			fmt.Printf("%s が脱落！（残り: %d人）\n", eliminated, queue.Size())
		}
	}

	if winner, ok := queue.Dequeue(); ok {
		fmt.Printf("\n勝者: %s！\n", winner)
	}
}
```

**動作の仕組み**:
```
プレイヤー: [Alice, Bob, Charlie, David, Eve]

ポテトを7回渡す:
  Dequeue → Enqueue を7回繰り返す
  [Alice, Bob, Charlie, David, Eve]
   ↓
  [Bob, Charlie, David, Eve, Alice]
   ↓
  （7回後）
  [Charlie, David, Eve, Alice, Bob]

脱落:
  Dequeue → Charlieが脱落
  残り: [David, Eve, Alice, Bob]
```

### 5. 円形キュー（リングバッファ）

固定サイズで効率的なキュー実装。

```go
// CircularQueue は円形キュー（リングバッファ）
type CircularQueue struct {
	items []interface{}
	front int
	rear  int
	size  int
	cap   int
}

// NewCircularQueue は指定サイズの円形キューを生成
func NewCircularQueue(capacity int) *CircularQueue {
	return &CircularQueue{
		items: make([]interface{}, capacity),
		front: 0,
		rear:  -1,
		size:  0,
		cap:   capacity,
	}
}

// Enqueue は円形キューに要素を追加
func (q *CircularQueue) Enqueue(item interface{}) bool {
	if q.IsFull() {
		return false
	}

	q.rear = (q.rear + 1) % q.cap
	q.items[q.rear] = item
	q.size++
	return true
}

// Dequeue は円形キューから要素を取り出す
func (q *CircularQueue) Dequeue() (interface{}, bool) {
	if q.IsEmpty() {
		return nil, false
	}

	item := q.items[q.front]
	q.front = (q.front + 1) % q.cap
	q.size--
	return item, true
}

// IsFull は円形キューが満杯かどうか
func (q *CircularQueue) IsFull() bool {
	return q.size == q.cap
}
```

**メリット**:
```
通常のキュー:
  Dequeue → 配列の先頭を削除 → すべての要素をシフト → O(n)
  [10, 20, 30, 40] → Dequeue(10) → [20, 30, 40]
                                      ↑ シフトが必要

円形キュー:
  Dequeue → frontインデックスを進める → O(1)
  [10, 20, 30, 40] → Dequeue(10) → frontを1つ進めるだけ
                                      シフト不要！
```

**動作イメージ**:
```
容量5の円形キュー:
インデックス: [0, 1, 2, 3, 4]

初期状態:
  front=0, rear=-1

Enqueue(10):
  rear=(0)%5=0 → items[0]=10
  front=0, rear=0, size=1

Enqueue(20), Enqueue(30), Enqueue(40), Enqueue(50):
  items=[10, 20, 30, 40, 50]
  front=0, rear=4, size=5 (満杯)

Dequeue():
  item=items[0]=10
  front=(1)%5=1, size=4

Enqueue(60):
  rear=(0)%5=0 → items[0]=60 （巡回して0に戻る）
  front=1, rear=0, size=5
```

## 計算量

### 時間計算量

| 操作 | 通常のキュー | 円形キュー | 説明 |
| --- | --- | --- | --- |
| Enqueue | O(1) | O(1) | 配列の末尾に追加 |
| Dequeue | O(n) | O(1) | 通常は要素のシフトが必要、円形は不要 |
| Peek | O(1) | O(1) | 配列の先頭を参照 |
| IsEmpty | O(1) | O(1) | サイズをチェック |
| Size | O(1) | O(1) | サイズを返す |

### Dequeueの計算量の違い

**通常のキュー（スライス実装）**:
```go
// O(n): すべての要素をシフト
item := q.items[0]
q.items = q.items[1:]  // 配列の再スライス
```

**円形キュー**:
```go
// O(1): インデックスを進めるだけ
item := q.items[q.front]
q.front = (q.front + 1) % q.cap  // インデックスの更新のみ
```

### 空間計算量

- **通常のキュー**: O(n) - 動的にサイズ変更
- **円形キュー**: O(k) - 固定サイズ（kは容量）

## キュー vs スタック

### 違い

| 特徴 | キュー（Queue） | スタック（Stack） |
| --- | --- | --- |
| 構造 | FIFO（先入れ先出し） | LIFO（後入れ先出し） |
| 取り出し | 最初に追加した要素 | 最後に追加した要素 |
| イメージ | 行列・待ち列 | 皿の積み重ね |
| 用途 | タスクキュー、BFS | 関数呼び出し、DFS、Undo |

### イメージ図

```
キュー（FIFO）:
  Enqueue → [1, 2, 3] → Dequeue
  追加       ↑     ↑    取り出し
            最初   最後

スタック（LIFO）:
  Push →  [1, 2, 3] → Pop
  追加     ↑     ↑   取り出し
          最初   最後
         （最後に取り出される）
```

### 探索アルゴリズムでの使い分け

```
BFS（幅優先探索）: キューを使用
  レベル0: [A]
  レベル1: [B, C]
  レベル2: [D, E, F, G]
  → 近い順に探索

DFS（深さ優先探索）: スタックを使用
  深さ0: [A]
  深さ1: [B]
  深さ2: [D]
  深さ3: [H]
  → 深い順に探索
```

## 使いどころ

### 向いている場面

1. **タスクの順次処理**
   - プリンタスプーラ
   - ジョブキュー
   - メッセージキュー

2. **幅優先探索（BFS）**
   - グラフの最短経路
   - ツリーのレベル順走査

3. **バッファリング**
   - ストリーミングデータ
   - ネットワークパケット
   - イベント処理

4. **順序保証が必要な処理**
   - リクエストの処理順序
   - トランザクションキュー

### 実例: Webサーバーのリクエスト処理

```go
type Request struct {
	id   int
	path string
}

type RequestQueue struct {
	queue *Queue
}

func NewRequestQueue() *RequestQueue {
	return &RequestQueue{
		queue: NewQueue(),
	}
}

// AddRequest はリクエストをキューに追加
func (rq *RequestQueue) AddRequest(req Request) {
	rq.queue.Enqueue(req)
	fmt.Printf("リクエスト追加: ID=%d, Path=%s （キューサイズ: %d）\n",
		req.id, req.path, rq.queue.Size())
}

// ProcessNext は次のリクエストを処理
func (rq *RequestQueue) ProcessNext() {
	if req, ok := rq.queue.Dequeue(); ok {
		r := req.(Request)
		fmt.Printf("処理中: ID=%d, Path=%s （残り: %d）\n",
			r.id, r.path, rq.queue.Size())
	}
}
```

**使用例**:
```go
reqQueue := NewRequestQueue()

// リクエストが届く
reqQueue.AddRequest(Request{1, "/api/users"})
reqQueue.AddRequest(Request{2, "/api/posts"})
reqQueue.AddRequest(Request{3, "/api/comments"})

// ワーカーが順番に処理
for !reqQueue.queue.IsEmpty() {
	reqQueue.ProcessNext()
}
```

## 注意点

### 1. メモリリーク（通常のキュー）

通常のスライス実装では、Dequeueでメモリが増え続ける可能性。

```go
// 問題: 内部配列がシュリンクしない
for i := 0; i < 1000000; i++ {
	queue.Enqueue(i)
}

for i := 0; i < 999999; i++ {
	queue.Dequeue()  // メモリは解放されない
}
// キューには1要素だけだが、内部配列は大きいまま
```

**対策**:
```go
// 定期的にキューを再作成
if queue.Size() < len(queue.items)/4 {
	// 新しいスライスにコピー
	newItems := make([]interface{}, queue.Size())
	copy(newItems, queue.items)
	queue.items = newItems
}
```

### 2. 空キューのDequeue/Peek

空のキューから取り出すとエラー。

```go
queue := NewQueue()

// 空キューからDequeue → エラー
if item, ok := queue.Dequeue(); !ok {
	fmt.Println("キューが空です")
}

// 空チェックを必ず行う
if !queue.IsEmpty() {
	item, _ := queue.Dequeue()
	// 処理
}
```

### 3. 円形キューの満杯チェック

円形キューは固定サイズなので満杯チェックが必要。

```go
cq := NewCircularQueue(5)

// 5個追加
for i := 0; i < 5; i++ {
	cq.Enqueue(i)
}

// 満杯なので追加失敗
if !cq.Enqueue(6) {
	fmt.Println("キューが満杯です")
}
```

### 4. スレッドセーフ性

基本実装はスレッドセーフではない。

```go
// マルチスレッド環境では mutex を使う
type SafeQueue struct {
	items []interface{}
	mu    sync.Mutex
}

func (q *SafeQueue) Enqueue(item interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

func (q *SafeQueue) Dequeue() (interface{}, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return nil, false
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}
```

## まとめ

### メリット

- 順序保証が確実
- 実装がシンプル
- 幅広い用途に使える
- BFSなど重要なアルゴリズムの基礎

### デメリット

- 通常のキューはDequeueがO(n)
- 途中の要素にアクセスできない
- 円形キューは固定サイズの制約

### 実装方法の選択

| 実装方法 | メリット | デメリット | 用途 |
| --- | --- | --- | --- |
| スライス | 実装が簡単、動的サイズ | DequeueがO(n) | 小規模、一時的な用途 |
| 円形キュー | Dequeue O(1)、効率的 | 固定サイズ | 高頻度処理、バッファ |
| リンクリスト | Dequeue O(1)、動的サイズ | メモリオーバーヘッド | 大規模、可変サイズ |

### 使うべき時

- **タスクの順次処理**: ジョブキュー、プリンタスプーラ
- **幅優先探索**: BFS、最短経路
- **バッファリング**: ストリーミング、パケット処理
- **順序保証**: リクエスト処理、トランザクション

### スタックとの使い分け

| 用途 | キュー（FIFO） | スタック（LIFO） |
| --- | --- | --- |
| タスク処理 | ✓ |  |
| BFS（幅優先） | ✓ |  |
| バッファリング | ✓ |  |
| 関数呼び出し |  | ✓ |
| DFS（深さ優先） |  | ✓ |
| Undo/Redo |  | ✓ |

キューはスタックと並ぶ基本的なデータ構造。FIFO（先入れ先出し）という単純な原理だが、タスク処理、BFS、バッファリングなど実用性が高い。実装方法によって計算量が変わるため、用途に応じて選択することが重要。プログラマーなら必ず理解すべきデータ構造。
