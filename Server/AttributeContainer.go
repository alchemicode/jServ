package main

import (
	"encoding/json"
	"fmt"
)

type AttributeContainer struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

//Default Constructor
func (ac *AttributeContainer) New(key string, value interface{}) {
	ac.Key = key
	ac.Value = value
}

func (ac AttributeContainer) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["key"] = ac.Key
	m["value"] = ac.Value
	return m
}

//Converts the container to Json text
func (ac AttributeContainer) ToJson() string {
	js, _ := json.Marshal(ac.ToMap())
	return string(js)
}

//Converts the container to a string
func (ac AttributeContainer) String() string {
	return fmt.Sprintf("{ \" %s \" : %v  }", ac.Key, ac.Value)
}
