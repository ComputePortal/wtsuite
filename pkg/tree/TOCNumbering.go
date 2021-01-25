package tree

import (
  "math"
  "strconv"
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type TOCNumbering struct {
  current int
  counts []int
}

func NewTOCNumbering() *TOCNumbering {
  return &TOCNumbering{0, nil}
}

func (n *TOCNumbering) Next(level int) {
  if n.counts == nil {
    n.current = level
    n.counts = []int{1}
  } else {
    diff := level - n.current 

    if diff > 0 {
      for i := 0; i < diff; i++ {
        n.counts = append(n.counts, 1)
      }
    } else {
      if (diff < 0) {
        if diff <= -len(n.counts) {
          n.counts = n.counts[0:1]
        } else {
          n.counts = n.counts[0: len(n.counts)+diff]
        }
      }

      n.counts[len(n.counts)-1] += 1
    }

    n.current = level
  }
}

func formatRoman(i int) string {
  if i >= 2000 {
    // fallback
    return strconv.Itoa(i)
  }

  var b strings.Builder

  if i >= 1000 {
    b.WriteString("M")
    i -= 1000
  }

  fnGroup := func(j int, one string, five string, ten string) {
    switch j {
    case 1:
      b.WriteString(one)
    case 2:
      b.WriteString(one)
      b.WriteString(one)
    case 3:
      b.WriteString(one)
      b.WriteString(one)
      b.WriteString(one)
    case 4:
      b.WriteString(one)
      b.WriteString(five)
    case 5:
      b.WriteString(five)
    case 6:
      b.WriteString(five)
      b.WriteString(one)
    case 7:
      b.WriteString(five)
      b.WriteString(one)
      b.WriteString(one)
    case 8:
      b.WriteString(five)
      b.WriteString(one)
      b.WriteString(one)
      b.WriteString(one)
    case 9: 
      b.WriteString(one)
      b.WriteString(ten)
    default:
      // dont write anything
    }
  }

  hundreds := int(math.Floor(float64(i)/100.0))

  fnGroup(hundreds, "C", "D", "M")

  i -= hundreds*100

  tens := int(math.Floor(float64(i)/10.0))

  fnGroup(tens, "X", "L", "C")

  i -= tens*10

  fnGroup(i, "I", "V", "X")

  return b.String()
}

func formatInt(scheme string, i int) string {
  switch scheme {
  case "decimal": 
    return strconv.Itoa(i)
  case "roman":
    return formatRoman(i)
  default:
    panic("unhandled")
  }
}

func (n *TOCNumbering) Write(scheme string) string {
  var b strings.Builder

  for _, count := range n.counts {
    b.WriteString(formatInt(scheme, count))
    b.WriteString(".")
  }

  b.WriteString(" ")

  return b.String()
}

func (n *TOCNumbering) CreateTag(scheme string, ctx context.Context) Tag {
  attr := tokens.NewEmptyStringDict(ctx)
  attr.Set("class", tokens.NewValueString("numbering", ctx))
  tag, err := BuildTag("span", attr, ctx)
  if err != nil {
    panic(err)
  }

  textTag := NewText(n.Write(scheme), ctx)

  tag.AppendChild(textTag)

  return tag
}
