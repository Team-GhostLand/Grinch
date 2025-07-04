package util

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"strings"
)

// Various file locations that should stay consistent within the app
const (
	Tempdir                 = ".temp"
	Backup                  = ".old"
	GitSvOvrrDir            = Tempdir + "/server-overrides"
	DevSvOvrrDir            = Tempdir + "/overrides/.SERVERSIDE"
	MrIndexFileLocation     = Tempdir + "/modrinth.index.json"
	GrWorkspaceFileLocation = ".gr-workspace"
	GrProjectFileLocation   = "grinch.kdl"
	ReasonablePerms         = 0755
)

func Hndl(err error, with string, cleanup bool) {
	if err != nil {
		if cleanup {
			_ = os.RemoveAll(Tempdir) //We don't check for errors; it doesn't really matter at this point. The app's about to crash, anyway.
		}
		log.Fatal("ERROR: ", with, ": ", err)
	}
}

func EnsureExtension(fname, ext string) string {
	if strings.HasSuffix(fname, "."+ext) {
		return fname
	} else {
		return fname + "." + ext
	}
}

func IsSafelyCreateable(path string) (bool, error) {
	_, err := os.Stat(path)

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) { //If we got a NotExists error - great! There's nothing. Can safely write.
			return true, nil
		}
		return false, err //...otherwise, either the direcotry/file DOES EXIST (but couldn't be Stat'd for some reason) or doesn't (but then likely wouldn't be safely createable, either - eg. because we have no permissions or the filesystem failed)
	}

	//If the error wasn't NIL, then Stat was succesful - ie. the direcotry/file must, logically, exist
	return false, fs.ErrExist

}

func StripFirstPathElement(path string) string {
	parts := strings.SplitN(path, "/", 2)
	return parts[1]
}
