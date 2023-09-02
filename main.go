package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
		return errors.New("missing collection - no place to save record")
	}
	if resource == "" {
		return errors.New("missing record - unable to save record (no name)")
	}
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	finalPath := filepath.Join(dir, resource+".json")
	tempPath := finalPath + ".tmp"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, _ := json.MarshalIndent(v, "", "")
	data = append(data, byte('\n'))
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return err
	}
	return os.Rename(tempPath, finalPath)
}

func (d *Driver) Read(collection, resource string, v interface{}) error {
	if collection == "" {
		return errors.New("missiong collection")
	}
	if resource == "" {
		return errors.New("missing resource")
	}

	recordPath := filepath.Join(d.dir, collection, resource)

	if _, err := stat(recordPath); err != nil {
		return err
	}

	record, err := os.ReadFile(recordPath + ".json")
	if err != nil {
		return err
	}
	return json.Unmarshal(record, &v)
}

func (d *Driver) ReadAll(collection string) ([]string, error){
	if collection == "" {
		return nil, errors.New("missing collection")
	}
	dir := filepath.Join(d.dir, collection)
	if _, err := stat(dir); err != nil {
		return nil, err
	}
	files, _ := os.ReadDir(dir)
	
	var records []string

	for _, file := range files {
		record, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		records = append(records, string(record))
	}
	return records, nil
}

func (d *Driver) Delete(collection, resource string) error {
	path := filepath.Join(collection, resource)
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch file, err := stat(dir); {
		case file == nil, err != nil:
			return fmt.Errorf("unable to find file or directory named %v",path)
		case file.Mode().IsDir():
			return os.RemoveAll(dir)
		case file.Mode().IsRegular():
			return os.RemoveAll(dir + ".json")
	}
	return nil
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	mut, ok := d.mutexes[collection]
	if !ok {
		mut = &sync.Mutex{}
		d.mutexes[collection] = mut
	}
	return mut
}

func stat(path string) (fi os.FileInfo, err error) {
	fi, err = os.Stat(path)
	if os.IsNotExist(err) {
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
		{"Vaibhav","20","8459492271", "Oracle", Address{"Kolpewadi", "Maharashtra", "India", "443102"}},
		{"Sarang","22","8459492271", "Oracle", Address{"Dhanora", "Maharashtra", "India", "443102"}},
		{"Pratham","21","8459492271", "Oracle", Address{"Takali", "Maharashtra", "India", "443102"}},
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
	fmt.Println(allusers)
	if err := db.Delete("users", "Hitesh"); err != nil {
		fmt.Println("Error: ", err)
	}

	if err := db.Delete("users", "");  err != nil {
		fmt.Println("Error: ", err)
	}

}