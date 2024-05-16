package goodbaddev

import (
    "fmt"
)

type MagicMnager struct{
    MagicName string
    MagicPoint int
    MagicPower int
    MagicTechnicalPoint int
}

func (mm *MagicMnager) NewMagic(mt MagicType) {
    switch mt {
    case Fire:
        mm.MagicName = "ファイヤ"
        mm.MagicPoint = 100
        mm.MagicPower = 50
        mm.MagicTechnicalPoint = 25
    case Shiden:
        mm.MagicName = "紫電"
        mm.MagicPoint = 200
        mm.MagicPower = 100
        mm.MagicTechnicalPoint = 50
    case HellFire:
        mm.MagicName = "地獄の業火"
        mm.MagicPoint = 400
        mm.MagicPower = 200
        mm.MagicTechnicalPoint = 100
    }
}

type MagicType string

const (
    Fire MagicType = "fire"
    Shiden MagicType = "shiden"
    HellFire MagicType = "hellFire"
)

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
