# Fate/Design Pattern

Fateシリーズのサーヴァントを題材に、Goのデザインパターンを学ぶデモ。

## Fateシリーズとは

TYPE-MOON制作のメディアミックス作品群。魔術師（マスター）が英霊（サーヴァント）を召喚し、聖杯を巡って戦う。

| 用語 | 説明 |
|-----|------|
| サーヴァント | 英霊が現界した存在。クラス適性に基づき召喚時にクラスが確定 |
| 宝具 | サーヴァントの必殺技 |
| 令呪 | マスターの絶対命令権（3画） |

## デザインパターンとFateの対応

| パターン | Fateの概念 | 実装 |
|---------|-----------|------|
| Factory | 召喚システム | `Summoner.Summon()` |
| Strategy | 宝具 | `NoblePhantasm` interface |
| Template Method | バトルアクション | `BattleTemplate.Run()` |
| Observer | マスター⇔サーヴァント | `Master.Command()` |

## 実行

```bash
cd demo/fate_servant
go run ./cmd/main.go
```

出力:
```
=== Fate/Design Pattern デモ ===

[Factory Pattern] 召喚システム
英霊はクラス適性を持ち、適性のあるクラスでのみ召喚可能
  アルトリアをSaberで召喚 → 宝具: Excalibur
  アルトリアをLancerで召喚 → 宝具: Rhongomyniad
  イスカンダルをSaberで召喚 → エラー: no aptitude for class

[Strategy Pattern] 宝具
同じ英霊でもクラスによって宝具が異なる
  Artoria (Saber): Excalibur → 9000ダメージ
  Artoria (Lancer): Rhongomyniad → 9500ダメージ

[Observer Pattern] マスターとサーヴァント
マスターの命令が契約サーヴァント全員に通知される
  マスター: 衛宮士郎（令呪: 3画）
  令呪を使用して宝具解放を命令:
    Artoria: Excalibur → 9000ダメージ
  残り令呪: 2画
```

## コード例

### Factory Pattern: 召喚

```go
summoner := fate.NewSummoner()

// クラス適性に基づく召喚
saber, _ := summoner.Summon("Artoria", fate.ClassSaber)   // OK
lancer, _ := summoner.Summon("Artoria", fate.ClassLancer) // OK

// 適性なし → エラー
_, err := summoner.Summon("Iskandar", fate.ClassSaber)
```

### Strategy Pattern: 宝具

```go
type NoblePhantasm interface {
    Name() string
    Damage() int
}

// 同じ英霊でもクラスによって宝具が異なる
saber.NoblePhantasm().Name()  // "Excalibur"
lancer.NoblePhantasm().Name() // "Rhongomyniad"
```

### Template Method Pattern: バトル

```go
type BattleAction interface {
    Prepare()
    Execute() int
    Cleanup()
}

// テンプレートが共通の流れを定義
template := fate.NewBattleTemplate(action)
damage := template.Run() // Prepare → Execute → Cleanup
```

### Observer Pattern: マスター

```go
master := fate.NewMaster("衛宮士郎")
master.Contract(fate.NewContractedServant(saber))

// 命令発行 → 全サーヴァントに通知
master.Command(fate.Command{Type: fate.CommandAttack})

// 令呪使用
master.UseCommandSpell(fate.Command{Type: fate.CommandUseNoblePhantasm})
```

## ファイル構成

```
fate_servant/
├── noble_phantasm.go  # Strategy: 宝具
├── servant.go         # サーヴァント基本定義
├── heroic_spirit.go   # Factory: 英霊・召喚
├── battle.go          # Template Method: バトル
├── master.go          # Observer: マスター
└── cmd/main.go        # デモ
```
