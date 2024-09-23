package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ripta/utilicue/pkg/cue2go"
)

func main() {
	gen := &cue2go.Generator{}
	root := &cobra.Command{
		Use:           "cue2go",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return gen.Run(args)
		},
	}

	if err := root.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
