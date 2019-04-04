package jsonrpc

import (
	"fmt"
	"net/http"
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
	Data    map[string]interface{} `json:"data"`
}

// Request - jrpc request object
type Request struct {
	JSONRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	ID      string                 `json:"id"`
}

// BatchRequest - jrpc batch request object
type BatchRequest []Request

// Response - jrpc response object
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   ErrorObject `json:"error,omitempty"`
	ID      string      `json:"id"`
}

func NewJSONRPC(fmap FuncMap) (JRPCServer, error) {
	return JRPCServer{
		funcMap: fmap,
	}, nil
}

func (srv JRPCServer) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := parseBody(w, r)
		if err != nil {
			return
		}

		fun, exists := srv.funcMap[req.Method]
		if !exists {
			writeJSON(w, http.StatusOK, Response{
				JSONRPC: "2.0",
				Error: ErrorObject{
					Code:    int(ErrorNoMethod),
					Message: fmt.Sprintf("Method %s is not supported", req.Method),
				},
			})
			return
		}

		result, err := fun(req.Params)
		if err != nil {
			writeJSON(w, http.StatusOK, Response{
				JSONRPC: "2.0",
				Error: ErrorObject{
					Code:    -32000,
					Message: err.Error(),
				},
			})
			return
		}

		writeJSON(w, http.StatusOK, Response{
			JSONRPC: "2.0",
			Result:  result,
			ID:      req.ID,
		})
	}
}
