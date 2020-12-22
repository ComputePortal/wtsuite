package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"../../../pkg/directives"
	"../../../pkg/files"
	"../../../pkg/tree"
)

// all the cmdargs that are needed for Config file reading
type CmdArgs struct {
	ConfigFile  string
	OutputDir   string
	IncludeDirs []string
	GlobalVars  map[string]string

	JsUrl           string
	CssUrl          string
	IncludeViews    []string
	ExcludeViews    []string
	IncludeControls []string
	ExcludeControls []string
	MathFontUrl     string
}

type SearchIndexConfig struct {
	TitleQuery     string `json:"title-query"`
	titleQuery     XQuery
	ContentQueries []string `json:"content-queries"`
	contentQueries []XQuery
	Pages          []string `json:"pages"` // XXX: is "views" a better name?
}

type SearchConfig struct {
	Indices map[string]SearchIndexConfig `json:"indices"`
	Ignore  []string                     `json:"ignore"`
}

type Config struct {
	Views    map[string]map[string]string `json:"views"` // first key is group name
	views    map[string]string
	Controls map[string]map[string][]string `json:"controls"` // input can be multiple paths/groups
	controls map[string][]string            // filtered using includeControls/excludeControls
	Files    map[string]string              `json:"files"` // rel paths in configFile, but abs paths here

	cssDst      string // for stylesheet bundling
	CssUrl      string `json:"css-url"` // for stylesheet link
	jsDst       string
	JsUrl       string `json:"js-url"`
	PxPerRem    int    `json:"px-per-rem"`    // for px(<X rem>) functin
	MathFontUrl string `json:"math-font-url"` // for math font woff2 file
	mathFontDst string

	Search SearchConfig `json:"search"`
}

func NewDefaultCmdArgs() CmdArgs {
	return CmdArgs{
		ConfigFile:  "",
		OutputDir:   "",
		IncludeDirs: make([]string, 0),
		GlobalVars:  make(map[string]string),

		JsUrl:           "",
		CssUrl:          "",
		IncludeViews:    make([]string, 0),
		ExcludeViews:    make([]string, 0),
		IncludeControls: make([]string, 0),
		ExcludeControls: make([]string, 0),
		MathFontUrl:     "",
	}
}

func PrintMessage(msg string) {
	fmt.Fprintf(os.Stderr, "\u001b[1m"+msg+"\u001b[0m\n\n")
}

func relToAbsConfigFileMap(configFname string, outputDir string, relMap map[string]string) (map[string]string, error) {
	absMap := make(map[string]string)

	for k, v := range relMap {
		kAbs, err := files.Search(configFname, k)
		if err != nil {
			return nil, errors.New("Error: bad config file src (" + k + ")")
		}

		vAbs, err := filepath.Abs(filepath.Join(outputDir, v))
		if err != nil {
			return nil, errors.New("Error: bad config file dst (" + v + ")")
		}

		absMap[kAbs] = vAbs
	}

	return absMap, nil
}

func viewIsIncluded(cmdArgs *CmdArgs, groupName string, viewName string) bool {
	if len(cmdArgs.IncludeViews) > 0 {
		ok := false
		for _, incl := range cmdArgs.IncludeViews {
			if incl == groupName || incl == viewName {
				ok = true
				break
			}
		}

		return ok
	} else {
		ok := true
		for _, excl := range cmdArgs.ExcludeViews {
			if excl == groupName || excl == viewName {
				ok = false
				break
			}
		}

		return ok
	}
}

func relToAbsConfigViewMap(cmdArgs *CmdArgs, relMap map[string]map[string]string) (map[string]string, error) {
	absMap := make(map[string]string)
	configFname := cmdArgs.ConfigFile
	outputDir := cmdArgs.OutputDir

	for groupName, group := range relMap {
		for k, v := range group {
			if viewIsIncluded(cmdArgs, groupName, k) {
				kAbs, err := files.Search(configFname, k)
				if err != nil {
					return nil, errors.New("Error: bad config file src (" + k + ")")
				}

				directives.RegisterURL(kAbs, v)

				vAbs, err := filepath.Abs(filepath.Join(outputDir, v))
				if err != nil {
					return nil, errors.New("Error: bad config file dst (" + v + ")")
				}

				absMap[kAbs] = vAbs
			}
		}
	}

	return absMap, nil
}

