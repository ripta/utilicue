package builder

import (
	"github.com/mitchellh/go-wordwrap"
	"go/ast"
	"strings"
)

const DocWidth = 80

func FormatDoc(comment string) *ast.CommentGroup {
	if comment == "" {
		return nil
	}

	groups := []string{}

	rawLines := strings.Split(comment, "\n\n")
	for _, rawLine := range rawLines {
		groups = append(groups, strings.ReplaceAll(rawLine, "\n", " "))
	}

	cg := &ast.CommentGroup{}
	for _, group := range groups {
		lines := strings.Split(wordwrap.WrapString(group, DocWidth), "\n")
		for _, line := range lines {
			cg.List = append(cg.List, &ast.Comment{
				Text: "// " + line,
			})
		}
	}

	return cg
}
