package transport

import "net/http"

// InletsHeader is used for internal connection-tracking
const InletsHeader = "x-inlets-id"

// CopyHeaders copies headers from one http.Header to another by value
func CopyHeaders(destination http.Header, source *http.Header) {
	for k, v := range *source {
		vClone := make([]string, len(v))
		copy(vClone, v)
		(destination)[k] = vClone
	}
}
