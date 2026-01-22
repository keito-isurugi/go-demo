// Package fate_servant demonstrates design patterns using Fate series Servants
package fate_servant

// NoblePhantasm は宝具のインターフェース（Strategy Pattern）
type NoblePhantasm interface {
	Name() string
	Damage() int
}

// Excalibur - アルトリア(Saber)の宝具
type Excalibur struct{}

func (e *Excalibur) Name() string  { return "Excalibur" }
func (e *Excalibur) Damage() int   { return 9000 }

// Rhongomyniad - アルトリア(Lancer)の宝具
type Rhongomyniad struct{}

func (r *Rhongomyniad) Name() string  { return "Rhongomyniad" }
func (r *Rhongomyniad) Damage() int   { return 9500 }

// GaeBolg - クー・フーリン(Lancer)の宝具
type GaeBolg struct{}

func (g *GaeBolg) Name() string  { return "Gae Bolg" }
func (g *GaeBolg) Damage() int   { return 7000 }

// WickerMan - クー・フーリン(Caster)の宝具
type WickerMan struct{}

func (w *WickerMan) Name() string  { return "Wicker Man" }
func (w *WickerMan) Damage() int   { return 7500 }

// UnlimitedBladeWorks - エミヤ(Archer)の宝具
type UnlimitedBladeWorks struct{}

func (u *UnlimitedBladeWorks) Name() string  { return "Unlimited Blade Works" }
func (u *UnlimitedBladeWorks) Damage() int   { return 8000 }

// GateOfBabylon - ギルガメッシュ(Archer)の宝具
type GateOfBabylon struct{}

func (g *GateOfBabylon) Name() string  { return "Gate of Babylon" }
func (g *GateOfBabylon) Damage() int   { return 10000 }

// EnumaElish - ギルガメッシュの最強宝具
type EnumaElish struct{}

func (e *EnumaElish) Name() string  { return "Enuma Elish" }
func (e *EnumaElish) Damage() int   { return 15000 }

// IoniouHetairoi - イスカンダル(Rider)の宝具
type IoniouHetairoi struct{}

func (i *IoniouHetairoi) Name() string  { return "Ionioi Hetairoi" }
func (i *IoniouHetairoi) Damage() int   { return 12000 }
