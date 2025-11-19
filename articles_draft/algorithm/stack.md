## 概要

- Goでスタック（Stack）を実装
- スタックはLIFO（Last In First Out：後入れ先出し）のデータ構造。最後に追加した要素が最初に取り出される。

### 特徴

- **LIFO構造**: 最後に入れたものが最初に出る
- **シンプルな操作**: Push（追加）とPop（取り出し）が基本
- **効率的**: すべての基本操作がO(1)
- **実用的**: 関数呼び出し、括弧マッチング、履歴管理など幅広く使われる

## スタックとは

### 基本的な考え方

スタックは、データを積み重ねるように管理するデータ構造。皿を積み重ねる様子を想像すると理解しやすい。

**イメージ**:
```
皿を積み重ねる（スタック）:

追加（Push）:          取り出し（Pop）:
  ↓                     ↑
[5]  ← 最後             [5]  ← 最初に取り出される
[4]                     [4]
[3]                     [3]
[2]                     [2]
[1]  ← 最初             [1]  ← 最後に取り出される
```

### なぜ重要か？

**プログラミングの基礎**:
- 関数呼び出しスタック（コールスタック）
- 再帰処理の実装
- 式の評価（逆ポーランド記法）
- Undo/Redo機能

**実用例**:
```
1. ブラウザの戻る/進む機能
   訪問履歴: [A] → [B] → [C]
   「戻る」でCをPop → Bに戻る

2. テキストエディタのUndo
   操作履歴: [操作1] → [操作2] → [操作3]
   「Undo」で操作3をPop → 操作2の状態に戻る

3. 関数呼び出し
   main() → funcA() → funcB()
   スタック: [main] → [funcA] → [funcB]
   funcBが終わると → Pop → funcAに戻る
```

## サンプルコード

### 基本実装（汎用型）

```go
package main

import "fmt"

// Stack はLIFO（後入れ先出し）のデータ構造
type Stack struct {
	items []interface{}
}

// NewStack は新しいスタックを生成
func NewStack() *Stack {
	return &Stack{
		items: []interface{}{},
	}
}

// Push はスタックに要素を追加
func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

// Pop はスタックから要素を取り出す
func (s *Stack) Pop() (interface{}, bool) {
	if s.IsEmpty() {
		return nil, false
	}

	lastIndex := len(s.items) - 1
	item := s.items[lastIndex]
	s.items = s.items[:lastIndex]

	return item, true
}

// Peek はスタックの最上位要素を確認（取り出さない）
func (s *Stack) Peek() (interface{}, bool) {
	if s.IsEmpty() {
		return nil, false
	}

	return s.items[len(s.items)-1], true
}

// IsEmpty はスタックが空かどうか
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

// Size はスタックのサイズを返す
func (s *Stack) Size() int {
	return len(s.items)
}

// Clear はスタックをクリア
func (s *Stack) Clear() {
	s.items = []interface{}{}
}
```

**使用例**:
```go
stack := NewStack()

// Push操作
stack.Push(10)
stack.Push(20)
stack.Push(30)
// スタック: [10, 20, 30]

// Peek操作（最上位要素の確認）
if top, ok := stack.Peek(); ok {
    fmt.Println(top) // 30
}

// Pop操作
if item, ok := stack.Pop(); ok {
    fmt.Println(item) // 30
}
// スタック: [10, 20]
```

### 整数専用スタック

型安全性を高めるため、整数専用のスタックも実装できる。

```go
// IntStack は整数専用のスタック
type IntStack struct {
	items []int
}

// NewIntStack は新しい整数スタックを生成
func NewIntStack() *IntStack {
	return &IntStack{
		items: []int{},
	}
}

// Push はスタックに要素を追加
func (s *IntStack) Push(item int) {
	s.items = append(s.items, item)
}

// Pop はスタックから要素を取り出す
func (s *IntStack) Pop() (int, bool) {
	if s.IsEmpty() {
		return 0, false
	}

	lastIndex := len(s.items) - 1
	item := s.items[lastIndex]
	s.items = s.items[:lastIndex]

	return item, true
}

// Peek はスタックの最上位要素を確認（取り出さない）
func (s *IntStack) Peek() (int, bool) {
	if s.IsEmpty() {
		return 0, false
	}

	return s.items[len(s.items)-1], true
}

// IsEmpty はスタックが空かどうか
func (s *IntStack) IsEmpty() bool {
	return len(s.items) == 0
}

// Size はスタックのサイズを返す
func (s *IntStack) Size() int {
	return len(s.items)
}
```

