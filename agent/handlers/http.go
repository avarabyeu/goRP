package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// StatusErr describes HTTP status error
type StatusErr struct {
	error
	statusCode int
}

// StackTrace returns stack trace of an error
func (se *StatusErr) StackTrace() errors.StackTrace {
	if st, ok := se.error.(stackTracer); ok {
		return st.StackTrace()
	}
	return nil
}

// NewStatusErr creates new status error
func NewStatusErr(statusCode int, err error) *StatusErr {
	return &StatusErr{statusCode: statusCode, error: err}
}

// HTTPHandlerFunc is a handler func for JSON/REST handlers
type HTTPHandlerFunc func(w http.ResponseWriter, rq *http.Request) error

// HTTPHandler handles HTTP requests
func HTTPHandler(f HTTPHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, rq *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		err := f(w, rq)
		if err != nil {
			var rs interface{}

			if se, ok := err.(*StatusErr); ok {
				w.WriteHeader(se.statusCode)
			}

			rs = map[string]string{"error": err.Error()}
			if st, ok := err.(stackTracer); ok {
				rs.(map[string]string)["stacktrace"] = fmt.Sprintf("%+v", st.StackTrace())
			}

			if wErr := json.NewEncoder(w).Encode(rs); wErr != nil {
				log.Error(wErr)
			}
		}

	}
}

// RestHandlerFunc handles REST requests
type RestHandlerFunc func(rq *http.Request) (interface{}, error)

// JSONHandler is a handler for JSON response content type handlers
func JSONHandler(f RestHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, rq *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		rs, err := f(rq)
		if err != nil {
			if se, ok := err.(*StatusErr); ok {
				w.WriteHeader(se.statusCode)
			}

			rs = map[string]string{"error": err.Error()}
			if st, ok := err.(stackTracer); ok {
				rs.(map[string]string)["stacktrace"] = fmt.Sprintf("%+v", st.StackTrace())
			}

		} else if rs == nil {
			rs = map[string]string{}
		}

		if wErr := json.NewEncoder(w).Encode(rs); wErr != nil {
			log.Error(wErr)
		}
	}
}

//ReadJSON reads and unmarshals JSON
func ReadJSON(r io.Reader, s interface{}) error {
	return json.NewDecoder(r).Decode(s)
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}
