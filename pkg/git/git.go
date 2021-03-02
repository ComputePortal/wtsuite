package git

import (
  "encoding/pem"
  "errors"
  "fmt"
  "path/filepath"
  "regexp"
  "strings"

  "github.com/computeportal/wtsuite/pkg/files"

  gitcore      "gopkg.in/src-d/go-git.v4"
  gitconfig    "gopkg.in/src-d/go-git.v4/config"
  gitplumbing  "gopkg.in/src-d/go-git.v4/plumbing"
  gittransport "gopkg.in/src-d/go-git.v4/plumbing/transport"
  gitssh       "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
  gitmemory    "gopkg.in/src-d/go-git.v4/storage/memory"
  billymemfs   "gopkg.in/src-d/go-billy.v4/memfs"
)

var (
  privateErrRe = regexp.MustCompile(`(auth|ssh)`)
)

func newAuthMethod(sshKey string) gittransport.AuthMethod {
  if sshKey == "" {
    return nil
  }

  // properly pem encode
  privPem, rest := pem.Decode([]byte(sshKey))
  if privPem == nil || len(rest) > 0 {
    fmt.Println(errors.New("bad key stored"))
    return nil
  }

  authMethod, err := gitssh.NewPublicKeys("wtaas", privPem.Bytes, "")
  if err != nil {
    panic(err)
  }

  return authMethod
}

func parseTag(ref_ gitplumbing.ReferenceName) (*files.SemVer, error) {
  ref := string(ref_)

  parts := strings.Split(ref, "/")

  if parts[1] != "tags" {
    return nil, errors.New("unhandled tag format")
  }

  tag := parts[2]
  if strings.HasPrefix(tag, "v") {
    tag = tag[1:]
  }

  return files.ParseSemVer(tag)
}

// returns nil if not found
func loopReferenceNames(url string, sshKey string, cond func(rn gitplumbing.ReferenceName) error) error {
  storer := gitmemory.NewStorage()

  remoteConfig := &gitconfig.RemoteConfig{
    Name: "origin",
    URLs: []string{url},
  }

  if err := remoteConfig.Validate(); err != nil {
    panic(err)
  }

  remote := gitcore.NewRemote(storer, remoteConfig)

  lstOptions := &gitcore.ListOptions{
    Auth: newAuthMethod(sshKey),
  }

  lst, err := remote.List(lstOptions)
  if err != nil {
    return err
  }

  for _, ref := range lst {
    if err := cond(ref.Name()); err != nil {
      return err
    }
  }

  return nil
}

func selectLatestTag(url string, svr *files.SemVerRange, sshKey string) (gitplumbing.ReferenceName, error) {
  found := false
  var bestSemVer *files.SemVer = nil
  var res gitplumbing.ReferenceName

  if err := loopReferenceNames(url, sshKey, func(rn gitplumbing.ReferenceName) error {
    if rn.IsTag() {
      semVer, err := parseTag(rn)
      if err == nil {
        if svr.Contains(semVer) {
          if !found || (found && semVer.After(bestSemVer))  {
            bestSemVer = semVer
            res = rn
            found = true
          } 
        }
      }
    }

    return nil
  }); err != nil {
    return res, err
  }

  if !found {
    return res, errors.New("no valid tags found")
  }

  return res, nil
}

func selectBranch(url string, branch string, sshKey string) (gitplumbing.ReferenceName, error) {
  found := false
  var res gitplumbing.ReferenceName

  fullName := "ref/head/" + branch

  if err := loopReferenceNames(url, sshKey, func(rn gitplumbing.ReferenceName) error {
    if rn.IsBranch() {
      if rn.String() == fullName {
        if found {
          return errors.New("duplicate branch?")
        } else {
          found = true
          res = rn
        }
      }
    }

    return nil
  }); err != nil {
    return res, err
  }

  return res, nil
}

