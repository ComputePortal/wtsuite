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

type HTMLCacheEntry struct {
	Deps         []string
	Control      string
	touched      bool
	lastModified time.Time
}

type HTMLCache struct {
	GitCommit    string                    // if this changes -> rebuild all
	CSSBundleURL string                    // if this changes -> rebuild all
	JSBundleURL  string                    // if this changes -> rebuild all
	PxPerRem     int                       // if this changes -> rebuild all
	Compact      bool                      // if this changes -> rebuild all
	GlobalVars   map[string]string         // if any of this changes -> rebuild all
	IndexMap     map[string]string         // abspath -> target abspath
	Data         map[string]HTMLCacheEntry // abs src path as key
}

func globalVarsNotEqual(gv1, gv2 map[string]string) bool {
	for k, v1 := range gv1 {
		if v2, ok := gv2[k]; !ok || v1 != v2 {
			return true
		}
	}

	for k, v2 := range gv2 {
		if v1, ok := gv1[k]; !ok || v1 != v2 {
			return true
		}
	}

	return false
}

func (c *HTMLCache) invalidateViews(indexMap map[string]string, viewControls map[string]string) {
	toDelete := make([]string, 0)

	// keep indexMap clean by removing any that dont exist anymore
	// also remove any that dont have the same target
	for k, oldTarget := range c.IndexMap {
		if newTarget, ok := indexMap[k]; !ok || oldTarget != newTarget {
			toDelete = append(toDelete, k)
		}
	}

	// also remove those that don't have the same controls
	for k, _ := range c.IndexMap {
		entry := c.Data[k]
		if _, ok := viewControls[k]; !ok {
			if entry.Control != "" {
				if VERBOSITY >= 2 {
					fmt.Printf("Info: invalidating cache for %s (old: %s, new: no control)\n", k, entry.Control)
				}

				toDelete = append(toDelete, k)
				continue
			}
		}

		if entry.Control != viewControls[k] {
			toDelete = append(toDelete, k)
			if VERBOSITY >= 2 {
				if entry.Control == "" {
					fmt.Printf("Info: invalidating cache for %s (old: no control, new: %s)\n", k, viewControls[k])
					//panic("block")
				} else {
					fmt.Printf("Info: invalidating cache for %s (old: %s, new: %s)\n", k, entry.Control, viewControls[k])
				}
			}

			continue
		}
	}

	for _, toD := range toDelete {
		if _, ok := c.Data[toD]; ok {
			delete(c.Data, toD)
		}
	}
}

