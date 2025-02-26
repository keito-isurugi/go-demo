package factory

// IPokemon インターフェース: ポケモンの基本的な振る舞いを定義
type IPokemon interface {
	setName(name string)
	setAttack(attack string)
	getAttack() string
	getName() string
}

// Pokemon 構造体: ポケモンの共通のプロパティを保持
type Pokemon struct {
	name   string
	attack string
}

// setName: ポケモンの名前を設定
func (p *Pokemon) setName(name string) {
	p.name = name
}

// setAttack: ポケモンの攻撃技を設定
func (p *Pokemon) setAttack(attack string) {
	p.attack = attack
}

// getAttack: ポケモンの攻撃技を取得
func (p *Pokemon) getAttack() string {
	return p.attack
}

// getName: ポケモンの名前を取得
func (p *Pokemon) getName() string {
	return p.name
}

// Pikachu 構造体: ポケモンの一種 (ピカチュウ)
type Pikachu struct {
	Pokemon
}

// NewPikachu: ピカチュウを生成するファクトリーメソッド
func NewPikachu() IPokemon {
	return &Pikachu{
		Pokemon: Pokemon{
			name:   "ピカチュウ",
			attack: "電気ショック",
		},
	}
}

// Eevee 構造体: ポケモンの一種 (イーブイ)
type Eevee struct {
	Pokemon
}

// NewEevee: イーブイを生成するファクトリーメソッド
func NewEevee() IPokemon {
	return &Eevee{
		Pokemon: Pokemon{
			name:   "イーブイ",
			attack: "たいあたり",
		},
	}
}

// GetPokemon: ポケモンの名前を指定してインスタンスを取得
func GetPokemon(name string) IPokemon {
	switch name {
	case "ピカチュウ":
		return NewPikachu()
	case "イーブイ":
		return NewEevee()
	default:
		return nil
	}
}

func FactoryExec() {
	pikachu := GetPokemon("ピカチュウ")
	eevee := GetPokemon("イーブイ")

	// 生成したポケモンの名前と技を出力
	println(pikachu.getName(), "の技: ", pikachu.getAttack())
	println(eevee.getName(), "の技: ", eevee.getAttack())
}