## 動作原理

### Push操作の流れ

```
初期状態: スタック = []

Push(10):
  ↓
[10]

Push(20):
  ↓
[10, 20]

Push(30):
  ↓
[10, 20, 30]  ← 30が最上位（top）
```

### Pop操作の流れ

```
スタック: [10, 20, 30]

Pop() → 30を取り出す
  ↓
[10, 20]  ← 20が新しい最上位

Pop() → 20を取り出す
  ↓
[10]  ← 10が最上位

Pop() → 10を取り出す
  ↓
[]  ← 空スタック
```

### Peek vs Pop

```
スタック: [10, 20, 30]

Peek() → 30を確認（取り出さない）
  スタックは変わらない: [10, 20, 30]

Pop() → 30を取り出す
  スタックが変わる: [10, 20]
```

## 実用例

### 1. 括弧のマッチング

括弧が正しく対応しているかをチェック。

```go
// isValidParentheses は括弧が正しくマッチしているかを判定
func isValidParentheses(s string) bool {
	stack := NewStack()
	pairs := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	for _, char := range s {
		// 開き括弧ならpush
		if char == '(' || char == '{' || char == '[' {
			stack.Push(char)
			continue
		}

		// 閉じ括弧の場合
		if expected, exists := pairs[char]; exists {
			if top, ok := stack.Pop(); !ok || top != expected {
				return false
			}
		}
	}

	return stack.IsEmpty()
}
```

**動作例**:
```
入力: "(())"

ステップ:
1. '(' → Push → スタック: ['(']
2. '(' → Push → スタック: ['(', '(']
3. ')' → Pop '(' → スタック: ['(']
4. ')' → Pop '(' → スタック: []

結果: 空スタック → 正しい ✓

入力: "(()"

ステップ:
1. '(' → Push → スタック: ['(']
2. '(' → Push → スタック: ['(', '(']
3. ')' → Pop '(' → スタック: ['(']
4. 終了 → スタックに '(' が残っている

結果: スタックに要素が残る → 不正 ✗
```

### 2. 逆ポーランド記法（RPN）の評価

後置記法の数式を評価する。

```go
// evaluateRPN は逆ポーランド記法を評価
func evaluateRPN(tokens []string) int {
	stack := NewIntStack()

	for _, token := range tokens {
		switch token {
		case "+":
			b, _ := stack.Pop()
			a, _ := stack.Pop()
			stack.Push(a + b)
		case "-":
			b, _ := stack.Pop()
			a, _ := stack.Pop()
			stack.Push(a - b)
		case "*":
			b, _ := stack.Pop()
			a, _ := stack.Pop()
			stack.Push(a * b)
		case "/":
			b, _ := stack.Pop()
			a, _ := stack.Pop()
			stack.Push(a / b)
		default:
			// 数値をpush
			num := 0
			fmt.Sscanf(token, "%d", &num)
			stack.Push(num)
		}
	}

	result, _ := stack.Pop()
	return result
}
```

**動作例**:
```
入力: ["2", "3", "+", "4", "*"]  → (2 + 3) * 4 = 20

ステップ:
1. "2" → Push(2) → スタック: [2]
2. "3" → Push(3) → スタック: [2, 3]
3. "+" → Pop(3), Pop(2), Push(2+3=5) → スタック: [5]
4. "4" → Push(4) → スタック: [5, 4]
5. "*" → Pop(4), Pop(5), Push(5*4=20) → スタック: [20]

結果: 20
```

**逆ポーランド記法とは**:
```
通常の記法（中置記法）: (2 + 3) * 4
逆ポーランド記法（後置記法）: 2 3 + 4 *

メリット:
- 括弧が不要
- 演算の優先順位が明確
- スタックで簡単に評価できる
```

