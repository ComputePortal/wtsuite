package glsl

type PrecisionType int

const (
  DEFAULTP PrecisionType = iota
  LOWP 
  MEDIUMP
  HIGHP
)

func PrecisionTypeToString(ptype PrecisionType) string {
  switch ptype {
  case DEFAULTP:
    return ""
  case LOWP:
    return "lowp"
  case MEDIUMP:
    return "mediump"
  case HIGHP:
    return "highp"
  default:
    panic("unhandled")
  }
}
