package algorithm

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

// DFSDemo は深さ優先探索のデモを実行
func DFSDemo() {
	fmt.Println("=== 基本的なグラフ探索 ===")
	// グラフの作成
	//     0
	//    / \
	//   1   2
	//  / \   \
	// 3   4   5
	g := NewGraph(6)
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 3)
	g.AddEdge(1, 4)
	g.AddEdge(2, 5)

	// 再帰版DFS
	g.DFS(0)

	// スタック版DFS
	g.DFSIterative(0)

	fmt.Println("\n=== 経路探索 ===")
	// 経路を見つける
	path := g.DFSWithPath(0, 5)
	if path != nil {
		fmt.Printf("経路 0 → 5: %v\n", path)
	}

	// 全ての経路を見つける
	allPaths := g.DFSAllPaths(0, 4)
	fmt.Printf("0 → 4 の全経路:\n")
	for i, p := range allPaths {
		fmt.Printf("  経路%d: %v\n", i+1, p)
	}

	fmt.Println("\n=== 連結成分 ===")
	// 非連結グラフの作成
	g2 := NewGraph(7)
	g2.AddEdge(0, 1)
	g2.AddEdge(1, 2)
	g2.AddEdge(3, 4)
	g2.AddEdge(5, 6)

	components := g2.DFSComponents()
	fmt.Printf("連結成分の数: %d\n", len(components))
	for i, comp := range components {
		fmt.Printf("  成分%d: %v\n", i+1, comp)
	}

	fmt.Println("\n=== 閉路検出 ===")
	// 閉路のないグラフ
	g3 := NewGraph(4)
	g3.AddEdge(0, 1)
	g3.AddEdge(1, 2)
	g3.AddEdge(2, 3)
	fmt.Printf("グラフ1に閉路: %v\n", g3.HasCycle())

	// 閉路のあるグラフ
	g4 := NewGraph(4)
	g4.AddEdge(0, 1)
	g4.AddEdge(1, 2)
	g4.AddEdge(2, 3)
	g4.AddEdge(3, 0)
	fmt.Printf("グラフ2に閉路: %v\n", g4.HasCycle())

	fmt.Println("\n=== トポロジカルソート ===")
	// 有向非巡回グラフ (DAG)
	//   5 → 2 → 3
	//   ↓    ↓
	//   0 → 1
	//   ↓
	//   4
	g5 := NewGraph(6)
	g5.AddDirectedEdge(5, 2)
	g5.AddDirectedEdge(5, 0)
	g5.AddDirectedEdge(2, 3)
	g5.AddDirectedEdge(2, 1)
	g5.AddDirectedEdge(0, 1)
	g5.AddDirectedEdge(0, 4)

	topoSort := g5.TopologicalSort()
	fmt.Printf("トポロジカルソート順序: %v\n", topoSort)
}