// age of the indexMap doesnt matter
func LoadHTMLCache(indexMap map[string]string,
	viewControls map[string]string,
	cssBundleURL string,
	jsBundleURL string,
	pxPerRem int,
	outputDir string,
	gitCommit string,
	compact bool,
	globalVars map[string]string,
	forceBuild bool) {
	src := cacheFile(outputDir + " html") // assume abspath

	c := &HTMLCache{
		gitCommit,
		cssBundleURL,
		jsBundleURL,
		pxPerRem,
		compact,
		make(map[string]string),
		make(map[string]string),
		make(map[string]HTMLCacheEntry),
	}

	if !forceBuild {
		if files.IsFile(src) {
			b, err := ioutil.ReadFile(src)
			if err == nil {
				buf := bytes.NewBuffer(b)
				decoder := gob.NewDecoder(buf)

				decodeErr := decoder.Decode(c)
				if decodeErr != nil ||
					c.GitCommit != gitCommit ||
					c.CSSBundleURL != cssBundleURL ||
					c.JSBundleURL != jsBundleURL ||
					c.PxPerRem != pxPerRem ||
					c.Compact != compact ||
					globalVarsNotEqual(c.GlobalVars, globalVars) {

					if decodeErr != nil {
						fmt.Fprintf(os.Stderr, "Warning: resetting view cache due to decode error (%s)\n", decodeErr.Error())
					} else if c.GitCommit != gitCommit {
						fmt.Fprintf(os.Stderr, "Warning: resetting view cache due to changed git commit (old: %s, new: %s)\n", c.GitCommit, gitCommit)
					} else if c.CSSBundleURL != cssBundleURL {
						fmt.Fprintf(os.Stderr, "Warning: resetting view cache due to changed css url (old: %s, new: %s)\n", c.CSSBundleURL, cssBundleURL)
					} else if c.JSBundleURL != jsBundleURL {
						fmt.Fprintf(os.Stderr, "Warning: resetting view cache due to changed js url (old: %s, new: %s)\n", c.JSBundleURL, jsBundleURL)
					} else if c.PxPerRem != pxPerRem {
						fmt.Fprintf(os.Stderr, "Warning: resetting view cache due to changed px/rem (old: %d, new: %d)\n", c.PxPerRem, pxPerRem)
					} else if c.Compact != compact {
						fmt.Fprintf(os.Stderr, "Warning: resetting view cache due to changed compact state output (old: %t, new: %t)\n", c.Compact, compact)
					} else if globalVarsNotEqual(c.GlobalVars, globalVars) {
						fmt.Fprint(os.Stderr, "Warning: resetting view cache due to changed compact state output (old: ", c.GlobalVars, ", new: ", globalVars, ")\n")
					}

					// reset
					c = &HTMLCache{
						gitCommit,
						cssBundleURL,
						jsBundleURL,
						pxPerRem,
						compact,
						globalVars,
						indexMap,
						make(map[string]HTMLCacheEntry),
					}
				} else {
					// remove any views that are no longer used,
					//  or no longer have the same target,
					//  or no longer have the same controls
					c.invalidateViews(indexMap, viewControls)
					c.IndexMap = indexMap
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
	files.HasUpstreamCacheDependency = _cache.HasUpstreamDependency
}

// for imports only, dont modify the css rules, because they are not needed
func (c *HTMLCache) StartUpdate(fname string) {
	//t := time.Time{}
	//if !c.RequiresUpdate(fname) {
	//t := c.Data[fname].lastModified
	//fmt.Println("starting partial update of ", fname)
	//} else {
	//fmt.Println("starting update of ", fname)
	//}

	c.Data[fname] = HTMLCacheEntry{make([]string, 0), "", true, time.Time{}}
}

func (c *HTMLCache) StartRootUpdate(fname string) {
	fmt.Println("starting root update of ", fname)
	c.Data[fname] = HTMLCacheEntry{make([]string, 0), "", true, time.Time{}}
}

func StartRootUpdate(fname string) {
	c, ok := _cache.(*HTMLCache)
	if !ok {
		panic("unexpected")
	}

	c.StartRootUpdate(fname)
}

func RollbackUpdate(fname string) {
	// simply remove the c.Data entry
	c, ok := _cache.(*HTMLCache)
	if !ok {
		panic("unexpected")
	}

	delete(c.Data, fname)
}

func (c *HTMLCache) AddDependency(fname string, dep string) {
	entry, ok := c.Data[fname]
	if !ok {
    c.StartUpdate(fname)
    entry = c.Data[fname]
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

func (c *HTMLCache) HasUpstreamDependency(thisPath string, upstreamPath string) bool {
	entry, ok := c.Data[thisPath]
	if !ok {
    c.StartUpdate(thisPath)
    entry = c.Data[thisPath]
	}


	// only append if not found
	for _, v := range entry.Deps {
		if v == upstreamPath || c.HasUpstreamDependency(v, upstreamPath) {
			return true
		}
	}

	return false
}

func SetViewControl(fname string, control string) {
	c, ok := _cache.(*HTMLCache)
	if !ok {
		panic("unexpected")
	}

	entry, ok := c.Data[fname]
	if !ok {
    c.StartUpdate(fname)
    entry = c.Data[fname]
	}

	entry.Control = control

	c.Data[fname] = entry
}

func (c *HTMLCache) requiresUpdate(fname string, age time.Time, m map[string]bool) bool {
	if prevVal, ok := m[fname]; ok {
		return prevVal
	} else {
		b := false
		entry, ok := c.Data[fname]
		if !ok {
			if VERBOSITY >= 3 {
				fmt.Printf("Info: %s requires update because it isn't available in the cache\n", fname)
			}
			b = true
		} else if t := c.Data[fname].lastModified; t.Equal(time.Time{}) || t.After(age) {
			if VERBOSITY >= 3 {
				if t.Equal(time.Time{}) {
					fmt.Printf("Info: %s requires update because its modification time is unknown\n", fname)
				} else {
					fmt.Printf("Info: %s requires update because its modification time is after the target modification time (%s > %s)\n", fname, t.Format(time.UnixDate), age.Format(time.UnixDate))
				}
			}
			b = true
		}

		m[fname] = true // if recursion by accident, then always true

		if !b {
			for _, dep := range entry.Deps {
				if c.requiresUpdate(dep, age, m) {
					b = true
					if VERBOSITY >= 3 {
						fmt.Printf("Info: %s requires update due to recursion (see %s)\n", fname, dep)
					}
					break
				}
			}
		}

		m[fname] = b

		return b
	}
}

func (c *HTMLCache) SyncLastModified() {
	for k, d := range c.Data {
		if t, err := lastModified(k); err == nil {
			if t.Equal(time.Time{}) {
				panic("last modified time can't be empty")
			}

			d.lastModified = t
			c.Data[k] = d

			if c.Data[k].lastModified.Equal(time.Time{}) {
				panic("set failed")
			}
		}
	}
}

func SyncHTMLLastModifiedTimes() {
	c, ok := _cache.(*HTMLCache)
	if !ok {
		panic("unexpected")
	}

	c.SyncLastModified()
}

// recursion must be on other function
func (c *HTMLCache) RequiresUpdate(fname string) bool {
	targetAge, err := lastModified(c.IndexMap[fname])

	// if err -> verbose fileError will be triggered later
	if err != nil {
		return true
	}

	// reuse results, otherwise we'll end up doing a lot of iterations
	m := make(map[string]bool)

	return c.requiresUpdate(fname, targetAge, m)
}

func (c *HTMLCache) touchUpwards(fname string) {
	entry, ok := c.Data[fname]
	if !ok {
		panic("all deps should be in cache data")
	}

  if !entry.touched { // avoid infinite loop in circular situations
    entry.touched = true

    c.Data[fname] = entry

    for _, dep := range entry.Deps {
      c.touchUpwards(dep)
    }
  }
}

func (c *HTMLCache) clean() {
	// touch all deps, moving upward from untouched index files
	// then remove any remaining untouched

	for k, _ := range c.IndexMap {
		entry, ok := c.Data[k]
		if !ok {
			// this is possible  if it just had an error and rolled back the update
			// (cache file doesnt need to be perfectly clean after an error anyway)
			delete(c.IndexMap, k)
			continue
		}

		if !entry.touched {
			c.touchUpwards(k)
		}
	}

	for k, v := range c.Data {
		if !v.touched {
			delete(c.Data, k)
		}
	}
}

func (c *HTMLCache) Save() []byte {
	c.clean()

	buf := bytes.Buffer{}

	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: "+err.Error())
	}

	return buf.Bytes()
}

func SaveHTMLCache(targetFile string) {
	SaveCache(targetFile + " html")
}
