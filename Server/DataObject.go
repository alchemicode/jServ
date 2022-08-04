package main

import (
	"encoding/json"
	"fmt"
)

type DataObject struct {
	Id   uint64                 `json:"id"`
	Data map[string]interface{} `json:"data"`
}

//Default Constructor
//Creates an empty Object with only an id
func (d *DataObject) WithoutData(id uint64) {
	d.Id = id
	d.Data = make(map[string]interface{})
}

//Map Constructor
//Creates an Object with given id and data map
func (d *DataObject) WithData(id uint64, data map[string]interface{}) {
	d.Id = id
	d.Data = data
}

//Json Constructor
//Creates an Object from given Json string
func (d *DataObject) FromJson(s string) {
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(s), &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat["id"].(float64))
	d.Id = uint64(dat["id"].(float64))
	d.Data = dat["data"].(map[string]interface{})
}

//Returns the Object as a map
func (d DataObject) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = d.Id
	m["data"] = d.Data
	return m
}

//Returns the Object in Json string format
func (d DataObject) ToJson() string {
	js, _ := json.Marshal(d)
	return string(js)
}

//Returns the Object as a string
func (d DataObject) String() string {
	return fmt.Sprintf(" \"id\" : %d\n \"data\" : %v", d.Id, d.Data)
}
