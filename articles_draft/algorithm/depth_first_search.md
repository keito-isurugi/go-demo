## 概要

- Goで深さ優先探索（DFS: Depth-First Search）を実装
- DFSはグラフやツリーを探索する基本的なアルゴリズム。できるだけ深く進んでから戻る探索戦略を取る。

### 特徴

- **深さ優先**: 一つの経路を可能な限り深く探索してから次の経路へ
- **2つの実装方法**: 再帰とスタック（反復）
- **空間計算量**: O(h) - hは探索の深さ
- **時間計算量**: O(V + E) - Vは頂点数、Eは辺数
- **応用範囲**: 経路探索、連結成分、閉路検出、トポロジカルソートなど

## 深さ優先探索とは

### アルゴリズムの基本動作

DFSは以下の手順で動作：

1. 開始頂点を訪問済みにする
2. 隣接する未訪問の頂点があれば、その頂点に移動して再帰的に探索
3. すべての隣接頂点を探索したら、一つ前の頂点に戻る（バックトラック）
4. 未訪問の頂点がなくなるまで繰り返す

### 具体例

```
グラフ:
    0
   / \
  1   2
 / \   \
3   4   5

DFS探索順序（頂点0から開始）:
0 → 1 → 3 → 4 → 2 → 5
```

**探索の流れ**:
1. 0から開始
2. 0の隣接頂点1を訪問
3. 1の隣接頂点3を訪問（深さ優先）
4. 3に隣接する未訪問頂点なし → バックトラックして1へ
5. 1の隣接頂点4を訪問
6. 4に隣接する未訪問頂点なし → バックトラックして0へ
7. 0の隣接頂点2を訪問
8. 2の隣接頂点5を訪問

### BFS（幅優先探索）との違い

| 特徴 | DFS（深さ優先） | BFS（幅優先） |
| --- | --- | --- |
| 探索戦略 | できるだけ深く | 同じ深さを先に |
| データ構造 | スタック（再帰） | キュー |
| 空間計算量 | O(h) 深さ | O(w) 幅 |
| 最短経路 | 保証なし | 保証あり |
| 用途 | 経路探索、閉路検出 | 最短経路、レベル探索 |

## サンプルコード

### グラフの基本構造

```go
package main

import "fmt"

// グラフを隣接リストで表現
type Graph struct {
	vertices int
	adjList  map[int][]int
}

// NewGraph はグラフを初期化
func NewGraph(vertices int) *Graph {
	return &Graph{
		vertices: vertices,
		adjList:  make(map[int][]int),
	}
}

// AddEdge は辺を追加（無向グラフ）
func (g *Graph) AddEdge(u, v int) {
	g.adjList[u] = append(g.adjList[u], v)
	g.adjList[v] = append(g.adjList[v], u)
}

// AddDirectedEdge は辺を追加（有向グラフ）
func (g *Graph) AddDirectedEdge(u, v int) {
	g.adjList[u] = append(g.adjList[u], v)
}
```

### 基本的なDFS実装（再帰版）

```go
// DFS は深さ優先探索（再帰版）
func (g *Graph) DFS(start int) {
	visited := make(map[int]bool)
	fmt.Println("DFS (再帰版):")
	g.dfsRecursive(start, visited)
	fmt.Println()
}

// dfsRecursive は再帰的な深さ優先探索の実装
func (g *Graph) dfsRecursive(v int, visited map[int]bool) {
	// 現在の頂点を訪問済みにする
	visited[v] = true
	fmt.Printf("%d ", v)

	// 隣接する頂点を再帰的に訪問
	for _, neighbor := range g.adjList[v] {
		if !visited[neighbor] {
			g.dfsRecursive(neighbor, visited)
		}
	}
}
```

**ポイント**:
- `visited`マップで訪問済み頂点を管理
- 再帰呼び出しで深さ優先探索を実現
- 各頂点は1回だけ訪問される

### スタックを使ったDFS（反復版）

