package fate_servant

// CommandType は命令の種類
type CommandType int

const (
	CommandAttack CommandType = iota
	CommandDefend
	CommandUseNoblePhantasm
)

// Command はマスターからの命令
type Command struct {
	Type CommandType
}

// ServantObserver はサーヴァント側のObserver
type ServantObserver interface {
	OnCommand(cmd Command) int
	Servant() Servant
}

// Master はマスター（Observer Pattern の Subject）
type Master struct {
	Name          string
	CommandSpells int
	servants      []ServantObserver
}

// NewMaster はマスターを生成
func NewMaster(name string) *Master {
	return &Master{
		Name:          name,
		CommandSpells: 3,
		servants:      make([]ServantObserver, 0),
	}
}

// Contract はサーヴァントと契約
func (m *Master) Contract(observer ServantObserver) {
	m.servants = append(m.servants, observer)
}

// Command は命令を発行
func (m *Master) Command(cmd Command) map[Servant]int {
	results := make(map[Servant]int)
	for _, s := range m.servants {
		results[s.Servant()] = s.OnCommand(cmd)
	}
	return results
}

// UseCommandSpell は令呪を使用
func (m *Master) UseCommandSpell(cmd Command) map[Servant]int {
	if m.CommandSpells <= 0 {
		return nil
	}
	m.CommandSpells--
	return m.Command(cmd)
}

// ContractedServant は契約されたサーヴァント（Observer）
type ContractedServant struct {
	servant Servant
}

// NewContractedServant は契約サーヴァントを生成
func NewContractedServant(servant Servant) *ContractedServant {
	return &ContractedServant{servant: servant}
}

// Servant はサーヴァントを返す
func (c *ContractedServant) Servant() Servant {
	return c.servant
}

// OnCommand は命令を受け取り結果を返す
func (c *ContractedServant) OnCommand(cmd Command) int {
	switch cmd.Type {
	case CommandAttack:
		return c.servant.Attack()
	case CommandUseNoblePhantasm:
		return c.servant.UseNoblePhantasm()
	default:
		return 0
	}
}
