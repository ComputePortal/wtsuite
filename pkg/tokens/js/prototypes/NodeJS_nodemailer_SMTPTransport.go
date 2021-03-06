package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type NodeJS_nodemailer_SMTPTransport struct {
  BuiltinPrototype
}

func NewNodeJS_nodemailer_SMTPTransportPrototype() values.Prototype {
  return &NodeJS_nodemailer_SMTPTransport{newBuiltinPrototype("SMTPTransport")}
}

func NewNodeJS_nodemailer_SMTPTransport(ctx context.Context) values.Value {
  return values.NewInstance(NewNodeJS_nodemailer_SMTPTransportPrototype(), ctx)
}

func (p *NodeJS_nodemailer_SMTPTransport) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*NodeJS_nodemailer_SMTPTransport); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *NodeJS_nodemailer_SMTPTransport) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)

  switch key {
  case "sendMail":
    opt := NewConfigObject(map[string]values.Value{
      "from": s,
      "to": s,
      "subject": s,
      "text": s,
      "html": s,
    }, ctx)

    return values.NewFunction([]values.Value{opt, NewVoidPromise(ctx)}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *NodeJS_nodemailer_SMTPTransport) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewNodeJS_nodemailer_SMTPTransportPrototype(), ctx), nil
}