```go
// DFSIterative は深さ優先探索（スタック版）
func (g *Graph) DFSIterative(start int) {
	visited := make(map[int]bool)
	stack := []int{start}

	fmt.Println("DFS (スタック版):")
	for len(stack) > 0 {
		// スタックから取り出す（後入れ先出し）
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if !visited[v] {
			visited[v] = true
			fmt.Printf("%d ", v)

			// 隣接する頂点をスタックに追加（逆順で追加して順序を保つ）
			for i := len(g.adjList[v]) - 1; i >= 0; i-- {
				neighbor := g.adjList[v][i]
				if !visited[neighbor] {
					stack = append(stack, neighbor)
				}
			}
		}
	}
	fmt.Println()
}
```

**ポイント**:
- スタック（LIFO: Last In First Out）を使用
- 再帰版と同じ動作を反復で実現
- スタックオーバーフローのリスクを回避

### 経路探索

```go
// DFSWithPath は経路を記録する深さ優先探索
func (g *Graph) DFSWithPath(start, goal int) []int {
	visited := make(map[int]bool)
	path := []int{}

	if g.dfsPathRecursive(start, goal, visited, &path) {
		return path
	}
	return nil
}

// dfsPathRecursive は経路を記録する再帰的DFS
func (g *Graph) dfsPathRecursive(v, goal int, visited map[int]bool, path *[]int) bool {
	visited[v] = true
	*path = append(*path, v)

	// ゴールに到達
	if v == goal {
		return true
	}

	// 隣接する頂点を探索
	for _, neighbor := range g.adjList[v] {
		if !visited[neighbor] {
			if g.dfsPathRecursive(neighbor, goal, visited, path) {
				return true
			}
		}
	}

	// この経路では見つからなかったので削除
	*path = (*path)[:len(*path)-1]
	return false
}
```

**使用例**:
```go
g := NewGraph(6)
g.AddEdge(0, 1)
g.AddEdge(1, 2)
g.AddEdge(2, 3)

path := g.DFSWithPath(0, 3)
// 結果: [0, 1, 2, 3]
```

### 全ての経路を見つける

```go
// DFSAllPaths は全ての経路を見つける深さ優先探索
func (g *Graph) DFSAllPaths(start, goal int) [][]int {
	visited := make(map[int]bool)
	path := []int{}
	allPaths := [][]int{}

	g.dfsAllPathsRecursive(start, goal, visited, path, &allPaths)
	return allPaths
}

// dfsAllPathsRecursive は全経路を記録する再帰的DFS
func (g *Graph) dfsAllPathsRecursive(v, goal int, visited map[int]bool, path []int, allPaths *[][]int) {
	visited[v] = true
	path = append(path, v)

	// ゴールに到達
	if v == goal {
		// 経路をコピーして保存
		pathCopy := make([]int, len(path))
		copy(pathCopy, path)
		*allPaths = append(*allPaths, pathCopy)
	} else {
		// 隣接する頂点を探索
		for _, neighbor := range g.adjList[v] {
			if !visited[neighbor] {
				g.dfsAllPathsRecursive(neighbor, goal, visited, path, allPaths)
			}
		}
	}

	// バックトラック
	visited[v] = false
}
```

**重要**: 全経路探索では`visited[v] = false`でバックトラックが必要。これにより同じ頂点を別の経路で再訪問できる。

### 連結成分の検出

```go
// DFSComponents は連結成分を見つける
func (g *Graph) DFSComponents() [][]int {
	visited := make(map[int]bool)
	components := [][]int{}

	for v := 0; v < g.vertices; v++ {
		if !visited[v] {
			component := []int{}
			g.dfsComponentRecursive(v, visited, &component)
			components = append(components, component)
		}
	}

	return components
}

// dfsComponentRecursive は連結成分を探索
func (g *Graph) dfsComponentRecursive(v int, visited map[int]bool, component *[]int) {
	visited[v] = true
	*component = append(*component, v)

	for _, neighbor := range g.adjList[v] {
		if !visited[neighbor] {
			g.dfsComponentRecursive(neighbor, visited, component)
		}
	}
}
```

**使用例**:
```go
// 非連結グラフ
g := NewGraph(7)
g.AddEdge(0, 1)
g.AddEdge(1, 2)
g.AddEdge(3, 4)
g.AddEdge(5, 6)

components := g.DFSComponents()
// 結果: [[0, 1, 2], [3, 4], [5, 6]]
```

