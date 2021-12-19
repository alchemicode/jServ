package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var requestTypes = map[string]bool{
	//False denotes that the server cannot receive that request type
	"GET":     true,
	"POST":    true,
	"PUT":     false,
	"HEAD":    false,
	"DELETE":  true,
	"PATCH":   false,
	"OPTIONS": false}

var requestPermissions = map[string]bool{
	//False denotes that an admin API key is required to make that request
	"QObject":        true,
	"QAttribute":     true,
	"QAllAttributes": true,
	"QByAttribute":   true,
	"QnewId":         false,
	"AEmpty":         true,
	"AObject":        true,
	"AAttribute":     true,
	"MObject":        true,
	"MAttribute":     true,
	"DObject":        true,
	"DAttribute":     true}

var ip string = "localhost"
var port int = 4040
var dbs []*Collection = make([]*Collection, 0)

func ReadConfig(ch chan bool) {
	content, err := os.ReadFile("config.json")
	if err != nil {
		ch <- false
		fmt.Println("Failed to read config file")
		return
	}
	var dat map[string]interface{}
	if err3 := json.Unmarshal([]byte(content), &dat); err3 != nil {
		//Channel returns false if there is any error
		ch <- false
		fmt.Println("Error when generating json data for config file")
		return
	}
	ip = dat["ip"].(string)
	port = int(dat["port"].(float64))

	rtemp := dat["Requests"].(map[string]interface{})
	for key, value := range rtemp {
		requestTypes[key] = value.(bool)
	}
	ptemp := dat["Permissions"].(map[string]interface{})
	for key, value := range ptemp {
		requestPermissions[key] = (value.(string) != "admin")
	}
	ch <- true
}

func ReadDatabases(ch chan bool) {
	files, err := ioutil.ReadDir("Databases/")
	if err != nil {
		ch <- false
		fmt.Println("Error when reading Database directory")
		return
	}
	for _, file := range files {
		if !file.IsDir() {
			name := strings.Split(file.Name(), ".")[0]
			col := new(Collection)
			col.New(name)
			dbs = append(dbs, col)
			fmt.Println("Loaded database \"" + name + "\"")
		}
	}
	ch <- true
}

func ReadFileAsLines(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error when opening " + filename)
		panic(err)
	}
	defer file.Close()
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func GenerateAdminApiKey(ch chan bool) {
	lines := ReadFileAsLines("data.jserv")
	if len(lines) > 1 {
		if lines[1] == "new" {
			//INSERT UUID GENERATION CODE HERE
			ch <- true
		} else {
			ch <- true
		}
	} else {
		fmt.Println("Failed to detect API Key. Write \"new\" on the second line of data.jserv")
		ch <- false
	}
}

func main() {
	fmt.Println(" * Starting...")
	ch := make(chan bool)
	go ReadConfig(ch)
	if result := <-ch; result {
		fmt.Println("Successfully read config")
	} else {
		fmt.Println("Error while reading config")
		os.Exit(1)
	}
	go ReadDatabases(ch)
	if result := <-ch; result {
		fmt.Println("Successfully generated Collections")
	} else {
		fmt.Println("Error while reading databases")
		os.Exit(1)
	}
	fmt.Printf(" * Starting server on %s:%d", ip, port)
}
