package fate_servant

import (
	"errors"
	"fmt"
)

// HeroicSpirit は英霊（クラス適性を持つ）
type HeroicSpirit struct {
	Name           string
	ClassAptitudes []ServantClass
	Configs        map[ServantClass]*ServantConfig
}

// ServantConfig はクラスごとの設定
type ServantConfig struct {
	Stats         Stats
	NoblePhantasm NoblePhantasm
	AttackPower   int
}

// HasAptitude は指定クラスへの適性があるか確認
func (h *HeroicSpirit) HasAptitude(class ServantClass) bool {
	for _, c := range h.ClassAptitudes {
		if c == class {
			return true
		}
	}
	return false
}

// Registry は英霊の座
type Registry struct {
	spirits map[string]*HeroicSpirit
}

// NewRegistry は英霊の座を生成
func NewRegistry() *Registry {
	r := &Registry{spirits: make(map[string]*HeroicSpirit)}
	r.registerDefaults()
	return r
}

func (r *Registry) registerDefaults() {
	// アルトリア（Saber, Lancer適性）
	r.spirits["Artoria"] = &HeroicSpirit{
		Name:           "Artoria",
		ClassAptitudes: []ServantClass{ClassSaber, ClassLancer},
		Configs: map[ServantClass]*ServantConfig{
			ClassSaber: {
				Stats:         Stats{"B", "B", "B", "A", "A+"},
				NoblePhantasm: &Excalibur{},
				AttackPower:   1500,
			},
			ClassLancer: {
				Stats:         Stats{"B", "A", "A", "A", "C"},
				NoblePhantasm: &Rhongomyniad{},
				AttackPower:   1400,
			},
		},
	}

	// クー・フーリン（Lancer, Caster適性）
	r.spirits["CuChulainn"] = &HeroicSpirit{
		Name:           "Cu Chulainn",
		ClassAptitudes: []ServantClass{ClassLancer, ClassCaster},
		Configs: map[ServantClass]*ServantConfig{
			ClassLancer: {
				Stats:         Stats{"B", "C", "A", "C", "E"},
				NoblePhantasm: &GaeBolg{},
				AttackPower:   1200,
			},
			ClassCaster: {
				Stats:         Stats{"E", "D", "C", "B", "D"},
				NoblePhantasm: &WickerMan{},
				AttackPower:   800,
			},
		},
	}

	// エミヤ（Archer適性）
	r.spirits["Emiya"] = &HeroicSpirit{
		Name:           "Emiya",
		ClassAptitudes: []ServantClass{ClassArcher},
		Configs: map[ServantClass]*ServantConfig{
			ClassArcher: {
				Stats:         Stats{"D", "C", "C", "B", "E"},
				NoblePhantasm: &UnlimitedBladeWorks{},
				AttackPower:   1000,
			},
		},
	}

	// ギルガメッシュ（Archer適性）
	r.spirits["Gilgamesh"] = &HeroicSpirit{
		Name:           "Gilgamesh",
		ClassAptitudes: []ServantClass{ClassArcher},
		Configs: map[ServantClass]*ServantConfig{
			ClassArcher: {
				Stats:         Stats{"B", "C", "C", "B", "A"},
				NoblePhantasm: &GateOfBabylon{},
				AttackPower:   1800,
			},
		},
	}

	// イスカンダル（Rider適性）
	r.spirits["Iskandar"] = &HeroicSpirit{
		Name:           "Iskandar",
		ClassAptitudes: []ServantClass{ClassRider},
		Configs: map[ServantClass]*ServantConfig{
			ClassRider: {
				Stats:         Stats{"B", "A", "D", "C", "A+"},
				NoblePhantasm: &IoniouHetairoi{},
				AttackPower:   1400,
			},
		},
	}
}

// Get は英霊を取得
func (r *Registry) Get(name string) (*HeroicSpirit, bool) {
	spirit, ok := r.spirits[name]
	return spirit, ok
}

// List は全英霊を返す
func (r *Registry) List() []*HeroicSpirit {
	list := make([]*HeroicSpirit, 0, len(r.spirits))
	for _, s := range r.spirits {
		list = append(list, s)
	}
	return list
}

// Summoner は召喚システム（Factory Pattern）
type Summoner struct {
	registry *Registry
}

// NewSummoner は召喚システムを生成
func NewSummoner() *Summoner {
	return &Summoner{registry: NewRegistry()}
}

var (
	ErrNotFound    = errors.New("heroic spirit not found")
	ErrNoAptitude  = errors.New("no aptitude for class")
)

// Summon は英霊を指定クラスで召喚
func (s *Summoner) Summon(name string, class ServantClass) (Servant, error) {
	spirit, ok := s.registry.Get(name)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, name)
	}

	if !spirit.HasAptitude(class) {
		return nil, fmt.Errorf("%w: %s -> %s", ErrNoAptitude, name, class)
	}

	config := spirit.Configs[class]
	return &BaseServant{
		trueName:      spirit.Name,
		class:         class,
		stats:         config.Stats,
		noblePhantasm: config.NoblePhantasm,
		attackPower:   config.AttackPower,
	}, nil
}

// Registry returns the heroic spirit registry
func (s *Summoner) Registry() *Registry {
	return s.registry
}
