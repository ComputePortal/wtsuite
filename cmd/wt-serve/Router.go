package main

import (
  "log"
  "net/http"
  "os"
  "path/filepath"

  "github.com/computeportal/wtsuite/pkg/files"
)

type Router struct {
  logger *log.Logger
  content *Tree
}

func NewRouter(root string) (*Router, error) {
  notFoundPath := filepath.Join(root, "404.html")
  if !files.IsFile(notFoundPath) {
    notFoundPath = ""
  }

  content, err := NewTree(
    root, 
    []string{"index.html"},
    map[string]string{
      ".bin": "application/octet-stream",
      ".css": "text/css",
      ".gif": "image/gif",
      ".html": "text/html",
      ".js": "application/javascript",
      ".json": "application/json",
      ".png": "image/png",
      ".svg": "image/svg+xml",
      ".txt": "text/plain",
      ".wasm": "application/wasm",
      ".woff2": "font/woff2",
    },
    notFoundPath,
  )

  if err != nil {
    return nil, err
  }

  return &Router{log.New(os.Stdout, "", log.Ltime), content}, nil
}

func (r *Router) ServeHTTP(resp_ http.ResponseWriter, req *http.Request) {
	// wrap resp_ so we have accecss to the returned status, size, etc
	resp := NewResponseWriter(resp_)

	if req.Method != "GET" {
		resp.WriteError("Error: not a GET request")
  } else {
    if err := r.content.Serve(resp, req); err != nil {
      r.LogError(err)
    }
  }

	r.LogAccess(resp, req)
}

func (r *Router) LogError(err error) {
  r.logger.Printf("Error: %s\n", err.Error())
}

func (r *Router) LogAccess(resp *ResponseWriter, req *http.Request) {
  r.logger.Printf("%s: %s (%d)\n", req.Method, req.URL.Path, resp.Status())
}
