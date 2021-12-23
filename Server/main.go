package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
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

var version string
var appname string
var ip string = "localhost"
var port int = 4040
var dbs []*Collection = make([]*Collection, 0)

var adminKey []string = make([]string, 0)
var userKeys []string = make([]string, 0)

func ReadConfig(ch chan bool) {
	//Reads the contents of the config file into a string
	content, err := os.ReadFile("config.json")
	if err != nil {
		//Channel returns false if there is any error
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
	//Reads IP and Port from the config file
	ip = dat["ip"].(string)
	port = int(dat["port"].(float64))

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
	appname = dat["appname"].(string)
	//Channel returns true if the read is successful
	ch <- true
}

func ReadDatabases(ch chan bool) {
	//Stores the files in the database directory in a list of files
	files, err := ioutil.ReadDir("Databases/")
	if err != nil {
		//Channel returns false if there is any error
		ch <- false
		fmt.Println("Error when reading Database directory")
		return
	}
	for _, file := range files {
		if !file.IsDir() {
			//Creates a new collection for each file in the directory
			name := strings.Split(file.Name(), ".")[0]
			col := new(Collection)
			col.New(name)
			dbs = append(dbs, col)
			fmt.Println(" * Loaded database \"" + name + "\"")
		}
	}
	//Channel returns true if the read is successful
	ch <- true
}

//Returns the contents of a file as a slice of strings
func ReadFileAsLines(filename string) []string {
	//Opens file given in filename
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error when opening " + filename)
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

// func StringToBinary(s string) string {
// 	buf := make([]byte, binary.MaxVarintLen64)
// 	var res string = ""
// 	for _, c := range s {
// 		n := binary.PutUvarint(buf, uint64(c))
// 		res += fmt.Sprintf("%x ", buf[:n])
// 	}
// 	fmt.Println(res)
// 	return res
// }
// func BinaryToString(s string) string {
// 	var bytes [][]byte = make([][]byte, len(s)/3)
// 	j := 0
// 	for i := 0; i < len(s) && j < len(s)/3; i += 3 {
// 		bytes[j] = []byte(s[i : i+2])
// 		j++
// 	}
// 	var res string = ""
// 	for _, b := range bytes {
// 		str := string(b[0]) + string(b[1])
// 		intVar, err := strconv.Atoi(str)
// 		x, n := binary.Uvarint([]byte(str))
// 		fmt.Sprintf("%v", n)
// 		res += fmt.Sprintf("%s", string(x))
// 	}
// 	fmt.Println(res)
// 	return res
// }

//Checksfor a 'new' keyword in the admin file and replaces it with a new uuid
func GenerateAdminApiKey(ch chan bool) {
	lines := ReadFileAsLines("admin.jserv")
	if len(lines) > 0 {
		if lines[len(lines)-1] == "new" {
			randomuuid := uuid.New()
			adminKey = append(adminKey, randomuuid.String())
			lines[len(lines)-1] = randomuuid.String()
			ioutil.WriteFile("admin.jserv", []byte(strings.Join(lines, "\n")), 0644)
			ch <- true
		} else {
			ch <- true
		}
	} else {
		//Channel returns false if there isn't a second line in the file
		fmt.Println("Failed to detect API Key. Write \"new\" on the last line of admin.jserv")
		ch <- false
	}
}

//Checks for a 'new' keyword in the admin file and replaces it with a new uuid
func GenerateUserApiKey(ch chan bool) {
	lines := ReadFileAsLines("keys.jserv")
	if len(lines) > 0 {
		if lines[len(lines)-1] == "new" {
			randomuuid := uuid.New()
			userKeys = append(userKeys, randomuuid.String())
			lines[len(lines)-1] = randomuuid.String()
			ioutil.WriteFile("keys.jserv", []byte(strings.Join(lines, "\n")), 0644)
			ch <- true
		} else {
			ch <- true
		}
	} else {
		//Channel returns false if there isn't a second line in the file
		fmt.Println("Failed to detect API Key. Write \"new\" on the last line of keys.jserv")
		ch <- false
	}
}

//Reads all api keys from admin and keys file
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

//Checks the validity of all the required jserv data files
func CheckFiles() {
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
	version = "0.2.0"
}

//The starting sequence to perform all the necessary checks before the server starts
func StartSequence() {
	fmt.Println(" * Starting...")
	CheckFiles()
	ch := make(chan bool)
	go ReadConfig(ch)
	if result := <-ch; result {
		fmt.Println(" * Successfully read config")
	} else {
		fmt.Println(" * Error while reading config")
		os.Exit(1)
	}
	go ReadDatabases(ch)
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
		fmt.Println(" * Failed to read API Keys")
		os.Exit(1)
	}
	fmt.Printf(" * Running jServ v%s for %s\n", version, appname)
}

//Checks if a string slice contains a string
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

//Checks for a collection of the given name
func FindCollection(c []*Collection, name string) *Collection {
	for _, v := range c {
		if v.name == name {
			return v
		}
	}
	return nil
}

//Checks for an object of the given id in a collection
func FindDataObject(c *Collection, id uint64) *DataObject {
	for _, v := range c.list {
		if v.Id == id {
			return &v
		}
	}
	return nil
}

//Checks for objects of a given attribute in a collection
func FindDataObjects(c *Collection, att string) []*DataObject {
	data := make([]*DataObject, 0)
	for _, v := range c.list {
		for k := range v.Data {
			if k == att {
				data = append(data, &v)
			}
		}
	}
	return data
}

//Checks if the given API key matches the permissions bool of a query type
func CheckApiKey(key string, permission bool) bool {
	if !permission {
		return contains(adminKey, key)
	} else {
		return (contains(userKeys, key) || contains(adminKey, key))
	}
}

func QObject(w http.ResponseWriter, req *http.Request) {
	var end string
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["QObject"]) {
		fmt.Printf("Object query from %s\n", req.RemoteAddr)
		var db string = req.URL.Query().Get("db")
		id, err := strconv.ParseUint(req.URL.Query().Get("id"), 10, 64)
		if err != nil {
			end = " > Failed to parse id query parameter"
		} else {
			fmt.Printf("Queried object %d from %s\n", id, db)
			C := FindCollection(dbs, db)
			if C != nil {
				data := FindDataObject(C, id)
				if data != nil {
					end = data.String()
				} else {
					end = fmt.Sprintf(" > Object %d could not be found in %s", id, db)
				}
			} else {
				end = " > Could not find collection " + db
			}
		}

	} else {
		end = " > Unauthorized Request from " + req.RemoteAddr
	}
	fmt.Println(end)
	fmt.Fprint(w, end)
}

