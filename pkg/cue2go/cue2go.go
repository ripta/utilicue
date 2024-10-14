package cue2go

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ripta/utilicue/pkg/goast/builder"
	"go/format"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

type Generator struct {
}

func (gen *Generator) Run(args []string) error {
	if len(args) == 0 {
		return errors.New("no input directories or files provided")
	}

	for _, arg := range args {
		files := []string{}

		abs, err := filepath.Abs(arg)
		if err != nil {
			return err
		}

		fi, err := os.Stat(abs)
		if err != nil {
			return err
		}

		if !fi.IsDir() {
			files = append(files, abs)
			continue
		}

		fs.WalkDir(os.DirFS(abs), ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			if filepath.Ext(path) == ".cue" {
				files = append(files, filepath.Join(abs, path))
			}
			return nil
		})

		ctx := cuecontext.New()
		insts := load.Instances(files, nil)
		if len(insts) != 1 {
			return fmt.Errorf("expected exactly one instance, got %d", len(insts))
		}

		v := ctx.BuildInstance(insts[0])
		if err := v.Err(); err != nil {
			log.Fatal(err)
		}

		pkg, err := processTopLevel(builder.NewFile(filepath.Base(abs)), v)
		if err != nil {
			return fmt.Errorf("error generating Go code: %w", err)
		}

		buf := &bytes.Buffer{}

		fmt.Fprintf(buf, "// Code generated by cue2go. DO NOT EDIT.\n")
		fmt.Fprintf(buf, "// Source: %s\n", abs)

		fset := token.NewFileSet()
		if err := format.Node(buf, fset, pkg.AsFileNode()); err != nil {
			return fmt.Errorf("error formatting output: %w", err)
		}

		if err := os.WriteFile(filepath.Join(abs, "generated.cue2go.go"), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("error writing output file: %w", err)
		}

		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			return fmt.Errorf("error formatting output: %w", err)
		}

		if err := os.WriteFile(filepath.Join(abs, "generated.cue2go.go"), formatted, 0644); err != nil {
			return fmt.Errorf("error writing output file: %w", err)
		}
	}

	return nil
}

func processValue(name cue.Selector, val cue.Value) (builder.Type, error) {
	ptr := IsOptional(name)
	kind := val.IncompleteKind()

	// Top values are implicitly optional, but since we already treat them as
	// `any` in Go, we don't want to add another indirection through pointers.
	//
	// However, in any other nullable case, do treat them as pointers, and then
	// clear the null bit from the kind so that it does not get taken into account.
	if kind != cue.TopKind {
		if kind&cue.NullKind == cue.NullKind {
			ptr = true
			kind = kind &^ cue.NullKind
		}
	}

	switch kind {

	case cue.StringKind, cue.IntKind, cue.FloatKind, cue.BoolKind:
		switch lt := name.LabelType(); lt {
		case cue.DefinitionLabel:
			ident := strings.TrimPrefix(name.String(), "#")
			expr := builder.FromKind(val.IncompleteKind()).WithPtr(ptr)
			return builder.NewType(ident).WithExpr(expr), nil
		case cue.StringLabel:
			ident := name.Unquoted()
			expr := builder.FromKind(val.IncompleteKind()).WithPtr(ptr)
			return builder.NewType(ident).WithExpr(expr), nil
		default:
			return builder.NoType, fmt.Errorf("unsupported label type %v at path %v", lt, name.String())
		}

	case cue.StructKind:
		if _, p := val.ReferencePath(); len(p.Selectors()) > 0 {
			ident := name.Unquoted()
			expr := builder.NewIdent(strings.TrimPrefix(p.String(), "#")).WithPtr(ptr)
			return builder.NewType(ident).WithExpr(expr), nil
		}

		return processStruct(name, val)

	case cue.ListKind:
		if _, p := val.ReferencePath(); len(p.Selectors()) > 0 {
			ident := strings.TrimPrefix(name.Unquoted(), "#")
			expr := builder.NewIdent(strings.TrimPrefix(p.String(), "#"))
			return builder.NewType(ident).WithComment(commentsFrom(val)).WithExpr(expr), nil
		}

		el := val.LookupPath(cue.MakePath(cue.AnyIndex))
		if _, p := el.ReferencePath(); len(p.Selectors()) > 0 {
			ident := strings.TrimPrefix(name.String(), "#")
			expr := builder.NewIdent("[]" + strings.TrimPrefix(p.String(), "#")).WithPtr(ptr)
			return builder.NewType(ident).WithComment(commentsFrom(val)).WithExpr(expr), nil
		}

	case cue.TopKind:
		ident := strings.TrimPrefix(name.String(), "#")
		return builder.NewType(ident).WithComment(commentsFrom(val)).WithExpr(builder.NewIdent("any")), nil

	case cue.BottomKind:
		return builder.NoType, fmt.Errorf("unsupported kind %v resolves to _|_ at path %v", kind, val.Path().String())

	default:
		return builder.NoType, fmt.Errorf("unsupported kind %v at path %v", kind, val.Path().String())
	}

	return builder.NoType, errors.New("unreachable")
}

func commentsFrom(val cue.Value) string {
	buf := &bytes.Buffer{}
	if cgs := val.Doc(); len(cgs) > 0 {
		for _, cg := range cgs {
			for _, c := range cg.List {
				text := strings.TrimPrefix(c.Text, "//")
				fmt.Fprintf(buf, "%s\n", strings.TrimSpace(text))
			}
		}
	}

	return buf.String()
}

// processStruct prints the top-level fields of a struct value
func processStruct(name cue.Selector, val cue.Value) (builder.Type, error) {
	ident := strings.TrimPrefix(name.String(), "#")
	expr := builder.NewStruct()

	// Iterate through the fields of the struct
	it, _ := val.Fields(cue.Optional(true))
	for it.Next() {
		field, err := processValue(it.Selector(), it.Value())
		if err != nil {
			return builder.NoType, fmt.Errorf("error processing field %v: %w", it.Selector(), err)
		}

		expr = expr.AddField(field)
	}

	return builder.NewType(ident).WithComment(commentsFrom(val)).WithExpr(expr), nil
}

func processTopLevel(pkg builder.File, val cue.Value) (builder.File, error) {
	it, err := val.Fields(cue.Definitions(true))
	if err != nil {
		return pkg, fmt.Errorf("error iterating over definitions: %w", err)
	}

	for it.Next() {
		v := it.Value()
		if !it.Selector().IsDefinition() {
			continue
		}

		decl, err := processValue(it.Selector(), v)
		if err != nil {
			return pkg, err
		}

		pkg = pkg.AddDecl(decl)
	}

	return pkg, nil
}

func IsOptional(sel cue.Selector) bool {
	sel.Optional()
	return sel.ConstraintType()&cue.OptionalConstraint == cue.OptionalConstraint
}

func IsRequired(sel cue.Selector) bool {
	return sel.ConstraintType()&cue.RequiredConstraint == cue.RequiredConstraint
}

func IsPattern(sel cue.Selector) bool {
	return sel.ConstraintType()&cue.PatternConstraint == cue.PatternConstraint
}
