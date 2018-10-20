// Package errorcode defines the customzied error code.
// https://confluence.tyo.invastsec.com/pages/resumedraft.action?draftId=5963807&draftShareId=0e50f541-3d27-42b9-a4a6-568bf5d73731
package errorcode

import (
	"encoding/json"
	"net/http"
)

// ids
const (
	pipe      = "passthru"
	generic   = "dongfeng.svc.generic"
	core      = "dongfeng.svc.core.server"
	coreproxy = "dongfeng.svc.core.proxy"
)

var (
	// Pipe means to passthru the customcode from other microservices
	Pipe = new(pipe, "PASSTHRU", "passthru", http.StatusBadGateway)
)

// Generic errors
var (
	// GenericHeIsDeadJim healthcheck dead
	GenericHeIsDeadJim = new(generic, "E0000", "he is dead, Jim", http.StatusBadGateway)
	// GenericNotAuthorized invalid token
	GenericNotAuthorized = new(generic, "E0001", "not authorized", http.StatusUnauthorized)
	// GenericInvalidToken invalid token
	GenericInvalidToken = new(generic, "E0002", "invalid token", http.StatusBadRequest)
	// GenericTokenExpired token expired
	GenericTokenExpired = new(generic, "E0003", "token expired", http.StatusBadRequest)
	// GenericInsufficientPrivileges insufficient privileges
	GenericInsufficientPrivileges = new(generic, "E0004", "insufficient privileges", http.StatusForbidden)
	// GenericInvalidMetaData invalid meta data
	GenericInvalidMetaData = new(generic, "E0005", "invalid meta data", http.StatusBadRequest)
)

// Core-Proxy
var ()

// Core
var ()

// Error implements the error interface.
type Error struct {
	ID         string `json:"id"`
	Code       int32  `json:"code"`
	CustomCode string `json:"custom_code"`
	Detail     string `json:"detail"`
	Status     string `json:"status"`
}

// new private constructor
func new(id, customcode, detail string, code int32) *Error {
	return &Error{
		ID:         id,
		CustomCode: customcode,
		Detail:     detail,
		Code:       code,
	}
}

// NewError new rpc error
func (e *Error) NewError(detail ...string) error {
	if len(detail) == 1 {
		e.Detail = detail[0]
	}
	return newError(e.ID, e.Detail, e.CustomCode, e.Code)
}

// Error error to string
func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func newError(id, detail, customcode string, httpcode int32) error {
	err := parse(detail)
	code := err.Code
	ccode := err.CustomCode

	// is errorcode is pipe, passthru
	if customcode != Pipe.CustomCode {
		code = httpcode
		ccode = customcode
	}

	return &Error{
		ID:         id,
		Code:       code,
		CustomCode: ccode,
		Detail:     detail,
		Status:     http.StatusText(int(code)),
	}
}

// parse tries to parse a JSON string into an error. If that fails, it will set the given string as the error detail.
func parse(err string) *Error {
	e := &Error{}
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Detail = err
	}
	return e
}
