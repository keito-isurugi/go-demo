// Package fate_servant demonstrates design patterns using Fate series Servants
package fate_servant

import "fmt"

// =============================================================================
// Strategy Pattern: 宝具（Noble Phantasm）の実装
// 各サーヴァントが異なる宝具を持ち、攻撃方法を動的に切り替え可能
// =============================================================================

// NoblePhantasm は宝具のインターフェース（Strategy）
type NoblePhantasm interface {
	GetName() string
	GetChant() string
	Activate() int // ダメージを返す
}

// Excalibur はセイバーの宝具（約束された勝利の剣）
type Excalibur struct{}

func (e *Excalibur) GetName() string {
	return "約束された勝利の剣（エクスカリバー）"
}

func (e *Excalibur) GetChant() string {
	return "誓いの言葉を紡ごう、世界を救う光となれ！エクスカリバー！"
}

func (e *Excalibur) Activate() int {
	fmt.Println(e.GetChant())
	fmt.Println("【" + e.GetName() + "】発動！")
	fmt.Println("黄金の光が敵を包み込む！")
	return 9000 // 高ダメージ
}

// GaeBolg はランサーの宝具（刺し穿つ死棘の槍）
type GaeBolg struct{}

func (g *GaeBolg) GetName() string {
	return "刺し穿つ死棘の槍（ゲイ・ボルク）"
}

func (g *GaeBolg) GetChant() string {
	return "その心臓、貰い受ける！"
}

func (g *GaeBolg) Activate() int {
	fmt.Println(g.GetChant())
	fmt.Println("【" + g.GetName() + "】発動！")
	fmt.Println("因果逆転の呪いが心臓を貫く！")
	return 7000 // 必中効果付き
}

// UnlimitedBladeWorks はアーチャーの宝具（無限の剣製）
type UnlimitedBladeWorks struct{}

func (u *UnlimitedBladeWorks) GetName() string {
	return "無限の剣製（アンリミテッドブレイドワークス）"
}

func (u *UnlimitedBladeWorks) GetChant() string {
	return `体は剣で出来ている
血潮は鉄で 心は硝子
幾たびの戦場を越えて不敗
ただの一度も敗走はなく
ただの一度も理解されない
彼の者は常に独り 剣の丘で勝利に酔う
故に、生涯に意味はなく
その体は、きっと剣で出来ていた`
}

func (u *UnlimitedBladeWorks) Activate() int {
	fmt.Println(u.GetChant())
	fmt.Println("【" + u.GetName() + "】発動！")
	fmt.Println("固有結界展開！無数の剣が降り注ぐ！")
	return 8000
}

// GateOfBabylon はギルガメッシュの宝具（王の財宝）
type GateOfBabylon struct{}

func (g *GateOfBabylon) GetName() string {
	return "王の財宝（ゲート・オブ・バビロン）"
}

func (g *GateOfBabylon) GetChant() string {
	return "雑種が...慢心せずに全力で叩き潰してやろう"
}

func (g *GateOfBabylon) Activate() int {
	fmt.Println(g.GetChant())
	fmt.Println("【" + g.GetName() + "】発動！")
	fmt.Println("無数の宝具が空間から射出される！")
	return 10000 // 最強クラス
}

// EnumaElish はギルガメッシュの最強宝具（天地乖離す開闘の星）
type EnumaElish struct{}

func (e *EnumaElish) GetName() string {
	return "天地乖離す開闘の星（エヌマ・エリシュ）"
}

func (e *EnumaElish) GetChant() string {
	return "目を開けろ、全てを見据え、ただ前を向け。世界を切り裂く！エヌマ・エリシュ！"
}

func (e *EnumaElish) Activate() int {
	fmt.Println(e.GetChant())
	fmt.Println("【" + e.GetName() + "】発動！")
	fmt.Println("乖離剣エアが世界を切り裂く！")
	return 15000 // 対界宝具
}

// Gordius はライダーの宝具（神威の車輪）
type Gordius struct{}

func (g *Gordius) GetName() string {
	return "神威の車輪（ゴルディアス・ホイール）"
}

func (g *Gordius) GetChant() string {
	return "蹂躙せよ！神威の車輪！"
}

func (g *Gordius) Activate() int {
	fmt.Println(g.GetChant())
	fmt.Println("【" + g.GetName() + "】発動！")
	fmt.Println("神牛が引く戦車が敵陣を蹂躙する！")
	return 6500
}

// IoniouHetairoi はイスカンダルの宝具（王の軍勢）
type IoniouHetairoi struct{}

func (i *IoniouHetairoi) GetName() string {
	return "王の軍勢（アイオニオン・ヘタイロイ）"
}

func (i *IoniouHetairoi) GetChant() string {
	return "我が覇道に、集え英雄！"
}

func (i *IoniouHetairoi) Activate() int {
	fmt.Println(i.GetChant())
	fmt.Println("【" + i.GetName() + "】発動！")
	fmt.Println("固有結界展開！数万の英霊が呼応し戦場を駆ける！")
	return 12000
}