func QAttribute(w http.ResponseWriter, req *http.Request) {
	var end string
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["QAttribute"]) {
		fmt.Printf("Attribute query from %s\n", req.RemoteAddr)
		db := req.URL.Query().Get("db")
		id, err := strconv.ParseUint(req.URL.Query().Get("id"), 10, 64)
		att := req.URL.Query().Get("a")
		if err != nil {
			end = " > Failed to parse id query parameter"
		} else {
			fmt.Printf("Queried attribute %s in %d from %s\n", att, id, db)
			C := FindCollection(dbs, db)
			if C != nil {
				data := FindDataObject(C, id)
				if data != nil {
					if val, ok := data.Data[att]; ok {
						attribute := new(AttributeContainer)
						attribute.New(att, val)
						end = attribute.ToJson()
					} else {
						end = fmt.Sprintf(" > Object %d in %s does not contain %s", id, db, att)
					}
				} else {
					end = fmt.Sprintf(" > Object %d could not be found in %s", id, db)
				}
			} else {
				end = " > Could not find collection " + db
			}
		}
	} else {
		end = " > Unauthorized Request from " + req.RemoteAddr
	}
	fmt.Println(end)
	fmt.Fprint(w, end)
}

func QAllAttributes(w http.ResponseWriter, req *http.Request) {
	var end string
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["QAllAttribute"]) {
		fmt.Printf("All Attributes query from %s\n", req.RemoteAddr)
		db := req.URL.Query().Get("db")
		att := req.URL.Query().Get("a")
		fmt.Printf("Queried objects with attribute %s from %s\n", att, db)
		C := FindCollection(dbs, db)
		if C != nil {
			data := FindDataObjects(C, att)
			endMap := make(map[uint64]interface{})
			if len(data) > 0 {
				for _, v := range data {
					endMap[v.Id] = v.Data[att]
				}
				js, _ := json.Marshal(endMap)
				end = string(js)
			} else {
				end = fmt.Sprintf(" > No objects with attribute %s could be found in %s", att, db)
			}
		} else {
			end = " > Could not find collection " + db
		}
	} else {
		end = " > Unauthorized Request from " + req.RemoteAddr
	}
	fmt.Println(end)
	fmt.Fprint(w, end)
}

func QByAttributes(w http.ResponseWriter, req *http.Request) {
	var end string
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["QByAttribute"]) {
		fmt.Printf("By Attributes query from %s\n", req.RemoteAddr)
		db := req.URL.Query().Get("db")
		att := req.URL.Query().Get("a")
		var attData map[string]interface{}
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&attData)
		if err != nil {
			end = " > Invalid JSON request body"
		} else {
			//Find a way to get the attribute value from the request body
			fmt.Printf("Queried objects with attribute %s from %s\n", att, db)
			C := FindCollection(dbs, db)
			if C != nil {
				data := FindDataObjects(C, att)
				if len(data) > 0 {
					js, _ := json.Marshal(data)
					end = string(js)
				} else {
					end = fmt.Sprintf(" > No objects with attribute %s could be found in %s", att, db)
				}
			} else {
				end = " > Could not find collection " + db
			}
		}
	} else {
		end = " > Unauthorized Request from " + req.RemoteAddr
	}
	fmt.Println(end)
	fmt.Fprint(w, end)
}

