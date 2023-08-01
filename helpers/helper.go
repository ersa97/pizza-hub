package helpers

import (
	"net/http"
)

func MethodHelper(r *http.Request, method string) bool {
	return r.Method == method
}
