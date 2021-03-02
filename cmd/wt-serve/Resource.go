package main

import (
  "net/http"
)

type Resource interface {
	Serve(resp *ResponseWriter, req *http.Request) error
}
