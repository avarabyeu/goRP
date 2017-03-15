package commons

import (
	"net/http"
	"goji.io/middleware"
)

//NoHandlerFound handler function that handles requests that doesn't match any handler
func NoHandlerFound(h func(w http.ResponseWriter, rq *http.Request)) (func(http.Handler) http.Handler) {
	return func(next http.Handler) http.Handler {
		mv := func(w http.ResponseWriter, rq *http.Request) {
			if nil == middleware.Handler(rq.Context()) {
				//not found
				h(w, rq)
			} else {
				next.ServeHTTP(w, rq)
			}

		}
		return http.HandlerFunc(mv)
	}
}