### 3. 文字列の反転

```go
// reverseString は文字列を反転
func reverseString(s string) string {
	stack := NewStack()

	// 文字を1つずつスタックに積む
	for _, char := range s {
		stack.Push(char)
	}

	// スタックから取り出して結果を構築
	result := ""
	for !stack.IsEmpty() {
		if char, ok := stack.Pop(); ok {
			result += string(char.(rune))
		}
	}

	return result
}
```

**動作例**:
```
入力: "hello"

Push段階:
h → [h]
e → [h, e]
l → [h, e, l]
l → [h, e, l, l]
o → [h, e, l, l, o]

Pop段階:
Pop() → o → result = "o"
Pop() → l → result = "ol"
Pop() → l → result = "oll"
Pop() → e → result = "olle"
Pop() → h → result = "olleh"

結果: "olleh"
```

### 4. 関数呼び出しのシミュレーション

```go
// simulateFunctionCalls は関数呼び出しスタックをシミュレーション
func simulateFunctionCalls() {
	callStack := NewStack()

	fmt.Println("関数呼び出し:")
	fmt.Println("main() → funcA() → funcB() → funcC()")

	// 関数呼び出し
	callStack.Push("main()")
	fmt.Printf("呼び出し: main()  → スタック: %v\n", callStack.items)

	callStack.Push("funcA()")
	fmt.Printf("呼び出し: funcA() → スタック: %v\n", callStack.items)

	callStack.Push("funcB()")
	fmt.Printf("呼び出し: funcB() → スタック: %v\n", callStack.items)

	callStack.Push("funcC()")
	fmt.Printf("呼び出し: funcC() → スタック: %v\n", callStack.items)

	fmt.Println("\n関数の終了（return）:")

	// 関数の終了
	for !callStack.IsEmpty() {
		if fn, ok := callStack.Pop(); ok {
			fmt.Printf("終了: %s → スタック: %v\n", fn, callStack.items)
		}
	}
}
```

**出力例**:
```
関数呼び出し:
呼び出し: main()  → スタック: [main()]
呼び出し: funcA() → スタック: [main(), funcA()]
呼び出し: funcB() → スタック: [main(), funcA(), funcB()]
呼び出し: funcC() → スタック: [main(), funcA(), funcB(), funcC()]

関数の終了（return）:
終了: funcC() → スタック: [main(), funcA(), funcB()]
終了: funcB() → スタック: [main(), funcA()]
終了: funcA() → スタック: [main()]
終了: main()  → スタック: []
```

### 5. Unixパスの簡略化

```go
// simplifyPath はUnixパスを簡略化
func simplifyPath(path string) string {
	stack := NewStack()
	components := []string{}

	// パスをコンポーネントに分割
	current := ""
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			if current != "" {
				components = append(components, current)
				current = ""
			}
		} else {
			current += string(path[i])
		}
	}
	if current != "" {
		components = append(components, current)
	}

	// スタックを使って処理
	for _, comp := range components {
		if comp == ".." {
			// 1つ上のディレクトリ
			if !stack.IsEmpty() {
				stack.Pop()
			}
		} else if comp != "." && comp != "" {
			// 通常のディレクトリ名
			stack.Push(comp)
		}
		// "." は現在のディレクトリなので無視
	}

	// スタックから結果を構築
	if stack.IsEmpty() {
		return "/"
	}

	result := ""
	temp := NewStack()
	for !stack.IsEmpty() {
		if item, ok := stack.Pop(); ok {
			temp.Push(item)
		}
	}

	for !temp.IsEmpty() {
		if item, ok := temp.Pop(); ok {
			result += "/" + item.(string)
		}
	}

	return result
}
```

**動作例**:
```
入力: "/home/user/documents/../pictures"

コンポーネント: ["home", "user", "documents", "..", "pictures"]

処理:
1. "home" → Push → スタック: [home]
2. "user" → Push → スタック: [home, user]
3. "documents" → Push → スタック: [home, user, documents]
4. ".." → Pop → スタック: [home, user]
5. "pictures" → Push → スタック: [home, user, pictures]

結果: "/home/user/pictures"
```

