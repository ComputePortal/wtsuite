// TODO: instead of transpiling the files, just look at the final files
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/computeportal/wtsuite/pkg/cache"
	"github.com/computeportal/wtsuite/pkg/directives"
	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/parsers"
	"github.com/computeportal/wtsuite/pkg/styles"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tree"
	"github.com/computeportal/wtsuite/pkg/tree/scripts"
)

var (
  VERBOSITY = 0
  cmdParser *parsers.CLIParser = nil
)

type CmdArgs struct {
  root string // can be a url
  configFile string // defaults to search.json in pwd
	searchIndexOutput string

	verbosity int
}

type SearchConfig struct {
	TitleQuery     string `json:"title-query"`
	titleQuery     styles.Selector
	ContentQuery   string `json:"content-query"`
	contentQueries []styles.Selector
	Ignore         []string `json:"ignore"`
}

func printMessageAndExit(msg string) {
	fmt.Fprintf(os.Stderr, "\u001b[1m"+msg+"\u001b[0m\n\n")
  os.Exit(1)
}

func printSyntaxErrorAndExit(err error) {
	os.Stderr.WriteString(err.Error())
	os.Exit(1)
}

func parseArgs() CmdArgs {
	// default args
	cmdArgs := CmdArgs{
		root: "",
    configFile: "search-config.json",
		searchIndexOutput: "search-index.json",

		verbosity: 0,
	}

	var positional []string = nil

  cmdParser = parsers.NewCLIParser(
    fmt.Sprintf("Usage: %s [options] <root>\n", os.Args[0]),
    "",
    []parsers.CLIOption{
      parsers.NewCLIString("c", "config", "-c, --config <config-file>   Defaults to ./search-config.json", &(cmdArgs.configFile)),
      parsers.NewCLIString("o", "output", "-o, --output <output-file>   Defaults to ./search-index.json", &(cmdArgs.searchIndexOutput)),
      parsers.NewCLICountFlag("v", "", "Verbosity", &(cmdArgs.verbosity)),
    },
    parsers.NewCLIRemaining(&positional),
  )

  if err := cmdParser.Parse(os.Args[1:]); err != nil {
    printMessageAndExit(err.Error())
  }

	if len(positional) != 1 {
		printMessageAndExit("Error: expected 1 positional arguments")
	}

	if !files.IsDir(positional[0]) {
    // TODO: might be url
		printMessageAndExit("Error: first argument is not a directory")
	} else {
    var err error
    cmdArgs.root, err = filepath.Abs(positional[0])
    if err != nil {
      printMessageAndExit("Error: root '"+positional[0]+"' "+err.Error())
    }
  }

	if err := files.AssertFile(cmdArgs.configFile); err != nil {
		printMessageAndExit("Error: configFile '"+cmdArgs.configFile+"' "+err.Error())
	}

	configFile, err := filepath.Abs(cmdArgs.configFile)
	if err != nil {
		printMessageAndExit("Error: configFile '"+cmdArgs.configFile+"' "+err.Error())
	} else {
		cmdArgs.configFile = configFile
	}

	absSearchIndexOutput, err := filepath.Abs(cmdArgs.searchIndexOutput)
	if err != nil {
		panic(err)
	}

	cmdArgs.searchIndexOutput = absSearchIndexOutput

	return cmdArgs
}

func setUpEnv(cmdArgs CmdArgs, cfg *SearchConfig) error {
	VERBOSITY = cmdArgs.verbosity
	directives.VERBOSITY = cmdArgs.verbosity
	tokens.VERBOSITY = cmdArgs.verbosity
	js.VERBOSITY = cmdArgs.verbosity
	parsers.VERBOSITY = cmdArgs.verbosity
	files.VERBOSITY = cmdArgs.verbosity
	cache.VERBOSITY = cmdArgs.verbosity
	tree.VERBOSITY = cmdArgs.verbosity
	//styles.VERBOSITY = cmdArgs.verbosity
	scripts.VERBOSITY = cmdArgs.verbosity

  return files.ResolvePackages(cmdArgs.configFile)
}

