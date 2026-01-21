package fate_servant

import (
	"errors"
	"fmt"
	"math/rand"
)

// =============================================================================
// 英霊（Heroic Spirit）とクラス適性の実装
// 英霊は複数のクラス適性を持ち、召喚時にいずれかのクラスに確定する
// =============================================================================

// HeroicSpirit は英霊（座にいる状態）を表す
type HeroicSpirit struct {
	TrueName       string                           // 真名
	ClassAptitudes []ServantClass                   // クラス適性
	ClassConfigs   map[ServantClass]*ServantConfig  // クラスごとの設定
}

// ServantConfig はクラスごとのサーヴァント設定
type ServantConfig struct {
	Stats           Stats
	NoblePhantasm   NoblePhantasm
	AttackPower     int
	AlternateNP     NoblePhantasm // 追加の宝具（ギルガメッシュのエヌマ・エリシュなど）
}

// HasAptitude は指定クラスへの適性があるか確認する
func (h *HeroicSpirit) HasAptitude(class ServantClass) bool {
	for _, c := range h.ClassAptitudes {
		if c == class {
			return true
		}
	}
	return false
}

// GetAptitudeString はクラス適性を文字列で返す
func (h *HeroicSpirit) GetAptitudeString() string {
	result := ""
	for i, c := range h.ClassAptitudes {
		if i > 0 {
			result += ", "
		}
		result += string(c)
	}
	return result
}

// =============================================================================
// 英霊のレジストリ（英霊の座）
// =============================================================================

// HeroicSpiritRegistry は英霊の座を表す
type HeroicSpiritRegistry struct {
	spirits map[string]*HeroicSpirit
}

// NewHeroicSpiritRegistry は英霊の座を生成する
func NewHeroicSpiritRegistry() *HeroicSpiritRegistry {
	registry := &HeroicSpiritRegistry{
		spirits: make(map[string]*HeroicSpirit),
	}
	registry.registerDefaultSpirits()
	return registry
}