## 計算量

### 時間計算量

| 操作 | 時間計算量 | 説明 |
| --- | --- | --- |
| Push | O(1) | 配列の末尾に追加 |
| Pop | O(1) | 配列の末尾から削除 |
| Peek | O(1) | 配列の末尾を参照 |
| IsEmpty | O(1) | 配列の長さをチェック |
| Size | O(1) | 配列の長さを返す |

### 空間計算量

- **O(n)**: n個の要素を格納
- スタックのサイズは格納する要素数に比例

### なぜO(1)なのか？

**配列の末尾操作は高速**:
```go
// Push: 配列の末尾に追加
s.items = append(s.items, item)  // O(1)*
// *償却計算量。稀に配列の再確保でO(n)

// Pop: 配列の末尾から削除
lastIndex := len(s.items) - 1
item := s.items[lastIndex]
s.items = s.items[:lastIndex]  // O(1)

// Peek: 配列の末尾を参照
return s.items[len(s.items)-1]  // O(1)
```

**重要なポイント**:
- 配列の末尾へのアクセスはインデックス計算だけ
- 要素をシフトする必要がない
- メモリの再確保は稀（償却O(1)）

## スタック vs キュー

### 違い

| 特徴 | スタック（Stack） | キュー（Queue） |
| --- | --- | --- |
| 構造 | LIFO（後入れ先出し） | FIFO（先入れ先出し） |
| 取り出し | 最後に追加した要素 | 最初に追加した要素 |
| イメージ | 皿の積み重ね | 行列・待ち列 |
| 用途 | 関数呼び出し、Undo | タスクキュー、BFS |

### イメージ図

```
スタック（LIFO）:
  Push →  [1, 2, 3] → Pop
  追加     ↑     ↑   取り出し
           最初   最後

キュー（FIFO）:
  Enqueue → [1, 2, 3] → Dequeue
  追加       ↑     ↑    取り出し
            最初   最後
```

## 使いどころ

### 向いている場面

1. **関数呼び出しの管理**
   - コールスタック
   - 再帰処理

2. **括弧・タグのマッチング**
   - 括弧の対応チェック
   - HTMLタグの検証

3. **式の評価**
   - 逆ポーランド記法
   - 中置記法→後置記法の変換

4. **履歴管理**
   - Undo/Redo機能
   - ブラウザの戻る/進む

5. **深さ優先探索（DFS）**
   - グラフ探索
   - 木構造の走査

### 実例: ブラウザの戻る/進む機能

```go
type Browser struct {
	backStack    *Stack  // 戻るスタック
	forwardStack *Stack  // 進むスタック
	currentPage  string
}

func NewBrowser() *Browser {
	return &Browser{
		backStack:    NewStack(),
		forwardStack: NewStack(),
		currentPage:  "home",
	}
}

// Visit は新しいページに移動
func (b *Browser) Visit(page string) {
	if b.currentPage != "" {
		b.backStack.Push(b.currentPage)
	}
	b.currentPage = page
	b.forwardStack.Clear() // 新規訪問で進む履歴をクリア
}

// Back は前のページに戻る
func (b *Browser) Back() string {
	if b.backStack.IsEmpty() {
		return b.currentPage
	}

	b.forwardStack.Push(b.currentPage)
	page, _ := b.backStack.Pop()
	b.currentPage = page.(string)
	return b.currentPage
}

// Forward は次のページに進む
func (b *Browser) Forward() string {
	if b.forwardStack.IsEmpty() {
		return b.currentPage
	}

	b.backStack.Push(b.currentPage)
	page, _ := b.forwardStack.Pop()
	b.currentPage = page.(string)
	return b.currentPage
}
```

**使用例**:
```go
browser := NewBrowser()

browser.Visit("google.com")
browser.Visit("github.com")
browser.Visit("stackoverflow.com")

// 戻る: stackoverflow.com → github.com
browser.Back()  // "github.com"

// 戻る: github.com → google.com
browser.Back()  // "google.com"

// 進む: google.com → github.com
browser.Forward()  // "github.com"
```

