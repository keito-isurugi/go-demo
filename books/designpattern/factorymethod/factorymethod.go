package factorymethod

import "fmt"

type Product interface {
	use()
	getOwner() string
}

type Factory interface {
	createProduct(owner string) Product
	registerProduct(product Product)
}
type AbstructFactory struct{
	Factory Factory
}
func NewAbstractFactry(f Factory) *AbstructFactory {
	return &AbstructFactory{Factory: f}
}
func (af AbstructFactory) create(owner string) Product{
	p := af.Factory.createProduct(owner)
	af.Factory.registerProduct(p)
	return p
}

type IDCard struct {
	owner string
}
func NewIDCard(owner string) *IDCard{
	fmt.Println(owner, "のカードを作ります。")
	return &IDCard{owner: owner}
}
func (idc IDCard) use() {
	fmt.Println(idc.owner, "のカードを使います。")
}
func (idc IDCard) getOwner() string{
	return idc.owner
}

type IDCardFactory struct {
	owners []string
}
func NewIDCardFactory() *IDCardFactory{
	return &IDCardFactory{owners: make([]string, 0)}
}
func (idcf IDCardFactory) createProduct(owner string) Product{
	return NewIDCard(owner)
}
func (idcf IDCardFactory) registerProduct(product Product) {
	idcf.owners = append(idcf.owners, product.getOwner())
}
func (idcf IDCardFactory) getOwners() []string {
	return idcf.owners
}

func Exec() {
	factory := NewIDCardFactory()
	abstractFactory := NewAbstractFactry(factory)

	card1 := abstractFactory.create("tanaka")
	card2 := abstractFactory.create("yamada")
	card3 := abstractFactory.create("watanabe")

	card1.use()
	card2.use()
	card3.use()
}
