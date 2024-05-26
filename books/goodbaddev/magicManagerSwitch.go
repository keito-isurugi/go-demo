package goodbaddev

import (
    "fmt"
)

type IMagic interface {
    GetName()
    CostMagicPoint()
    AttackPower()
    CostTechnicalPoint()
}

type MagicMnager struct{
    MagicName string
    MagicPoint int
    MagicPower int
    MagicTechnicalPoint int
}

type MagicType string

const (
    MagicFire MagicType = "fire"
    MagicShiden MagicType = "shiden"
    MagicHellFire MagicType = "hellFire"
)

// switchを一箇所にした実装
func (mm *MagicMnager) NewMagic(mt MagicType) {
    switch mt {
    case MagicFire:
        mm.MagicName = "ファイヤ"
        mm.MagicPoint = 100
        mm.MagicPower = 50
        mm.MagicTechnicalPoint = 25
    case MagicShiden:
        mm.MagicName = "紫電"
        mm.MagicPoint = 200
        mm.MagicPower = 100
        mm.MagicTechnicalPoint = 50
    case MagicHellFire:
        mm.MagicName = "地獄の業火"
        mm.MagicPoint = 400
        mm.MagicPower = 200
        mm.MagicTechnicalPoint = 100
    }
}
func (mm *MagicMnager) GetName() {
    fmt.Println("魔法名: ", mm.MagicName)
}
func (mm *MagicMnager) CostMagicPoint() {
    fmt.Println("消費MP: ", mm.MagicPoint)
}
func (mm *MagicMnager) AttackPower() {
    fmt.Println("攻撃力: ", mm.MagicPower)
}
func (mm *MagicMnager) CostTechnicalPoint() {
    fmt.Println("消費TP: ", mm.MagicTechnicalPoint)
}

func MgickSwitch() {
    magicGouka := &MagicMnager{}
	magicGouka.NewMagic(MagicHellFire)
	magicGouka.GetName()
	magicGouka.CostMagicPoint()
	magicGouka.AttackPower()
	magicGouka.CostTechnicalPoint()
	fmt.Println("=================")
	magicShiden := &MagicMnager{}
	magicShiden.NewMagic(MagicShiden)
	magicShiden.GetName()
	magicShiden.CostMagicPoint()
	magicShiden.AttackPower()
	magicShiden.CostTechnicalPoint()
	fmt.Println("=================")
}