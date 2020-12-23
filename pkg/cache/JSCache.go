package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/computeportal/wtsuite/pkg/files"
)

type JSCacheEntry struct {
	Deps    []string // exported
	touched bool     // if still false when saving -> delete this entry
}

type JSCache struct {
	// abs filepath -> details
	age  time.Time // == time.Unix(0, 0) if file doesnt yet exist
	data map[string]JSCacheEntry
}

func LoadJSCache(targetFile string, forceBuild bool) {
	src := cacheFile(targetFile)

	age := time.Unix(0, 0) // very old dummy time
	data := make(map[string]JSCacheEntry)

	if !forceBuild {
		if files.IsFile(src) {
			b, err := ioutil.ReadFile(src)
			if err == nil {
				buf := bytes.NewBuffer(b)
				decoder := gob.NewDecoder(buf)

				statAge, statErr := lastModified(targetFile)

				decodeErr := decoder.Decode(&data)
				if decodeErr != nil || statErr != nil {
					data = make(map[string]JSCacheEntry)
				} else {
					age = statAge
				}
			}
		} else if files.IsDir(src) {
			fmt.Fprintf(os.Stderr, "Error: cache file is a directory, this shouldn't be possible")
			os.Exit(1)
		}
	}

	_cache = &JSCache{age, data}

	files.StartCacheUpdate = _cache.StartUpdate
	files.AddCacheDependency = _cache.AddDependency
}

// also use this for adding files
func (c *JSCache) StartUpdate(fname string) {
	c.data[fname] = JSCacheEntry{make([]string, 0), true}
}

func (c *JSCache) AddDependency(fname string, dep string) {
	entry, ok := c.data[fname]
	if !ok {
		panic(fname + " not found in JSCache")
	}

	// only append if not found
	for _, v := range entry.Deps {
		if v == dep {
			return
		}
	}

	entry.Deps = append(entry.Deps, dep)

	c.data[fname] = entry
}

func (c *JSCache) HasUpstreamDependency(thisPath string, upstreamPath string) bool {
	entry, ok := c.data[thisPath]
	if !ok {
		panic(thisPath + " not found in HTMLCache")
	}

	// only append if not found
	for _, v := range entry.Deps {
		if v == upstreamPath || c.HasUpstreamDependency(v, upstreamPath) {
			return true
		}
	}

	return false
}

func (c *JSCache) RequiresUpdate(fname string) bool {
	entry, ok := c.data[fname]
	if !ok {
		return true
	}

	// if err -> verbose fileError will be triggered later
	if t, err := lastModified(fname); err != nil || t.After(c.age) {
		return true
	}

	for _, dep := range entry.Deps {
		if c.RequiresUpdate(dep) {
			return true
		}
	}

	return false
}

func (c *JSCache) Save() []byte {
	// delete all untouched data entries
	for k, v := range c.data {
		if !v.touched {
			delete(c.data, k)
		}
	}

	buf := bytes.Buffer{}

	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(c.data)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: "+err.Error())
	}

	return buf.Bytes()
}
