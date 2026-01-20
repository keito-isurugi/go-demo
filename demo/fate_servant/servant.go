package fate_servant

import "fmt"

// =============================================================================
// Factory Pattern: サーヴァントの生成
// クラスに応じたサーヴァントを生成するファクトリー
// =============================================================================

// ServantClass はサーヴァントのクラスを表す
type ServantClass string

const (
	ClassSaber    ServantClass = "Saber"
	ClassArcher   ServantClass = "Archer"
	ClassLancer   ServantClass = "Lancer"
	ClassRider    ServantClass = "Rider"
	ClassCaster   ServantClass = "Caster"
	ClassAssassin ServantClass = "Assassin"
	ClassBerserker ServantClass = "Berserker"
)

// Stats はサーヴァントのステータスを表す
type Stats struct {
	Strength     string // 筋力
	Endurance    string // 耐久
	Agility      string // 敏捷
	Mana         string // 魔力
	Luck         string // 幸運
	NoblePhantasm string // 宝具
}

// Servant はサーヴァントの基本インターフェース
type Servant interface {
	GetTrueName() string
	GetClass() ServantClass
	GetStats() Stats
	GetNoblePhantasm() NoblePhantasm
	SetNoblePhantasm(np NoblePhantasm)
	Attack() int
	UseNoblePhantasm() int
	Introduce()
}

// BaseServant はサーヴァントの基本実装
type BaseServant struct {
	TrueName      string
	Class         ServantClass
	Stats         Stats
	NoblePhantasmWeapon NoblePhantasm
	AttackPower   int
}

func (s *BaseServant) GetTrueName() string {
	return s.TrueName
}

func (s *BaseServant) GetClass() ServantClass {
	return s.Class
}

func (s *BaseServant) GetStats() Stats {
	return s.Stats
}

func (s *BaseServant) GetNoblePhantasm() NoblePhantasm {
	return s.NoblePhantasmWeapon
}

func (s *BaseServant) SetNoblePhantasm(np NoblePhantasm) {
	s.NoblePhantasmWeapon = np
}

func (s *BaseServant) Attack() int {
	fmt.Printf("%s（%s）の通常攻撃！\n", s.TrueName, s.Class)
	return s.AttackPower
}

func (s *BaseServant) UseNoblePhantasm() int {
	if s.NoblePhantasmWeapon == nil {
		fmt.Println("宝具が設定されていません")
		return 0
	}
	fmt.Printf("\n%s、宝具解放！\n", s.TrueName)
	return s.NoblePhantasmWeapon.Activate()
}

func (s *BaseServant) Introduce() {
	fmt.Println("================================================")
	fmt.Printf("真名: %s\n", s.TrueName)
	fmt.Printf("クラス: %s\n", s.Class)
	fmt.Println("ステータス:")
	fmt.Printf("  筋力: %s | 耐久: %s | 敏捷: %s\n", s.Stats.Strength, s.Stats.Endurance, s.Stats.Agility)
	fmt.Printf("  魔力: %s | 幸運: %s | 宝具: %s\n", s.Stats.Mana, s.Stats.Luck, s.Stats.NoblePhantasm)
	if s.NoblePhantasmWeapon != nil {
		fmt.Printf("宝具名: %s\n", s.NoblePhantasmWeapon.GetName())
	}
	fmt.Println("================================================")
}

// =============================================================================
// 具体的なサーヴァント実装
// =============================================================================

// Artoria はセイバーのアルトリア・ペンドラゴン
type Artoria struct {
	BaseServant
}

// NewArtoria はアルトリアを生成する
func NewArtoria() *Artoria {
	return &Artoria{
		BaseServant: BaseServant{
			TrueName: "アルトリア・ペンドラゴン",
			Class:    ClassSaber,
			Stats: Stats{
				Strength:     "B",
				Endurance:    "B",
				Agility:      "B",
				Mana:         "A",
				Luck:         "A+",
				NoblePhantasm: "A++",
			},
			NoblePhantasmWeapon: &Excalibur{},
			AttackPower:   1500,
		},
	}
}

// CuChulainn はランサーのクー・フーリン
type CuChulainn struct {
	BaseServant
}

