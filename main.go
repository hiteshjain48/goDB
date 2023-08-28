package main

import (
	"fmt"
	"os"
	"encoding/json"
	"sync"
	"github.com/jcelliott/lumber"
)

const version = "1.0.0"
type Address struct {
	City string
	State string
	Country string
	Pincode json.Number
}
type User struct {
	Name string
	Age json.Number
	Contact string
	Company string
	Address Address
}
func main() {
	dir := "./"

	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	employees := []User{
		{"Hitesh","21","8459492271", "Oracle", Address{"Datala", "Maharashtra", "India", "443102"}},
	}

	for _, value := range employees{
		db.Write("users", value.Name, User{
			Name: value.Name,
			Age: value.Age,
			Contact: value.Contact,
			Company: value.Company,
			Address: value.Address,
		})
	}

	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println(records)

	allusers := []User{}

	for _,f := range records {
		employeeFound := User{}
		if err := json.Unmarshal([]byte(f), &employeeFound); err != nil {
			fmt.Println("Error", err)
		}
		allusers = append(allusers, employeeFound)
	}
}