package fate_servant

// BattleAction はバトルアクションのインターフェース（Template Method Pattern）
type BattleAction interface {
	Prepare()
	Execute() int
	Cleanup()
}

// BattleTemplate はバトルのテンプレート
type BattleTemplate struct {
	action BattleAction
}

// NewBattleTemplate はバトルテンプレートを生成
func NewBattleTemplate(action BattleAction) *BattleTemplate {
	return &BattleTemplate{action: action}
}

// Run はテンプレートメソッド
func (b *BattleTemplate) Run() int {
	b.action.Prepare()
	damage := b.action.Execute()
	b.action.Cleanup()
	return damage
}

// AttackAction は通常攻撃アクション
type AttackAction struct {
	servant Servant
}

func NewAttackAction(s Servant) *AttackAction {
	return &AttackAction{servant: s}
}

func (a *AttackAction) Prepare() {}
func (a *AttackAction) Execute() int { return a.servant.Attack() }
func (a *AttackAction) Cleanup() {}

// NoblePhantasmAction は宝具アクション
type NoblePhantasmAction struct {
	servant Servant
}

func NewNoblePhantasmAction(s Servant) *NoblePhantasmAction {
	return &NoblePhantasmAction{servant: s}
}

func (a *NoblePhantasmAction) Prepare() {}
func (a *NoblePhantasmAction) Execute() int { return a.servant.UseNoblePhantasm() }
func (a *NoblePhantasmAction) Cleanup() {}
