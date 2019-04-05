package jsonrpc

import (
	"encoding/json"
	"fmt"
)

func parseBatch(body []byte) ([]Request, error) {
	var requests []Request

	if err := json.Unmarshal(body, &requests); err != nil {
		return nil, fmt.Errorf("unable to parse request body: %s", err.Error())
	}
	return requests, nil
}

func parseIndividual(body []byte) (Request, error) {
	var request Request
	if err := json.Unmarshal(body, &request); err != nil {
		return Request{}, fmt.Errorf("unable to parse request body: %s", err.Error())
	}
	return request, nil
}
