package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"alchemicode.com/jServ/data"
	core "github.com/alchemicode/jserv-core"
	"github.com/google/uuid"

	msg "github.com/vmihailenco/msgpack/v5"
)

func ReadConfig(ch chan bool) {
	//Reads the contents of the config file into a string
	content, err := os.ReadFile("config.json")
	if err != nil {
		//Channel returns false if there is any error
		ch <- false
		fmt.Println(" > ! Failed to read config file")
		return
	}
	var dat map[string]interface{}
	if err3 := json.Unmarshal([]byte(content), &dat); err3 != nil {
		//Channel returns false if there is any error
		ch <- false
		fmt.Println(" > ! Error when generating json data for config file")
		return
	}
	//Reads IP and Port from the config file
	ip = dat["ip"].(string)
	port = int(dat["port"].(float64))
	debug = dat["debug"].(bool)
	write_interval = dat["write-interval"].(float64)
	//Reads in the values of the Requests list
	rtemp := dat["Requests"].(map[string]interface{})
	for key, value := range rtemp {
		requestTypes[key] = value.(bool)
	}
	//Reads in the values of the Permissions list
	ptemp := dat["Permissions"].(map[string]interface{})
	for key, value := range ptemp {
		requestPermissions[key] = (value.(string) != "admin")
	}
	atemp := dat["Aliases"].(map[string]interface{})
	for key, value := range atemp {
		aliases[key] = (value.(string))
	}
	stemp := dat["Services"].(map[string]interface{})
	for key, value := range stemp {
		services[key] = (value.(string))
	}
	appname = dat["appname"].(string)
	//Channel returns true if the read is successful
	ch <- true
}

func LoadDatabase(ch chan bool) {
	//Stores the files in the database directory in a list of files
	files, err := os.ReadDir("../Collections")
	if err != nil {
		//Channel returns false if there is any error
		ch <- false
		fmt.Println(" > ! Error when reading Collections directory")
		return
	}
	for _, file := range files {
		if !file.IsDir() {
			//Creates a new collection for each file in the directory
			name := strings.Split(file.Name(), ".")[0]
			col := new(core.Collection)
			col.New(name)
			ReadCollection(col)
		}
	}
	//Channel returns true if the read is successful
	ch <- true
}

// Returns the contents of a file as a slice of strings
func ReadFileAsLines(filename string) []string {
	//Opens file given in filename
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(" > ! Error when opening " + filename)
		panic(err)
	}
	defer file.Close()
	//Makes a string slice and adds each line in the file
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func ReadCollection(c *core.Collection) {
	ch := make(chan bool)
	go readCollection(c, ch)
	if ok := <-ch; !ok {
		fmt.Println(" > ! Failed to read " + c.Name + ".dat")
	} else {
		cols = append(cols, c)
		fmt.Println(" * Loaded the collection \"" + c.Name + "\"")
	}
}

func readCollection(c *core.Collection, ch chan bool) {
	path := filepath.Join("../Collections", c.Name+".dat")
	//Reads the contents of the file as a []byte
	content, err := os.ReadFile(path)
	if err != nil {
		//Channel returns false if there is any error
		ch <- false
		fmt.Println(" > ! Error while reading file " + c.Name + ".dat")
	} else {
		if len(content) == 0 {
			b, _ := msg.Marshal([]map[string]interface{}{})
			os.WriteFile(path, b, 0777)
			ch <- c.FromMsgPack(b)
		} else {
			ch <- c.FromMsgPack(content)
		}
	}
}

func UpdateCollection(c *core.Collection) {
	ch := make(chan bool)
	go updateCollection(c, ch)
	if ok := <-ch; !ok {
		fmt.Println(" > ! Failed to update " + c.Name + ".dat")
	}
}

func updateCollection(c *core.Collection, ch chan bool) {
	path := filepath.Join("../Collections", c.Name+".dat")
	//Reads the contents of the file as a []byte
	bytes, ok := c.ToMsgPack()
	if !ok {
		//Channel returns false if there is any error
		ch <- false
		fmt.Println(" > ! Error while serializing " + c.Name + ".dat")
	} else {
		if err := os.Truncate(path, 0); err != nil {
			ch <- false
			fmt.Println(" > ! Error when opening file " + c.Name + ".dat")
		} else {

			if err := os.WriteFile(path, bytes, 0777); err != nil {
				//Channel returns false if there is any error
				ch <- false
				fmt.Println(" > ! Error when opening file " + c.Name + ".dat")
			} else {
				//Channel returns true if the update was successful
				ch <- true
			}
		}

	}
}

