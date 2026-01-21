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
	// 1. 召喚システムとクラス適性デモ
	// ==========================================================================
	fmt.Println("\n\n【1. 召喚システム】英霊とクラス適性")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	summonSystem := servant.NewSummoningSystem()
	registry := summonSystem.GetRegistry()

	// 英霊の座に登録されている英霊を確認
	fmt.Println("\n--- 英霊の座に登録されている英霊 ---")
	for _, spirit := range registry.ListSpirits() {
		fmt.Printf("【%s】クラス適性: %s\n", spirit.TrueName, spirit.GetAptitudeString())
	}

	// ==========================================================================
	// 2. クラス適性に基づく召喚
	// ==========================================================================
	fmt.Println("\n\n【2. クラス適性に基づく召喚】")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	// アルトリアをSaberとして召喚
	fmt.Println("\n--- アルトリアをSaberとして召喚 ---")
	artoriaSaber, err := summonSystem.Summon("アルトリア", servant.ClassSaber)
	if err != nil {
		fmt.Printf("召喚失敗: %v\n", err)
	} else {
		artoriaSaber.Introduce()
	}

	// アルトリアをLancerとして召喚（別クラス）
	fmt.Println("\n--- アルトリアをLancerとして召喚 ---")
	artoriaLancer, err := summonSystem.Summon("アルトリア", servant.ClassLancer)
	if err != nil {
		fmt.Printf("召喚失敗: %v\n", err)
	} else {
		artoriaLancer.Introduce()
	}

	// クー・フーリンをCasterとして召喚
	fmt.Println("\n--- クー・フーリンをCasterとして召喚 ---")
	cuCaster, err := summonSystem.Summon("クー・フーリン", servant.ClassCaster)
	if err != nil {
		fmt.Printf("召喚失敗: %v\n", err)
	} else {
		cuCaster.Introduce()
	}

	// ==========================================================================
	// 3. 適性のないクラスへの召喚（エラーケース）
	// ==========================================================================
	fmt.Println("\n\n【3. 適性のないクラスへの召喚】")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	// イスカンダルをSaberとして召喚しようとする（失敗）
	fmt.Println("\n--- イスカンダルをSaberとして召喚（適性なし）---")
	_, err = summonSystem.Summon("イスカンダル", servant.ClassSaber)
	if err != nil {
		fmt.Printf("召喚失敗: %v\n", err)
	}

	// エミヤをBerserkerとして召喚しようとする（失敗）
	fmt.Println("\n--- エミヤをBerserkerとして召喚（適性なし）---")
	_, err = summonSystem.Summon("エミヤ", servant.ClassBerserker)
	if err != nil {
		fmt.Printf("召喚失敗: %v\n", err)
	}

	// ==========================================================================
	// 4. Strategy Pattern - クラスによる宝具の違い
	// ==========================================================================
	fmt.Println("\n\n【4. Strategy Pattern】クラスによる宝具の違い")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	fmt.Println("\n--- セイバー・アルトリアの宝具（エクスカリバー）---")
	artoriaSaber.UseNoblePhantasm()

	fmt.Println("\n--- ランサー・アルトリアの宝具（ロンゴミニアド）---")
	artoriaLancer.UseNoblePhantasm()

	fmt.Println("\n--- キャスター・クー・フーリンの宝具（ウィッカーマン）---")
	cuCaster.UseNoblePhantasm()

	// クー・フーリンをLancerとして召喚
	fmt.Println("\n--- ランサー・クー・フーリンの宝具（ゲイ・ボルク）---")
	cuLancer, _ := summonSystem.Summon("クー・フーリン", servant.ClassLancer)
	cuLancer.UseNoblePhantasm()

	// ==========================================================================
	// 5. ギルガメッシュの複数宝具
	// ==========================================================================
	fmt.Println("\n\n【5. 複数宝具】ギルガメッシュ")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	gilgamesh, _ := summonSystem.Summon("ギルガメッシュ", servant.ClassArcher)
	gilgamesh.Introduce()

	fmt.Println("\n--- 王の財宝 ---")
	gilgamesh.UseNoblePhantasm()

	// 追加宝具の使用
	if summoned, ok := gilgamesh.(*servant.SummonedServant); ok {
		if summoned.HasAlternateNoblePhantasm() {
			fmt.Println("\n--- エヌマ・エリシュ（追加宝具）---")
			summoned.UseAlternateNoblePhantasm()
		}
	}

	// ==========================================================================
	// 6. 聖杯戦争シミュレーション（同一クラス重複不可）
	// ==========================================================================
	fmt.Println("\n\n【6. 聖杯戦争】同一クラス重複不可ルール")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	// 新しい召喚システム（聖杯戦争用）
	holyGrailWar := servant.NewSummoningSystem()

	fmt.Println("\n--- 第五次聖杯戦争 サーヴァント召喚 ---")

	// Saberの枠を埋める
	_, err = holyGrailWar.SummonForHolyGrailWar("アルトリア", servant.ClassSaber)
	if err != nil {
		fmt.Printf("召喚失敗: %v\n", err)
	}

	// 同じSaberの枠に別の英霊を召喚しようとする
	fmt.Println("\n--- 別のSaberを召喚しようとする ---")
	_, err = holyGrailWar.SummonForHolyGrailWar("クー・フーリン", servant.ClassSaber)
	if err != nil {
		fmt.Printf("召喚失敗: %v\n", err)
	}

	// Lancerの枠は空いている
	fmt.Println("\n--- Lancerの枠は空いているので召喚可能 ---")
	_, err = holyGrailWar.SummonForHolyGrailWar("クー・フーリン", servant.ClassLancer)
	if err != nil {
		fmt.Printf("召喚失敗: %v\n", err)
	}

	// ==========================================================================
	// 7. Template Method Pattern & Observer Pattern（既存デモ）
	// ==========================================================================
	fmt.Println("\n\n【7. バトルシステム】Template Method & Observer Pattern")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	// マスターを作成
	shirou := servant.NewMaster("衛宮士郎")
	tokiomi := servant.NewMaster("遠坂時臣")

	// サーヴァントを契約
	loyalSaber := servant.NewLoyalServant(artoriaSaber)
	proudArcher := servant.NewProudServant(gilgamesh)

	shirou.RegisterServant(loyalSaber)
	tokiomi.RegisterServant(proudArcher)

	// バトル
	fmt.Println("\n--- バトル開始 ---")
	battleSystem := servant.NewBattleSystem(artoriaSaber, gilgamesh)
	battleSystem.ExecuteTurn(false)
	battleSystem.SwapRoles()
	battleSystem.ExecuteTurn(true)

	// 令呪使用
	fmt.Println("\n--- 令呪使用 ---")
	shirou.UseCommandSpell(servant.Command{
		Message: "セイバー、宝具を解放せよ！",
	})

	fmt.Println("\n\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    デモ終了                                ║")
	fmt.Println("║  実装された設定:                                           ║")
	fmt.Println("║  • 英霊ごとのクラス適性                                    ║")
	fmt.Println("║  • クラスに応じた宝具・ステータスの変化                    ║")
	fmt.Println("║  • 聖杯戦争での同一クラス重複不可ルール                    ║")
	fmt.Println("║                                                            ║")
	fmt.Println("║  実装されたデザインパターン:                               ║")
	fmt.Println("║  • Strategy Pattern - 宝具の動的切り替え                   ║")
	fmt.Println("║  • Factory Pattern - 召喚システム                          ║")
	fmt.Println("║  • Template Method Pattern - バトルアクション              ║")
	fmt.Println("║  • Observer Pattern - マスター⇔サーヴァント通信           ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
}
