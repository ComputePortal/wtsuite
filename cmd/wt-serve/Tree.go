package main

import (
  "errors"
  "net/http"
  "path/filepath"
  "os"
  "strings"

  "github.com/computeportal/wtsuite/pkg/files"
)

type Tree struct {
  root     string
  indexNames []string
  mimeTypes map[string]string
	files    map[string]Resource // url to resource map
	notFound *HTML
}

func NewTree(root string, indexNames []string, mimeTypes map[string]string,
	notFoundPath string) (*Tree, error) {
	files := make(map[string]Resource)

	reuseMap := make(map[string]Resource) // path -> resource, so we can reuse resource for different urls

	rootLen := len(root)

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.New("unable to walk file tree at root \""+root+"\"")
		}

		var url string
		if path == root {
			url = "/"
		} else {
			// abbreviate the path by cutting off the root
			url = path[rootLen:]
		}

		resource, err := pathToResource(path, info.IsDir(), indexNames, mimeTypes, reuseMap)
		if err != nil {
			return err
		}

		if resource != nil {
			files[url] = resource

			if info.IsDir() {
				if strings.HasSuffix(url, "/") {
					files[strings.TrimSuffix(url, "/")] = resource
				} else {
					files[url+"/"] = resource
				}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	notFoundResource, err := pathToResource(notFoundPath, false, indexNames, mimeTypes, reuseMap)
	if err != nil {
		return nil, err
	}

	if notFoundResource == nil {
    // create a default
    notFoundResource = DefaultNotFoundHTML(notFoundPath)
	}

	notFoundFile, ok := notFoundResource.(*HTML)
	if !ok {
		return nil, errors.New("404 file \""+notFoundPath+"\" isn't an html resource")
	}

  // remember root/indexNames/mimeTypes if files are added later
	return &Tree{
    root, 
    indexNames,
    mimeTypes,
    files, 
    notFoundFile,
  }, nil
}

// returns nil if no (valid) resource found at path
func pathToResource(path string, isDir bool, indexNames []string, mimeTypes map[string]string, reuseMap map[string]Resource) (Resource, error) {
	resourcePath := path
	if isDir {
		indexFound := false
		// look for index, first one found wins
		for _, testIndex := range indexNames {
			testPath := filepath.Join(path, testIndex)
			if files.IsFile(testPath) {
				resourcePath = testPath
				indexFound = true
				break
			}
		}

		if !indexFound {
			return nil, nil
		}
	}

	if reuseMap != nil {
    if prev, ok := reuseMap[resourcePath]; ok {
      // assume it already passed the mimetype test
      return prev, nil
    }
	} 

  // test mime types
  ext := filepath.Ext(resourcePath)

  if mimeType, ok := mimeTypes[ext]; ok {
    var resource Resource

    if mimeType == "text/html" {
      var err error
      resource, err = NewHTML(resourcePath)
      if err != nil {
        return nil, err
      }
    } else {
      var err error
      resource, err = NewFile(resourcePath, mimeType)
      if err != nil {
        return nil, err
      }
    }

    if reuseMap != nil {
      reuseMap[resourcePath] = resource
    }

    return resource, nil
  }

  return nil, nil
}

// all errors result in false
func (t *Tree) updateIfFile(url string) (Resource, bool) {
  path := filepath.Join(t.root, url)

  info, err := os.Stat(path)
  if err != nil {
    return nil, false
  }

  resource, err := pathToResource(path, info.IsDir(), t.indexNames, t.mimeTypes, nil)
  if resource == nil || err != nil {
    return nil, false
  }

  t.files[url] = resource

  if info.IsDir() {
    if strings.HasSuffix(url, "/") {
      t.files[strings.TrimSuffix(url, "/")] = resource
    } else {
      t.files[url+"/"] = resource
    }
  }

  return resource, true
}

func (t *Tree) Serve(resp *ResponseWriter, req *http.Request) error {
	if resource, ok := t.files[req.URL.Path]; ok {
		return resource.Serve(resp, req)
	} else {
    if resource, ok := t.updateIfFile(req.URL.Path); ok {
      return resource.Serve(resp, req)
    } else {
      return t.notFound.ServeStatus(resp, req, http.StatusNotFound)
    }
	}
}
