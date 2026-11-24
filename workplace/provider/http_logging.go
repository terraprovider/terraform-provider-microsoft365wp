package provider

import (
	"net/http"
	"net/http/httputil"
)

// sensitiveRequestHeaders are redacted in logged request dumps (they carry access tokens).
var sensitiveRequestHeaders = []string{"Authorization"}

// dumpRequestRedacted returns the dump of req like httputil.DumpRequestOut(req, true), but with sensitive headers
// redacted. req itself keeps its headers and (still readable) body.
func dumpRequestRedacted(req *http.Request) ([]byte, error) {
	r := *req
	r.Header = req.Header.Clone()
	for _, h := range sensitiveRequestHeaders {
		if r.Header.Get(h) != "" {
			r.Header.Set(h, "<redacted>")
		}
	}
	dump, err := httputil.DumpRequestOut(&r, true)
	// DumpRequestOut has consumed the body and replaced it (on the copy) with a new reader over the same bytes
	req.Body = r.Body
	return dump, err
}
