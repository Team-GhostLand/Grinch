package trans

import (
	"cmp"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/Team-GhostLand/Grinch/util"
)

var ImportTransformPredicates = map[string]util.MrIndexModSideSupportPair{
	"CLO_":     {Client: util.MssOptional, Server: util.MssRequired},
	"SVO_":     {Client: util.MssRequired, Server: util.MssOptional},
	"CLO+SVO_": {Client: util.MssOptional, Server: util.MssOptional},
	"CLX_":     {Client: util.MssUnsupported, Server: util.MssRequired},
	"SVX_":     {Client: util.MssRequired, Server: util.MssUnsupported},
	"CLX+SVX_": {Client: util.MssUnsupported, Server: util.MssUnsupported},
	"CLO+SVX_": {Client: util.MssOptional, Server: util.MssUnsupported},
	"CLX+SVO_": {Client: util.MssUnsupported, Server: util.MssOptional},
}

func SwapServerDevToGit() error {
	if _, err := os.Stat(filepath.FromSlash(util.DevServerOvrrDir)); errors.Is(err, fs.ErrNotExist) {
		return nil //Do nothing if there were no server-overrides in the 1st place
	}
	return os.Rename(filepath.FromSlash(util.DevServerOvrrDir), filepath.FromSlash(util.ServerOvrrDir))
}

func ApplyJsonParamsOnImport(mi *util.MrIndex, params util.PackDefParams) {
	if params.Names.Git != "" {
		mi.Name = params.Names.Git
	}

	if params.Description != "" {
		mi.Desc = params.Description
	}
}

func SortMrIndexOnImport(mi *util.MrIndex) {
	slices.SortFunc(mi.Mods, func(a, b util.MrIndexModInstance) int {
		return cmp.Compare(a.Path, b.Path)
	})
}