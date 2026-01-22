package fate_servant

// ServantClass はサーヴァントのクラス
type ServantClass string

const (
	ClassSaber     ServantClass = "Saber"
	ClassArcher    ServantClass = "Archer"
	ClassLancer    ServantClass = "Lancer"
	ClassRider     ServantClass = "Rider"
	ClassCaster    ServantClass = "Caster"
	ClassAssassin  ServantClass = "Assassin"
	ClassBerserker ServantClass = "Berserker"
)

// Stats はサーヴァントのステータス
type Stats struct {
	Strength  string
	Endurance string
	Agility   string
	Mana      string
	Luck      string
}

// Servant はサーヴァントのインターフェース
type Servant interface {
	TrueName() string
	Class() ServantClass
	Stats() Stats
	NoblePhantasm() NoblePhantasm
	Attack() int
	UseNoblePhantasm() int
}

// BaseServant はサーヴァントの基本実装
type BaseServant struct {
	trueName      string
	class         ServantClass
	stats         Stats
	noblePhantasm NoblePhantasm
	attackPower   int
}

func (s *BaseServant) TrueName() string           { return s.trueName }
func (s *BaseServant) Class() ServantClass        { return s.class }
func (s *BaseServant) Stats() Stats               { return s.stats }
func (s *BaseServant) NoblePhantasm() NoblePhantasm { return s.noblePhantasm }
func (s *BaseServant) Attack() int                { return s.attackPower }
func (s *BaseServant) UseNoblePhantasm() int {
	if s.noblePhantasm == nil {
		return 0
	}
	return s.noblePhantasm.Damage()
}
