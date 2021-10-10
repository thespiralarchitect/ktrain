package httputil

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

// type Context struct {
// 	R *http.Request
// 	W http.ResponseWriter
// }
type JsonBinder struct{}
type QueryURLBinder struct{}

type HTTPBinder interface {
	BindRequest(i interface{}, req *http.Request) error
}

func (JsonBinder) BindRequest(i interface{}, req *http.Request) error {
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return errors.New("Error reading body request")
	}
	err = json.Unmarshal(b, i)
	if err != nil {
		return errors.New("Error unmarshal")
	}
	return nil
}
func (QueryURLBinder) BindRequest(obj interface{}, req *http.Request) error {
	values := req.URL.Query()
	if err := mapstructure.Decode(values, obj); err != nil {
		return err
	}
	return nil
}
