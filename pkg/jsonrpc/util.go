package jsonrpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	output, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)

	return nil
}

func readBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusOK, Response{
			JSONRPC: "2.0",
			Error: ErrorObject{
				Code:    int(ErrorInvalidRequest),
				Message: fmt.Sprintf("Unable to read request body: %s", err.Error()),
			},
		})
		return nil, err
	}
	return body, nil
}
