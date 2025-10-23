package tyozetsu

import  "fmt"

type User struct {
	Name string
}

func (user *User) GetUserPublic() {
	fmt.Println(user.Name)
}

// 多様性
type Pet interface {
	reaction() 
}

type PetshopCustomer struct {}
func (pc *PetshopCustomer) touch(pet Pet) {
	pet.reaction()
}

type Dog struct{
	Voice string
}
func (dog *Dog) reaction() {
	fmt.Println(dog.Voice)
}
type Cat struct{
	Voice string
}
func (cat *Cat) reaction() {
	fmt.Println(cat.Voice)
}

func Exec() {
	customer := &PetshopCustomer{}
	customer.touch(&Dog{Voice: "wan!"})
	customer.touch(&Cat{Voice: "nyan!"})
}