package trans

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Team-GhostLand/Grinch/util"
)

type ExportMode int

const (
	EmDefault ExportMode = iota
	EmQuick
	EmDev
	EmSlim
	EmTweakable
)

func SwapServerGitToDev() error {
	if _, err := os.Stat(filepath.FromSlash(util.GitSvOvrrDir)); errors.Is(err, fs.ErrNotExist) {
		return nil //Do nothing if there were no server-overrides in the 1st place
	}
	return os.Rename(filepath.FromSlash(util.GitSvOvrrDir), filepath.FromSlash(util.DevSvOvrrDir))
}

func ResolveServerRemovals() error {
	return errors.New("default export mode not yet implemented, please use --quick for now (and - if needed - create a server-pack out of it using the make_serverpack.sh script available on our GitHub)") //TODO
}

func DoExportJsonTransforms(em ExportMode) error {
	if em == EmDefault || em == EmQuick {
		return nil //They don't need any JSON transforms - early-return
	}

	mi, err := util.GetMrIndexJson(util.MrIndexFileLocation)
	if err != nil {
		return err
	}

	switch em {
	case EmDev:
		util.DoClientsideSupportJsonTransforms(&mi, util.MssUnsupported, util.MssRequired, true)
	case EmSlim:
		util.DoClientsideSupportJsonTransforms(&mi, util.MssOptional, util.MssUnsupported, false)
	case EmTweakable:
		util.DoClientsideSupportJsonTransforms(&mi, util.MssOptional, util.MssOptional, true)
	}

	return util.SetMrIndexJson(mi, util.MrIndexFileLocation)
}