func controlIsIncluded(cmdArgs *CmdArgs, groupName string, controlName string) bool {
	if len(cmdArgs.IncludeControls) > 0 {
		ok := false
		for _, incl := range cmdArgs.IncludeControls {
			if incl == groupName || incl == controlName {
				ok = true
				break
			}
		}

		return ok
	} else {
		ok := true
		for _, excl := range cmdArgs.ExcludeControls {
			if excl == groupName || excl == controlName {
				ok = false
				break
			}
		}

		return ok
	}
}

func expandViewGlobs(globs []string, cmdArgs *CmdArgs, viewsRelMap map[string]map[string]string) ([]string, bool, error) {
	hasViews := false
	absViews := make([]string, 0) // the result

	// loop the globs to find the views
	for _, glob := range globs {
		for viewGroupName, viewGroupViews := range viewsRelMap {
			for view, _ := range viewGroupViews {
				if viewGroupName == glob || glob == "*" || view == glob {
					hasViews = true
					// add all group views

					if viewIsIncluded(cmdArgs, viewGroupName, view) {
						absView, err := files.Search(cmdArgs.ConfigFile, view)
						if err != nil {
							return nil, false, errors.New("Error: bad config file view src (" + view + ")")
						}

						// only append if unique
						unique := true
						for _, testView := range absViews {
							if testView == absView {
								unique = false
							}
						}
						if unique {
							absViews = append(absViews, absView)
						}
					}
				}
			}
		}
	}

	return absViews, hasViews, nil
}

func relAndGlobToAbsControlsFileMap(cmdArgs *CmdArgs, controlGroups map[string]map[string][]string,
	viewsRelMap map[string]map[string]string) (map[string][]string, error) {
	absMap := make(map[string][]string)
	configFname := cmdArgs.ConfigFile

	for groupName, controlsMap := range controlGroups {
		for k, globs := range controlsMap {
			if controlIsIncluded(cmdArgs, groupName, k) {
				kAbs, err := files.Search(configFname, k)
				if err != nil {
					return nil, errors.New("Error: bad config controls src (" + k + ", relative to " + configFname + ")")
				}

				// absViews might be empty due to exclusion flags, but that might not be a problem
				absViews, hasViews, err := expandViewGlobs(globs, cmdArgs, viewsRelMap)
				if err != nil {
					return nil, err
				}

				if !hasViews {
					PrintMessage("Warning: no views found for '" + k + "'\n")
				}

				if len(absViews) > 0 {
					absMap[kAbs] = absViews
				}
			}
		}
	}

	return absMap, nil
}

func expandSearchViews(cmdArgs *CmdArgs, viewsRelMap map[string]map[string]string, cfg *Config) error {
	// check of max one per view is done later
	for key, indexConfig := range cfg.Search.Indices {
		var hasViews bool = false
		var err error = nil
		indexConfig.Pages, hasViews, err = expandViewGlobs(indexConfig.Pages, cmdArgs, viewsRelMap)
		if err != nil {
			return err
		}

		if !hasViews {
			PrintMessage("Warning: no views found for search index " + key + "'\n'")
		}

		// also compile the queries
		indexConfig.titleQuery, err = ParseXQuery(indexConfig.TitleQuery)
		if err != nil {
			return err
		}

		indexConfig.contentQueries = make([]XQuery, len(indexConfig.ContentQueries))
		for i, querySource := range indexConfig.ContentQueries {
			indexConfig.contentQueries[i], err = ParseXQuery(querySource)
			if err != nil {
				return err
			}
		}

		// save strategy struct back into list
		cfg.Search.Indices[key] = indexConfig
	}

	// sort the ignore words, (just to be sure, needed for BinarySearch)
	sort.Strings(cfg.Search.Ignore)

	return nil
}

