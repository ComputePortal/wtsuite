package files

type Source interface {
  Search(callerPath string, srcPath string) (string, error)
  Read(path string) ([]byte, error)
}
