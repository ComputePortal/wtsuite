package scripts

type SrcScript struct {
	src string
}

func NewSrcScript(src string) (*SrcScript, error) {
	return &SrcScript{src}, nil
}

func (s *SrcScript) Write() (string, error) {
	return "", nil
}

func (s *SrcScript) Dependencies() []string {
	return []string{s.src}
}
