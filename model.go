package main

import "fmt"

type User struct {
	Age  int
	Id   int
	Name Name
}

func NewUser(age int, id int, name Name) (User, error) {
	u := User{
		Age:  age,
		Id:   id,
		Name: name,
	}
	return u, u.Valdation()
}

func (u *User) Valdation() error {
	if u.Age < 10 && u.Name.Middle != nil {
		return fmt.Errorf("middle name is from 10 years old")
	}

	return nil
}

type Name struct {
	First  string
	Last   string
	Middle *string
}

func NewName(first, last string, middle *string) (Name, error) {
	n := Name{
		First:  first,
		Last:   last,
		Middle: middle,
	}

	return n, n.Valdation()
}

func (n *Name) Valdation() error {
	if n.Middle != nil {
		if n.First == "" {
			return fmt.Errorf("first name is required when middle name is provided")
		}
		if n.Last == "" {
			return fmt.Errorf("last name is required when middle name is provided")
		}
	}

	return nil
}
