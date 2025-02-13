package main

import "fmt"

// このインターフェースを使って異なる攻撃方法を実装する
// ポケモンの攻撃方法はこのインターフェースを実装した内容によって決まる
type AttackStrategy interface {
	Attack() string
}

// 物理攻撃(AttackStrategyインターフェースを実装)
type PhysicalAttack struct{}
func (pa *PhysicalAttack) Attack() string {
	return "Physical Attack!"
}

// 特殊攻撃(AttackStrategyインターフェースを実装)
type SpecialAttack struct{}
func (sa *SpecialAttack) Attack() string {
	return "Special Attack!"
}

// ポケモン構造体
// AttackStrategy を差し替えることで、攻撃方法を動的に変更できる
type Pokemon struct {
	name   string
	attack AttackStrategy
}

// ポケモンを生成するコンストラクタ
func NewPokemon(name string, attack AttackStrategy) *Pokemon {
	return &Pokemon{
		name:   name,
		attack: attack,
	}
}

// ポケモンの攻撃を表示
func (p *Pokemon) PerformAttack() {
	fmt.Printf("%s uses: %s\n", p.name, p.attack.Attack())
}

// ポケモンの攻撃方法を変更
func (p *Pokemon) ChangeStrategy(newStrategy AttackStrategy) {
	p.attack = newStrategy
}

// 実行
func storategyExec() {
	// 物理攻撃を持つポケモン「ピカチュウ」を作成
	pikachu := NewPokemon("Pikachu", &PhysicalAttack{})
	pikachu.PerformAttack()

	// 特殊攻撃に変更
	pikachu.ChangeStrategy(&SpecialAttack{})
	pikachu.PerformAttack()
}

// 出力
// Pikachu uses: Physical Attack!
// Pikachu uses: Special Attack!