package builder

import (
	"go/ast"
	"go/token"
)

type Type struct {
	name    string
	comment string
	expr    ast.Expr
}

func NewType(name string) *Type {
	return &Type{
		name: name,
	}
}

func (t Type) AsDecl() ast.Decl {
	return &ast.GenDecl{
		Doc: FormatDoc(t.comment),
		Tok: token.TYPE,
		Specs: []ast.Spec{
			t.AsSpec(),
		},
	}
}

func (t Type) AsField() *ast.Field {
	if t.expr == nil {
		panic("expression on type builder " + t.name + " not set")
	}

	return &ast.Field{
		Names: []*ast.Ident{
			ast.NewIdent(t.name),
		},
		Type: t.expr,
	}
}

func (t Type) AsSpec() ast.Spec {
	if t.expr == nil {
		panic("expression on type builder " + t.name + " not set")
	}

	return &ast.TypeSpec{
		Name: ast.NewIdent(t.name),
		Type: t.expr,
	}
}

func (t Type) WithComment(comment string) Type {
	if t.comment != "" {
		panic("comment on type builder " + t.name + " already set")
	}

	t.comment = comment
	return t
}

func (t Type) WithExpr(expr Expresser) Type {
	if t.expr != nil {
		panic("expression on type builder " + t.name + " already set")
	}

	t.expr = expr.AsExpr()
	return t
}

func (t Type) WithStarExpr(expr Expresser) Type {
	if t.expr != nil {
		panic("expression on type builder " + t.name + " already set")
	}

	t.expr = &ast.StarExpr{
		X: expr.AsExpr(),
	}
	return t
}
