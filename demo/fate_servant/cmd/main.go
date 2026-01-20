package main

import (
	"fmt"
	servant "fate_servant"
)

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║     Fate/Design Pattern - サーヴァントシミュレーター      ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	// ==========================================================================
	// 1. Factory Pattern デモ
	// ==========================================================================
	fmt.Println("\n\n【1. Factory Pattern】サーヴァントの召喚")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	factory := &servant.ServantFactory{}

	// 名前でサーヴァントを召喚
	artoria := factory.CreateServant("アルトリア")
	gilgamesh := factory.CreateServant("ギルガメッシュ")
	emiya := factory.CreateServant("エミヤ")

	// クラスでサーヴァントを召喚
	lancer := factory.CreateServantByClass(servant.ClassLancer)

	fmt.Println("\n召喚されたサーヴァント:")
	artoria.Introduce()
	gilgamesh.Introduce()
	emiya.Introduce()
	lancer.Introduce()

	// ==========================================================================
	// 2. Strategy Pattern デモ
	// ==========================================================================
	fmt.Println("\n\n【2. Strategy Pattern】宝具の切り替え")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	fmt.Println("\n--- アルトリアの宝具 ---")
	artoria.UseNoblePhantasm()

	fmt.Println("\n--- エミヤの宝具 ---")
	emiya.UseNoblePhantasm()

	// ギルガメッシュは複数の宝具を持つ
	fmt.Println("\n--- ギルガメッシュの宝具（王の財宝）---")
	gilgamesh.UseNoblePhantasm()

	// ギルガメッシュの真の力
	if gil, ok := gilgamesh.(*servant.Gilgamesh); ok {
		fmt.Println("\n--- ギルガメッシュの真の宝具（エヌマ・エリシュ）---")
		gil.UseEnumaElish()
	}

	// ==========================================================================
	// 3. Template Method Pattern デモ
	// ==========================================================================
	fmt.Println("\n\n【3. Template Method Pattern】バトルアクション")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	// バトルシステムを作成
	battleSystem := servant.NewBattleSystem(artoria, gilgamesh)

	// 通常攻撃
	fmt.Println("\n--- 通常攻撃ターン ---")
	battleSystem.ExecuteTurn(false)

	// 攻守交代
	battleSystem.SwapRoles()

	// 宝具攻撃
	fmt.Println("\n--- 宝具攻撃ターン ---")
	battleSystem.ExecuteTurn(true)

	// コンボ攻撃のデモ
	fmt.Println("\n--- コンボ攻撃デモ ---")
	comboAction := servant.NewComboAttackAction(lancer, 5)
	template := servant.NewBattleTemplate(comboAction)
	template.PerformAction()

	// ==========================================================================
	// 4. Observer Pattern デモ
	// ==========================================================================
	fmt.Println("\n\n【4. Observer Pattern】マスターとサーヴァントの契約")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	// マスターを作成
	shirou := servant.NewMaster("衛宮士郎")
	tokiomi := servant.NewMaster("遠坂時臣")

	// サーヴァントを契約（性格によって反応が変わる）
	loyalSaber := servant.NewLoyalServant(artoria)
	proudArcher := servant.NewProudServant(gilgamesh)
	normalLancer := servant.NewContractedServant(lancer)

	// 契約を結ぶ
	shirou.RegisterServant(loyalSaber)
	shirou.RegisterServant(normalLancer)
	tokiomi.RegisterServant(proudArcher)

	// 命令を発行
	fmt.Println("\n--- 士郎が攻撃命令を出す ---")
	shirou.IssueCommand(servant.Command{
		Type:    servant.CommandAttack,
		Message: "あの敵を倒してくれ！",
	})

	fmt.Println("\n--- 時臣が宝具解放を命じる ---")
	tokiomi.IssueCommand(servant.Command{
		Type:    servant.CommandUseNoblePhantasm,
		Message: "全力で敵を殲滅せよ",
	})

	fmt.Println("\n--- 時臣が撤退を命じる（誇り高きギルガメッシュの反応）---")
	tokiomi.IssueCommand(servant.Command{
		Type:    servant.CommandRetreat,
		Message: "一時撤退だ",
	})

	// 令呪を使用
	fmt.Println("\n--- 士郎が令呪を使用 ---")
	shirou.UseCommandSpell(servant.Command{
		Message: "セイバー、全力で守れ！",
	})

	fmt.Printf("\n士郎の残り令呪: %d画\n", shirou.GetRemainingCommandSpells())

	// ==========================================================================
	// 5. 複合デモ: 聖杯戦争シミュレーション
	// ==========================================================================
	fmt.Println("\n\n【5. 複合デモ】聖杯戦争バトル")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	fmt.Println("\n========== 第四次聖杯戦争 最終決戦 ==========")
	fmt.Println("セイバー vs ギルガメッシュ")
	fmt.Println("============================================")

	// イスカンダルを召喚して参戦
	iskandar := factory.CreateServant("イスカンダル")
	iskandar.Introduce()

	fmt.Println("\n--- イスカンダルの王の軍勢 ---")
	if rider, ok := iskandar.(*servant.Iskandar); ok {
		rider.UseIoniouHetairoi()
	}

	fmt.Println("\n\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    デモ終了                                ║")
	fmt.Println("║  実装されたデザインパターン:                               ║")
	fmt.Println("║  • Strategy Pattern - 宝具の動的切り替え                   ║")
	fmt.Println("║  • Factory Pattern - サーヴァントの生成                    ║")
	fmt.Println("║  • Template Method Pattern - バトルアクションの共通処理    ║")
	fmt.Println("║  • Observer Pattern - マスター⇔サーヴァント通信           ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
}
