package trans

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Team-GhostLand/Grinch/util"
)

func SwapServerGitToDev() error {
	if _, err := os.Stat(filepath.FromSlash(util.ServerOvrrDir)); errors.Is(err, fs.ErrNotExist) {
		return nil //Do nothing if there were no server-overrides in the 1st place
	}
	return os.Rename(filepath.FromSlash(util.ServerOvrrDir), filepath.FromSlash(util.DevServerOvrrDir))
}

func ResolveServerRemovals() error {
	sourceFile := filepath.FromSlash(util.RemovalsFileLocation)

	data_bin, err := os.ReadFile(sourceFile)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) { //To see if we got ANY OTHER KIND of error than „not exists” (that's different from what directly checking fs.ErrExist does) - these errors are irrecoverable, so we crash
			return err
		}

		return nil //If REMOVALS.txt doesn't exist, there's nothing to resolve
	}

	data := string(data_bin)
	lines := strings.Split(data, "\n")

	for _, l := range lines {
		path := util.NormalOvrrDir + "/" + l
		newPath := util.ClientOvrrDir + "/" + l
		err = os.Rename(filepath.FromSlash(path), filepath.FromSlash(newPath))
		if err != nil {
			//Recovery attempt: Create any missing directories in the target path, then try again
			err = os.MkdirAll(filepath.FromSlash(newPath), util.ReasonableDirPerms)
			if err != nil {
				return err
			}
			err = os.Remove(filepath.FromSlash(newPath)) //The code above always creates a directory. We remove the last path element, in case it was supposed to be a file.
			if err != nil {
				return err
			}
			err = os.Rename(filepath.FromSlash(path), filepath.FromSlash(newPath))
			if err != nil {
				return err
			}
		}
	}

	err = os.Remove(sourceFile)
	return err
}

func DoExportJsonTransforms(mp *util.ModpackDefinition, em util.ExportMode) (util.MrIndex, error) {

	mi, err := util.GetMrIndexJson(util.MrIndexFileLocation)
	if err != nil {
		return mi, err
	}

	mi.Name = util.ResolveTemplateString(&mi, util.GetIngameExportName(mp, em))

	switch em {
	case util.EmDev:
		util.DoClientsideSupportJsonTransforms(&mi, util.MssUnsupported, util.MssRequired, true)
	case util.EmSlim:
		util.DoClientsideSupportJsonTransforms(&mi, util.MssOptional, util.MssUnsupported, false)
	case util.EmTweakable:
		util.DoClientsideSupportJsonTransforms(&mi, util.MssOptional, util.MssOptional, true)
	}

	return mi, util.SetMrIndexJson(mi, util.MrIndexFileLocation)
}