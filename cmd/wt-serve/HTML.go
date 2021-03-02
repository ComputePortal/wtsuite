package main

import (
  "errors"
	"io/ioutil"
	"net/http"
  
  "github.com/computeportal/wtsuite/pkg/files"
)

type HTML struct {
  tryPath string // for DefaultNotFound, in case an actual NotFound html is found
  FileData
}

func NewHTML(path string) (*HTML, error) {
  if !files.IsFile(path) {
    return nil, errors.New("\""+path+"\" not found")
  }

  return &HTML{"", newFileData(path, "text/html")}, nil
}

func DefaultNotFoundHTML(tryPath string) *HTML {
  h := &HTML{tryPath, newFileData("", "text/html")}

  h.FileData.cache([]byte(`<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>Not Found</title></head><body><h1>Not Found</h1></body></html>"`))

  return h
}

func (h *HTML) cache() error {
  if h.path != "" && (h.buf == nil || !h.FileData.isUpToDate()) {
		b, err := ioutil.ReadFile(h.path)
		if err != nil {
			return errors.New("unable to read file \""+h.path+"\" at serve time")
		}

    h.FileData.cache(b)
    h.FileData.grabLatestModTime()
  } else if h.path == "" && h.tryPath != "" && files.IsFile(h.tryPath) {
    h.path = h.tryPath
    return h.cache()
  }

  return nil
}

func (h *HTML) Serve(resp *ResponseWriter, req *http.Request) error {
  return h.ServeStatus(resp, req, http.StatusOK)
}

func (h *HTML) ServeStatus(resp *ResponseWriter, req *http.Request, status int) error {
	if err := h.cache(); err != nil {
		return err
	}

  return h.FileData.ServeStatus(resp, req, status)
}
