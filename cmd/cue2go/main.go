package main

import (
	"context"
	"fmt"
	"github.com/thediveo/enumflag"
	"os"

	"github.com/spf13/cobra"

	"github.com/ripta/utilicue/pkg/cue2go"
)

func main() {
	gen := &cue2go.Generator{
		ExportMode: cue2go.ExportModeRespectSource,
	}

	root := &cobra.Command{
		Use:           "cue2go",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return gen.Run(args)
		},
	}

	modeValue := enumflag.New(&gen.ExportMode, "export-mode", cue2go.ExportModeIds, enumflag.EnumCaseInsensitive)
	root.Flags().VarP(modeValue, "export-mode", "e", "export mode: respect-source, all)")

	if err := root.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
