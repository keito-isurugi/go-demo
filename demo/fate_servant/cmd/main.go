package main

import (
	"fmt"

	fate "fate_servant"
)

func main() {
	fmt.Println("=== Fate/Design Pattern デモ ===")

	summoner := fate.NewSummoner()

	// 1. Factory Pattern: 召喚システム
	fmt.Println("\n[Factory Pattern] 召喚システム")
	fmt.Println("英霊はクラス適性を持ち、適性のあるクラスでのみ召喚可能")

	saber, _ := summoner.Summon("Artoria", fate.ClassSaber)
	fmt.Printf("  アルトリアをSaberで召喚 → 宝具: %s\n", saber.NoblePhantasm().Name())

	lancer, _ := summoner.Summon("Artoria", fate.ClassLancer)
	fmt.Printf("  アルトリアをLancerで召喚 → 宝具: %s\n", lancer.NoblePhantasm().Name())

	archer, _ := summoner.Summon("Gilgamesh", fate.ClassArcher)
	fmt.Printf("  ギルガメッシュをArcherで召喚 → 宝具: %s\n", archer.NoblePhantasm().Name())

	_, err := summoner.Summon("Iskandar", fate.ClassSaber)
	fmt.Printf("  イスカンダルをSaberで召喚 → エラー: %v\n", err)

	// 2. Strategy Pattern: 宝具
	fmt.Println("\n[Strategy Pattern] 宝具")
	fmt.Println("同じ英霊でもクラスによって宝具が異なる")

	fmt.Printf("  %s (%s): %s → %dダメージ\n",
		saber.TrueName(), saber.Class(),
		saber.NoblePhantasm().Name(), saber.UseNoblePhantasm())

	fmt.Printf("  %s (%s): %s → %dダメージ\n",
		lancer.TrueName(), lancer.Class(),
		lancer.NoblePhantasm().Name(), lancer.UseNoblePhantasm())

	// 3. Template Method Pattern: バトルアクション
	fmt.Println("\n[Template Method Pattern] バトルアクション")
	fmt.Println("Prepare → Execute → Cleanup の流れを共通化")

	attackAction := fate.NewAttackAction(saber)
	npAction := fate.NewNoblePhantasmAction(archer)

	attackTemplate := fate.NewBattleTemplate(attackAction)
	npTemplate := fate.NewBattleTemplate(npAction)

	fmt.Printf("  通常攻撃: %dダメージ\n", attackTemplate.Run())
	fmt.Printf("  宝具攻撃: %dダメージ\n", npTemplate.Run())

	// 4. Observer Pattern: マスターとサーヴァント
	fmt.Println("\n[Observer Pattern] マスターとサーヴァント")
	fmt.Println("マスターの命令が契約サーヴァント全員に通知される")

	master := fate.NewMaster("衛宮士郎")
	master.Contract(fate.NewContractedServant(saber))
	master.Contract(fate.NewContractedServant(lancer))

	fmt.Printf("  マスター: %s（令呪: %d画）\n", master.Name, master.CommandSpells)

	fmt.Println("  攻撃命令を発行:")
	results := master.Command(fate.Command{Type: fate.CommandAttack})
	for servant, damage := range results {
		fmt.Printf("    %s → %dダメージ\n", servant.TrueName(), damage)
	}

	fmt.Println("  令呪を使用して宝具解放を命令:")
	results = master.UseCommandSpell(fate.Command{Type: fate.CommandUseNoblePhantasm})
	for servant, damage := range results {
		fmt.Printf("    %s: %s → %dダメージ\n",
			servant.TrueName(), servant.NoblePhantasm().Name(), damage)
	}
	fmt.Printf("  残り令呪: %d画\n", master.CommandSpells)
}
