package algorithm

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

// DemoQueue はキューの動作デモ
func DemoQueue() {
	fmt.Println("\n1. 基本的なキュー操作:")
	fmt.Println("=================================")

	queue := NewQueue()

	// Enqueue操作
	fmt.Println("\nEnqueue操作:")
	queue.Enqueue(10)
	fmt.Printf("Enqueue(10) → キュー: %v\n", queue.items)

	queue.Enqueue(20)
	fmt.Printf("Enqueue(20) → キュー: %v\n", queue.items)

	queue.Enqueue(30)
	fmt.Printf("Enqueue(30) → キュー: %v\n", queue.items)

	queue.Enqueue(40)
	fmt.Printf("Enqueue(40) → キュー: %v\n", queue.items)

	// Peek操作
	fmt.Println("\nPeek操作（先頭の確認）:")
	if front, ok := queue.Peek(); ok {
		fmt.Printf("先頭の要素: %v （キューは変わらない: %v）\n", front, queue.items)
	}

	// Dequeue操作
	fmt.Println("\nDequeue操作:")
	if item, ok := queue.Dequeue(); ok {
		fmt.Printf("Dequeue() → %v （キュー: %v）\n", item, queue.items)
	}

	if item, ok := queue.Dequeue(); ok {
		fmt.Printf("Dequeue() → %v （キュー: %v）\n", item, queue.items)
	}

	// サイズ確認
	fmt.Printf("\nキューのサイズ: %d\n", queue.Size())

	fmt.Println("\n2. プリンタキューのシミュレーション:")
	fmt.Println("=================================")
	simulatePrinterQueue()

	fmt.Println("\n3. タスクキューの処理:")
	fmt.Println("=================================")
	processTaskQueue()

	fmt.Println("\n4. BFS（幅優先探索）のデモ:")
	fmt.Println("=================================")
	demoBFS()

	fmt.Println("\n5. ホットポテトゲーム:")
	fmt.Println("=================================")
	demoHotPotato()

	fmt.Println("\n6. 円形キュー（リングバッファ）:")
	fmt.Println("=================================")
	demoCircularQueue()
}

// simulatePrinterQueue はプリンタキューをシミュレート
func simulatePrinterQueue() {
	type PrintJob struct {
		id   int
		name string
	}

	printerQueue := NewQueue()

	fmt.Println("\n印刷ジョブの追加:")
	jobs := []PrintJob{
		{1, "document1.pdf"},
		{2, "photo.jpg"},
		{3, "report.docx"},
		{4, "invoice.pdf"},
	}

	for _, job := range jobs {
		printerQueue.Enqueue(job)
		fmt.Printf("ジョブ追加: ID=%d, ファイル=%s\n", job.id, job.name)
	}

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

	fmt.Println("\nすべての印刷ジョブが完了しました")
}

// processTaskQueue はタスクキューを処理
func processTaskQueue() {
	type Task struct {
		id       int
		priority string
		name     string
	}

	taskQueue := NewQueue()

	fmt.Println("\nタスクの追加:")
	tasks := []Task{
		{1, "高", "データベースバックアップ"},
		{2, "中", "ログ分析"},
		{3, "高", "セキュリティスキャン"},
		{4, "低", "レポート生成"},
	}

	for _, task := range tasks {
		taskQueue.Enqueue(task)
		fmt.Printf("タスク追加: [%s] %s\n", task.priority, task.name)
	}

	fmt.Println("\nタスクの実行:")
	for !taskQueue.IsEmpty() {
		if task, ok := taskQueue.Dequeue(); ok {
			t := task.(Task)
			fmt.Printf("実行中: [%s] %s （残りタスク: %d）\n",
				t.priority, t.name, taskQueue.Size())
		}
	}
}

// demoBFS は幅優先探索のデモ
func demoBFS() {
	// グラフの隣接リスト表現
	// 0 → 1, 2
	// 1 → 3, 4
	// 2 → 5
	// 3, 4, 5 → なし
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

	fmt.Printf("\n開始ノード: %d\n", start)
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

// demoHotPotato はホットポテトゲームのデモ
func demoHotPotato() {
	players := []string{"Alice", "Bob", "Charlie", "David", "Eve"}
	queue := NewQueue()

	fmt.Println("\nプレイヤー:")
	for _, player := range players {
		queue.Enqueue(player)
		fmt.Printf("- %s\n", player)
	}

	passes := 7 // ポテトを渡す回数
	fmt.Printf("\nポテトを%d回渡します...\n\n", passes)

	round := 1
	for queue.Size() > 1 {
		fmt.Printf("ラウンド %d:\n", round)

		// ポテトを渡す
		for i := 0; i < passes; i++ {
			if player, ok := queue.Dequeue(); ok {
				queue.Enqueue(player)
			}
		}

		// 現在ポテトを持っている人が脱落
		if eliminated, ok := queue.Dequeue(); ok {
			fmt.Printf("  %s が脱落！（残り: %d人）\n", eliminated, queue.Size())
		}

		round++
	}

	if winner, ok := queue.Dequeue(); ok {
		fmt.Printf("\n勝者: %s！\n", winner)
	}
}

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

// IsEmpty は円形キューが空かどうか
func (q *CircularQueue) IsEmpty() bool {
	return q.size == 0
}

// IsFull は円形キューが満杯かどうか
func (q *CircularQueue) IsFull() bool {
	return q.size == q.cap
}

// Size は円形キューの現在のサイズ
func (q *CircularQueue) Size() int {
	return q.size
}

// demoCircularQueue は円形キューのデモ
func demoCircularQueue() {
	cq := NewCircularQueue(5)

	fmt.Println("\n円形キュー（容量: 5）の操作:")

	// Enqueue
	fmt.Println("\nEnqueue操作:")
	for i := 1; i <= 5; i++ {
		if cq.Enqueue(i * 10) {
			fmt.Printf("Enqueue(%d) 成功 （サイズ: %d/%d）\n", i*10, cq.Size(), cq.cap)
		}
	}

	// 満杯時のEnqueue
	fmt.Println("\n満杯時のEnqueue:")
	if !cq.Enqueue(60) {
		fmt.Println("Enqueue(60) 失敗: キューが満杯です")
	}

	// Dequeue
	fmt.Println("\nDequeue操作:")
	for i := 0; i < 3; i++ {
		if item, ok := cq.Dequeue(); ok {
			fmt.Printf("Dequeue() → %v （サイズ: %d/%d）\n", item, cq.Size(), cq.cap)
		}
	}

	// 空きができたのでEnqueue
	fmt.Println("\n空きができたのでEnqueue:")
	for i := 6; i <= 8; i++ {
		if cq.Enqueue(i * 10) {
			fmt.Printf("Enqueue(%d) 成功 （サイズ: %d/%d）\n", i*10, cq.Size(), cq.cap)
		}
	}

	// 残りをすべてDequeue
	fmt.Println("\n残りをすべてDequeue:")
	for !cq.IsEmpty() {
		if item, ok := cq.Dequeue(); ok {
			fmt.Printf("Dequeue() → %v （サイズ: %d/%d）\n", item, cq.Size(), cq.cap)
		}
	}

	fmt.Println("\nメリット:")
	fmt.Println("- 固定サイズで効率的なメモリ使用")
	fmt.Println("- 配列の先頭削除によるシフトが不要")
	fmt.Println("- EnqueueとDequeueの両方がO(1)")
}
