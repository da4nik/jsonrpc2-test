package procedures

import (
	"encoding/json"
	"fmt"

	"github.com/da4nik/jrpc2_try/internal/log"
	"github.com/da4nik/jrpc2_try/pkg/jsonrpc"
	"github.com/da4nik/jrpc2_try/pkg/services/auth"
)

// Map returns procedures map for json rpc
func Map() jsonrpc.FuncMap {
	return jsonrpc.FuncMap{
		"auth.login": login,
	}
}

type loginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func login(rawParams []byte) (interface{}, error) {
	var params loginParams
	err := json.Unmarshal(rawParams, &params)
	if err != nil {
		return nil, fmt.Errorf("bad params")
	}

	token, err := auth.Authenticate(params.Email, params.Password)
	if err != nil {
		return nil, err
	}

	log.Debugf("Procedure: login")
	log.Debugf("params %+v", params)

	return map[string]string{"token": token}, nil
}
