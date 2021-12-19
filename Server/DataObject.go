package main

import (
	"encoding/json"
	"fmt"
)

type DataObject struct {
	Id   int                    `json:"id"`
	Data map[string]interface{} `json:"data"`
}

func (d *DataObject) WithEmptyMap(id int) {
	d.Id = id
	d.Data = make(map[string]interface{})
}

func (d *DataObject) WithData(id int, data map[string]interface{}) {
	d.Id = id
	d.Data = data
}

func (d *DataObject) FromJson(s string) {
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(s), &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat["id"].(float64))
	d.Id = int(dat["id"].(float64))
	d.Data = dat["data"].(map[string]interface{})
}

func (d DataObject) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = d.Id
	m["data"] = d.Data
	return m
}

func (d DataObject) ToJson() string {
	m := make(map[string]interface{})
	m["id"] = d.Id
	m["data"] = d.Data
	js, _ := json.Marshal(m)
	return string(js)
}

func (d DataObject) String() string {
	return fmt.Sprintf(" \"id\" : %d\n \"data\" : %v", d.Id, d.Data)
}
