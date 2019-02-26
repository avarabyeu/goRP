package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

type StatusErr struct {
	error
	statusCode int
}

func (se *StatusErr) StackTrace() errors.StackTrace {
	if st, ok := se.error.(stackTracer); ok {
		return st.StackTrace()
	}
	return nil
}

func NewStatusErr(statusCode int, error error) *StatusErr {
	return &StatusErr{statusCode: statusCode, error: error}
}

type HTTPHandlerFunc func(w http.ResponseWriter, rq *http.Request) error

func HTTPHandler(f HTTPHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, rq *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		err := f(w, rq)
		if nil != err {
			var rs interface{}

			if se, ok := err.(*StatusErr); ok {
				w.WriteHeader(se.statusCode)
			}

			rs = map[string]string{"error": err.Error()}
			if st, ok := err.(stackTracer); ok {
				rs.(map[string]string)["stacktrace"] = fmt.Sprintf("%+v", st.StackTrace())
			}

			if wErr := json.NewEncoder(w).Encode(rs); nil != wErr {
				log.Error(wErr)
			}
		}

	}
}

type RestHandlerFunc func(rq *http.Request) (interface{}, error)

func JSONHandler(f RestHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, rq *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		rs, err := f(rq)
		if nil != err {
			if se, ok := err.(*StatusErr); ok {
				w.WriteHeader(se.statusCode)
			}

			rs = map[string]string{"error": err.Error()}
			if st, ok := err.(stackTracer); ok {
				rs.(map[string]string)["stacktrace"] = fmt.Sprintf("%+v", st.StackTrace())
			}

		} else if nil == rs {
			rs = map[string]string{}
		}

		if wErr := json.NewEncoder(w).Encode(rs); nil != wErr {
			log.Error(wErr)
		}
	}
}

func ReadJSON(io io.Reader, s interface{}) error {
	bytes, err := ioutil.ReadAll(io)
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return json.Unmarshal(bytes, s)
	//return json.NewDecoder(bytes).Decode(s)
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}