type SearchIndexPage struct {
	Url     string   `json:"url"`     // used as key
	Title   string   `json:"title"`   // should be unique for each indexed page
	Content []string `json:"content"` // each string is a paragraph
}

type SearchIndex struct {
	Pages   []SearchIndexPage      `json:"pages"`   // sorted
	Ignore  map[string]string      `json:"ignore"`  // key is same as value
	Index   map[string]interface{} `json:"index"`   // nested character tree, leaves are indices into pages array
	Partial map[string]interface{} `json:"partial"` // nested character tree which doesn't start at beginning of word
}

func NewSearchIndex() *SearchIndex {
	return &SearchIndex{
		Pages:   make([]SearchIndexPage, 0),
		Ignore:  make(map[string]string),
		Index:   make(map[string]interface{}), // XXX: can be left empty for initial test
		Partial: make(map[string]interface{}), // XXX: can be left empty for initial test
	}
}

func (si *SearchIndex) AddPage(url string, title string, content []string) {
	si.Pages = append(si.Pages, SearchIndexPage{url, title, content})
}

func findRootParagraph(xpath []tree.Tag) tree.Tag {
	for _, t := range xpath {
		if t.Name() == "p" {
			return t
		}
	}

	return nil
}

func extractTagText(tags []tree.Tag) []string {
  str := make([]string, 0)

  for _, t := range tags {
    if err := tree.WalkText(t, []tree.Tag{}, func(_ []tree.Tag, s string) error {
      str = append(str, s)

      return nil
    }); err != nil {
      panic("unexpected")
    }
  }

  return str
}

func parseHTMLFile(cmdArgs CmdArgs, cfg *SearchConfig, path string, si *SearchIndex) error {
  url := path[len(cmdArgs.root):]

  p, err := parsers.NewXMLParser(path)
  if err != nil {
    return err
  }

	rawTags, err := p.BuildTags()
	if err != nil {
		return err
	}

	root := tree.NewRoot(p.NewContext(0, 1))
	node := directives.NewRootNode(root, directives.HTML)
  // the source isn't really used, because the html file doesnt contain import statements
	fileScope := directives.NewFileScope(false, directives.NewFileCache())

	for _, tag := range rawTags {
		if err := directives.BuildTag(fileScope, node, tag); err != nil {
			return err
		}
	}

  tree.RegisterParents(root)
  if err := root.Validate(); err != nil {
    return err
  }

  // now we can apply the search queries
  titleTags := cfg.titleQuery.Match(root)
  if len(titleTags) > 1 {
    return errors.New("Error: multiple titles found in " + path)
  }

  titleParts := extractTagText(titleTags)
  title := strings.Join(titleParts, " ")

  contentTags := make([]tree.Tag, 0)

  for _, sel := range cfg.contentQueries {
    contentTags = append(contentTags, sel.Match(root)...)
  }

  content := extractTagText(contentTags)

  fmt.Println("indexing ", url, "(title=", title, ")")
  si.AddPage(url, title, content)

  return nil
}

func registerSearchableContent(cmdArgs CmdArgs, cfg *SearchConfig) (*SearchIndex, error) {
	searchIndex := NewSearchIndex()

  if err := filepath.Walk(cmdArgs.root, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }

    // only look at html files
    if filepath.Ext(path) != ".html" {
      return nil
    }

    // now read the html file
    if err := parseHTMLFile(cmdArgs, cfg, path, searchIndex); err != nil {
      return err
    }

    return nil
  }); err != nil {
    return nil, err
  }

	return searchIndex, nil
}

func (si *SearchIndex) indexWord(m map[string]interface{}, pageID int, f string) error {
	chars := strings.Split(f, "")

	for _, char := range chars {
		if mInner, ok := m[char]; ok {
			m = mInner.(map[string]interface{})
		} else {
			mInner := make(map[string]interface{})
			m[char] = mInner
			m = mInner
		}
	}

	if pages_, ok := m["pages"]; ok {
		pages := pages_.([]float64)
		unique := true
		for _, page := range pages {
			if int(page) == pageID {
				unique = false
				break
			}
		}

		if unique {
			pages = append(pages, float64(pageID))
			m["pages"] = pages
		}
	} else {
		m["pages"] = []float64{float64(pageID)}
	}

	return nil
}

