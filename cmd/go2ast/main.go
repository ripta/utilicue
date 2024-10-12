package main

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"os"
)

// go run ./cmd/go2ast ./examples/db/generated.cue2go.go
func main() {
	fset := token.NewFileSet()
	flags := parser.ParseComments | parser.SkipObjectResolution | parser.AllErrors

	file, err := parser.ParseFile(fset, os.Args[1], nil, flags)
	if err != nil {
		panic(err)
	}

	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
