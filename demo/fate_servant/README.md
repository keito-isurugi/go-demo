# Fate/Design Pattern - サーヴァントシミュレーター

Fateシリーズのサーヴァントを題材に、Goのデザインパターンを学ぶデモプロジェクトです。

## Fateシリーズとは

**Fate/stay night** を原点とするTYPE-MOON制作のメディアミックス作品群です。

物語の中心となるのは「聖杯戦争」と呼ばれる儀式。7人の魔術師（マスター）が、過去の英雄や伝説上の人物を「サーヴァント」として召喚し、万能の願望機「聖杯」を巡って戦います。

### 主要な概念

| 用語 | 説明 |
|-----|------|
| **サーヴァント** | 英霊（過去の英雄）が現界した存在。英霊ごとに生前の逸話に基づく「クラス適性」があり、召喚時に適性のあるクラス（Saber, Archer, Lancer, Rider, Caster, Assassin, Berserker）のいずれかに確定する |
| **マスター** | サーヴァントを召喚・契約した魔術師。令呪でサーヴァントに絶対命令を下せる |
| **宝具（Noble Phantasm）** | サーヴァントの切り札。生前の逸話や伝説が具現化した必殺技 |
| **令呪** | マスターに与えられる3画の絶対命令権。使い切ると契約が弱まる |

このプロジェクトでは、これらの概念をデザインパターンで表現しています。

---

## 実装されたデザインパターン

### 1. Strategy Pattern（戦略パターン）

**ファイル**: `noble_phantasm.go`

宝具（Noble Phantasm）の実装に使用。各サーヴァントが異なる宝具を持ち、動的に切り替え可能です。

```go
// NoblePhantasm インターフェース
type NoblePhantasm interface {
    GetName() string
    GetChant() string
    Activate() int
}

// 具体的な宝具実装
type Excalibur struct{}      // エクスカリバー
type GaeBolg struct{}        // ゲイ・ボルク
type UnlimitedBladeWorks struct{} // 無限の剣製
```

| 宝具 | サーヴァント | 説明 |
|-----|------------|------|
| エクスカリバー | アルトリア | 約束された勝利の剣 |
| ゲイ・ボルク | クー・フーリン | 刺し穿つ死棘の槍 |
| 無限の剣製 | エミヤ | 固有結界 |
| 王の財宝 | ギルガメッシュ | ゲート・オブ・バビロン |
| エヌマ・エリシュ | ギルガメッシュ | 天地乖離す開闘の星 |
| 王の軍勢 | イスカンダル | アイオニオン・ヘタイロイ |

---

### 2. Factory Pattern（ファクトリーパターン）

**ファイル**: `servant.go`

サーヴァントの生成を一元管理。名前またはクラスでサーヴァントを召喚できます。

```go
factory := &ServantFactory{}

// 名前で召喚
artoria := factory.CreateServant("アルトリア")
gilgamesh := factory.CreateServant("ギルガメッシュ")

// クラスで召喚
saber := factory.CreateServantByClass(ClassSaber)
lancer := factory.CreateServantByClass(ClassLancer)
```

| クラス | デフォルトサーヴァント |
|-------|---------------------|
| Saber | アルトリア・ペンドラゴン |
| Archer | エミヤ |
| Lancer | クー・フーリン |
| Rider | イスカンダル |

---

### 3. Template Method Pattern（テンプレートメソッドパターン）

**ファイル**: `battle.go`

バトルアクションの共通処理を定義し、具体的な行動はサブクラスで実装します。

```go
// バトルの流れ: Prepare → Execute → Finalize
type BattleAction interface {
    Prepare()
    Execute() int
    Finalize()
}

// 具体的なアクション
type NormalAttackAction struct{}    // 通常攻撃
type NoblePhantasmAction struct{}   // 宝具攻撃
type ComboAttackAction struct{}     // コンボ攻撃
```

