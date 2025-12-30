package util

import (
	"errors"
	"io/fs"
	"log"
	"os"
)

const (
	Tempdir                 = ".temp"
	Backup                  = ".old"
	ServerOvrrDir           = Tempdir + "/server-overrides"
	ClientOvrrDir           = Tempdir + "/client-overrides"
	NormalOvrrDir           = Tempdir + "/overrides"
	DevServerOvrrDir        = NormalOvrrDir + "/.SERVERSIDE"
	MrIndexFileLocation     = Tempdir + "/modrinth.index.json"
	RemovalsFileLocation    = ServerOvrrDir + "/REMOVALS.txt"
	GrWorkspaceFileLocation = ".gr-workspace"
	GrProjectFileLocation   = "grinch.kdl"
	DisabledExtension       = "disabled"
	ReasonableFilePerms     = 0644
	ReasonableDirPerms      = 0755
)

func Hndl(err error, with string, cleanup bool) {
	if err != nil {
		if cleanup {
			_ = os.RemoveAll(Tempdir) //We don't check for errors; it doesn't really matter at this point. The app's about to crash, anyway.
		}
		log.Fatal("ERROR: ", with, ": ", err)
	}
}

func IsSafelyCreatable(path string) (bool, error) {
	_, err := os.Stat(path)

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) { //If we got a NotExists error - great! There's nothing. Can safely write.
			return true, nil
		}
		return false, err //...otherwise, either the directory/file DOES EXIST (but couldn't be Stat'd for some reason) or doesn't (but then likely wouldn't be safely creatable, either - eg. because we have no permissions or the filesystem failed)
	}

	//If the error wasn't NIL, then Stat was successful - ie. the directory/file must, logically, exist
	return false, fs.ErrExist

}