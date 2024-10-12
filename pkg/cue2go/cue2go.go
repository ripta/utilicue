package cue2go

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
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

		buf := &bytes.Buffer{}

		fmt.Fprintf(buf, "// Code generated by cue2go. DO NOT EDIT.\n")
		fmt.Fprintf(buf, "// Source: %s\n", abs)
		fmt.Fprintf(buf, "package %s\n", filepath.Base(abs))

		it, err := v.Fields(cue.Definitions(true))
		if err != nil {
			return fmt.Errorf("error iterating over definitions: %w", err)
		}

		for it.Next() {
			v := it.Value()
			if !it.Selector().IsDefinition() {
				continue
			}

			fmt.Fprint(buf, "type ")
			valueToGo(buf, it.Selector(), v)
		}

		raw := buf.Bytes()

		if err := os.WriteFile(filepath.Join(abs, "generated.cue2go.go"), raw, 0644); err != nil {
			return fmt.Errorf("error writing output file: %w", err)
		}

		formatted, err := format.Source(raw)
		if err != nil {
			return fmt.Errorf("error formatting output: %w", err)
		}

		if err := os.WriteFile(filepath.Join(abs, "generated.cue2go.go"), formatted, 0644); err != nil {
			return fmt.Errorf("error writing output file: %w", err)
		}
	}

	return nil
}

func valueToGo(buf *bytes.Buffer, name cue.Selector, val cue.Value) {
	ptr := ""
	if name.ConstraintType()&cue.OptionalConstraint == cue.OptionalConstraint {
		ptr = "*"
	}

	switch k := val.IncompleteKind(); k {

	case cue.StringKind, cue.IntKind, cue.FloatKind, cue.BoolKind:
		switch lt := name.LabelType(); lt {
		case cue.DefinitionLabel:
			fmt.Fprintf(buf, "%v %v\n", strings.TrimPrefix(name.String(), "#"), val.IncompleteKind())
		case cue.StringLabel:
			fmt.Fprintf(buf, "\t%v %s%v\n", name.Unquoted(), ptr, val.IncompleteKind())
		default:
			panic(fmt.Sprintf("unsupported label type %v at path %v", lt, name.String()))
		}

	case cue.StructKind:
		if _, p := val.ReferencePath(); len(p.Selectors()) > 0 {
			fmt.Fprintf(buf, "\t%v %s%v\n", name.Unquoted(), ptr, strings.TrimPrefix(p.String(), "#"))
			return
		}

		structToType(buf, name, val)

	case cue.ListKind:
		el := val.LookupPath(cue.MakePath(cue.AnyIndex))
		if _, p := el.ReferencePath(); len(p.Selectors()) > 0 {
			fmt.Fprint(buf, "\n")
			copyComments(buf, val)
			fmt.Fprintf(buf, "%v []%v\n", strings.TrimPrefix(name.String(), "#"), strings.TrimPrefix(p.String(), "#"))
			return
		}

	case cue.TopKind:
		fmt.Fprintf(buf, "%v any\n", strings.TrimPrefix(name.String(), "#"))

	default:
		panic(fmt.Sprintf("unsupported kind %v at path %v", k, val.Path().String()))
	}
}

func copyComments(buf *bytes.Buffer, val cue.Value) {
	if cgs := val.Doc(); len(cgs) > 0 {
		for _, cg := range cgs {
			for _, c := range cg.List {
				fmt.Fprintf(buf, "%s\n", c.Text)
			}
		}
	}
}

// structToType prints the top-level fields of a struct value
func structToType(buf *bytes.Buffer, name cue.Selector, val cue.Value) {
	fmt.Fprint(buf, "\n")
	copyComments(buf, val)
	fmt.Fprintf(buf, "%v struct {\n", strings.TrimPrefix(name.String(), "#"))

	// Iterate through the fields of the struct
	it, _ := val.Fields(cue.Optional(true))
	for it.Next() {
		valueToGo(buf, it.Selector(), it.Value())
	}

	fmt.Fprintf(buf, "}\n")
}
