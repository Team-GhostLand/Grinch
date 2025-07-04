package trans

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Team-GhostLand/Grinch/util"
)

var ImportTransformPredicates = map[string]util.MrIndexModSideSupportPair{
	"GR_CLO_":     {Client: util.MssOptional, Server: util.MssRequired},
	"GR_SVO_":     {Client: util.MssRequired, Server: util.MssOptional},
	"GR_CLO_SVO_": {Client: util.MssOptional, Server: util.MssOptional},
	"GR_CLX_":     {Client: util.MssUnsupported, Server: util.MssRequired},
	"GR_SVX_":     {Client: util.MssRequired, Server: util.MssUnsupported},
	"GR_CLX_SVX_": {Client: util.MssUnsupported, Server: util.MssUnsupported},
	"GR_CLO_SVX_": {Client: util.MssOptional, Server: util.MssUnsupported},
	"GR_CLX_SVO_": {Client: util.MssUnsupported, Server: util.MssOptional},
}

func SwapServerDevToGit() error {
	if _, err := os.Stat(filepath.FromSlash(util.DevSvOvrrDir)); errors.Is(err, fs.ErrNotExist) {
		return nil //Do nothing if there were no server-overrides in the 1st place
	}
	return os.Rename(filepath.FromSlash(util.DevSvOvrrDir), filepath.FromSlash(util.GitSvOvrrDir))
}
