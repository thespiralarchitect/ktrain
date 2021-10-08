package binding

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator"
)

type Context struct {
	R *http.Request
	W http.ResponseWriter
}
type jsonBinding struct{}

var (
	JSON = jsonBinding{}
)

type HTTPBinder interface {
	BindJSONRequest(i interface{}, req *http.Request) error
}
type HTTPShouldBind interface {
	ShouldBind(i interface{}) error
}

func decodeJSON(r io.Reader, i interface{}) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(i); err != nil {
		return err
	}
	return Validator(i)
}
func (jsonBinding) BindJSONRequest(i interface{}, req Context) error {
	if req.R == nil || req.R.Body == nil {
		return fmt.Errorf("invalid request")
	}
	return decodeJSON(req.R.Body, i)
}

func (req Context) ShouldBindJSON(i interface{}, j jsonBinding) error {

	return j.BindJSONRequest(i, req)
}
func (req Context) ShouldBind(i interface{}) error {
	return req.ShouldBindJSON(i, JSON)
}
func Validator(i interface{}) error {
	validate := validator.New()
	return validate.Struct(i)
}
