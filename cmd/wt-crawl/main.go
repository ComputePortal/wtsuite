package main

import (
  "fmt"
  "io/ioutil"
  "math"
  "net/http"
  "os"
  "path/filepath"
  "regexp"
  "sort"
  "strings"
  "sync"

	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/parsers"
)

var (
  VERBOSITY = 0
  cmdParser *parsers.CLIParser = nil
  doctypeRE = regexp.MustCompile(`^[\s]*<!(doctype)|(DOCTYPE)[\s]*html[\s]*>`)
  hrefRE = regexp.MustCompile(`<a.*href=['"](.*?)['"]`)
)

type CmdArgs struct {
  rootURL string
  dstDir string
  nCores int

  verbosity int
}

type CrawlStatePage struct {
  url string
  done bool
  pending bool
}

func (c *CrawlStatePage) SetDone() {
  c.done = true
  c.pending = false
}

func (c *CrawlStatePage) IsDone() bool {
  return c.done
}

func (c *CrawlStatePage) IsPending() bool {
  return c.pending
}

func NewCrawlStatePage(url string) *CrawlStatePage {
  return &CrawlStatePage{url, false, true}
}

type CrawlState struct {
  pages map[string]*CrawlStatePage
  lock *sync.RWMutex
}

func NewCrawlState() *CrawlState {
  return &CrawlState{make(map[string]*CrawlStatePage), &sync.RWMutex{}}
}

func (cs *CrawlState) HasPending() bool {
  for _, p := range cs.pages {
    if p.IsPending() {
      return true
    }
  }

  return false
}

func (cs *CrawlState) GetAllPending() []string {
  res := []string{}

  for url, p := range cs.pages {
    if p.IsPending() {
      res = append(res, url)
    }
  }

  sort.Strings(res)

  return res
}

func (cs *CrawlState) AddPending(url string) {
  cs.lock.Lock()

  if _, ok := cs.pages[url]; !ok {
    cs.pages[url] = NewCrawlStatePage(url)
  }

  cs.lock.Unlock()
}

func (cs *CrawlState) SetDone(url string) {
  cs.lock.Lock()

  cs.pages[url].SetDone()

  cs.lock.Unlock()
}

func (cs *CrawlState) Len() int {
  cs.lock.RLock()

  defer cs.lock.RUnlock()

  return len(cs.pages)
}

func (cs *CrawlState) CountDone() int {
  cs.lock.RLock()

  count := 0

  for _, pageState := range cs.pages {
    if pageState.IsDone() {
      count += 1
    }
  }

  cs.lock.RUnlock()

  return count
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
		rootURL: "",
    dstDir: "", // based on rootURL by default
    nCores: 1,

		verbosity: 0,
	}

	var positional []string = nil

  cmdParser = parsers.NewCLIParser(
    fmt.Sprintf("Usage: %s [options] <root>\n", os.Args[0]),
    "<root>   an URL (i.e. starts with <scheme>://)",
    []parsers.CLIOption{
      parsers.NewCLIString("o", "output", "-o, --output <output-dir>   Defaults to ./<host-name>_<uri>", &(cmdArgs.dstDir)),
      parsers.NewCLIInt("j", "", "-j <n-processes>   Defaults to 1", &(cmdArgs.nCores)),
      parsers.NewCLICountFlag("v", "", "Verbosity", &(cmdArgs.verbosity)),
    },
    parsers.NewCLIRemaining(&positional),
  )

  if err := cmdParser.Parse(os.Args[1:]); err != nil {
    printMessageAndExit(err.Error())
  }

	if len(positional) != 1 {
		printMessageAndExit("Error: expected 1 positional argument")
	}

  if files.IsURL(positional[0]) {
    cmdArgs.rootURL = strings.TrimRight(positional[0], "/")
  } else {
		printMessageAndExit("Error: first argument is not an url")
	} 

  if cmdArgs.dstDir == "" {
    // remove the scheme, replace all forward slashes with underscores
    dstDir := strings.Split(cmdArgs.rootURL, "://")[1]
    dstDir = strings.Replace(dstDir, "/", "_", -1)
    dstDir = strings.Replace(dstDir, ".", "_", -1) // so dots dont confuse the ext
    var err error
    cmdArgs.dstDir, err = filepath.Abs(dstDir)
    if err != nil {
      printMessageAndExit("Error: dstDir '" + dstDir + "' " + err.Error())
    }
  }

  if err := os.MkdirAll(cmdArgs.dstDir, 0755); err != nil {
    printMessageAndExit("Error: dstDir '" + cmdArgs.dstDir + "' " + err.Error())
  }

  if cmdArgs.nCores < 1 {
    printMessageAndExit("Error: -j <n-cores> must be larger than 0")
  }

	return cmdArgs
}

func setUpEnv(cmdArgs CmdArgs) error {
	VERBOSITY = cmdArgs.verbosity

  return nil
}

func urlToPath(rootURL string, dstDir string, url string) string {
  path := filepath.Join(dstDir, strings.TrimPrefix(url, rootURL))

  if filepath.Ext(path) != ".html" {
    path = filepath.Join(path, "index.html")
  }

  return path
}

