package trans

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Team-GhostLand/Grinch/util"
)

func SwapServerDevToGit() error {
	if _, err := os.Stat(filepath.FromSlash(util.DevSvOvrrDir)); errors.Is(err, fs.ErrNotExist) {
		return nil //Do nothing if there were no server-overrides in the 1st place
	}
	return os.Rename(filepath.FromSlash(util.DevSvOvrrDir), filepath.FromSlash(util.GitSvOvrrDir))
}

func DoImportJsonTransforms( /*???*/ ) error {
	return nil //TODO: Implement
}
