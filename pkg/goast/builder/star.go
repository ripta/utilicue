package builder

import "go/ast"

type Star struct {
	X ast.Expr
}

func NewStar(x Expresser) Star {
	return Star{
		X: x.AsExpr(),
	}
}

func (s Star) AsExpr() ast.Expr {
	return &ast.StarExpr{
		X: s.X,
	}
}
