package scripts

type InlineScript struct {
	content string
}

func NewInlineScript(content string) (*InlineScript, error) {
	return &InlineScript{content}, nil
}

func (s *InlineScript) Write() (string, error) {
	return s.content, nil
}

func (s *InlineScript) Dependencies() []string {
	return []string{}
}
