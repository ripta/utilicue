package builder

import "go/ast"

type File struct {
	name  string
	decls []ast.Decl
}

func NewFile(name string) File {
	return File{
		name: name,
	}
}

func (f File) AddDecl(d Declarer) File {
	f.decls = append(f.decls, d.AsDecl())
	return f
}

func (f File) AsFileNode() *ast.File {
	return &ast.File{
		Name:  ast.NewIdent(f.name),
		Decls: f.decls,
	}
}
