package jsonrpc

import (
	"fmt"
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

type RPCFunc func(params map[string]interface{}) (interface{}, error)
type FuncMap map[string]RPCFunc

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
	JSONRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	ID      *string                `json:"id"`
}

// BatchRequest - jrpc batch request object
type BatchRequest []Request

// Response - jrpc response object
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   ErrorObject `json:"error,omitempty"`
	ID      *string     `json:"id"`
}

func NewJSONRPC(fmap FuncMap) JRPCServer {
	return JRPCServer{
		funcMap: fmap,
	}
}

func (srv JRPCServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(w, r)
	if err != nil {
		return
	}

	req, errInd := parseIndividual(body)
	if errInd == nil {
		result := srv.processRequest(&req)
		writeJSON(w, http.StatusOK, result)
		return
	}

	reqs, errBatch := parseBatch(body)
	if errBatch == nil {
		var result = srv.processBatch(reqs)
		writeJSON(w, http.StatusOK, result)
		return
	}

	writeJSON(w, http.StatusOK, Response{
		JSONRPC: "2.0",
		Error: ErrorObject{
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
			Error: ErrorObject{
				Code:    int(ErrorNoMethod),
				Message: fmt.Sprintf("Method %s is not supported", req.Method),
			},
			ID: req.ID,
		}
	}

	log.Debugf("Processing procedure \"%s\" with %s", req.Method, req.Params)
	result, err := fun(req.Params)
	if err != nil {
		return Response{
			JSONRPC: "2.0",
			Error: ErrorObject{
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
