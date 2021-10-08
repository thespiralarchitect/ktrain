package httputil

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// type Context struct {
// 	R *http.Request
// 	W http.ResponseWriter
// }
type JsonBinding struct{}

type HTTPBinder interface {
	BindJSONRequest(i interface{}, req *http.Request) error
}
// func decodeJSON(r io.Reader, i interface{}, w http.ResponseWriter) error {
// 	decoder := json.NewDecoder(r)
// 	if err := decoder.Decode(i); err != nil {
// 		RespondError(w, http.StatusInternalServerError, "Error unmarshal body request")
// 		return err
// 	}
// 	return nil
// }
func(JsonBinding) BindJSONRequest(i interface{}, req *http.Request) error {
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return errors.New("Error reading body request")
	}
	err = json.Unmarshal(b,i)
	if err != nil {
		return errors.New("Error unmarshal")
	}
	return nil
}