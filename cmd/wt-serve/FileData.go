package main

import (
  "bytes"
  //"compress/gzip"
  "crypto/md5"
	"encoding/base64"
  "fmt"
  "io"
  "net/http"
  "os"
  "strings"
  "time"
)

var COMPRESS = true

type FileData struct {
  path       string
  modTime    time.Time
  mimeType   string
	eTag       string
	buf        *bytes.Reader
	compressed *bytes.Reader
}

func newFileData(path string, mimeType string) FileData {
  return FileData{path, time.Time{}, mimeType, "", nil, nil}
}

// error also returns false
func (f *FileData) isUpToDate() bool {
  if f.modTime.Equal(time.Time{}) {
    return false
  }

  fInfo, err := os.Stat(f.path)
  if err != nil {
    return false
  }

  if f.modTime.After(fInfo.ModTime()) {
    return true
  } else {
    return false
  }
}

func (f *FileData) grabLatestModTime() {
  fInfo, err := os.Stat(f.path)
  if err == nil {
    f.modTime = fInfo.ModTime()
  }
}

func (f *FileData) cache(b []byte) {
  f.buf = bytes.NewReader(b)

  sum_ := md5.Sum(b)
  sum := make([]byte, 16)
  for i, _ := range sum {
    sum[i] = sum_[i]
  }

  f.eTag = "\"" + base64.StdEncoding.EncodeToString(sum) + "\""

  // XXX: does compression make sense for local development?
  /*if len(b) > 1400 {
    cBytes := &bytes.Buffer{}

    cWriter := gzip.NewWriter(cBytes)
    cWriter.Write(b)
    cWriter.Close()

    f.compressed = bytes.NewReader(cBytes.Bytes())
  }*/
}

func (f *FileData) ServeStatus(resp *ResponseWriter, req *http.Request, status int) error {
	if eTag := req.Header.Get("If-None-Match"); eTag != "" {
		if eTag == f.eTag {
			resp.WriteHeader(http.StatusNotModified)
			fmt.Fprintf(resp, "")
			return nil
		}
	}

	resp.Header().Set("Content-Type", f.mimeType)
	resp.Header().Set("ETag", f.eTag)

	if status != http.StatusOK {
		resp.WriteHeader(status)
	}

	acceptEncoding := req.Header.Get("Accept-Encoding")

	if COMPRESS && strings.Contains(acceptEncoding, "gzip") && f.compressed != nil {
		resp.Header().Set("Content-Encoding", "gzip")
		f.compressed.Seek(0, io.SeekStart)
		io.Copy(resp, f.compressed)
	} else {
		f.buf.Seek(0, io.SeekStart)
		io.Copy(resp, f.buf)
	}

	return nil
}
