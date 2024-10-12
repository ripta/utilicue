package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
)

// go run ./cmd/ast2go
func main() {
	fset := token.NewFileSet()
	node := &ast.File{
		Name: ast.NewIdent("db"),
		Decls: []ast.Decl{
			// type Title string
			&ast.GenDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "// Title is a record of a person's title.",
						},
					},
				},
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("Title"),
						Type: ast.NewIdent("string"),
					},
				},
			},
			// type Identity struct {
			// 	First  string
			// 	Middle string
			// 	Last   string
			// 	Nick   *string
			// }
			&ast.GenDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "// Identity is a record of a person's name, including first, middle,",
						},
						{
							Text: "// last, and nick names.",
						},
					},
				},
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("Identity"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("First"),
										},
										Type: &ast.Ident{Name: "string"},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Middle"),
										},
										Type: &ast.Ident{Name: "string"},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Last"),
										},
										Type: &ast.Ident{Name: "string"},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Nick"),
										},
										Type: &ast.StarExpr{
											X: &ast.Ident{Name: "string"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	buf := &bytes.Buffer{}
	buf.WriteString("// Code generated by ast2go. DO NOT EDIT.\n")
	if err := format.Node(buf, fset, node); err != nil {
		panic(err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	fmt.Println(string(formatted))
}