```
========== バトルアクション開始 ==========
アルトリア・ペンドラゴン（Saber）が構えを取る...  ← Prepare
アルトリア・ペンドラゴン（Saber）の通常攻撃！    ← Execute
攻撃完了。体勢を整える。                        ← Finalize
========== 与ダメージ: 1500 ==========
```

---

### 4. Observer Pattern（オブザーバーパターン）

**ファイル**: `master.go`

マスターとサーヴァントの契約関係を表現。マスターからの命令をサーヴァントに通知します。

```go
// マスター（Subject）
type Master struct {
    Name          string
    CommandSpells int  // 令呪
    servants      []ServantObserver
}

// サーヴァント（Observer）
type ServantObserver interface {
    OnCommand(cmd Command)
    GetServant() Servant
}
```

#### サーヴァントの性格による反応の違い

| タイプ | 実装 | 特徴 |
|-------|-----|------|
| 通常 | `ContractedServant` | 標準的な反応 |
| 忠誠心が高い | `LoyalServant` | 騎士道精神に基づく反応（セイバー向け） |
| 誇り高い | `ProudServant` | 王としての威厳ある反応（ギルガメッシュ向け） |

```
[マスター: 遠坂時臣] 命令を発行: 撤退
[ギルガメッシュ] 撤退だと？この王に逃げろと？...貴様、後で話がある
```

---

## 実行方法

```bash
cd demo/fate_servant
go run ./cmd/main.go
```

## クラス図

```mermaid
classDiagram
    class Servant {
        <<interface>>
        +GetTrueName() string
        +GetClass() ServantClass
        +GetStats() Stats
        +GetNoblePhantasm() NoblePhantasm
        +Attack() int
        +UseNoblePhantasm() int
    }

    class NoblePhantasm {
        <<interface>>
        +GetName() string
        +GetChant() string
        +Activate() int
    }

    class ServantObserver {
        <<interface>>
        +OnCommand(cmd Command)
        +GetServant() Servant
    }

    class Master {
        +Name string
        +CommandSpells int
        +RegisterServant(observer)
        +IssueCommand(cmd)
        +UseCommandSpell(cmd)
    }

    class ServantFactory {
        +CreateServant(name) Servant
        +CreateServantByClass(class) Servant
    }

    class BattleTemplate {
        +PerformAction() int
    }

    class BattleAction {
        <<interface>>
        +Prepare()
        +Execute() int
        +Finalize()
    }

    Servant --> NoblePhantasm : has
    Master --> ServantObserver : notifies
    ServantObserver --> Servant : wraps
    ServantFactory --> Servant : creates
    BattleTemplate --> BattleAction : uses
```

## 各パターンの関係

```mermaid
flowchart TB
    subgraph "Factory Pattern"
        F[ServantFactory]
        F --> S1[Artoria]
        F --> S2[Gilgamesh]
        F --> S3[Emiya]
    end

    subgraph "Strategy Pattern"
        S1 --> NP1[Excalibur]
        S2 --> NP2[Gate of Babylon]
        S2 --> NP3[Enuma Elish]
        S3 --> NP4[Unlimited Blade Works]
    end

    subgraph "Observer Pattern"
        M[Master]
        M -->|command| O1[LoyalServant]
        M -->|command| O2[ProudServant]
        O1 --> S1
        O2 --> S2
    end

    subgraph "Template Method Pattern"
        BT[BattleTemplate]
        BT --> BA1[NormalAttack]
        BT --> BA2[NoblePhantasm]
        BT --> BA3[ComboAttack]
    end
```

## ファイル構成

```
demo/fate_servant/
├── go.mod              # モジュール定義
├── README.md           # このファイル
├── noble_phantasm.go   # Strategy Pattern - 宝具
├── servant.go          # Factory Pattern - サーヴァント
├── battle.go           # Template Method Pattern - バトル
├── master.go           # Observer Pattern - マスター
└── cmd/
    └── main.go         # デモ実行ファイル
```