### 実例: テキストエディタのUndo/Redo

```go
type Editor struct {
	text       string
	undoStack  *Stack
	redoStack  *Stack
}

func NewEditor() *Editor {
	return &Editor{
		text:      "",
		undoStack: NewStack(),
		redoStack: NewStack(),
	}
}

// Type はテキストを追加
func (e *Editor) Type(newText string) {
	e.undoStack.Push(e.text)
	e.text += newText
	e.redoStack.Clear() // 新規入力でRedo履歴をクリア
}

// Undo は操作を取り消す
func (e *Editor) Undo() {
	if e.undoStack.IsEmpty() {
		return
	}

	e.redoStack.Push(e.text)
	prevText, _ := e.undoStack.Pop()
	e.text = prevText.(string)
}

// Redo は操作をやり直す
func (e *Editor) Redo() {
	if e.redoStack.IsEmpty() {
		return
	}

	e.undoStack.Push(e.text)
	nextText, _ := e.redoStack.Pop()
	e.text = nextText.(string)
}
```

**使用例**:
```go
editor := NewEditor()

editor.Type("Hello")     // text = "Hello"
editor.Type(" World")    // text = "Hello World"
editor.Type("!")         // text = "Hello World!"

editor.Undo()            // text = "Hello World"
editor.Undo()            // text = "Hello"

editor.Redo()            // text = "Hello World"
```

## 注意点

### 1. スタックオーバーフロー

スタックに要素を追加しすぎるとメモリ不足になる。

```go
// 注意: 無限にPushするとメモリ不足
stack := NewStack()
for i := 0; i < 10000000000; i++ {
	stack.Push(i)  // メモリ不足のリスク
}
```

**対策**:
- スタックのサイズに上限を設ける
- 定期的に不要な要素をクリア

### 2. 空スタックのPop/Peek

空のスタックからPopやPeekすると問題が起きる。

```go
stack := NewStack()

// 空スタックからPop → エラー
if item, ok := stack.Pop(); !ok {
	fmt.Println("スタックが空です")
}

// 空チェックを必ず行う
if !stack.IsEmpty() {
	item, _ := stack.Pop()
	// 処理
}
```

### 3. スレッドセーフ性

基本実装はスレッドセーフではない。

```go
// マルチスレッド環境では mutex を使う
type SafeStack struct {
	items []interface{}
	mu    sync.Mutex
}

func (s *SafeStack) Push(item interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, item)
}

func (s *SafeStack) Pop() (interface{}, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.items) == 0 {
		return nil, false
	}

	lastIndex := len(s.items) - 1
	item := s.items[lastIndex]
	s.items = s.items[:lastIndex]

	return item, true
}
```

## まとめ

### メリット

- すべての基本操作がO(1)で高速
- 実装がシンプルで理解しやすい
- メモリ効率が良い（必要な分だけ確保）
- 幅広い用途に使える

### デメリット

- 途中の要素にアクセスできない
- LIFO以外の操作には向かない
- サイズ制限がない場合メモリ不足のリスク

### 使うべき時

- **関数呼び出しの管理**: コールスタック、再帰
- **括弧マッチング**: 括弧検証、タグ検証
- **式の評価**: 逆ポーランド記法、計算機
- **履歴管理**: Undo/Redo、ブラウザ履歴
- **深さ優先探索**: DFS、バックトラッキング

### キューとの使い分け

| 用途 | スタック（LIFO） | キュー（FIFO） |
| --- | --- | --- |
| 関数呼び出し | ✓ |  |
| Undo/Redo | ✓ |  |
| 括弧マッチング | ✓ |  |
| DFS（深さ優先） | ✓ |  |
| BFS（幅優先） |  | ✓ |
| タスクキュー |  | ✓ |
| 待ち行列 |  | ✓ |

スタックは最も基本的なデータ構造の一つで、プログラミングのあらゆる場面で使われる。LIFO（後入れ先出し）という単純な原理だが、関数呼び出し、式の評価、履歴管理など実用性が高い。すべての操作がO(1)で効率的。プログラマーなら必ず理解すべき重要なデータ構造。
