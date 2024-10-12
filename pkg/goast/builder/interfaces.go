package builder

import "go/ast"

type Declarer interface {
	AsDecl() ast.Decl
}

type Expresser interface {
	AsExpr() ast.Expr
}

type Fielder interface {
	AsField() *ast.Field
}

type Specifier interface {
	AsSpec() ast.Spec
}
