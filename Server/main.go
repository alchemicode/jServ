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
	"QAllObjects":    true,
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
	files, err := ioutil.ReadDir("Databases")
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
			fmt.Println(" * Loaded the database \"" + name + "\"")
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
		if v.Name == name {
			return v
		}
	}
	return nil
}

//Checks for an object of the given id in a collection
func FindDataObject(c *Collection, id uint64) *DataObject {
	for _, v := range c.List {
		if v.Id == id {
			return &v
		}
	}
	return nil
}

//Checks for objects of a given attribute in a collection
func FindDataObjects(c *Collection, att string) []*DataObject {
	data := make([]*DataObject, 0)
	for _, v := range c.List {
		for k := range v.Data {
			if k == att {
				data = append(data, &v)
			}
		}
	}
	return data
}

func RemoveDataObject(c *Collection, id uint64) {
	for i, v := range c.List {
		if v.Id == id {
			c.List[i] = c.List[len(c.List)-1]    // Copy last element to index i.
			c.List[len(c.List)-1] = DataObject{} // Erase last element (write zero value).
			c.List = c.List[:len(c.List)-1]
			break
		}
	}
}

func RemoveAttribute(c *Collection, id uint64, att string) {
	for _, v := range c.List {
		if v.Id == id {
			for j := range v.Data {
				if j == att {
					delete(v.Data, j)
					break
				}
			}
		}
	}
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
	end := Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["QObject"]) {
		fmt.Printf("Object query from %s\n", req.RemoteAddr)
		//Gets necessary query parameters
		var db string = req.URL.Query().Get("db")
		id, err := strconv.ParseUint(req.URL.Query().Get("id"), 10, 64)
		if err != nil {
			end.WithoutData("error", "Failed to parse id query parameter")
		} else {
			fmt.Printf("Queried object %d from %s\n", id, db)
			//Gets reference to collection
			C := FindCollection(dbs, db)
			if C != nil {
				//Gets reference to DataObject with specified id
				data := FindDataObject(C, id)
				if data != nil {
					//Returns DataObject as a JSON object in response
					end.WithData("ok", fmt.Sprintf("Successfully queried object %d from %s\n", id, db), data.ToJson())
				} else {
					end.WithoutData("error", fmt.Sprintf("Object %d could not be found in %s", id, db))
				}
			} else {
				end.WithoutData("error", "Could not find collection "+db)
			}
		}
	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}
	//Changes console message to add '>' prefix if it is an error message
	if end.Status == "ok" {
		fmt.Println(end.Message)
	} else {
		fmt.Println(" > " + end.Message)
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func QAllObjects(w http.ResponseWriter, req *http.Request) {
	end := Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["QAllObjects"]) {
		fmt.Printf("Objects query from %s\n", req.RemoteAddr)
		//Gets necessary query parameters
		var db string = req.URL.Query().Get("db")
		fmt.Printf("Queried objects from %s\n", db)
		//Gets reference to collection
		C := FindCollection(dbs, db)
		if C != nil {
			//Gets reference to DataObject with specified id
			var data []map[string]interface{}
			for _, v := range C.List {
				data = append(data, v.ToMap())
			}
			if data != nil {
				if js, err := json.Marshal(data); err != nil {
					end.WithoutData("error", "Failed to parse list to JSON")
				} else {
					//Returns DataObject as a JSON object in response
					end.WithData("ok", fmt.Sprintf("Successfully queried objects from %s\n", db), string(js))
				}
			} else {
				end.WithoutData("error", fmt.Sprintf("No objects could not be found in %s", db))
			}
		} else {
			end.WithoutData("error", "Could not find collection "+db)
		}

	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}
	//Changes console message to add '>' prefix if it is an error message
	if end.Status == "ok" {
		fmt.Println(end.Message)
	} else {
		fmt.Println(" > " + end.Message)
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func QAttribute(w http.ResponseWriter, req *http.Request) {
	end := Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["QAttribute"]) {
		fmt.Printf("Attribute query from %s\n", req.RemoteAddr)
		//Gets necessary query parameters
		db := req.URL.Query().Get("db")
		id, err := strconv.ParseUint(req.URL.Query().Get("id"), 10, 64)
		att := req.URL.Query().Get("a")
		if err != nil {
			end.WithoutData("error", "Failed to parse id query parameter")
		} else {
			fmt.Printf("Queried attribute %s in %d from %s\n", att, id, db)
			//Gets reference to collection
			C := FindCollection(dbs, db)
			if C != nil {
				//Gets reference to DataObject with specified id
				data := FindDataObject(C, id)
				if data != nil {
					//Gets attribute from the DataObject
					if val, ok := data.Data[att]; ok {
						attribute := new(AttributeContainer)
						attribute.New(att, val)
						//Returns AttributeContainer as a JSON Object in response
						end.WithData("ok", fmt.Sprintf("Successfully queried object %d from %s\n", id, db), attribute.ToJson())
					} else {
						end.WithoutData("error", fmt.Sprintf("Object %d in %s does not contain %s", id, db, att))
					}
				} else {
					end.WithoutData("error", fmt.Sprintf("Object %d could not be found in %s", id, db))
				}
			} else {
				end.WithoutData("error", "Could not find collection "+db)
			}
		}
	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}
	//Changes console message to add '>' prefix if it is an error message
	if end.Status == "ok" {
		fmt.Println(end.Message)
	} else {
		fmt.Println(" > " + end.Message)
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func QAllAttributes(w http.ResponseWriter, req *http.Request) {
	end := Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["QAllAttribute"]) {
		fmt.Printf("All Attributes query from %s\n", req.RemoteAddr)
		//Gets necessary query parameters
		db := req.URL.Query().Get("db")
		att := req.URL.Query().Get("a")
		fmt.Printf("Queried objects with attribute %s from %s\n", att, db)
		//Gets reference to collection
		C := FindCollection(dbs, db)
		if C != nil {
			//Gets reference to DataObjects with specified attribute
			data := FindDataObjects(C, att)
			//Makes list for all the ids of the objects
			var list []uint64
			if len(data) > 0 {
				for _, v := range data {
					//Adds all the object ids to the list
					list = append(list, v.Id)
				}
				if js, err := json.Marshal(list); err != nil {
					end.WithoutData("error", "Failed to parse list to JSON")
				} else {
					//Returns the list of DataObject ids in the response
					end.WithData("ok", fmt.Sprintf("Successfully queried objects with attribute %s from %s\n", att, db), string(js))
				}
			} else {
				end.WithoutData("error", fmt.Sprintf("No objects with attribute %s could be found in %s", att, db))
			}
		} else {
			end.WithoutData("error", "Could not find collection "+db)
		}
	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}
	//Changes console message to add '>' prefix if it is an error message
	if end.Status == "ok" {
		fmt.Println(end.Message)
	} else {
		fmt.Println(" > " + end.Message)
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func QByAttributes(w http.ResponseWriter, req *http.Request) {
	end := Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["QByAttribute"]) {
		fmt.Printf("By Attributes query from %s\n", req.RemoteAddr)
		//Gets necessary query parameters
		db := req.URL.Query().Get("db")
		att := req.URL.Query().Get("a")
		//Reads JSON from request body
		var attData map[string]interface{}
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&attData)
		if err != nil {
			end.WithoutData("error", "Invalid JSON request body")
		} else {
			fmt.Printf("Queried objects with attribute %s from %s\n", att, db)
			//Gets reference to collection
			C := FindCollection(dbs, db)
			if C != nil {
				//Gets reference to DataObjects with specified attribute
				data := FindDataObjects(C, att)
				var list []map[string]interface{}
				if len(data) > 0 {
					for _, v := range data {
						list = append(list, v.ToMap())
					}
					if js, err := json.Marshal(list); err != nil {
						end.WithoutData("error", "Failed to parse list to JSON")
					} else {
						//Returns list of DataObjects in the response
						end.WithData("ok", fmt.Sprintf("Queried objects with attribute %s from %s\n", att, db), string(js))
					}

				} else {
					end.WithoutData("error", fmt.Sprintf("No objects with attribute %s could be found in %s", att, db))
				}
			} else {
				end.WithoutData("error", "Could not find collection "+db)
			}
		}
	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}
	//Changes console message to add '>' prefix if it is an error message
	if end.Status == "ok" {
		fmt.Println(end.Message)
	} else {
		fmt.Println(" > " + end.Message)
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func QNewId(w http.ResponseWriter, req *http.Request) {
	end := Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["QNewId"]) {
		fmt.Printf("New ID query from %s\n", req.RemoteAddr)
		//Gets necessary query parameters
		db := req.URL.Query().Get("db")
		fmt.Printf("Queried %s for new id\n", db)
		//Gets reference to collection
		C := FindCollection(dbs, db)
		if C != nil {
			//Finds the next unused id
			maxID := uint64(0)
			for _, v := range C.List {
				if v.Id > maxID {
					maxID = v.Id
				}
			}
			maxID += 1
			//Creates AttributeContainer to return the new id
			ac := AttributeContainer{}
			ac.New("id", maxID)
			//Returns id AttributeContainer in response
			end.WithData("ok", fmt.Sprintf("Queried %s for new ID\n", db), ac.ToJson())
		} else {
			end.WithoutData("error", "Could not find collection "+db)
		}
	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}
	//Changes console message to add '>' prefix if it is an error message
	if end.Status == "ok" {
		fmt.Println(end.Message)
	} else {
		fmt.Println(" > " + end.Message)
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func AEmpty(w http.ResponseWriter, req *http.Request) {
	end := Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["AEmpty"]) {
		fmt.Printf("Empty object add request from %s\n", req.RemoteAddr)
		//Gets necessary query parameters
		db := req.URL.Query().Get("db")
		id, err := strconv.ParseUint(req.URL.Query().Get("id"), 10, 64)
		if err != nil {
			end.WithoutData("error", "Failed to parse id query parameter")
		} else {
			fmt.Printf("Requested to add object %d to %s\n", id, db)
			//Gets reference to collection
			C := FindCollection(dbs, db)
			if C != nil {
				//Gets reference to DataObject of specified id
				data := FindDataObject(C, id)
				if data == nil {
					//Creates new empty DataObject
					obj := DataObject{}
					obj.WithoutData(id)
					//Updates collection
					C.List = append(C.List, obj)
					C.UpdateFile()
					//Returns confirmation message
					end.WithoutData("ok", fmt.Sprintf("Successfully added object %d to %s", id, db))
				} else {
					end.WithoutData("error", fmt.Sprintf("Object %d already exists in %s", id, db))
				}
			} else {
				end.WithoutData("error", "Could not find collection "+db)
			}
		}
	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}
	//Changes console message to add '>' prefix if it is an error message
	if end.Status == "ok" {
		fmt.Println(end.Message)
	} else {
		fmt.Println(" > " + end.Message)
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func AObject(w http.ResponseWriter, req *http.Request) {
	end := Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["AObject"]) {
		fmt.Printf("Object add request from %s\n", req.RemoteAddr)
		//Gets necessary query parameters
		db := req.URL.Query().Get("db")
		//Reads JSON from request body
		var objData map[string]interface{}
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&objData)
		if err != nil {
			end.WithoutData("error", "Failed to parse Request body")
		} else {
			//Creates new DataObject to be added
			obj := DataObject{}
			obj.WithData(uint64(objData["id"].(float64)), objData["data"].(map[string]interface{}))
			fmt.Printf("Requested to add object %d to %s\n", obj.Id, db)
			//Gets reference to collection
			C := FindCollection(dbs, db)
			if C != nil {
				//Tries getting reference to DataObject with new id
				//If the id already exists in the collection the request returns an error
				data := FindDataObject(C, obj.Id)
				if data == nil {
					//Adds new DataObject to collection
					C.List = append(C.List, obj)
					//Updates collection
					C.UpdateFile()
					//Returns confirmation message
					end.WithoutData("ok", fmt.Sprintf("Successfully added object %d to %s", obj.Id, db))
				} else {
					end.WithoutData("error", fmt.Sprintf(" > Object %d already exists in %s", obj.Id, db))
				}
			} else {
				end.WithoutData("error", "Could not find collection "+db)
			}
		}
	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}
	//Changes console message to add '>' prefix if it is an error message
	if end.Status == "ok" {
		fmt.Println(end.Message)
	} else {
		fmt.Println(" > " + end.Message)
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func AAttribute(w http.ResponseWriter, req *http.Request) {
	end := Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["AAttribute"]) {
		fmt.Printf("Attribute add request from %s\n", req.RemoteAddr)
		//Gets necessary query parameters
		db := req.URL.Query().Get("db")
		id, err := strconv.ParseUint(req.URL.Query().Get("id"), 10, 64)
		if err != nil {
			end.WithoutData("error", "Failed to parse id query parameter")
		} else {
			att := req.URL.Query().Get("a")
			//Reads JSON from request body
			var attData map[string]interface{}
			decoder := json.NewDecoder(req.Body)
			decoder.DisallowUnknownFields()
			err := decoder.Decode(&attData)
			if err != nil {
				end.WithoutData("error", "Failed to parse JSON Response body")
			} else {
				fmt.Printf("Requested to add attribute %s to object %d in %s\n", att, id, db)
				//Gets reference to collection
				C := FindCollection(dbs, db)
				if C != nil {
					//Gets reference to DataObject with specified id
					data := FindDataObject(C, id)
					if data != nil {
						//Adds attribute to the DataObject
						data.Data[att] = attData[att]
						//Updates collection
						C.UpdateFile()
						//Returns confirmation message
						end.WithoutData("ok", fmt.Sprintf("Successfully added attribute %s to object %d in %s", att, id, db))
					} else {
						end.WithoutData("error", fmt.Sprintf("Object %d could not be found in %s", id, db))
					}
				} else {
					end.WithoutData("error", "Could not find collection "+db)
				}
			}
		}
	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}
	//Changes console message to add '>' prefix if it is an error message
	if end.Status == "ok" {
		fmt.Println(end.Message)
	} else {
		fmt.Println(" > " + end.Message)
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func MObject(w http.ResponseWriter, req *http.Request) {
	end := Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["MObject"]) {
		fmt.Printf("Modify object request from %s\n", req.RemoteAddr)
		//Gets necessary query parameters
		db := req.URL.Query().Get("db")
		id, err := strconv.ParseUint(req.URL.Query().Get("id"), 10, 64)
		newId, err2 := strconv.ParseUint(req.URL.Query().Get("n"), 10, 64)
		if err != nil || err2 != nil {
			end.WithoutData("error", "Failed to parse id query parameter")
		} else {
			fmt.Printf("Requested to mod object %d to %d\n", id, newId)
			//Gets reference to collection
			C := FindCollection(dbs, db)
			if C != nil {
				//Gets reference to DataObject with specified id
				data := FindDataObject(C, id)
				//Tries getting reference to DataObject with new id
				sameData := FindDataObject(C, newId)
				if sameData != nil {
					end.WithoutData("error", fmt.Sprintln("Object %d already exists in %s", newId, db))
				} else if data != nil {
					//Changes DataObject id
					FindDataObject(C, id).Id = newId
					//Updates collection
					C.UpdateFile()
					//Returns confirmation message
					end.WithoutData("error", fmt.Sprintf("Successfully modded object %d to %d", id, newId))
				} else {
					end.WithoutData("error", fmt.Sprintf("Object %d could not be found in %s", id, db))
				}
			} else {
				end.WithoutData("error", "Could not find collection "+db)
			}
		}
	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}
	//Changes console message to add '>' prefix if it is an error message
	if end.Status == "ok" {
		fmt.Println(end.Message)
	} else {
		fmt.Println(" > " + end.Message)
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func MAttribute(w http.ResponseWriter, req *http.Request) {
	end := Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["MAttribute"]) {
		fmt.Printf("Modify Attribute request from %s\n", req.RemoteAddr)
		//Gets necessary query parameters
		db := req.URL.Query().Get("db")
		id, err := strconv.ParseUint(req.URL.Query().Get("id"), 10, 64)
		if err != nil {
			end.WithoutData("error", "Failed to parse id query parameter")
		} else {
			att := req.URL.Query().Get("a")
			//Reads JSON from request body
			var attData map[string]interface{}
			decoder := json.NewDecoder(req.Body)
			decoder.DisallowUnknownFields()
			err := decoder.Decode(&attData)
			if err != nil {
				end.WithoutData("error", "Failed to parse JSON Response body")
			} else {
				fmt.Printf("Requested to modify attribute %s in object %d in %s\n", att, id, db)
				//Gets reference to collection
				C := FindCollection(dbs, db)
				if C != nil {
					//Gets reference to DataObject with specified id
					data := FindDataObject(C, id)
					if data != nil {
						//Changes attribute of DataObject
						data.Data[att] = attData[att]
						//Updates collection
						C.UpdateFile()
						//Returns confirmation message
						end.WithoutData("ok", fmt.Sprintf("Successfully modified attribute %s in object %d in %s", att, id, db))
					} else {
						end.WithoutData("error", fmt.Sprintf("Object %d could not be found in %s", id, db))
					}
				} else {
					end.WithoutData("error", "Could not find collection "+db)
				}
			}
		}
	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}
	//Changes console message to add '>' prefix if it is an error message
	if end.Status == "ok" {
		fmt.Println(end.Message)
	} else {
		fmt.Println(" > " + end.Message)
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func DObject(w http.ResponseWriter, req *http.Request) {
	end := Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["DObject"]) {
		fmt.Printf("Delete object request from %s\n", req.RemoteAddr)
		//Gets necessary query parameters
		db := req.URL.Query().Get("db")
		id, err := strconv.ParseUint(req.URL.Query().Get("id"), 10, 64)
		if err != nil {
			end.WithoutData("error", "Failed to parse id query parameter")
		} else {
			fmt.Printf("Requested to delete object %d\n", id)
			//Gets reference to collection
			C := FindCollection(dbs, db)
			if C != nil {
				//Gets reference to DataObject with specified id
				data := FindDataObject(C, id)
				if data != nil {
					//Removes DataObject from collection
					RemoveDataObject(C, id)
					//Update collection
					C.UpdateFile()
					//Returns confirmation message
					end.WithoutData("ok", fmt.Sprintf("Successfully deleted object %d", id))
				} else {
					end.WithoutData("error", fmt.Sprintf("Object %d could not be found in %s", id, db))
				}
			} else {
				end.WithoutData("error", "Could not find collection "+db)
			}
		}
	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}
	//Changes console message to add '>' prefix if it is an error message
	if end.Status == "ok" {
		fmt.Println(end.Message)
	} else {
		fmt.Println(" > " + end.Message)
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func DAttribute(w http.ResponseWriter, req *http.Request) {
	end := Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["DAttribute"]) {
		fmt.Printf("Delete attribute request from %s\n", req.RemoteAddr)
		//Gets necessary query parameters
		db := req.URL.Query().Get("db")
		id, err := strconv.ParseUint(req.URL.Query().Get("id"), 10, 64)
		if err != nil {
			end.WithoutData("error", "Failed to parse id query parameter")
		} else {
			att := req.URL.Query().Get("a")
			fmt.Printf("Requested to delete attribute %s in object %d in %s\n", att, id, db)
			//Gets reference to collection
			C := FindCollection(dbs, db)
			if C != nil {
				//Gets reference to DataObject with specified id
				data := FindDataObject(C, id)
				if data != nil {
					if data.Data[att] != nil {
						//Removes attribute from DataObject
						RemoveAttribute(C, id, att)
						//Updates collection
						C.UpdateFile()
						//Sends confirmation message
						end.WithoutData("ok", fmt.Sprintf("Successfully modified attribute %s in object %d in %s", att, id, db))
					} else {
						end.WithoutData("error", fmt.Sprintf("Attribute %s does not exist in %d in %s", att, id, db))
					}
				} else {
					end.WithoutData("error", fmt.Sprintf("Object %d could not be found in %s", id, db))
				}
			} else {
				end.WithoutData("error", "Could not find collection "+db)
			}
		}
	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}
	//Changes console message to add '>' prefix if it is an error message
	if end.Status == "ok" {
		fmt.Println(end.Message)
	} else {
		fmt.Println(" > " + end.Message)
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func main() {
	StartSequence()
	http.HandleFunc("/query", QObject)
	http.HandleFunc("/query/objects", QAllObjects)
	http.HandleFunc("/query/attribute", QAttribute)
	http.HandleFunc("/query/allAttributes", QAllAttributes)
	http.HandleFunc("/query/byAttribute", QByAttributes)
	http.HandleFunc("/query/newId", QNewId)
	http.HandleFunc("/add", AEmpty)
	http.HandleFunc("/add/object", AObject)
	http.HandleFunc("/add/attribute", AAttribute)
	http.HandleFunc("/mod/object", MObject)
	http.HandleFunc("/mod/attribute", MAttribute)
	http.HandleFunc("/delete/object", DObject)
	http.HandleFunc("/delete/attribute", DAttribute)
	fmt.Printf(" * Server bound to %s:%d\n", ip, port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
