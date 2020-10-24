package cache

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type Cache interface {
	StartUpdate(fname string)
	AddDependency(fname string, dep string)
	HasUpstreamDependency(thisPath string, upstreamPath string) bool
	RequiresUpdate(fname string) bool
	Save() []byte
}

const (
	HTMLPPCACHEDIR = ".htmlppcache"
)

var (
	VERBOSITY = 0
)

var _cache Cache = nil

func cacheDir() string {
	return filepath.Join(os.Getenv("HOME"), HTMLPPCACHEDIR)
}

func cacheFile(targetFile string) string {
	key := base64.StdEncoding.EncodeToString([]byte(targetFile))

	return filepath.Join(cacheDir(), key)
}

func lastModified(path string) (time.Time, error) {
	status, err := os.Stat(path)
	if err != nil {
		return time.Now(), err
	}

	return status.ModTime(), nil
}

func RequiresUpdate(fname string) bool {
	if !filepath.IsAbs(fname) {
		panic("expected absolute fname")
	}

	// also assume file exists (if it doesn't then error will soon be given)

	return _cache.RequiresUpdate(fname)
}

func SaveCache(targetFile string) {
	dst := cacheFile(targetFile)

	// make sure the dst directory exists
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Internal Error: unable to make .htmlppcache dir in HOME\n")
		os.Exit(1)
	}

	if err := ioutil.WriteFile(dst, _cache.Save(), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Internal Error: when writing cache to  %s\n", dst)
		os.Exit(1)
	}
}