func crawlPage(cmdArgs CmdArgs, url string, rawBytes []byte, cs *CrawlState) error {
  // simply convert to string, and search for `<a.*?href=["']\(.*?\)['"]`

  raw := string(rawBytes)

  // must find <!DOCTYPE html> in raw
  if !doctypeRE.MatchString(raw) {
    return nil
  }

  res := hrefRE.FindAllStringSubmatch(raw, -1)

  hostURL := strings.Join(strings.Split(cmdArgs.rootURL, "/")[0:3], "/")
  urlBase := url
  if filepath.Ext(urlBase) != "" {
    parts := strings.Split(urlBase, "/")
    urlBase = strings.Join(parts[0:len(parts)-1], "/")
  } else {
    urlBase = strings.TrimRight(urlBase, "/")
  }

  for _, match := range res {
    if len(match) == 2 {
      href := match[1]
      if strings.Contains(href, "?") {
        href = strings.Split(href, "?")[0]
      }

      if strings.Contains(href, "#") {
        href = strings.Split(href, "#")[0]
      }

      if href == "." || href == "/" || href == "" || strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") {
        continue
      } else if strings.HasSuffix(href, ".json") || strings.HasSuffix(href, ".txt") || strings.HasSuffix(href, ".docx") || strings.HasSuffix(href, ".pdf") {
        continue
      } else if strings.Contains(href, "/../") {
        // probably something malformed, best ignore for now
        continue
      }

      pendingURL := ""
      if (strings.HasPrefix(href, "https://") || strings.HasPrefix(href, "http://")) && !strings.HasPrefix(href, cmdArgs.rootURL) {
        // other domain, ignore
      } else if strings.HasPrefix(href, cmdArgs.rootURL) {
        pendingURL = href
      } else if strings.HasPrefix(href, "/") {
        pendingURL = hostURL + href
        if !strings.HasPrefix(strings.ToLower(pendingURL), strings.ToLower(cmdArgs.rootURL)) {
          continue
        }
      } else { // assume relative path
        pendingURL = urlBase + "/" + href
      }

      // can't check suffix because could be php or all kinds of other crap
      if pendingURL != "" {
        cs.AddPending(pendingURL)
      }
    }
  }

  return nil
}

func loadAndCrawlPage(cmdArgs CmdArgs, url string, cs *CrawlState) error {
  path := urlToPath(cmdArgs.rootURL, cmdArgs.dstDir, url)

  var rawBytes []byte = nil
  fromDisc := false
  if files.IsFile(path) {
    fromDisc = true
    var err error
    rawBytes, err = ioutil.ReadFile(path)
    if err != nil {
      return err
    }
  } else {
    resp, err := http.Get(url)
    if err != nil {
      return err // probably time-out error
    }

    if resp.StatusCode == 200 {
      if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
        rawBytes, err = ioutil.ReadAll(resp.Body)
        if err != nil {
          return err
        }
      } else {
        return nil // not a html page, but not a problem
      }
    } else {
      return nil // unable to load, but not a problem
    }
  }

  if !fromDisc {
    // save !
    dir := filepath.Dir(path)
    if !files.IsDir(dir) {
      if err := os.MkdirAll(dir, 0755); err != nil {
        return err
      }
    }

    fmt.Fprintf(os.Stdout, "Saving %d/%d (%s)\n", cs.CountDone(), cs.Len(), strings.TrimPrefix(url, cmdArgs.rootURL))

    if err := ioutil.WriteFile(path, rawBytes, 0644); err != nil {
      return err
    }
  }

  if err := crawlPage(cmdArgs, url, rawBytes, cs); err != nil {
    return err
  }

  return nil
}

func splitForProcesses(n int, urls []string) [][]string {
  nPart := int(math.Ceil(float64(len(urls))/float64(n))) // trunc

  res := [][]string{}

  for i := 0; i < len(urls); i += nPart {
    start := i
    stop := int(math.Min(float64(start+nPart), float64(len(urls))))

    res = append(res, urls[start:stop])
  }

  return res
}

func crawl(cmdArgs CmdArgs) error {
  cs := NewCrawlState()
  cs.AddPending(cmdArgs.rootURL)

  gen := 0
  for cs.HasPending() {
    urls := cs.GetAllPending()

    fmt.Fprintf(os.Stdout, "%d: crawling %d pages\n", gen, len(urls))

    urlGroups := splitForProcesses(cmdArgs.nCores, urls)

    var wg sync.WaitGroup

    for _, urlGroup := range urlGroups {
      wg.Add(1)

      go func(urls_ []string) {
        defer wg.Done()

        for _, url := range urls_ {
          if err := loadAndCrawlPage(cmdArgs, url, cs); err != nil {
            printMessageAndExit(err.Error()+"\n")
          }

          cs.SetDone(url)
        }
      }(urlGroup)

      wg.Wait()
    }

    gen += 1
  }

  return nil
}

func main() {
  cmdArgs := parseArgs()

  if err := setUpEnv(cmdArgs); err != nil {
		printMessageAndExit(err.Error()+"\n")
  }

  if err := crawl(cmdArgs); err != nil {
		printMessageAndExit(err.Error()+"\n")
  }
}
