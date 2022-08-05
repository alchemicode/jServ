package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Collection struct {
	Name string
	List []DataObject
}

//Named Constructor Constructor
//Creates an empty Collection object with just a name
//and reads its file
func (c *Collection) New(name string) {
	c.Name = name
	c.List = make([]DataObject, 0)
	c.ReadFile()
}

//Asynchronously reads a collection's Json file
func (c *Collection) ReadFile() {
	//Channel provides success data
	ch := make(chan bool)
	go c.readFile(ch)
	result := <-ch
	if !result {
		fmt.Println("Failed to read database: " + c.Name + ".json")
	}
}

//Reads the collection's Json file
func (c *Collection) readFile(ch chan bool) {
	path := filepath.Join("Databases", c.Name+".json")
	//Reads the contents of the file as a string
	content, err2 := os.ReadFile(path)
	if err2 != nil {
		//Channel returns false if there is any error
		ch <- false
		fmt.Println("Error while reading file " + c.Name + ".json")
	} else {
		//Generates map data from Json file
		var dat []map[string]interface{}
		if err3 := json.Unmarshal([]byte(content), &dat); err3 != nil {
			//Channel returns false if there is any error
			ch <- false
			fmt.Println("Error when generating json data for " + c.Name + ".json")
		} else {
			//Reads each object in the generated data from the Json file
			//and populates the collection's list
			for i := 0; i < len(dat); i++ {
				obj := new(DataObject)
				obj.WithData(uint64(dat[i]["id"].(float64)), dat[i]["data"].(map[string]interface{}))
				c.List = append(c.List, *obj)
			}
			//Channel returns true if the read was successful
			ch <- true
		}

	}

}

//Asynchronously updates the collection's Json file
func (c *Collection) UpdateFile() {
	//Channel provides success data
	ch := make(chan bool)
	go c.updateFile(ch)
	result := <-ch
	if result {
		fmt.Println(" * Updated " + c.Name + ".json")
	} else {
		fmt.Println(" > Failed to update " + c.Name + ".json")
	}
}

func (c *Collection) updateFile(ch chan bool) {
	path := filepath.Join("Databases", c.Name+".json")
	//Converts the collection's list data into Json data
	js, _ := json.Marshal(c.List)
	//Overwrites Json file with new data
	if err := os.WriteFile(path, []byte(js), 0644); err != nil {
		//Channel returns false if there is any error
		ch <- false
		fmt.Println(" > Error when opening file " + c.Name + ".json")
	} else {
		//Channel returns true if the update was successful
		ch <- true
	}
}