### 閉路検出（無向グラフ）

```go
// HasCycle は閉路が存在するか判定（無向グラフ）
func (g *Graph) HasCycle() bool {
	visited := make(map[int]bool)

	for v := 0; v < g.vertices; v++ {
		if !visited[v] {
			if g.hasCycleRecursive(v, -1, visited) {
				return true
			}
		}
	}

	return false
}

// hasCycleRecursive は閉路検出の再帰実装
func (g *Graph) hasCycleRecursive(v, parent int, visited map[int]bool) bool {
	visited[v] = true

	for _, neighbor := range g.adjList[v] {
		if !visited[neighbor] {
			if g.hasCycleRecursive(neighbor, v, visited) {
				return true
			}
		} else if neighbor != parent {
			// 訪問済みかつ親でない = 閉路発見
			return true
		}
	}

	return false
}
```

**ポイント**:
- `parent`パラメータで直前の頂点を記録
- 訪問済みの頂点が親でない場合、閉路が存在

### トポロジカルソート（有向グラフ）

```go
// TopologicalSort はトポロジカルソートを実行（有向グラフ）
func (g *Graph) TopologicalSort() []int {
	visited := make(map[int]bool)
	stack := []int{}

	for v := 0; v < g.vertices; v++ {
		if !visited[v] {
			g.topologicalSortRecursive(v, visited, &stack)
		}
	}

	// スタックを逆順にして返す
	result := make([]int, len(stack))
	for i := 0; i < len(stack); i++ {
		result[i] = stack[len(stack)-1-i]
	}

	return result
}

// topologicalSortRecursive はトポロジカルソートの再帰実装
func (g *Graph) topologicalSortRecursive(v int, visited map[int]bool, stack *[]int) {
	visited[v] = true

	for _, neighbor := range g.adjList[v] {
		if !visited[neighbor] {
			g.topologicalSortRecursive(neighbor, visited, stack)
		}
	}

	// 後処理でスタックに追加
	*stack = append(*stack, v)
}
```

**使用例**:
```go
// 有向非巡回グラフ (DAG)
//   5 → 2 → 3
//   ↓    ↓
//   0 → 1
//   ↓
//   4
g := NewGraph(6)
g.AddDirectedEdge(5, 2)
g.AddDirectedEdge(5, 0)
g.AddDirectedEdge(2, 3)
g.AddDirectedEdge(2, 1)
g.AddDirectedEdge(0, 1)
g.AddDirectedEdge(0, 4)

topoSort := g.TopologicalSort()
// 結果: [5, 2, 0, 3, 1, 4] または [5, 0, 2, 4, 3, 1] など
```

## 計算量

### 時間計算量

| 操作 | 時間計算量 | 説明 |
| --- | --- | --- |
| DFS（全頂点） | O(V + E) | 各頂点と各辺を1回ずつ訪問 |
| 経路探索 | O(V + E) | 最悪でも全グラフを探索 |
| 連結成分 | O(V + E) | 全頂点を探索 |
| 閉路検出 | O(V + E) | 全頂点を探索 |
| トポロジカルソート | O(V + E) | 全頂点と辺を処理 |

- V: 頂点数（Vertices）
- E: 辺数（Edges）

### 空間計算量

| 実装方法 | 空間計算量 | 説明 |
| --- | --- | --- |
| 再帰版DFS | O(h) | 再帰スタックの深さh |
| スタック版DFS | O(V) | 明示的なスタック |
| visited配列 | O(V) | 全頂点分の記録 |

- h: グラフの最大深さ（最悪でV）

### パフォーマンス特性

**DFSが有利な場面**:
- グラフが疎（辺が少ない）
- 深い探索が必要
- メモリが限られている

**BFSが有利な場面**:
- 最短経路が必要
- レベル単位の探索
- 幅が狭いグラフ

## 動作原理の詳細

### 再帰とスタックの関係

DFSの再帰版とスタック版は本質的に同じ：

**再帰版**:
```
dfs(0):
  visit 0
  dfs(1):
    visit 1
    dfs(3):
      visit 3
    dfs(4):
      visit 4
  dfs(2):
    visit 2
```

