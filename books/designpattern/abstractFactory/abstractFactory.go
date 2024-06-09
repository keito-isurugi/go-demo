package abstractfactory

import "fmt"

type IBikeFactory interface {
	makeSportsBike() ISportsBike
	makeNekidBike() INekidBike
}

func GetBikeFactory(maker string) (IBikeFactory, error) {
	if maker == "Honda" {
		return &Honda{}, nil
	}
	if maker == "Yamaha" {
		return &Yamaha{}, nil
	}

	return nil, fmt.Errorf("maker name is invalid")
}

type ISportsBike interface {
	setLogo(logo string)
	setName(name string)
	setCC(cc int)
	getLogo() string
	getName() string
	getCC() int
}
type SportsBike struct {
	logo string
	name string
	cc   int
}

func (sb *SportsBike) setLogo(logo string) {
	sb.logo = logo
}

func (sb *SportsBike) getLogo() string {
	return sb.logo
}
func (sb *SportsBike) setName(name string) {
	sb.name = name
}

func (sb *SportsBike) getName() string {
	return sb.name
}
func (sb *SportsBike) setCC(cc int) {
	sb.cc = cc
}
func (sb *SportsBike) getCC() int {
	return sb.cc
}

type INekidBike interface {
	setLogo(logo string)
	setName(name string)
	setCC(cc int)
	getLogo() string
	getName() string
	getCC() int
}
type NekidBike struct {
	logo string
	name string
	cc   int
}

func (nb *NekidBike) setLogo(logo string) {
	nb.logo = logo
}

func (nb *NekidBike) getLogo() string {
	return nb.logo
}
func (nb *NekidBike) setName(name string) {
	nb.name = name
}

func (nb *NekidBike) getName() string {
	return nb.name
}
func (nb *NekidBike) setCC(cc int) {
	nb.cc = cc
}
func (nb *NekidBike) getCC() int {
	return nb.cc
}

type HondaSportsBike struct {
	SportsBike
}
type HondaNekidBike struct {
	NekidBike
}
type YamahaSportsBike struct {
	SportsBike
}
type YamahaNekidBike struct {
	NekidBike
}

type Honda struct{}

func (h *Honda) makeSportsBike() ISportsBike {
	return &HondaSportsBike{
		SportsBike: SportsBike{
			logo: "Honda",
			name: "NSR250R",
			cc:   250,
		},
	}
}
func (h *Honda) makeNekidBike() INekidBike {
	return &HondaNekidBike{
		NekidBike: NekidBike{
			logo: "Honda",
			name: "CB400SF",
			cc:   400,
		},
	}
}

type Yamaha struct{}

func (h *Yamaha) makeSportsBike() ISportsBike {
	return &YamahaSportsBike{
		SportsBike: SportsBike{
			logo: "Yamaha",
			name: "YZF-R1",
			cc:   1000,
		},
	}
}
func (h *Yamaha) makeNekidBike() INekidBike {
	return &YamahaNekidBike{
		NekidBike: NekidBike{
			logo: "Yamaha",
			name: "XJR400",
			cc:   400,
		},
	}
}

func Exec() {
	hondaFactory, _ := GetBikeFactory("Honda")
	yamahaFactory, _ := GetBikeFactory("Yamaha")

	hondaSportsBike := hondaFactory.makeSportsBike()
	hondaNekidBike := hondaFactory.makeNekidBike()

	yamahaSportsBike := yamahaFactory.makeSportsBike()
	yamahaNekidBike := yamahaFactory.makeNekidBike()

	printSportsBikeDetails(hondaSportsBike)
	printNekidBikeDetails(hondaNekidBike)

	printSportsBikeDetails(yamahaSportsBike)
	printNekidBikeDetails(yamahaNekidBike)
}

func printSportsBikeDetails(sb ISportsBike) {
	fmt.Printf("Logo: %s\n", sb.getLogo())
	fmt.Printf("Name: %s\n", sb.getName())
	fmt.Printf("CC: %d\n", sb.getCC())
}

func printNekidBikeDetails(sb INekidBike) {
	fmt.Printf("Logo: %s\n", sb.getLogo())
	fmt.Printf("Name: %s\n", sb.getName())
	fmt.Printf("CC: %d\n", sb.getCC())
}
