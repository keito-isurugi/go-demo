package fate_servant

import "fmt"

// =============================================================================
// Template Method Pattern: バトルアクションの共通処理
// 戦闘の流れを定義し、具体的な行動はサブクラスで実装
// =============================================================================

// BattleAction はバトルアクションのインターフェース
type BattleAction interface {
	Prepare()
	Execute() int
	Finalize()
}

// BattleTemplate はバトルのテンプレート（Template Method パターン）
type BattleTemplate struct {
	action BattleAction
}

// NewBattleTemplate はバトルテンプレートを生成する
func NewBattleTemplate(action BattleAction) *BattleTemplate {
	return &BattleTemplate{action: action}
}

// PerformAction はテンプレートメソッド（共通の流れを定義）
func (b *BattleTemplate) PerformAction() int {
	fmt.Println("\n========== バトルアクション開始 ==========")
	b.action.Prepare()
	damage := b.action.Execute()
	b.action.Finalize()
	fmt.Printf("========== 与ダメージ: %d ==========\n\n", damage)
	return damage
}

// =============================================================================
// 具体的なバトルアクション実装
// =============================================================================

// NormalAttackAction は通常攻撃アクション
type NormalAttackAction struct {
	servant Servant
}

// NewNormalAttackAction は通常攻撃アクションを生成
func NewNormalAttackAction(servant Servant) *NormalAttackAction {
	return &NormalAttackAction{servant: servant}
}

func (a *NormalAttackAction) Prepare() {
	fmt.Printf("%s（%s）が構えを取る...\n", a.servant.GetTrueName(), a.servant.GetClass())
}

func (a *NormalAttackAction) Execute() int {
	return a.servant.Attack()
}

func (a *NormalAttackAction) Finalize() {
	fmt.Println("攻撃完了。体勢を整える。")
}

// NoblePhantasmAction は宝具使用アクション
type NoblePhantasmAction struct {
	servant Servant
}

// NewNoblePhantasmAction は宝具アクションを生成
func NewNoblePhantasmAction(servant Servant) *NoblePhantasmAction {
	return &NoblePhantasmAction{servant: servant}
}

func (a *NoblePhantasmAction) Prepare() {
	fmt.Printf("%sが魔力を集中させる...\n", a.servant.GetTrueName())
	fmt.Println("令呪を以て命ずる！宝具を解放せよ！")
}

func (a *NoblePhantasmAction) Execute() int {
	return a.servant.UseNoblePhantasm()
}

func (a *NoblePhantasmAction) Finalize() {
	fmt.Println("魔力が霧散していく...")
}

// ComboAttackAction はコンボ攻撃アクション
type ComboAttackAction struct {
	servant Servant
	hits    int
}

// NewComboAttackAction はコンボ攻撃アクションを生成
func NewComboAttackAction(servant Servant, hits int) *ComboAttackAction {
	return &ComboAttackAction{servant: servant, hits: hits}
}

func (a *ComboAttackAction) Prepare() {
	fmt.Printf("%sがコンボ攻撃の構えを取る！\n", a.servant.GetTrueName())
}

func (a *ComboAttackAction) Execute() int {
	totalDamage := 0
	for i := 1; i <= a.hits; i++ {
		fmt.Printf("【%d Hit】", i)
		damage := a.servant.Attack() / 2 // コンボは1発あたりのダメージ減
		totalDamage += damage
	}
	fmt.Printf("コンボ完了！%d Hits!\n", a.hits)
	return totalDamage
}

func (a *ComboAttackAction) Finalize() {
	fmt.Println("見事なコンボ攻撃だった！")
}

// =============================================================================
// バトルシステム（複数のパターンを組み合わせ）
// =============================================================================

// BattleSystem は戦闘を管理するシステム
type BattleSystem struct {
	attacker Servant
	defender Servant
}

// NewBattleSystem はバトルシステムを生成
func NewBattleSystem(attacker, defender Servant) *BattleSystem {
	return &BattleSystem{
		attacker: attacker,
		defender: defender,
	}
}

// ExecuteTurn はターンを実行する
func (b *BattleSystem) ExecuteTurn(useNoblePhantasm bool) int {
	var action BattleAction
	if useNoblePhantasm {
		action = NewNoblePhantasmAction(b.attacker)
	} else {
		action = NewNormalAttackAction(b.attacker)
	}

	template := NewBattleTemplate(action)
	damage := template.PerformAction()

	fmt.Printf(">>> %s が %s に %d ダメージ！ <<<\n",
		b.attacker.GetTrueName(),
		b.defender.GetTrueName(),
		damage)

	return damage
}

// SwapRoles は攻守を入れ替える
func (b *BattleSystem) SwapRoles() {
	b.attacker, b.defender = b.defender, b.attacker
	fmt.Printf("\n--- 攻守交代: %s の番 ---\n", b.attacker.GetTrueName())
}
