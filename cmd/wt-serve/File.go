package main

import (
  "errors"
	"io/ioutil"
	"net/http"
  
  "github.com/computeportal/wtsuite/pkg/files"
)

type File struct {
  FileData
}

// path here is absolute on this filesystem
func NewFile(path string, mimeType string) (*File, error) {
	if !files.IsFile(path) {
		return nil, errors.New("\""+path+"\" not found")
	}

	return &File{newFileData(path, mimeType)}, nil
}

// dont cache the files in debug mode
func (f *File) cache() error {
	if f.buf == nil || !f.FileData.isUpToDate() {
		b, err := ioutil.ReadFile(f.path)
		if err != nil {
			return errors.New("unable to read file \""+f.path+"\" at serve time")
		}

    f.FileData.cache(b)
    f.FileData.grabLatestModTime()
	}

	return nil
}

func (f *File) Serve(resp *ResponseWriter, req *http.Request) error {
	return f.ServeStatus(resp, req, http.StatusOK)
}

func (f *File) ServeStatus(resp *ResponseWriter, req *http.Request, status int) error {
	if err := f.cache(); err != nil {
		return err
	}

  return f.FileData.ServeStatus(resp, req, status)
}
