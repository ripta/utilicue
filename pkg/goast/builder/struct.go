package builder

import (
	"go/ast"
)

type Struct struct {
	fields []*ast.Field
}

func NewStruct() Struct {
	return Struct{}
}

func (s Struct) AddField(f Fielder) Struct {
	s.fields = append(s.fields, f.AsField())
	return s
}

func (s Struct) AsExpr() ast.Expr {
	return &ast.StructType{
		Fields: &ast.FieldList{
			List: s.fields,
		},
	}
}