func ReadConfigFile(cmdArgs *CmdArgs) (*Config, error) {
	cfg := &Config{
		Views:       make(map[string]map[string]string),
		views:       make(map[string]string),
		Controls:    make(map[string]map[string][]string),
		controls:    make(map[string][]string),
		Files:       make(map[string]string),
		cssDst:      "",
		CssUrl:      "",
		jsDst:       "",
		JsUrl:       "",
		PxPerRem:    0,
		MathFontUrl: "",
		mathFontDst: "",
		Search: SearchConfig{
			Indices: make(map[string]SearchIndexConfig),
			Ignore:  make([]string, 0),
		},
	}

	b, err := ioutil.ReadFile(cmdArgs.ConfigFile)
	if err != nil {
		return cfg, errors.New("Error: problem reading the config file")
	}

	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, errors.New("Error: bad config file syntax (" + err.Error() + ")")
	}

	// need include dirs now because rel paths are converted to abs paths
	files.AppendIncludeDirs(cmdArgs.IncludeDirs)

	cfg.controls, err = relAndGlobToAbsControlsFileMap(cmdArgs, cfg.Controls, cfg.Views)
	if err != nil {
		return cfg, err
	}

	cfg.Files, err = relToAbsConfigFileMap(cmdArgs.ConfigFile, cmdArgs.OutputDir, cfg.Files)
	if err != nil {
		return cfg, err
	}

	cfg.views, err = relToAbsConfigViewMap(cmdArgs, cfg.Views)
	if err != nil {
		return cfg, err
	}

	if err := expandSearchViews(cmdArgs, cfg.Views, cfg); err != nil {
		return cfg, err
	}

	if cmdArgs.CssUrl != "" {
		cfg.CssUrl = cmdArgs.CssUrl
	}

	if cmdArgs.JsUrl != "" {
		cfg.JsUrl = cmdArgs.JsUrl
	}

	if cmdArgs.MathFontUrl != "" {
		cfg.MathFontUrl = cmdArgs.MathFontUrl
	}

	if cfg.CssUrl == "" {
		return cfg, errors.New("Error: empty css-url in config file")
	}

	cfg.cssDst, err = filepath.Abs(filepath.Join(cmdArgs.OutputDir, cfg.CssUrl))
	if err != nil {
		return cfg, errors.New("Error: bad css-dst path (" + cfg.CssUrl + ")")
	}

	if cfg.PxPerRem != 0 {
		if cfg.PxPerRem < 0 {
			return cfg, errors.New("Error: negative px-per-rem not allowed")
		}
	}

	if cfg.JsUrl == "" {
		return cfg, errors.New("Error: empty js-url in config file")
	}

	cfg.jsDst, err = filepath.Abs(filepath.Join(cmdArgs.OutputDir, cfg.JsUrl))
	if err != nil {
		return cfg, errors.New("Error: bad js-dst path (" + cfg.JsUrl + ")")
	}

	if cfg.MathFontUrl != "" {
		cfg.mathFontDst, err = filepath.Abs(filepath.Join(cmdArgs.OutputDir, cfg.MathFontUrl))
		if err != nil {
			return cfg, errors.New("Error: bad math-font-dst path (" + cfg.MathFontUrl + ")")
		}
	}

	return cfg, nil
}

func (c *Config) GetControls() map[string][]string {
	return c.controls
}

func (c *Config) GetViews() map[string]string {
	return c.views
}

func (c *Config) GetJsDst() string {
	return c.jsDst
}

func (c *Config) GetCssDst() string {
	return c.cssDst
}

func (c *Config) GetMathFontDst() string {
	return c.mathFontDst
}

func (s *SearchIndexConfig) TitleMatch(xpath []tree.Tag) bool {
	return s.titleQuery.Match(xpath)
}

func (s *SearchIndexConfig) ContentMatch(xpath []tree.Tag) bool {
	for _, cq := range s.contentQueries {
		if cq.Match(xpath) {
			return true
		}
	}

	return false
}
