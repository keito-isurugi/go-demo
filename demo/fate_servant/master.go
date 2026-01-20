package fate_servant

import "fmt"

// =============================================================================
// Observer Pattern: マスターとサーヴァントの通信
// マスターからの命令をサーヴァントが受け取る
// =============================================================================

// Command はマスターからの命令を表す
type Command struct {
	Type    CommandType
	Message string
}

// CommandType は命令の種類
type CommandType int

const (
	CommandAttack CommandType = iota
	CommandDefend
	CommandUseNoblePhantasm
	CommandRetreat
	CommandCommandSpell // 令呪使用
)

func (c CommandType) String() string {
	switch c {
	case CommandAttack:
		return "攻撃"
	case CommandDefend:
		return "防御"
	case CommandUseNoblePhantasm:
		return "宝具解放"
	case CommandRetreat:
		return "撤退"
	case CommandCommandSpell:
		return "令呪発動"
	default:
		return "不明"
	}
}

// ServantObserver はサーヴァント側のObserverインターフェース
type ServantObserver interface {
	OnCommand(cmd Command)
	GetServant() Servant
}

// MasterSubject はマスター側のSubjectインターフェース
type MasterSubject interface {
	RegisterServant(observer ServantObserver)
	RemoveServant(observer ServantObserver)
	IssueCommand(cmd Command)
	UseCommandSpell(cmd Command) // 令呪使用
}

// =============================================================================
// Master の実装
// =============================================================================

// Master はマスターを表す
type Master struct {
	Name          string
	CommandSpells int // 残り令呪数
	servants      []ServantObserver
}

// NewMaster はマスターを生成する
func NewMaster(name string) *Master {
	return &Master{
		Name:          name,
		CommandSpells: 3, // 初期令呪は3画
		servants:      make([]ServantObserver, 0),
	}
}

// RegisterServant はサーヴァントを登録する
func (m *Master) RegisterServant(observer ServantObserver) {
	m.servants = append(m.servants, observer)
	fmt.Printf("[契約] %sが%sと契約を結んだ\n", m.Name, observer.GetServant().GetTrueName())
}

// RemoveServant はサーヴァントを解除する
func (m *Master) RemoveServant(observer ServantObserver) {
	for i, s := range m.servants {
		if s == observer {
			m.servants = append(m.servants[:i], m.servants[i+1:]...)
			fmt.Printf("[契約解除] %sとの契約が解除された\n", observer.GetServant().GetTrueName())
			return
		}
	}
}

// IssueCommand は命令を発行する（全サーヴァントに通知）
func (m *Master) IssueCommand(cmd Command) {
	fmt.Printf("\n[マスター: %s] 命令を発行: %s\n", m.Name, cmd.Type)
	if cmd.Message != "" {
		fmt.Printf("「%s」\n", cmd.Message)
	}
	for _, servant := range m.servants {
		servant.OnCommand(cmd)
	}
}

// UseCommandSpell は令呪を使用する（強制命令）
func (m *Master) UseCommandSpell(cmd Command) {
	if m.CommandSpells <= 0 {
		fmt.Printf("[%s] 令呪が残っていない！\n", m.Name)
		return
	}

	m.CommandSpells--
	fmt.Printf("\n[令呪発動] %sが令呪を使用！（残り%d画）\n", m.Name, m.CommandSpells)
	fmt.Println("令呪を以て命ずる！")

	cmd.Type = CommandCommandSpell
	for _, servant := range m.servants {
		servant.OnCommand(cmd)
	}
}

// GetRemainingCommandSpells は残り令呪数を返す
func (m *Master) GetRemainingCommandSpells() int {
	return m.CommandSpells
}

// =============================================================================
// ContractedServant の実装（Observer側）
// =============================================================================

// ContractedServant は契約されたサーヴァント
type ContractedServant struct {
	servant Servant
}

// NewContractedServant は契約サーヴァントを生成する
func NewContractedServant(servant Servant) *ContractedServant {
	return &ContractedServant{servant: servant}
}