**スタック版**:
```
スタック: [0]
→ pop 0, visit 0, push [1, 2]
スタック: [1, 2]
→ pop 2, visit 2, push [5]
スタック: [1, 5]
→ pop 5, visit 5
スタック: [1]
→ pop 1, visit 1, push [3, 4]
...
```

両方とも同じ順序で頂点を訪問するが、スタック版は明示的にスタックを管理する。

### バックトラックの仕組み

全経路探索でのバックトラック：

```go
visited[v] = true    // 訪問済みにする
path = append(path, v)

// 探索...

visited[v] = false   // 未訪問に戻す（バックトラック）
```

**なぜ必要か**:
- 一つの経路で訪問済みにした頂点を、別の経路でも訪問できるようにする
- 全ての可能な経路を探索するため

**例**:
```
グラフ: 0 - 1 - 3
        |   |
        +- 2 +

0から3への全経路:
1. 0 → 1 → 3
2. 0 → 2 → 1 → 3

経路1で1を訪問済みにしても、
バックトラックで未訪問に戻すことで経路2も見つけられる
```

## 使いどころ

### 向いてる場面

- **迷路探索**: 一つの経路を深く探索
- **パズルソルバー**: 数独、8クイーン問題など
- **依存関係の解決**: トポロジカルソート
- **グラフ連結性**: 連結成分の検出
- **閉路検出**: デッドロック検出、DAG判定

### 実例1: 迷路探索

```go
type Maze struct {
	grid [][]int
	rows, cols int
}

func (m *Maze) solveDFS(r, c int, visited [][]bool, path *[][2]int) bool {
	// ゴールに到達
	if r == m.rows-1 && c == m.cols-1 {
		*path = append(*path, [2]int{r, c})
		return true
	}

	// 訪問済みまたは壁
	if visited[r][c] || m.grid[r][c] == 1 {
		return false
	}

	visited[r][c] = true
	*path = append(*path, [2]int{r, c})

	// 4方向を探索
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for _, d := range dirs {
		nr, nc := r+d[0], c+d[1]
		if nr >= 0 && nr < m.rows && nc >= 0 && nc < m.cols {
			if m.solveDFS(nr, nc, visited, path) {
				return true
			}
		}
	}

	// バックトラック
	*path = (*path)[:len(*path)-1]
	return false
}
```

### 実例2: ファイルシステムの探索

```go
type FileNode struct {
	name     string
	isDir    bool
	children []*FileNode
}

func (f *FileNode) findFile(target string) *FileNode {
	if f.name == target {
		return f
	}

	if f.isDir {
		for _, child := range f.children {
			if result := child.findFile(target); result != nil {
				return result
			}
		}
	}

	return nil
}

func (f *FileNode) printTree(depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	fmt.Printf("%s%s\n", indent, f.name)

	if f.isDir {
		for _, child := range f.children {
			child.printTree(depth + 1)
		}
	}
}
```

### 実例3: 島の数を数える

2Dグリッド上で連結した陸地（島）の数を数える：

```go
func numIslands(grid [][]byte) int {
	if len(grid) == 0 {
		return 0
	}

	rows, cols := len(grid), len(grid[0])
	count := 0

	var dfs func(r, c int)
	dfs = func(r, c int) {
		// 範囲外または水
		if r < 0 || r >= rows || c < 0 || c >= cols || grid[r][c] == '0' {
			return
		}

		// 訪問済みにする（陸地を水に変更）
		grid[r][c] = '0'

		// 4方向を探索
		dfs(r-1, c)
		dfs(r+1, c)
		dfs(r, c-1)
		dfs(r, c+1)
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == '1' {
				count++
				dfs(r, c)
			}
		}
	}

	return count
}

// 使用例
grid := [][]byte{
	{'1', '1', '0', '0', '0'},
	{'1', '1', '0', '0', '0'},
	{'0', '0', '1', '0', '0'},
	{'0', '0', '0', '1', '1'},
}
// 結果: 3つの島
```

## 関連する問題

### 1. 強連結成分（Strongly Connected Components）

