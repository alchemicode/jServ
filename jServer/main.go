package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	core "github.com/alchemicode/jserv-core"

	"alchemicode.com/jServ/data"
)

var requestTypes = map[string]bool{
	//False denotes that the server cannot receive that request type
	"GET":     true,
	"POST":    true,
	"PUT":     false,
	"HEAD":    false,
	"DELETE":  true,
	"PATCH":   false,
	"OPTIONS": false,
}

var requestPermissions = map[string]bool{
	//False denotes that an admin API key is required to make that request
	"Query":  true,
	"Add":    true,
	"Mod":    true,
	"Delete": true,
	"Purge":  true,
}

var aliases = map[string]string{
	"127.0.0.1": "localhost",
}

var services = map[string]string{
	"/": "index",
}

var version string
var appname string
var debug bool
var ip string = "localhost"
var port int = 4040
var write_interval float64 = 10
var cols []*core.Collection = make([]*core.Collection, 0)

var db_free = true

var adminKey []string = make([]string, 0)
var userKeys []string = make([]string, 0)

func Query(w http.ResponseWriter, req *http.Request) {
	end := core.Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["Query"]) {
		query := new(core.Query)
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(query)
		if err != nil {
			end.WithoutData("error", "Failed to parse Request Body")
		} else {
			fmt.Printf("(%s) Query request to %s\n", CheckAlias(req.RemoteAddr), strings.Join(query.Collections, ", "))
			if debug {
				fmt.Println(" ~ \n" + query.ToJson())
			}
			ret := data.ProcessQuery(cols, query)
			data := map[string]interface{}{"documents": ret}
			if err != nil {
				end.WithoutData("error", "Failed to parse queried data")
			}
			end.WithData("ok", "Successful Query", data)
		}
	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}

	if end.Status != "ok" {
		fmt.Println(" > ! Request failed")
		if debug {
			fmt.Println(" ~ " + end.ToJson())
		}
	} else {
		fmt.Println(" > Query Successful")
		if debug {
			fmt.Println(" ~ \n" + end.ToJson())
		}
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())

}

func Add(w http.ResponseWriter, req *http.Request) {
	end := core.Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["Add"]) {
		//Gets necessary query parameters
		add := new(core.Add)
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(add); err != nil {
			end.WithoutData("error", "Failed to parse Request Body")
		} else {
			db_free = false
			fmt.Printf("(%s) Add request to %s\n", CheckAlias(req.RemoteAddr), add.Collection)
			if debug {
				fmt.Println(" ~ \n" + add.ToJson())
			}
			data.ProcessAdd(cols, add, &end)
			db_free = true
		}

	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}

	if end.Status != "ok" {
		fmt.Println(" > ! Request failed")
		if debug {
			fmt.Println(" ~ " + end.ToJson())
		}
	} else {
		fmt.Println(" > Add Successful")
		if debug {
			fmt.Println(" ~ \n" + end.ToJson())
		}
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func Mod(w http.ResponseWriter, req *http.Request) {
	end := core.Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["Mod"]) {
		//Gets necessary query parameters
		mod := new(core.Mod)
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(mod)
		if err != nil {
			end.WithoutData("error", "Failed to parse Request Body")
		} else {
			db_free = false
			fmt.Printf("(%s) Mod request to %s\n", CheckAlias(req.RemoteAddr), mod.Collection+"->doc_"+mod.Document)
			if debug {
				fmt.Println(" ~ \n" + mod.ToJson())
			}
			data.ProcessModify(cols, mod, &end)
			db_free = true
		}
	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}
	//Changes console message to add '>' prefix if it is an error message
	if end.Status != "ok" {
		fmt.Println(" > ! Request failed")
		if debug {
			fmt.Println(" ~ " + end.ToJson())
		}
	} else {
		fmt.Println(" > Mod Successful")
		if debug {
			fmt.Println(" ~ \n" + end.ToJson())
		}
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func Delete(w http.ResponseWriter, req *http.Request) {
	end := core.Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["Delete"]) {
		del := new(core.Delete)
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(del)
		if err != nil {
			end.WithoutData("error", "Failed to parse Request Body")
		} else {
			db_free = false
			fmt.Printf("(%s) Delete request to %s\n", CheckAlias(req.RemoteAddr), del.Collection+"->doc_"+del.Document)
			if debug {
				fmt.Println(" ~ \n" + del.ToJson())
			}
			data.ProcessDelete(cols, del, &end)
			db_free = true
		}
	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}
	//Changes console message to add '>' prefix if it is an error message
	if end.Status != "ok" {
		fmt.Println(" > ! Request failed")
		if debug {
			fmt.Println(" ~ " + end.ToJson())
		}
	} else {
		fmt.Println(" > Delete Successful")
		if debug {
			fmt.Println(" ~ \n" + end.ToJson())
		}
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func Purge(w http.ResponseWriter, req *http.Request) {
	end := core.Response{}
	if CheckApiKey(req.Header.Get("x-api-key"), requestPermissions["Purge"]) {
		del := new(core.Query)
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(del)
		if err != nil {
			end.WithoutData("error", "Failed to parse Request Body")
		} else {
			db_free = false
			fmt.Printf("(%s) Purge request to %s\n", CheckAlias(req.RemoteAddr), strings.Join(del.Collections, ", "))
			if debug {
				fmt.Println(" ~ \n" + del.ToJson())
			}
			data.ProcessPurge(cols, del, &end)
			db_free = true
		}
	} else {
		end.WithoutData("error", "Unauthorized Request from "+req.RemoteAddr)
	}
	//Changes console message to add '>' prefix if it is an error message
	if end.Status != "ok" {
		fmt.Println(" > ! Request failed")
		if debug {
			fmt.Println(" ~ " + end.ToJson())
		}
	} else {
		fmt.Println(" > Purge Successful")
		if debug {
			fmt.Println(" ~ \n" + end.ToJson())
		}
	}
	//Writes response to http response
	fmt.Fprint(w, end.ToJson())
}

func main() {
	// ep, err := python.NewEmbeddedPython("example")
	// if err != nil {
	// 	panic(err)
	// }
	// cmd := ep.PythonCmd("-c", "print('hello')")
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err = cmd.Run()
	StartSequence(&version, &appname)
	http.HandleFunc("/j/db/query", Query)
	http.HandleFunc("/j/db/add", Add)
	http.HandleFunc("/j/db/mod", Mod)
	http.HandleFunc("/j/db/delete", Delete)
	http.HandleFunc("/j/db/purge", Purge)
	//https.HandleFunc("/", )

	server_done := false
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer func() { server_done = true }()
		fmt.Printf(" * Server bound to %s:%d\n", ip, port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Duration(write_interval) * time.Second)
		for !server_done {
			<-ticker.C
			if db_free {
				for _, col := range cols {
					UpdateCollection(col)
				}
				if debug {
					fmt.Println(" ~ Updated Collection")
				}
			}
		}
	}()
	wg.Wait()
}
