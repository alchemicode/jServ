package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Collection struct {
	name string
	list []DataObject
}

func (c *Collection) New(name string) {
	c.name = name
	c.list = make([]DataObject, 0)
	c.ReadFile()
}

func (c *Collection) ReadFile() {
	ch := make(chan bool)
	go c.readFile(ch)
	result := <-ch
	if result {
		fmt.Println("Succeeded to read database: " + c.name + ".json")
	} else {
		fmt.Println("Failed to read database: " + c.name + ".json")
	}
}

func (c *Collection) readFile(ch chan bool) {
	path := "Databases/" + c.name + ".json"
	content, err2 := os.ReadFile(path)
	if err2 != nil {
		ch <- false
		fmt.Println("Error while reading file " + c.name + ".json")
	} else {
		var dat []map[string]interface{}
		if err3 := json.Unmarshal([]byte(content), &dat); err3 != nil {
			ch <- false
			fmt.Println("Error when generating json data for " + c.name + ".json")
		} else {
			for i := 0; i < len(dat); i++ {
				obj := new(DataObject)
				obj.WithData(int(dat[i]["id"].(float64)), dat[i]["data"].(map[string]interface{}))
				c.list = append(c.list, *obj)
			}
			ch <- true
		}

	}

}

func (c *Collection) UpdateFile() {
	ch := make(chan bool)
	go c.updateFile(ch)
	result := <-ch
	if result {
		fmt.Println("Succeeded to update database: " + c.name + ".json")
	} else {
		fmt.Println("Failed to update database: " + c.name + ".json")
	}
}

func (c *Collection) updateFile(ch chan bool) {
	path := "Databases/" + c.name + "1.json"
	js, _ := json.Marshal(c.list)
	if err := os.WriteFile(path, []byte(js), 0644); err != nil {
		ch <- false
		fmt.Println("Error when opening file " + c.name + ".json")
	} else {
		ch <- true
	}
}
