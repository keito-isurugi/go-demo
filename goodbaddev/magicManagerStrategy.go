package goodbaddev

import (
    "fmt"
)

type MagicManager struct {
    strategy IMagic
}

func NewMagicManager(strategy IMagic) *MagicManager {
    return &MagicManager{strategy: strategy}
}

func (m *MagicManager) SetStrategy(strategy IMagic) {
    m.strategy = strategy
}

func (m *MagicManager) Execute() {
    m.strategy.AttackPower()
    m.strategy.CostMagicPoint()
    m.strategy.AttackPower()
    m.strategy.CostTechnicalPoint()
}

type Fire struct {}
func (f *Fire) GetName() {
    fmt.Println("魔法名: ファイヤ")
}
func (f *Fire) CostMagicPoint() {
    fmt.Println("消費MP: 100")
}
func (f *Fire) AttackPower() {
    fmt.Println("攻撃力: 50")
}
func (f *Fire) CostTechnicalPoint() {
    fmt.Println("消費TP: 25")
}

type Shiden struct {}
func (s *Shiden) GetName() {
    fmt.Println("魔法名: 紫電")
}
func (s *Shiden) CostMagicPoint() {
    fmt.Println("消費MP: 200")
}
func (s *Shiden) AttackPower() {
    fmt.Println("攻撃力: 100")
}
func (s *Shiden) CostTechnicalPoint() {
    fmt.Println("消費TP: 50")
}

type HellFire struct {}
func (hf *HellFire) GetName() {
    fmt.Println("魔法名: 地獄の業火")
}
func (hf *HellFire) CostMagicPoint() {
    fmt.Println("消費MP: 400")
}
func (hf *HellFire) AttackPower() {
    fmt.Println("攻撃力: 200")
}
func (hf *HellFire) CostTechnicalPoint() {
    fmt.Println("消費TP: 100")
}