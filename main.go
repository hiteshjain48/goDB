package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

type (
	Logger interface{
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct {
		mutex sync.Mutex
		mutexes map[string]*sync.Mutex
		dir string
		log Logger
	}
)

type Options struct {
	Logger
}

func New(dir string, options *Options) (*Driver, error){
	dir = filepath.Clean(dir)
	opts := Options{}
	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}

	driver := Driver{
		dir: dir,
		mutexes: make(map[string]*sync.Mutex),
		log: opts.Logger,
	}

	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return &driver, nil
	}
	opts.Logger.Debug("Creating the database at '%s'...\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

func (d * Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" {
		return errors.New("Missing collection - no place to save record")
	}
	if resource == "" {
		return errors.New("Missing record - unable to save record (no name)!")
	}
	mutex := d.getOrCreateMutex()
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	finalPath := filepath.Join(dir, resource+".json")
	tempPath := finalPath + ".tmp"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data := json.MarshalIndent(v, "", "")
	data = append(data, byte('\n'))
	if err := os.WriteFile(tempPath, b, 0644); err != nil {
		return err
	}

}

func (d *Driver) Read() (){

}

func (d *Driver) ReadAll(collection string) (){

}

func (d *Driver) Delete() error {

}

func (d *Driver) getOrCreateMutex() *sync.Mutex {

}

func stat(path string) (os.FileInfo, error) {
	fi, err := os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
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

	// if err := db.Delete("user", "Hitesh"); err != nil {
	// 	fmt.Println("Error: ", err)
	// }

	// if err := db.Delete("user", "");  err != nil {
	// 	fmt.Println("Error: ", err)
	// }

}