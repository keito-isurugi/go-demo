package goodbaddev

import (
	"fmt"
)

type Person struct {
	Name string
	Age int
}

type Persons struct {
	people []Person
}

func (ps *Persons) Add(p Person) {
	ps.people = append(ps.people, p)
}

func (ps *Persons) GetAll() []Person{
	return ps.people
}

func (ps *Persons) FilterByAge(minAge int) Persons{
	var filterd []Person
	for _, p := range ps.people {
		if p.Age >= minAge {
			filterd = append(filterd, p)
		}
	}

	return Persons{people: filterd}
}

func FirstCollection() {
	persons := &Persons{}

	persons.Add(Person{Name: "Yamada", Age: 10})
	persons.Add(Person{Name: "Suzuki", Age: 30})
	persons.Add(Person{Name: "Watanabe", Age: 40})
	persons.Add(Person{Name: "Sato", Age: 50})
	
	allPersons := persons.GetAll()
	fmt.Printf("%+v", allPersons)
	fmt.Println()
	
	adultsPersons := persons.FilterByAge(20)
	fmt.Printf("%+v", adultsPersons)
}