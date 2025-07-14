package trans

import (
	"cmp"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

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
	if _, err := os.Stat(filepath.FromSlash(util.DevSvOvrrDir)); errors.Is(err, fs.ErrNotExist) {
		return nil //Do nothing if there were no server-overrides in the 1st place
	}
	return os.Rename(filepath.FromSlash(util.DevSvOvrrDir), filepath.FromSlash(util.GitSvOvrrDir))
}

func SolveJsonImportConstraints(mi *util.MrIndex, constr util.PackDefConstraints) error {
	if constr.Name != "" {
		mi.Name = constr.Name
	}

	if constr.Description != "" {
		mi.Desc = constr.Description
	}

	if !strings.HasPrefix(mi.Ver, constr.Version) { //no need for a != "" check - every string starts out empty, after all
		return errors.New("modpack version " + mi.Ver + " doesn't start with the expected " + constr.Version)
	}

	return nil
}

func SolveFileImportConstraints(fset util.PackDefConstrFilterSet) error {
	if fset.Allow == nil && fset.Disallowed == nil && fset.Expect == nil {
		return nil
	}
	if fset.Allow == nil || fset.Disallowed == nil || fset.Expect == nil {
		return errors.New("if using file constraints, all 3 filters (allow+expect+disallowed) must be present (they CAN be empty, if you don't want to use them), but you only specified 2 or 1 out of them")
	}

	return errors.New("not yet implemented") //TODO
}

func SortMrIndexOnImport(mi *util.MrIndex) {
	slices.SortFunc(mi.Mods,
		func(a, b util.MrIndexModInstance) int {
			return cmp.Compare(a.Path, b.Path)
		})
}
