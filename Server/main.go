package main

import (
	"fmt"
	//"net/http"
	//"encoding/json"
	//"os"
)

func main() {
	ac := new(AttributeContainer)
	ac.New("Yeet", "Yaw")
	fmt.Println(ac)
	fmt.Println(ac.ToJson())

	fmt.Println("From empty map")
	do := new(DataObject)
	do.WithEmptyMap(1)
	do.Data["Yah"] = "Yeet"
	fmt.Println(do)

	fmt.Println("From JSON data")
	thing := ` {"id":2, "data" : { "Yah" : "Yeet" } } `
	d2 := new(DataObject)
	d2.FromJson(thing)
	fmt.Println(d2)
	fmt.Println(d2.ToJson())

	fmt.Println("Trying Collection")
	collection := new(Collection)
	collection.New("example")
	fmt.Println("Collection " + collection.name + ":")
	for i := 0; i < len(collection.list); i++ {
		fmt.Println(collection.list[i])
	}
	collection.list = append(collection.list, *d2)
	collection.UpdateFile()
}