func QNewId(w http.ResponseWriter, req *http.Request) {
	var end string
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["QByAttribute"]) {
		fmt.Printf("New ID query from %s\n", req.RemoteAddr)
		db := req.URL.Query().Get("db")
		var attData map[string]interface{}
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&attData)
		if err != nil {
			end = " > Invalid JSON request body"
		} else {
			//Find a way to get the attribute value from the request body
			fmt.Printf("Queried %s for new ID\n", db)
			C := FindCollection(dbs, db)
			if C != nil {
				maxID := uint64(0)
				for _, v := range C.list {
					if v.Id > maxID {
						maxID = v.Id
					}
				}
				maxID += 1
				js, _ := json.Marshal(maxID)
				end = string(js)
			} else {
				end = " > Could not find collection " + db
			}
		}
	} else {
		end = " > Unauthorized Request from " + req.RemoteAddr
	}
	fmt.Println(end)
	fmt.Fprint(w, end)
}

func AEmpty(w http.ResponseWriter, req *http.Request) {
	var end string
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["AEmpty"]) {
		fmt.Printf("Empty object add request from %s\n", req.RemoteAddr)
		db := req.URL.Query().Get("db")
		id, err := strconv.ParseUint(req.URL.Query().Get("id"), 10, 64)
		if err != nil {
			end = " > Failed to parse id query parameter"
		} else {
			fmt.Printf("Requested to add object %d to %s\n", id, db)
			C := FindCollection(dbs, db)
			if C != nil {
				data := FindDataObject(C, id)
				if data == nil {
					obj := DataObject{}
					obj.WithEmptyMap(id)
					C.list = append(C.list, obj)
					C.UpdateFile()
					end = fmt.Sprintf("Successfully added object %d to %s", id, db)
				} else {
					end = fmt.Sprintf(" > Object %d already exists in %s", id, db)
				}
			} else {
				end = " > Could not find collection " + db
			}
		}

	} else {
		end = " > Unauthorized Request from " + req.RemoteAddr
	}
	fmt.Println(end)
	fmt.Fprint(w, end)
}

func AObject(w http.ResponseWriter, req *http.Request) {
	var end string
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["AObject"]) {
		fmt.Printf("Object add request from %s\n", req.RemoteAddr)
		db := req.URL.Query().Get("db")
		var objData map[string]interface{}
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&objData)
		if err != nil {
			end = " > Failed to parse Request body"
		} else {
			obj := DataObject{}
			obj.WithData(uint64(objData["id"].(float64)), objData["data"].(map[string]interface{}))
			fmt.Printf("Requested to add object %d to %s\n", obj.Id, db)
			C := FindCollection(dbs, db)
			if C != nil {
				data := FindDataObject(C, obj.Id)
				if data == nil {
					C.list = append(C.list, obj)
					C.UpdateFile()
					end = fmt.Sprintf("Successfully added object %d to %s", obj.Id, db)
				} else {
					end = fmt.Sprintf(" > Object %d already exists in %s", obj.Id, db)
				}
			} else {
				end = " > Could not find collection " + db
			}
		}

	} else {
		end = " > Unauthorized Request from " + req.RemoteAddr
	}
	fmt.Println(end)
	fmt.Fprint(w, end)
}

func AAttribute(w http.ResponseWriter, req *http.Request) {
	var end string
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["AAttribute"]) {
		fmt.Printf("Attribute add request from %s\n", req.RemoteAddr)
		db := req.URL.Query().Get("db")
		id, err := strconv.ParseUint(req.URL.Query().Get("id"), 10, 64)
		if err != nil {
			end = " > Failed to parse id query parameter"
		} else {
			att := req.URL.Query().Get("a")
			var attData map[string]interface{}
			decoder := json.NewDecoder(req.Body)
			decoder.DisallowUnknownFields()
			err := decoder.Decode(&attData)
			if err != nil {
				end = " > Failed to parse JSON Response body"
			} else {
				fmt.Printf("Requested to add attribute %s to object %d in %s\n", att, id, db)
				C := FindCollection(dbs, db)
				if C != nil {
					data := FindDataObject(C, id)
					if data != nil {
						data.Data[att] = attData[att]
						C.UpdateFile()
						end = fmt.Sprintf("Successfully added attribute %s to object %d in %s", att, id, db)
					} else {
						end = fmt.Sprintf(" > Object %d does not exist in %s", id, db)
					}
				} else {
					end = " > Could not find collection " + db
				}
			}
		}
	} else {
		end = " > Unauthorized Request from " + req.RemoteAddr
	}
	fmt.Println(end)
	fmt.Fprint(w, end)
}

func main() {
	StartSequence()
	http.HandleFunc("/query", QObject)
	http.HandleFunc("/query/attribute", QAttribute)
	http.HandleFunc("/query/allAttributes", QAllAttributes)
	http.HandleFunc("/query/byAttributes", QByAttributes)
	http.HandleFunc("/query/newId", QNewId)
	http.HandleFunc("/add", AEmpty)
	http.HandleFunc("/add/object", AObject)
	http.HandleFunc("/add/attribute", AAttribute)
	fmt.Printf(" * Server bound to %s:%d\n", ip, port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), nil)
	if err != nil {
		log.Fatal(err)
	}
}