package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"../files"
	"../tokens/js"
)

type ControlCacheEntry struct {
	Deps    []string
	touched bool // if still false when saving -> delete this entry
}

type ControlCache struct {
	age      time.Time
	Compact  bool
	Controls map[string][]js.ViewInterface // Data contains all js files, not just controls
	Data     map[string]ControlCacheEntry
}

func (c *ControlCache) reset() {
	c.Controls = make(map[string][]js.ViewInterface)
	c.Data = make(map[string]ControlCacheEntry)
}

func (c *ControlCache) invalidateControls(controlViews map[string][]string,
	viewInterfaces map[string]*js.ViewInterface) {

	// if the controls differ, then update everything
	if len(controlViews) != len(c.Controls) {
		c.reset()
		return
	}

	for control, origViewInterfs := range c.Controls {
		views, ok := controlViews[control]
		if !ok {
			c.reset()
			return
		}

		for _, view := range views {
			vif, ok := viewInterfaces[view]
			if !ok {
				continue
				//panic("unexpected")
			}

			found := false
			for _, origViewInterf := range origViewInterfs {
				if vif.IsSame(&origViewInterf) {
					found = true
				}
			}

			if !found {
				c.reset()
				return
			}
		}
	}
}

func LoadControlCache(controls map[string][]string, jsDst string, viewInterfaces map[string]*js.ViewInterface,
	compact bool, forceBuild bool) {
	// the cache file names is based on jsDst
	src := cacheFile(jsDst)

	age := time.Unix(0, 0) // very old dummy time)
	c := &ControlCache{
		age,
		compact,
		make(map[string][]js.ViewInterface),
		make(map[string]ControlCacheEntry),
	}

	if !forceBuild {
		if files.IsFile(src) {
			b, err := ioutil.ReadFile(src)
			if err == nil {
				buf := bytes.NewBuffer(b)
				decoder := gob.NewDecoder(buf)
				decodeErr := decoder.Decode(c)

				statAge, statErr := lastModified(jsDst)

				if decodeErr != nil ||
					statErr != nil ||
					c.Compact != compact {
					c = &ControlCache{
						age,
						compact,
						make(map[string][]js.ViewInterface),
						make(map[string]ControlCacheEntry),
					}
				} else {
					// remove everything if:
					//  any controls dont match
					//  any views dont match
					c.invalidateControls(controls, viewInterfaces)

					c.age = statAge
				}
			}
		} else if files.IsDir(src) {
			fmt.Fprintf(os.Stderr, "Error: cache file is a directory, this shouldn't be possible")
			os.Exit(1)
		}
	}

	_cache = c

	files.StartCacheUpdate = _cache.StartUpdate
	files.AddCacheDependency = _cache.AddDependency
}

func (c *ControlCache) StartUpdate(fname string) {
	c.Data[fname] = ControlCacheEntry{make([]string, 0), true}
}

func (c *ControlCache) AddDependency(fname string, dep string) {
	entry, ok := c.Data[fname]
	if !ok {
		panic(fname + " not found ControlCache")
	}

	// only append if not found
	for _, v := range entry.Deps {
		if v == dep {
			return
		}
	}

	entry.Deps = append(entry.Deps, dep)

	c.Data[fname] = entry
}

func (c *ControlCache) HasUpstreamDependency(thisPath string, upstreamPath string) bool {
	entry, ok := c.Data[thisPath]
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

func AddControl(control string, views []string, viewInterfaces map[string]*js.ViewInterface) {
	c, ok := _cache.(*ControlCache)
	if !ok {
		panic("unexpected")
	}

	uniqueViewInterfs := make([]js.ViewInterface, 0)

	for _, view := range views {
		vif, ok := viewInterfaces[view]

		if !ok {
			for testViewName, _ := range viewInterfaces {
				fmt.Println("found: ", testViewName)
			}

			panic("view interface for " + view + " not found")
		}

		if vif == nil {
			panic("view interface cant be nil")
		}

		unique := true
		for _, other := range uniqueViewInterfs {
			if vif.IsSame(&other) {
				unique = false
			}
		}

		if unique {
			uniqueViewInterfs = append(uniqueViewInterfs, *vif)
		}
	}

	c.Controls[control] = uniqueViewInterfs
}

func (c *ControlCache) requiresUpdate(fname string, m map[string]bool) bool {
	if prev, ok := m[fname]; ok {
		return prev
	} else {
		res := false
		entry, ok := c.Data[fname]
		if !ok {
			res = true
		} else if t, err := lastModified(fname); err != nil || t.After(c.age) {
			res = true
		} else {
			for _, dep := range entry.Deps {
				if c.requiresUpdate(dep, m) {
					res = true
					break
				}
			}
		}

		m[fname] = res
		return res
	}
}

func (c *ControlCache) RequiresUpdate(fname string) bool {
	return c.requiresUpdate(fname, make(map[string]bool))
}

func (c *ControlCache) Save() []byte {
	// delete all untouched data entries
	for k, v := range c.Data {
		if !v.touched {
			delete(c.Data, k)
		}
	}

	buf := bytes.Buffer{}

	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: "+err.Error())
	}

	return buf.Bytes()
}