// Checksfor a 'new' keyword in the admin file and replaces it with a new uuid
func GenerateAdminApiKey(ch chan bool) {
	lines := ReadFileAsLines("admin.jserv")
	if len(lines) > 0 {
		if lines[len(lines)-1] == "new" {
			randomuuid := uuid.New()
			adminKey = append(adminKey, randomuuid.String())
			lines[len(lines)-1] = randomuuid.String()
			os.WriteFile("admin.jserv", []byte(strings.Join(lines, "\n")), 0644)
			ch <- true
		} else {
			ch <- true
		}
	} else {
		//Channel returns false if there isn't a second line in the file
		fmt.Println(" > ! Failed to detect API Key. Write \"new\" on the last line of admin.jserv")
		ch <- false
	}
}

// Checks for a 'new' keyword in the admin file and replaces it with a new uuid
func GenerateUserApiKey(ch chan bool) {
	lines := ReadFileAsLines("keys.jserv")
	if len(lines) > 0 {
		if lines[len(lines)-1] == "new" {
			randomuuid := uuid.New()
			userKeys = append(userKeys, randomuuid.String())
			lines[len(lines)-1] = randomuuid.String()
			os.WriteFile("keys.jserv", []byte(strings.Join(lines, "\n")), 0644)
			ch <- true
		} else {
			ch <- true
		}
	} else {
		//Channel returns false if there isn't a second line in the file
		fmt.Println(" > ! Failed to detect API Key. Write \"new\" on the last line of keys.jserv")
		ch <- false
	}
}

// Reads all api keys from admin and keys file
func ReadKeys(ch chan bool) {
	//reads the lines that aren't 'new' '' or ' '
	lines := ReadFileAsLines("admin.jserv")
	for i := 0; i < len(lines); i++ {
		if lines[i] != "new" && lines[i] != "" && lines[i] != " " {
			adminKey = append(adminKey, lines[i])
		}
	}
	lines = ReadFileAsLines("keys.jserv")
	for i := 0; i < len(lines); i++ {
		if lines[i] != "new" && lines[i] != "" && lines[i] != " " {
			userKeys = append(userKeys, lines[i])
		}
	}
	ch <- true
}

// Checks the validity of all the required jserv data files
func CheckFiles(version *string) {
	//Checks if the files exist, create them if not, and panic if there is any error
	if _, err := os.Stat("admin.jserv"); os.IsNotExist(err) {
		f, err := os.Create("admin.jserv")
		if err != nil {
			panic(err)
		}
		f.WriteString("new")
	}
	if _, err := os.Stat("keys.jserv"); os.IsNotExist(err) {
		f, err := os.Create("keys.jserv")
		if err != nil {
			panic(err)
		}
		f.WriteString("new")
	}
	*version = "1.0.0"
}

// Checks if the given API key matches the permissions bool of a query type
func CheckApiKey(key string, permission bool) bool {
	if !permission {
		c, _ := data.Contains(adminKey, key)
		return c
	} else {
		uc, _ := data.Contains(userKeys, key)
		ac, _ := data.Contains(adminKey, key)
		return uc || ac
	}
}

func CheckAlias(address string) string {
	v, found := aliases[address]
	if found {
		return v
	} else {
		return address
	}
}

// The starting sequence to perform all the necessary checks before the server starts
func StartSequence(version *string, appname *string) {
	fmt.Println(" * Starting...")
	CheckFiles(version)
	ch := make(chan bool)
	go ReadConfig(ch)
	if result := <-ch; result {
		fmt.Println(" * Successfully read config")
	} else {
		fmt.Println(" * Error while reading config")
		os.Exit(1)
	}
	go LoadDatabase(ch)
	if result := <-ch; result {
		fmt.Println(" * Successfully generated Collections")
	} else {
		fmt.Println(" * Error while reading databases")
		os.Exit(1)
	}
	go GenerateAdminApiKey(ch)
	if result := <-ch; !result {
		fmt.Println(" * Error while reading/generating admin API Key")
		os.Exit(1)
	}
	go GenerateUserApiKey(ch)
	if result := <-ch; !result {
		fmt.Println(" * Error while reading/generating user API Key")
		os.Exit(1)
	}
	go ReadKeys(ch)
	if result := <-ch; !result {
		fmt.Println(" * ! Failed to read API Keys")
		os.Exit(1)
	}
	fmt.Printf(" * Running jServ v%s for %s\n", *version, *appname)
}
