package main

import "encoding/json"

type Response struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func (r *Response) WithoutData(status string, message string) {
	r.Status = status
	r.Message = message
	r.Data = make(map[string]interface{})
}

func (r *Response) WithData(status string, message string, data map[string]interface{}) {
	r.Status = status
	r.Message = message
	r.Data = data
}

//Returns the Object in Json string format
func (r Response) ToJson() string {
	js, _ := json.Marshal(r)
	return string(js)
}
