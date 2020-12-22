package math

type LevelScope interface {
	Level() float64
	FracDepth() int

	Scale(f float64) LevelScope
	IncrFracDepth() LevelScope
}

type LevelScopeData struct {
	level     float64
	fracDepth int
}

func NewLevelScope() LevelScope {
	return &LevelScopeData{1.0, 0}
}

func (l *LevelScopeData) Level() float64 {
	return l.level
}

func (l *LevelScopeData) FracDepth() int {
	return l.fracDepth
}

func (l *LevelScopeData) Scale(f float64) LevelScope {
	return &LevelScopeData{l.Level() * f, l.FracDepth()}
}

func (l *LevelScopeData) IncrFracDepth() LevelScope {
	return &LevelScopeData{l.Level(), l.FracDepth() + 1}
}