// registerDefaultSpirits はデフォルトの英霊を登録する
func (r *HeroicSpiritRegistry) registerDefaultSpirits() {
	// アルトリア・ペンドラゴン（Saber, Lancer, Rider適性）
	r.spirits["アルトリア"] = &HeroicSpirit{
		TrueName:       "アルトリア・ペンドラゴン",
		ClassAptitudes: []ServantClass{ClassSaber, ClassLancer, ClassRider},
		ClassConfigs: map[ServantClass]*ServantConfig{
			ClassSaber: {
				Stats: Stats{
					Strength: "B", Endurance: "B", Agility: "B",
					Mana: "A", Luck: "A+", NoblePhantasm: "A++",
				},
				NoblePhantasm: &Excalibur{},
				AttackPower:   1500,
			},
			ClassLancer: {
				Stats: Stats{
					Strength: "B", Endurance: "A", Agility: "A",
					Mana: "A", Luck: "C", NoblePhantasm: "A++",
				},
				NoblePhantasm: &Rhongomyniad{},
				AttackPower:   1400,
			},
			ClassRider: {
				Stats: Stats{
					Strength: "B", Endurance: "A", Agility: "A",
					Mana: "A", Luck: "B", NoblePhantasm: "A+",
				},
				NoblePhantasm: &Llamrei{},
				AttackPower:   1300,
			},
		},
	}

	// クー・フーリン（Lancer, Caster, Berserker適性）
	r.spirits["クー・フーリン"] = &HeroicSpirit{
		TrueName:       "クー・フーリン",
		ClassAptitudes: []ServantClass{ClassLancer, ClassCaster, ClassBerserker},
		ClassConfigs: map[ServantClass]*ServantConfig{
			ClassLancer: {
				Stats: Stats{
					Strength: "B", Endurance: "C", Agility: "A",
					Mana: "C", Luck: "E", NoblePhantasm: "B",
				},
				NoblePhantasm: &GaeBolg{},
				AttackPower:   1200,
			},
			ClassCaster: {
				Stats: Stats{
					Strength: "E", Endurance: "D", Agility: "C",
					Mana: "B", Luck: "D", NoblePhantasm: "B",
				},
				NoblePhantasm: &WickerMan{},
				AttackPower:   800,
			},
			ClassBerserker: {
				Stats: Stats{
					Strength: "A", Endurance: "B+", Agility: "A+",
					Mana: "C", Luck: "E", NoblePhantasm: "A",
				},
				NoblePhantasm: &GaeBolgBeast{},
				AttackPower:   1600,
			},
		},
	}

	// エミヤ（Archer, Assassin適性）
	r.spirits["エミヤ"] = &HeroicSpirit{
		TrueName:       "エミヤ",
		ClassAptitudes: []ServantClass{ClassArcher, ClassAssassin},
		ClassConfigs: map[ServantClass]*ServantConfig{
			ClassArcher: {
				Stats: Stats{
					Strength: "D", Endurance: "C", Agility: "C",
					Mana: "B", Luck: "E", NoblePhantasm: "E~A++",
				},
				NoblePhantasm: &UnlimitedBladeWorks{},
				AttackPower:   1000,
			},
			ClassAssassin: {
				Stats: Stats{
					Strength: "D", Endurance: "D", Agility: "B",
					Mana: "B", Luck: "E", NoblePhantasm: "E~A++",
				},
				NoblePhantasm: &UnlimitedBladeWorks{},
				AttackPower:   900,
			},
		},
	}

	// ギルガメッシュ（Archer, Caster適性）
	r.spirits["ギルガメッシュ"] = &HeroicSpirit{
		TrueName:       "ギルガメッシュ",
		ClassAptitudes: []ServantClass{ClassArcher, ClassCaster},
		ClassConfigs: map[ServantClass]*ServantConfig{
			ClassArcher: {
				Stats: Stats{
					Strength: "B", Endurance: "C", Agility: "C",
					Mana: "B", Luck: "A", NoblePhantasm: "EX",
				},
				NoblePhantasm: &GateOfBabylon{},
				AttackPower:   1800,
				AlternateNP:   &EnumaElish{},
			},
			ClassCaster: {
				Stats: Stats{
					Strength: "D", Endurance: "C", Agility: "C",
					Mana: "A", Luck: "A", NoblePhantasm: "EX",
				},
				NoblePhantasm: &GateOfBabylon{},
				AttackPower:   1200,
			},
		},
	}

	// イスカンダル（Rider適性）
	r.spirits["イスカンダル"] = &HeroicSpirit{
		TrueName:       "イスカンダル",
		ClassAptitudes: []ServantClass{ClassRider},
		ClassConfigs: map[ServantClass]*ServantConfig{
			ClassRider: {
				Stats: Stats{
					Strength: "B", Endurance: "A", Agility: "D",
					Mana: "C", Luck: "A+", NoblePhantasm: "A++",
				},
				NoblePhantasm: &Gordius{},
				AttackPower:   1400,
				AlternateNP:   &IoniouHetairoi{},
			},
		},
	}
}

// GetSpirit は英霊を取得する
func (r *HeroicSpiritRegistry) GetSpirit(name string) (*HeroicSpirit, bool) {
	spirit, ok := r.spirits[name]
	return spirit, ok
}

// ListSpirits は登録されている英霊の一覧を返す
func (r *HeroicSpiritRegistry) ListSpirits() []*HeroicSpirit {
	spirits := make([]*HeroicSpirit, 0, len(r.spirits))
	for _, s := range r.spirits {
		spirits = append(spirits, s)
	}
	return spirits
}

// =============================================================================
// 召喚システム（SummoningSystem）
// =============================================================================

var (
	ErrSpiritNotFound    = errors.New("英霊が見つかりません")
	ErrNoAptitude        = errors.New("指定されたクラスへの適性がありません")
	ErrClassSlotOccupied = errors.New("そのクラスの枠は既に埋まっています")
)

// SummoningSystem は召喚システムを表す
type SummoningSystem struct {
	registry      *HeroicSpiritRegistry
	occupiedSlots map[ServantClass]bool // 聖杯戦争での枠管理
}

