package recursion

import "fmt"

type HanoiMove struct {
	From string
	To   string
	Disk int
}

func TowerOfHanoi(n int, from, to, aux string) []HanoiMove {
	moves := []HanoiMove{}
	hanoiHelper(n, from, to, aux, &moves)
	return moves
}

func hanoiHelper(n int, from, to, aux string, moves *[]HanoiMove) {
	if n == 1 {
		*moves = append(*moves, HanoiMove{From: from, To: to, Disk: n})
		return
	}
	
	hanoiHelper(n-1, from, aux, to, moves)
	
	*moves = append(*moves, HanoiMove{From: from, To: to, Disk: n})
	
	hanoiHelper(n-1, aux, to, from, moves)
}

func TowerOfHanoiWithSteps(n int, from, to, aux string, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	
	fmt.Printf("%sHanoi(%d, %s, %s, %s) を呼び出し\n", indent, n, from, to, aux)
	
	if n == 1 {
		fmt.Printf("%s  ディスク1を%sから%sへ移動\n", indent, from, to)
		return
	}
	
	fmt.Printf("%s  ステップ1: %d枚のディスクを%sから%sへ（%s経由）\n", indent, n-1, from, aux, to)
	TowerOfHanoiWithSteps(n-1, from, aux, to, depth+1)
	
	fmt.Printf("%s  ステップ2: ディスク%dを%sから%sへ移動\n", indent, n, from, to)
	
	fmt.Printf("%s  ステップ3: %d枚のディスクを%sから%sへ（%s経由）\n", indent, n-1, aux, to, from)
	TowerOfHanoiWithSteps(n-1, aux, to, from, depth+1)
}

func CountHanoiMoves(n int) int {
	return (1 << n) - 1
}

func PrintHanoiState(towers map[string][]int) {
	maxHeight := 0
	for _, tower := range towers {
		if len(tower) > maxHeight {
			maxHeight = len(tower)
		}
	}
	
	for level := maxHeight - 1; level >= 0; level-- {
		for _, peg := range []string{"A", "B", "C"} {
			if level < len(towers[peg]) {
				fmt.Printf("  %d  ", towers[peg][level])
			} else {
				fmt.Print("  |  ")
			}
		}
		fmt.Println()
	}
	fmt.Println("-----+-----+-----")
	fmt.Println("  A     B     C  ")
	fmt.Println()
}

func SimulateTowerOfHanoi(n int) {
	towers := map[string][]int{
		"A": make([]int, 0),
		"B": make([]int, 0),
		"C": make([]int, 0),
	}
	
	for i := n; i >= 1; i-- {
		towers["A"] = append(towers["A"], i)
	}
	
	fmt.Printf("初期状態（%d枚のディスク）:\n", n)
	PrintHanoiState(towers)
	
	moves := TowerOfHanoi(n, "A", "C", "B")
	
	for i, move := range moves {
		fmt.Printf("移動 %d: ディスク%dを%sから%sへ\n", i+1, move.Disk, move.From, move.To)
		
		fromTower := towers[move.From]
		disk := fromTower[len(fromTower)-1]
		towers[move.From] = fromTower[:len(fromTower)-1]
		towers[move.To] = append(towers[move.To], disk)
		
		PrintHanoiState(towers)
	}
	
	fmt.Printf("\n完了！ 総移動回数: %d\n", len(moves))
}

func TowerOfHanoiIterative(n int) []HanoiMove {
	moves := []HanoiMove{}
	
	if n%2 == 0 {
		moves = iterativeHanoi(n, "A", "B", "C")
	} else {
		moves = iterativeHanoi(n, "A", "C", "B")
	}
	
	return moves
}

func iterativeHanoi(n int, source, auxiliary, destination string) []HanoiMove {
	moves := []HanoiMove{}
	totalMoves := (1 << n) - 1
	
	towers := map[string][]int{
		source:      make([]int, 0),
		auxiliary:   make([]int, 0),
		destination: make([]int, 0),
	}
	
	for i := n; i >= 1; i-- {
		towers[source] = append(towers[source], i)
	}
	
	pegs := []string{source, auxiliary, destination}
	
	for i := 1; i <= totalMoves; i++ {
		from := pegs[(i&(i-1))%3]
		to := pegs[((i|(i-1))+1)%3]
		
		if len(towers[from]) > 0 && 
		   (len(towers[to]) == 0 || towers[from][len(towers[from])-1] < towers[to][len(towers[to])-1]) {
			disk := towers[from][len(towers[from])-1]
			towers[from] = towers[from][:len(towers[from])-1]
			towers[to] = append(towers[to], disk)
			moves = append(moves, HanoiMove{From: from, To: to, Disk: disk})
		} else {
			disk := towers[to][len(towers[to])-1]
			towers[to] = towers[to][:len(towers[to])-1]
			towers[from] = append(towers[from], disk)
			moves = append(moves, HanoiMove{From: to, To: from, Disk: disk})
		}
	}
	
	return moves
}