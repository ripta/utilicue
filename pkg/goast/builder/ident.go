package builder

import "go/ast"

type Ident struct {
	name string
}

func NewIdent(name string) Ident {
	return Ident{
		name: name,
	}
}

func (i Ident) AsExpr() ast.Expr {
	return ast.NewIdent(i.name)
}