func (si *SearchIndex) IndexWord(pageID int, f string) error {
	if err := si.indexWord(si.Index, pageID, f); err != nil {
		return err
	}

	// partial versions of the word are also indexed (include the full word itself)
	for i := 0; i < len(f); i++ {
		fPart := f[i:]

		if err := si.indexWord(si.Partial, pageID, fPart); err != nil {
			return err
		}
	}

	return nil
}

func isIgnoredWord(cfg *SearchConfig, w string) bool {
	i := sort.SearchStrings(cfg.Ignore, w)

	if i > -1 && i < len(cfg.Ignore) {
		return cfg.Ignore[i] == w
	} else {
		return false
	}
}

func indexSentence(cfg *SearchConfig, si *SearchIndex, pageID int, sentence string) error {
	fields := strings.FieldsFunc(strings.Trim(sentence, "."), func(r rune) bool {
		return r < 46 || // keep period as decimal separator
			r == 47 || // forward slash
			r == 58 || // :
			r == 59 || // ;
			r == 63 || // ?
			r == 95 // _
	})

	for _, field := range fields {
		f := strings.ToLower(strings.TrimSpace(field))
		if f != "" {
			if !isIgnoredWord(cfg, f) {
				if err := si.IndexWord(pageID, f); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// actually fill the index/partial nested trees
func buildSearchIndex(cfg *SearchConfig, searchIndex *SearchIndex) error {
	// loop each word of each page
	for i, page := range searchIndex.Pages {
		if err := indexSentence(cfg, searchIndex, i, page.Title); err != nil {
			return err
		}
		for _, paragraph := range page.Content {
			if err := indexSentence(cfg, searchIndex, i, paragraph); err != nil {
				return err
			}
		}
	}

	// add the ignored values
	for _, w := range cfg.Ignore {
		searchIndex.Ignore[w] = w
	}

	return nil
}

func ReadConfigFile(cmdArgs *CmdArgs) (*SearchConfig, error) {
  cfg := &SearchConfig{
    TitleQuery: "",
    titleQuery: nil,
    ContentQuery: "",
    contentQueries: []styles.Selector{},
    Ignore:         []string{},
  }

	b, err := ioutil.ReadFile(cmdArgs.configFile)
	if err != nil {
		return cfg, errors.New("Error: problem reading the config file")
	}

	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, errors.New("Error: bad config file syntax (" + err.Error() + ")")
	}

  // parse
  titleQueries, err := styles.ParseSelectorList(tokens.NewValueString(cfg.TitleQuery, context.NewContext(context.NewSource(cfg.TitleQuery), cmdArgs.configFile)))
  if err != nil {
    return nil, err
  }

  if len(titleQueries) != 1 {
    return cfg, errors.New("Error: expected only one title query")
  }

  cfg.titleQuery = titleQueries[0]

  cfg.contentQueries, err = styles.ParseSelectorList(tokens.NewValueString(cfg.ContentQuery, context.NewContext(context.NewSource(cfg.ContentQuery), cmdArgs.configFile)))
  if err != nil {
    return nil, err
  }

  return cfg, nil
}

func main() {
	cmdArgs := parseArgs()

	// age of the configFile doesn't matter
	cfg, err := ReadConfigFile(&(cmdArgs))
	if err != nil {
		printMessageAndExit(err.Error()+"\n")
	}

  if err := setUpEnv(cmdArgs, cfg); err != nil {
		printMessageAndExit(err.Error()+"\n")
  }

	searchIndex, err := registerSearchableContent(cmdArgs, cfg)
	if err != nil {
		printMessageAndExit(err.Error()+"\n")
	}

	if err := buildSearchIndex(cfg, searchIndex); err != nil {
		printMessageAndExit(err.Error()+"\n")
	}

	b, err := json.Marshal(searchIndex)
	if err != nil {
		printMessageAndExit(err.Error()+"\n")
	}

	if err := ioutil.WriteFile(cmdArgs.searchIndexOutput, b, 0644); err != nil {
		printMessageAndExit(err.Error()+"\n")
	}
}
