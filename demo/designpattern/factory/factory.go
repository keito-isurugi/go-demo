package factory

type IPokemon interface {
	setName(name string)
	setAttack(attack string)
	getAttack() string
	getName() string
}

type Pokemon struct {
	name   string
	attack string
}
func (p *Pokemon) setName(name string) {
	p.name = name
}
func (p *Pokemon) setAttack(attack string) {
	p.attack = attack
}
func (p *Pokemon) getAttack() string {
	return p.attack
}
func (p *Pokemon) getName() string {
	return p.name
}

type Pikachu struct {
	Pokemon
}
func NewPikachu() IPokemon {
	return &Pikachu{
		Pokemon: Pokemon{
			name:   "ピカチュウ",
			attack: "電気ショック",
		},
	}
}

type Eevee struct {
	Pokemon
}
func NewEevee() IPokemon {
	return &Eevee{
		Pokemon: Pokemon{
			name:   "イーブイ",
			attack: "たいあたり",
		},
	}
}

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

	println(pikachu.getName(), "の技: ", pikachu.getAttack())
	println(eevee.getName(), "の技: ", eevee.getAttack())
}