// GetServant はサーヴァントを返す
func (c *ContractedServant) GetServant() Servant {
	return c.servant
}

// OnCommand はマスターからの命令を受け取る
func (c *ContractedServant) OnCommand(cmd Command) {
	servantName := c.servant.GetTrueName()

	switch cmd.Type {
	case CommandAttack:
		fmt.Printf("[%s] 了解した、マスター。攻撃する！\n", servantName)
		c.servant.Attack()

	case CommandDefend:
		fmt.Printf("[%s] 防御態勢を取る！\n", servantName)

	case CommandUseNoblePhantasm:
		fmt.Printf("[%s] 宝具を解放する！\n", servantName)
		c.servant.UseNoblePhantasm()

	case CommandRetreat:
		fmt.Printf("[%s] 撤退する...不本意だが仕方あるまい\n", servantName)

	case CommandCommandSpell:
		fmt.Printf("[%s] 令呪の力が体を駆け巡る...！\n", servantName)
		if cmd.Message != "" {
			fmt.Printf("[%s] 「%s」...命令を遂行する！\n", servantName, cmd.Message)
		}
		// 令呪の強制力でステータスアップなどの効果を表現可能
	}
}

// =============================================================================
// 特殊なサーヴァントObserver（性格による反応の違い）
// =============================================================================

// ProudServant は誇り高いサーヴァント（ギルガメッシュなど）
type ProudServant struct {
	servant Servant
}

// NewProudServant は誇り高いサーヴァントを生成
func NewProudServant(servant Servant) *ProudServant {
	return &ProudServant{servant: servant}
}

func (p *ProudServant) GetServant() Servant {
	return p.servant
}

func (p *ProudServant) OnCommand(cmd Command) {
	servantName := p.servant.GetTrueName()

	switch cmd.Type {
	case CommandAttack:
		fmt.Printf("[%s] フン、雑種どもを蹴散らしてやろう\n", servantName)
		p.servant.Attack()

	case CommandDefend:
		fmt.Printf("[%s] 王に防御など必要ない。だが...今回だけは聞いてやる\n", servantName)

	case CommandUseNoblePhantasm:
		fmt.Printf("[%s] 慢心せずに全力で叩き潰してやる\n", servantName)
		p.servant.UseNoblePhantasm()

	case CommandRetreat:
		fmt.Printf("[%s] 撤退だと？この王に逃げろと？...貴様、後で話がある\n", servantName)

	case CommandCommandSpell:
		fmt.Printf("[%s] 令呪か...この王を強制するとは度胸がある\n", servantName)
		fmt.Printf("[%s] 良いだろう、その願い叶えてやる\n", servantName)
	}
}

// LoyalServant は忠誠心の高いサーヴァント（セイバーなど）
type LoyalServant struct {
	servant Servant
}

// NewLoyalServant は忠誠心の高いサーヴァントを生成
func NewLoyalServant(servant Servant) *LoyalServant {
	return &LoyalServant{servant: servant}
}

func (l *LoyalServant) GetServant() Servant {
	return l.servant
}

func (l *LoyalServant) OnCommand(cmd Command) {
	servantName := l.servant.GetTrueName()

	switch cmd.Type {
	case CommandAttack:
		fmt.Printf("[%s] 承知しました、マスター。騎士として敵を討ちます！\n", servantName)
		l.servant.Attack()

	case CommandDefend:
		fmt.Printf("[%s] マスターを守るのは騎士の務め。防御を固めます\n", servantName)

	case CommandUseNoblePhantasm:
		fmt.Printf("[%s] 宝具、解放します！\n", servantName)
		l.servant.UseNoblePhantasm()

	case CommandRetreat:
		fmt.Printf("[%s] 了解しました。戦略的撤退を行います\n", servantName)

	case CommandCommandSpell:
		fmt.Printf("[%s] 令呪の力が...！マスターの意志に従います！\n", servantName)
	}
}