// NewSummoningSystem は召喚システムを生成する
func NewSummoningSystem() *SummoningSystem {
	return &SummoningSystem{
		registry:      NewHeroicSpiritRegistry(),
		occupiedSlots: make(map[ServantClass]bool),
	}
}

// Summon は英霊を指定クラスで召喚する
func (s *SummoningSystem) Summon(spiritName string, class ServantClass) (Servant, error) {
	spirit, ok := s.registry.GetSpirit(spiritName)
	if !ok {
		return nil, ErrSpiritNotFound
	}

	if !spirit.HasAptitude(class) {
		return nil, fmt.Errorf("%w: %sは%sの適性を持っていません（適性: %s）",
			ErrNoAptitude, spirit.TrueName, class, spirit.GetAptitudeString())
	}

	config, ok := spirit.ClassConfigs[class]
	if !ok {
		return nil, fmt.Errorf("クラス設定が見つかりません: %s", class)
	}

	fmt.Printf("\n【召喚】%sを%sとして召喚します...\n", spirit.TrueName, class)

	servant := &SummonedServant{
		BaseServant: BaseServant{
			TrueName:            spirit.TrueName,
			Class:               class,
			Stats:               config.Stats,
			NoblePhantasmWeapon: config.NoblePhantasm,
			AttackPower:         config.AttackPower,
		},
		alternateNP: config.AlternateNP,
	}

	fmt.Printf("「%sのサーヴァント、召喚に応じ参上した」\n", class)

	return servant, nil
}

// SummonForHolyGrailWar は聖杯戦争用の召喚（同一クラス重複不可）
func (s *SummoningSystem) SummonForHolyGrailWar(spiritName string, class ServantClass) (Servant, error) {
	if s.occupiedSlots[class] {
		return nil, fmt.Errorf("%w: %s", ErrClassSlotOccupied, class)
	}

	servant, err := s.Summon(spiritName, class)
	if err != nil {
		return nil, err
	}

	s.occupiedSlots[class] = true
	return servant, nil
}

// SummonRandom は適性クラスからランダムに召喚する
func (s *SummoningSystem) SummonRandom(spiritName string) (Servant, error) {
	spirit, ok := s.registry.GetSpirit(spiritName)
	if !ok {
		return nil, ErrSpiritNotFound
	}

	if len(spirit.ClassAptitudes) == 0 {
		return nil, fmt.Errorf("クラス適性がありません")
	}

	// ランダムにクラスを選択
	randomIndex := rand.Intn(len(spirit.ClassAptitudes))
	selectedClass := spirit.ClassAptitudes[randomIndex]

	fmt.Printf("【触媒なし召喚】クラスはランダムに決定されます...\n")
	return s.Summon(spiritName, selectedClass)
}

// GetRegistry は英霊の座を返す
func (s *SummoningSystem) GetRegistry() *HeroicSpiritRegistry {
	return s.registry
}

// ResetSlots は聖杯戦争の枠をリセットする
func (s *SummoningSystem) ResetSlots() {
	s.occupiedSlots = make(map[ServantClass]bool)
}

// =============================================================================
// SummonedServant（召喚されたサーヴァント）
// =============================================================================

// SummonedServant は召喚されたサーヴァントを表す
type SummonedServant struct {
	BaseServant
	alternateNP NoblePhantasm
}

// UseAlternateNoblePhantasm は追加の宝具を使用する
func (s *SummonedServant) UseAlternateNoblePhantasm() int {
	if s.alternateNP == nil {
		fmt.Println("追加の宝具はありません")
		return 0
	}
	fmt.Printf("\n%s、真の力を解放！\n", s.TrueName)
	return s.alternateNP.Activate()
}

// HasAlternateNoblePhantasm は追加の宝具を持っているか確認する
func (s *SummonedServant) HasAlternateNoblePhantasm() bool {
	return s.alternateNP != nil
}

// GetAlternateNoblePhantasm は追加の宝具を返す
func (s *SummonedServant) GetAlternateNoblePhantasm() NoblePhantasm {
	return s.alternateNP
}