func cloneRef(libURL string, ref gitplumbing.ReferenceName, dst string, sshKey string) error {
  wt := billymemfs.New()

  storer := gitmemory.NewStorage()

  cloneOptions := &gitcore.CloneOptions{
    URL: libURL,
    Auth: newAuthMethod(sshKey),
    ReferenceName: ref,
    SingleBranch: true,
    NoCheckout: true, // checkout follows further down
    RecurseSubmodules: gitcore.NoRecurseSubmodules,
    Progress: nil,
  }

  if err := cloneOptions.Validate(); err != nil {
    return err
  }

  repo, err := gitcore.Clone(storer, wt, cloneOptions)
  if err != nil {
    return err
  }

  worktree, err := repo.Worktree()
  if err != nil {
    return err
  }

  if err := worktree.Checkout(&gitcore.CheckoutOptions{
    Branch: ref,
  }); err != nil {
    return err
  }

  return writeWorktree(wt, dst)
}

func correctURL(url string) string {
  if !strings.HasSuffix(url, ".git") {
    url += ".git"
  }
  
  if !strings.HasSuffix(url, "https://") {
    url = "https://" + url
  }

  return url
}

func FetchRangedTag(libURL string, libDst string, svr *files.SemVerRange, sshKey string) error {
  libURL = correctURL(libURL)

  tagRef, err := selectLatestTag(libURL, svr, sshKey)
  if err != nil {
    return err
  }

  semVer, err := parseTag(tagRef)
  if err != nil {
    return err
  }

  if semVer.Write() == "0.0.0" {
    panic("something went wrong")
  }

  dst := filepath.Join(libDst, semVer.Write())

  if files.IsFile(dst) {
    return errors.New("Error: destination " + dst + " is a file")
  }

  if !files.IsDir(dst) {
    if err := cloneRef(libURL, tagRef, dst, sshKey); err != nil {
      return err
    }
  } // else: assume it is still the same

  return nil
}

// empty sshKey for public
func FetchLatestTag(libURL string, dstPath string, sshKey string) error {
  return FetchRangedTag(libURL, dstPath, files.NewSemVerRange(nil, nil), sshKey)
}

func FetchBranch(repoURL string, branch string, dstPath string, sshKey string) error {
  repoURL = correctURL(repoURL)

  branchRef, err := selectBranch(repoURL, branch, sshKey)
  if err != nil {
    return err
  }

  if files.IsFile(dstPath) {
    return errors.New("Error: destination " + dstPath + " is a file")
  }

  // always fetch, regardless of local state
  if err := cloneRef(repoURL, branchRef, dstPath, sshKey); err != nil {
    return err
  }

  return nil
}

func ForcePush(srcDir string, dstURL string, sshKey string) error {
  dstURL = correctURL(dstURL)

  wt, err := readWorktree(srcDir)
  if err != nil {
    return err
  }

  // first init
  storer := gitmemory.NewStorage()
  repo, err := gitcore.Init(storer, wt)
  if err != nil {
    return err
  }

  // set the remote info
  // assume branch-name is "main"
  cfg, err := repo.Config()
  if err != nil {
    return err
  }

  remoteConfig := &gitconfig.RemoteConfig{
    Name: "origin", 
    URLs: []string{dstURL},
  }

  if err := remoteConfig.Validate(); err != nil {
    panic(err)
  }

  cfg.Remotes["origin"] = remoteConfig

  // TODO: check that remote is actually setable this way (otherwise we first have to clone, and then upload the differences, which is slightly more hassle)

  // first we must clone, and then we can update
  pushOptions := &gitcore.PushOptions{
    RemoteName: "origin",
    Auth: newAuthMethod(sshKey),
  }

  if err := pushOptions.Validate(); err != nil {
    return err
  }

  if err := repo.Push(pushOptions); err != nil {
    return err
  }

  return nil
}

func FetchPublicOrPrivate(url string, svr *files.SemVerRange) (string, error) {
  dstBase := files.PkgInstallDir(url)

  // assume public (i.e. empty sshKey)
  if err := FetchRangedTag(url, dstBase, svr, ""); err != nil {
    if privateErrRe.MatchString(strings.ToLower(err.Error())) {
      dstBase = files.PrivatePkgInstallDir(url)

      // read SSH key from path
      sshKey, err := files.ReadPrivateSSHKey()
      if err != nil {
        return "", err
      }

      if err := FetchRangedTag(url, dstBase, svr, sshKey); err != nil {
        return "", err
      }
    } else {
      return "", err
    }
  }

  return dstBase, nil
}

func RegisterFetchPublicOrPrivate() {
  files.FetchPublicOrPrivate = FetchPublicOrPrivate
}
