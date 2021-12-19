package main

import (
	"fmt"
	"encoding/json"
)

type AttributeContainer struct {
	key string
	value interface{}
	m map[string]interface{}
}

func (ac *AttributeContainer) New(key string, value interface{}) {
	ac.m = make(map[string]interface{})
	ac.key = key
	ac.value = value
	ac.m[key] = value
}

func (ac AttributeContainer) ToJson() string{
	js, _ := json.Marshal(ac.m)
    return string(js)
}


func (ac AttributeContainer) String() string{
	return fmt.Sprintf("{ \" %s \" : %v  }", ac.key, ac.value)
}