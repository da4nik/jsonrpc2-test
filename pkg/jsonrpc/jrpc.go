package jsonrpc

// ErrorCode - error code type
type ErrorCode int

// JSON RPC 2.0 error codes
const (
	ErrorParse          ErrorCode = 32700
	ErrorInvalidRequest ErrorCode = -32600
	ErrorNoMethod       ErrorCode = -32601
	ErrorBadParams      ErrorCode = -32602
	ErrorInternalError  ErrorCode = -32603
	ErrorServerError    ErrorCode = -32000
)

// M - map for json
type M map[string]interface{}

// ErrorObject - jrpc error object
type ErrorObject struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    M      `json:"data"`
}

// Request - jrpc request object
type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      string      `json:"id"`
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