有向グラフの強連結成分を見つける（Kosaraju's Algorithm）：

```go
func (g *Graph) SCCKosaraju() [][]int {
	// ステップ1: DFSで終了時刻順にスタックに積む
	visited := make(map[int]bool)
	stack := []int{}

	for v := 0; v < g.vertices; v++ {
		if !visited[v] {
			g.fillOrder(v, visited, &stack)
		}
	}

	// ステップ2: グラフを転置
	transpose := g.getTranspose()

	// ステップ3: スタックの順序でDFS
	visited = make(map[int]bool)
	scc := [][]int{}

	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if !visited[v] {
			component := []int{}
			transpose.dfsComponentRecursive(v, visited, &component)
			scc = append(scc, component)
		}
	}

	return scc
}
```

### 2. 橋と関節点

グラフの橋（bridge）や関節点（articulation point）を見つける：

```go
func (g *Graph) findBridges() [][2]int {
	visited := make(map[int]bool)
	disc := make(map[int]int)    // 発見時刻
	low := make(map[int]int)     // 到達可能な最小時刻
	parent := make(map[int]int)
	bridges := [][2]int{}
	time := 0

	var dfs func(u int)
	dfs = func(u int) {
		visited[u] = true
		disc[u] = time
		low[u] = time
		time++

		for _, v := range g.adjList[u] {
			if !visited[v] {
				parent[v] = u
				dfs(v)

				// 子の最小到達時刻で更新
				if low[v] < low[u] {
					low[u] = low[v]
				}

				// 橋の条件: low[v] > disc[u]
				if low[v] > disc[u] {
					bridges = append(bridges, [2]int{u, v})
				}
			} else if v != parent[u] {
				// バックエッジ
				if disc[v] < low[u] {
					low[u] = disc[v]
				}
			}
		}
	}

	for v := 0; v < g.vertices; v++ {
		if !visited[v] {
			dfs(v)
		}
	}

	return bridges
}
```

### 3. 有向グラフの閉路検出

有向グラフでの閉路検出（バックエッジを検出）：

```go
func (g *Graph) hasCycleDirected() bool {
	visited := make(map[int]bool)
	recStack := make(map[int]bool) // 再帰スタック

	var dfs func(v int) bool
	dfs = func(v int) bool {
		visited[v] = true
		recStack[v] = true

		for _, neighbor := range g.adjList[v] {
			if !visited[neighbor] {
				if dfs(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				// 再帰スタック上の頂点 → 閉路発見
				return true
			}
		}

		recStack[v] = false
		return false
	}

	for v := 0; v < g.vertices; v++ {
		if !visited[v] {
			if dfs(v) {
				return true
			}
		}
	}

	return false
}
```

## 最適化テクニック

### 1. メモ化（Memoization）

重複計算を避けるためのメモ化：

```go
func (g *Graph) pathCountMemo(start, goal int, memo map[int]int) int {
	if start == goal {
		return 1
	}

	if count, ok := memo[start]; ok {
		return count
	}

	total := 0
	for _, neighbor := range g.adjList[start] {
		total += g.pathCountMemo(neighbor, goal, memo)
	}

	memo[start] = total
	return total
}
```

### 2. 早期終了

条件を満たしたら即座に終了：

```go
func (g *Graph) hasPath(start, goal int) bool {
	visited := make(map[int]bool)

	var dfs func(v int) bool
	dfs = func(v int) bool {
		if v == goal {
			return true // 早期終了
		}

		visited[v] = true

		for _, neighbor := range g.adjList[v] {
			if !visited[neighbor] {
				if dfs(neighbor) {
					return true
				}
			}
		}

		return false
	}

	return dfs(start)
}
```

### 3. 反復深化深さ優先探索（IDDFS）

深さ制限を徐々に増やすことで、BFSの最短経路とDFSの省メモリを両立：

