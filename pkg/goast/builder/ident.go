package builder

import (
	"cuelang.org/go/cue"
	"go/ast"
)

type Ident struct {
	name string
	ptr  bool
}

func NewIdent(name string) Ident {
	return Ident{
		name: name,
	}
}

func (i Ident) WithPtr(ptr bool) Ident {
	i.ptr = ptr
	return i
}

func (i Ident) AsExpr() ast.Expr {
	if i.ptr {
		return &ast.StarExpr{
			X: ast.NewIdent(i.name),
		}
	}

	return ast.NewIdent(i.name)
}

func FromKind(kind cue.Kind) Ident {
	return NewIdent(kind.String())
}
