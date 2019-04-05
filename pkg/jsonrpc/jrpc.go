package jsonrpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/da4nik/jrpc2_try/internal/log"
)

// ErrorCode - error code type
type ErrorCode int

// JSON RPC 2.0 error codes
const (
	ErrorParse          ErrorCode = -32700
	ErrorInvalidRequest ErrorCode = -32600
	ErrorNoMethod       ErrorCode = -32601
	ErrorBadParams      ErrorCode = -32602
	ErrorInternalError  ErrorCode = -32603
	ErrorServerError    ErrorCode = -32000
)

// RPCFunc - function signature to process RPC call
type RPCFunc func(rawParams []byte) (interface{}, error)

// FuncMap map remote procedure names and functions
type FuncMap map[string]RPCFunc

// JRPCServer - json rpc server instance
type JRPCServer struct {
	funcMap FuncMap
}

// ErrorObject - jrpc error object
type ErrorObject struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// Request - jrpc request object
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      *string         `json:"id"`
}

// BatchRequest - jrpc batch request object
type BatchRequest []Request

// Response - jrpc response object
type Response struct {
	JSONRPC string       `json:"jsonrpc"`
	Result  interface{}  `json:"result,omitempty"`
	Error   *ErrorObject `json:"error,omitempty"`
	ID      *string      `json:"id"`
}

// NewJSONRPC - creates new json rpc server instance
func NewJSONRPC(fmap FuncMap) JRPCServer {
	return JRPCServer{
		funcMap: fmap,
	}
}

func (srv JRPCServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := srv.readBody(w, r)
	if err != nil {
		return
	}

	req, errInd := srv.parseIndividual(body)
	if errInd == nil {
		result := srv.processRequest(&req)
		writeJSON(w, http.StatusOK, result)
		return
	}

	reqs, errBatch := srv.parseBatch(body)
	if errBatch == nil {
		result := srv.processBatch(reqs)
		writeJSON(w, http.StatusOK, result)
		return
	}

	writeJSON(w, http.StatusOK, Response{
		JSONRPC: "2.0",
		Error: &ErrorObject{
			Code:    int(ErrorParse),
			Message: fmt.Sprintf("Unable to parse request body: %s%s", errInd.Error(), errBatch.Error()),
		},
	})
}

func (srv JRPCServer) processRequest(req *Request) Response {
	fun, exists := srv.funcMap[req.Method]
	if !exists {
		return Response{
			JSONRPC: "2.0",
			Error: &ErrorObject{
				Code:    int(ErrorNoMethod),
				Message: fmt.Sprintf("Method %s is not supported", req.Method),
			},
			ID: req.ID,
		}
	}

	log.Debugf("Processing procedure \"%s\" with %s", req.Method, req.Params)
	result, err := fun([]byte(req.Params))
	if err != nil {
		return Response{
			JSONRPC: "2.0",
			Error: &ErrorObject{
				Code:    -32000,
				Message: err.Error(),
			},
			ID: req.ID,
		}
	}

	return Response{
		JSONRPC: "2.0",
		Result:  result,
		ID:      req.ID,
		Error:   nil,
	}
}

func (srv JRPCServer) processBatch(reqs []Request) []Response {
	var result []Response

	var respChan = make(chan Response, len(reqs))
	defer close(respChan)

	var wg sync.WaitGroup
	wg.Add(len(reqs))

	for _, req := range reqs {
		go func(req Request, wg *sync.WaitGroup) {
			respChan <- srv.processRequest(&req)
			wg.Done()
		}(req, &wg)
	}

	wg.Wait()
	for i := 0; i < len(reqs); i++ {
		result = append(result, <-respChan)
	}

	return result
}

func (srv JRPCServer) readBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusOK, Response{
			JSONRPC: "2.0",
			Error: &ErrorObject{
				Code:    int(ErrorInvalidRequest),
				Message: fmt.Sprintf("Unable to read request body: %s", err.Error()),
			},
		})
		return nil, err
	}
	return body, nil
}

func (srv JRPCServer) parseBatch(body []byte) ([]Request, error) {
	var requests []Request

	if err := json.Unmarshal(body, &requests); err != nil {
		return nil, fmt.Errorf("unable to parse request body: %s", err.Error())
	}
	return requests, nil
}

func (srv JRPCServer) parseIndividual(body []byte) (Request, error) {
	var request Request
	if err := json.Unmarshal(body, &request); err != nil {
		return Request{}, fmt.Errorf("unable to parse request body: %s", err.Error())
	}
	return request, nil
}