```go
func (g *Graph) IDDFS(start, goal, maxDepth int) bool {
	for depth := 0; depth <= maxDepth; depth++ {
		visited := make(map[int]bool)
		if g.dfsLimited(start, goal, depth, visited) {
			return true
		}
	}
	return false
}

func (g *Graph) dfsLimited(v, goal, depth int, visited map[int]bool) bool {
	if v == goal {
		return true
	}

	if depth <= 0 {
		return false
	}

	visited[v] = true

	for _, neighbor := range g.adjList[v] {
		if !visited[neighbor] {
			if g.dfsLimited(neighbor, goal, depth-1, visited) {
				return true
			}
		}
	}

	return false
}
```

## デバッグとテスト

### テストケース

```go
func TestDFS() {
	// テスト1: 基本的な探索
	g1 := NewGraph(4)
	g1.AddEdge(0, 1)
	g1.AddEdge(0, 2)
	g1.AddEdge(1, 3)

	// 再帰版とスタック版の結果が一致するか
	// ...

	// テスト2: 閉路検出
	g2 := NewGraph(4)
	g2.AddEdge(0, 1)
	g2.AddEdge(1, 2)
	g2.AddEdge(2, 3)
	g2.AddEdge(3, 0)

	if !g2.HasCycle() {
		fmt.Println("FAIL: 閉路を検出できない")
	}

	// テスト3: 連結成分
	g3 := NewGraph(6)
	g3.AddEdge(0, 1)
	g3.AddEdge(2, 3)
	g3.AddEdge(4, 5)

	components := g3.DFSComponents()
	if len(components) != 3 {
		fmt.Printf("FAIL: 期待3個、実際%d個\n", len(components))
	}
}
```

### デバッグ用可視化

```go
// DFSの探索過程を可視化
func (g *Graph) DFSVisualize(start int) {
	visited := make(map[int]bool)
	depth := 0

	var dfs func(v, d int)
	dfs = func(v, d int) {
		visited[v] = true

		indent := ""
		for i := 0; i < d; i++ {
			indent += "  "
		}
		fmt.Printf("%s訪問: %d (深さ%d)\n", indent, v, d)

		for _, neighbor := range g.adjList[v] {
			if !visited[neighbor] {
				fmt.Printf("%s→ %d に移動\n", indent, neighbor)
				dfs(neighbor, d+1)
				fmt.Printf("%s← %d に戻る\n", indent, v)
			} else {
				fmt.Printf("%s× %d は訪問済み\n", indent, neighbor)
			}
		}
	}

	dfs(start, 0)
}
```

**出力例**:
```
訪問: 0 (深さ0)
→ 1 に移動
  訪問: 1 (深さ1)
  → 3 に移動
    訪問: 3 (深さ2)
  ← 1 に戻る
  → 4 に移動
    訪問: 4 (深さ2)
  ← 1 に戻る
← 0 に戻る
→ 2 に移動
  訪問: 2 (深さ1)
```

## まとめ

### メリット

- 実装がシンプル（特に再帰版）
- 空間計算量O(h)で省メモリ
- 経路探索、閉路検出など多様な応用
- バックトラックとの相性が良い
- トポロジカルソートに適している

### デメリット

- 最短経路の保証なし
- 深いグラフでスタックオーバーフローのリスク
- 訪問順序が直感的でない場合がある

### 使うべき時

- **経路探索**: 一つでも経路が見つかればよい
- **連結性判定**: グラフが連結か、成分はいくつか
- **閉路検出**: 有向・無向グラフの閉路
- **トポロジカルソート**: タスクの依存関係
- **迷路・パズル**: バックトラックが必要な問題

### アルゴリズム選択のガイド

| 要件 | 推奨アルゴリズム | 理由 |
| --- | --- | --- |
| 最短経路 | BFS | 重み無しグラフで最短保証 |
| 経路の存在 | DFS | 省メモリ、実装簡単 |
| 全経路列挙 | DFS + バックトラック | バックトラックとの相性 |
| レベル探索 | BFS | 同じ深さを先に処理 |
| トポロジカルソート | DFS | 後処理でスタックに積む |
| 連結成分 | DFS または BFS | どちらでも可 |
| 閉路検出 | DFS | バックエッジ検出が容易 |

深さ優先探索は、グラフアルゴリズムの基礎であり、様々な応用問題の解法の核となる重要なアルゴリズム。再帰とスタックの関係、バックトラックの仕組みを理解することで、複雑な問題にも対応できる。
