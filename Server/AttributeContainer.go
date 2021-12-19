package main

import (
	"encoding/json"
	"fmt"
)

type AttributeContainer struct {
	key   string
	value interface{}
	m     map[string]interface{}
}

//Default Constructor
func (ac *AttributeContainer) New(key string, value interface{}) {
	ac.m = make(map[string]interface{})
	ac.key = key
	ac.value = value
	ac.m[key] = value
}

//Converts the container to Json text
func (ac AttributeContainer) ToJson() string {
	js, _ := json.Marshal(ac.m)
	return string(js)
}

//Converts the container to a string
func (ac AttributeContainer) String() string {
	return fmt.Sprintf("{ \" %s \" : %v  }", ac.key, ac.value)
}
