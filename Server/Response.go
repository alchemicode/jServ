package main

import "encoding/json"

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (r *Response) WithoutData(status string, message string) {
	r.Status = status
	r.Message = message
	r.Data = ""
}

func (r *Response) WithData(status string, message string, data string) {
	r.Status = status
	r.Message = message
	r.Data = data
}

//Returns the Object in Json string format
func (r Response) DataToMap() map[string]interface{} {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(r.Data), &data); err != nil {
		panic(err)
	}
	return data
}

func (r Response) ToJson() string {
	data := make(map[string]interface{})
	data["status"] = r.Status
	data["message"] = r.Message
	data["data"] = r.DataToMap()
	js, _ := json.Marshal(data)
	return string(js)
}
