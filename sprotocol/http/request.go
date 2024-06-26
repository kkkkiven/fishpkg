package http

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httputil"

	"github.com/pkg/errors"
)

//go:generate msgp

type HTTPRequest struct {
	RemoteAddr string
	Raw        []byte
}

func (this *HTTPRequest) ReadRequest() (*http.Request, error) {
	m := bytes.NewBuffer(this.Raw)
	r, err := http.ReadRequest(bufio.NewReader(m))
	if err != nil {
		return nil, errors.WithMessage(err, "read request failed")
	}

	r.RemoteAddr = this.RemoteAddr
	return r, nil
}

func (this *HTTPRequest) DumpRequest(r *http.Request) error {
	r.URL.Scheme = "http"
	r.URL.Host = r.Host

	raw, err := httputil.DumpRequest(r, true)
	if err != nil {
		return errors.WithMessage(err, "dump request failed")
	}

	this.Raw = raw
	this.RemoteAddr = r.RemoteAddr
	return nil
}
