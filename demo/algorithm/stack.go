package algorithm

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

// Clear はスタックをクリア
func (s *IntStack) Clear() {
	s.items = []int{}
}

// DemoStack はスタックのデモを実行
func DemoStack() {
	fmt.Println("=================================")
	fmt.Println("1. 基本的なスタック操作")
	fmt.Println("=================================")

	stack := NewIntStack()

	// Push操作
	fmt.Println("\n■ Push操作:")
	values := []int{10, 20, 30, 40, 50}
	for _, v := range values {
		stack.Push(v)
		fmt.Printf("Push(%d) → スタック: %v\n", v, stack.items)
	}

	// Peek操作
	fmt.Println("\n■ Peek操作（最上位要素の確認）:")
	if top, ok := stack.Peek(); ok {
		fmt.Printf("最上位要素: %d (サイズ: %d)\n", top, stack.Size())
	}

	// Pop操作
	fmt.Println("\n■ Pop操作:")
	for !stack.IsEmpty() {
		if item, ok := stack.Pop(); ok {
			fmt.Printf("Pop() → %d (残り: %v)\n", item, stack.items)
		}
	}

	// 空スタックの操作
	fmt.Println("\n■ 空スタックの操作:")
	if _, ok := stack.Pop(); !ok {
		fmt.Println("Pop() → スタックが空です")
	}
	if _, ok := stack.Peek(); !ok {
		fmt.Println("Peek() → スタックが空です")
	}

	fmt.Println("\n=================================")
	fmt.Println("2. スタックの実用例: 括弧のマッチング")
	fmt.Println("=================================")

	testCases := []string{
		"(())",
		"(()())",
		"(()",
		"())",
		"{[()]}",
		"{[(])}",
	}

	for _, tc := range testCases {
		result := isValidParentheses(tc)
		status := "✓ 正しい"
		if !result {
			status = "✗ 不正"
		}
		fmt.Printf("'%s' → %s\n", tc, status)
	}

	fmt.Println("\n=================================")
	fmt.Println("3. スタックの実用例: 逆ポーランド記法")
	fmt.Println("=================================")

	rpnExpressions := [][]string{
		{"2", "3", "+"},                    // 2 + 3 = 5
		{"2", "3", "+", "4", "*"},          // (2 + 3) * 4 = 20
		{"5", "1", "2", "+", "4", "*", "+"}, // 5 + ((1 + 2) * 4) = 17
	}

	for _, expr := range rpnExpressions {
		result := evaluateRPN(expr)
		fmt.Printf("%v = %d\n", expr, result)
	}

	fmt.Println("\n=================================")
	fmt.Println("4. スタックの実用例: 文字列の反転")
	fmt.Println("=================================")

	testStrings := []string{
		"hello",
		"world",
		"Go言語",
	}

	for _, str := range testStrings {
		reversed := reverseString(str)
		fmt.Printf("'%s' → '%s'\n", str, reversed)
	}

	fmt.Println("\n=================================")
	fmt.Println("5. スタックの実用例: 関数呼び出しのシミュレーション")
	fmt.Println("=================================")

	simulateFunctionCalls()

	fmt.Println("\n=================================")
	fmt.Println("6. スタックの実用例: パスの簡略化")
	fmt.Println("=================================")

	paths := []string{
		"/home/user/documents/../pictures",
		"/a/./b/../../c/",
		"/a//b////c/d//././/..",
	}

	for _, path := range paths {
		simplified := simplifyPath(path)
		fmt.Printf("'%s'\n  → '%s'\n", path, simplified)
	}
}

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

// simulateFunctionCalls は関数呼び出しスタックをシミュレーション
func simulateFunctionCalls() {
	callStack := NewStack()

	fmt.Println("関数呼び出し:")
	fmt.Println("main() → funcA() → funcB() → funcC()")
	fmt.Println()

	// 関数呼び出し
	callStack.Push("main()")
	fmt.Printf("呼び出し: main()  → スタック: %v\n", getStackItems(callStack))

	callStack.Push("funcA()")
	fmt.Printf("呼び出し: funcA() → スタック: %v\n", getStackItems(callStack))

	callStack.Push("funcB()")
	fmt.Printf("呼び出し: funcB() → スタック: %v\n", getStackItems(callStack))

	callStack.Push("funcC()")
	fmt.Printf("呼び出し: funcC() → スタック: %v\n", getStackItems(callStack))

	fmt.Println("\n関数の終了（return）:")

	// 関数の終了
	for !callStack.IsEmpty() {
		if fn, ok := callStack.Pop(); ok {
			fmt.Printf("終了: %s     → スタック: %v\n", fn, getStackItems(callStack))
		}
	}
}

// getStackItems はスタックの中身を取得（デバッグ用）
func getStackItems(s *Stack) []interface{} {
	return s.items
}

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
