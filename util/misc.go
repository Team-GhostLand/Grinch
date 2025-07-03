package util

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	Tempdir = ".temp"
)

func Hndl(err error, with string, cleanup bool) {
	if err != nil {
		if cleanup {
			_ = os.RemoveAll(Tempdir) //We don't check for errors; it doesn't really matter at this point. The app's about to crash, anyway.
		}
		log.Fatal("ERROR: ", with, ": ", err)
	}
}

func GetExportName(mp *ModpackDefinition, nameOverride string) string {
	ext := "mrpack"
	if nameOverride != "" {
		return EnsureExtension(nameOverride, ext)
	} else if mp.NameOut != "" {
		return EnsureExtension(mp.NameOut, ext)
	} else {
		return mp.Name + "." + ext
	}
}

func EnsureExtension(fname, ext string) string {
	if strings.HasSuffix(fname, "."+ext) {
		return fname
	} else {
		return fname + "." + ext
	}
}

func IsSafelyCreateable(name string) (bool, error) {
	_, err := os.Stat(name)

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) { //If we got a NotExists error - great! There's nothing. Can safely write.
			return true, nil
		}
		return false, err //...otherwise, either the direcotry/file DOES EXIST (but couldn't be opened for some reason) or doesn't (but then likely wouldn't be safely createable, either - eg. because we have no permissions or the filesystem failed)
	}

	//If the error wasn't NIL, then Stat was succesful - ie. the direcotry/file must, logically, exist
	return false, fs.ErrExist

}

func MakeZipFile(src, dest string) error {
	//thx, https://stackoverflow.com/questions/37869793/how-do-i-zip-a-directory-containing-sub-directories-or-files-in-golang
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Ensure that `path` is not absolute; it should not start with "/".
		// This snippet happens to work because I don't use
		// absolute paths, but ensure your real-world code
		// transforms path into a zip-root relative path.
		f, err := w.Create(StripFirstPathElement(path)) //Modification from stackoverflow anwser - without it, we'd copy names with the „.temp” prefix, which we don't want
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	err = filepath.Walk(src, walker)
	if err != nil {
		return err
	}

	return nil
}

func StripFirstPathElement(path string) string {
	parts := strings.SplitN(path, "/", 2)
	return parts[1]
}
