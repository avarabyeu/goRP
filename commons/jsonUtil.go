package commons

import (
	"encoding/json"
	"net/http"
	"fmt"
)

const CONTENT_TYPE_HEADER string = "Content-Type"

var jsonContentTypeValue = []string{"application/json; charset=utf-8"}
var jsContentTypeValue = []string{"application/javascript; charset=utf-8"}

//WriteJSON serializes body to provided writer
func WriteJSON(status int, body interface{}, w http.ResponseWriter) error {
	header := w.Header()
	if val := header[CONTENT_TYPE_HEADER]; len(val) == 0 {
		header[CONTENT_TYPE_HEADER] = jsonContentTypeValue
	}
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(body)
}

//WriteJSONP serializes body as JSONP
func WriteJSONP(status int, body interface{}, callback string, w http.ResponseWriter) error {
	header := w.Header()
	if val := header[CONTENT_TYPE_HEADER]; len(val) == 0 {
		header[CONTENT_TYPE_HEADER] = jsContentTypeValue
	}
	jsonArr, err := json.Marshal(body)
	if nil != err {
		return err
	}

	w.WriteHeader(status)
	_, err = w.Write([]byte(fmt.Sprintf("%s(%s);", callback, jsonArr)))
	return err
}
