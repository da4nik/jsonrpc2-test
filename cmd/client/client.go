package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type m map[string]interface{}

func post(url, contentType string, data interface{}) {
	buf, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("error marshaling data: %s\n", err.Error())
		return
	}

	resp, err := http.Post(url, contentType, bytes.NewReader(buf))
	if err != nil {
		fmt.Printf("error sending request: %s\n", err.Error())
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading body %s\n", err.Error())
		return
	}

	fmt.Printf("Body: %s\n", b)
	var restBody interface{}
	err = json.Unmarshal(b, &restBody)
	if err != nil {
		fmt.Printf("error parsing body %s\n", err.Error())
		return
	}

	fmt.Println("Response: %+v\n", resp)
	fmt.Println("Response body: %+v\n", restBody)
}

func main() {
	requestBatch := []m{
		{
			"jsonrpc": "2.0",
			"method":  "auth.login",
			"params": m{
				"email":    "user@example.com",
				"password": "password",
			},
			"id": "1",
		},
		{
			"jsonrpc": "2.0",
			"method":  "auth.logout",
			"params":  nil,
			"id":      "2",
		},
	}

	requestSingle := m{
		"jsonrpc": "2.0",
		"method":  "auth.login",
		"params": m{
			"email":    "user@example.com",
			"password": "password",
		},
		"id": "1",
	}

	fmt.Printf("\n-------------- Batch\n\n")
	post("http://localhost:8080/rpc", "application/json", requestBatch)
	fmt.Printf("\n-------------- Single\n\n")
	post("http://localhost:8080/rpc", "application/json", requestSingle)
}
