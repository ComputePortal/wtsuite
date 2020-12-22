package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"

	"../files"
)

type FileCacheEntry struct {
}

type FileCache struct {
	TargetMap map[string]string
	Data      map[string]FileCacheEntry
}

func (c *FileCache) invalidateFiles(filesMap map[string]string) {
	toDelete := make([]string, 0)

	// keep indexMap clean by removing any that dont exist anymore
	for k, oldTarget := range c.TargetMap {
		if newTarget, ok := filesMap[k]; !ok || oldTarget != newTarget {
			toDelete = append(toDelete, k)
		}
	}

	for _, toD := range toDelete {
		if _, ok := c.Data[toD]; ok {
			delete(c.Data, toD)
		}
	}
}

func LoadFileCache(filesMap map[string]string, outputDir string, forceBuild bool) {
	src := cacheFile(outputDir + " file")

	c := &FileCache{
		make(map[string]string),
		make(map[string]FileCacheEntry),
	}

	if !forceBuild {
		if files.IsFile(src) {
			b, err := ioutil.ReadFile(src)
			if err == nil {
				buf := bytes.NewBuffer(b)
				decoder := gob.NewDecoder(buf)

				decodeErr := decoder.Decode(c)
				if decodeErr != nil {
					// reset
					c = &FileCache{
						make(map[string]string),
						make(map[string]FileCacheEntry),
					}
				} else {
					// remove any files that are no longer used
					c.invalidateFiles(filesMap)

					c.TargetMap = filesMap
				}
			}
		} else if files.IsDir(src) {
			fmt.Fprintf(os.Stderr, "Error: cache file is directory, this shouldn't be possible")
			os.Exit(1)
		}
	}

	_cache = c

	files.StartCacheUpdate = _cache.StartUpdate
	files.AddCacheDependency = nil
	files.HasUpstreamCacheDependency = nil
}

func (c *FileCache) StartUpdate(fname string) {
	c.Data[fname] = FileCacheEntry{}
}

func (c *FileCache) AddDependency(fname string, dep string) {
	panic("regular files cannot have dependencies")
}

func (c *FileCache) HasUpstreamDependency(thisPath string, upstreamPath string) bool {
	panic("regular files cannot have dependencies")
}

func (c *FileCache) RequiresUpdate(fname string) bool {
	if _, ok := c.Data[fname]; !ok {
		return true
	}

	targetAge, err := lastModified(c.TargetMap[fname])

	if err != nil {
		return true
	}

	// if err -> verbose fileError will be triggered later
	if t, err := lastModified(fname); err != nil || t.After(targetAge) {
		return true
	}

	return false
}

func (c *FileCache) Save() []byte {
	buf := bytes.Buffer{}

	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: "+err.Error())
	}

	return buf.Bytes()
}

func SaveFileCache(targetFile string) {
	SaveCache(targetFile + " file")
}