// NewCuChulainn はクー・フーリンを生成する
func NewCuChulainn() *CuChulainn {
	return &CuChulainn{
		BaseServant: BaseServant{
			TrueName: "クー・フーリン",
			Class:    ClassLancer,
			Stats: Stats{
				Strength:     "B",
				Endurance:    "C",
				Agility:      "A",
				Mana:         "C",
				Luck:         "E",
				NoblePhantasm: "B",
			},
			NoblePhantasmWeapon: &GaeBolg{},
			AttackPower:   1200,
		},
	}
}

// Emiya はアーチャーのエミヤ
type Emiya struct {
	BaseServant
}

// NewEmiya はエミヤを生成する
func NewEmiya() *Emiya {
	return &Emiya{
		BaseServant: BaseServant{
			TrueName: "エミヤ",
			Class:    ClassArcher,
			Stats: Stats{
				Strength:     "D",
				Endurance:    "C",
				Agility:      "C",
				Mana:         "B",
				Luck:         "E",
				NoblePhantasm: "E~A++",
			},
			NoblePhantasmWeapon: &UnlimitedBladeWorks{},
			AttackPower:   1000,
		},
	}
}

// Gilgamesh は英雄王ギルガメッシュ
type Gilgamesh struct {
	BaseServant
	alternateNoblePhantasm NoblePhantasm // エヌマ・エリシュ
}

// NewGilgamesh はギルガメッシュを生成する
func NewGilgamesh() *Gilgamesh {
	return &Gilgamesh{
		BaseServant: BaseServant{
			TrueName: "ギルガメッシュ",
			Class:    ClassArcher,
			Stats: Stats{
				Strength:     "B",
				Endurance:    "C",
				Agility:      "C",
				Mana:         "B",
				Luck:         "A",
				NoblePhantasm: "EX",
			},
			NoblePhantasmWeapon: &GateOfBabylon{},
			AttackPower:   1800,
		},
		alternateNoblePhantasm: &EnumaElish{},
	}
}

// UseEnumaElish は乖離剣エアを使用する
func (g *Gilgamesh) UseEnumaElish() int {
	fmt.Printf("\n%s、真の力を解放！\n", g.TrueName)
	return g.alternateNoblePhantasm.Activate()
}

// Iskandar はライダーのイスカンダル
type Iskandar struct {
	BaseServant
	ultimateNoblePhantasm NoblePhantasm // 王の軍勢
}

// NewIskandar はイスカンダルを生成する
func NewIskandar() *Iskandar {
	return &Iskandar{
		BaseServant: BaseServant{
			TrueName: "イスカンダル",
			Class:    ClassRider,
			Stats: Stats{
				Strength:     "B",
				Endurance:    "A",
				Agility:      "D",
				Mana:         "C",
				Luck:         "A+",
				NoblePhantasm: "A++",
			},
			NoblePhantasmWeapon: &Gordius{},
			AttackPower:   1400,
		},
		ultimateNoblePhantasm: &IoniouHetairoi{},
	}
}

// UseIoniouHetairoi は王の軍勢を使用する
func (i *Iskandar) UseIoniouHetairoi() int {
	fmt.Printf("\n%s、固有結界を展開！\n", i.TrueName)
	return i.ultimateNoblePhantasm.Activate()
}

// =============================================================================
// Factory パターン実装
// =============================================================================

// ServantFactory はサーヴァントを生成するファクトリー
type ServantFactory struct{}

// CreateServant は指定された名前のサーヴァントを生成する
func (f *ServantFactory) CreateServant(name string) Servant {
	switch name {
	case "Artoria", "アルトリア":
		return NewArtoria()
	case "CuChulainn", "クー・フーリン":
		return NewCuChulainn()
	case "Emiya", "エミヤ":
		return NewEmiya()
	case "Gilgamesh", "ギルガメッシュ":
		return NewGilgamesh()
	case "Iskandar", "イスカンダル":
		return NewIskandar()
	default:
		return nil
	}
}

// CreateServantByClass はクラスに基づいてデフォルトのサーヴァントを生成する
func (f *ServantFactory) CreateServantByClass(class ServantClass) Servant {
	switch class {
	case ClassSaber:
		return NewArtoria()
	case ClassLancer:
		return NewCuChulainn()
	case ClassArcher:
		return NewEmiya()
	case ClassRider:
		return NewIskandar()
	default:
		return nil
	}
}
