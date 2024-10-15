package cue2go

import (
	"github.com/thediveo/enumflag"
)

type ExportModeFlag enumflag.Flag

const (
	ExportModeRespectSource ExportModeFlag = iota
	ExportModeAll
)

var ExportModeIds = map[ExportModeFlag][]string{
	ExportModeRespectSource: {"respect-source", "default"},
	ExportModeAll:           {"all"},
}
