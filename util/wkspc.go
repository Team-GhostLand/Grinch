package util

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

type WorkspaceConfigFile = []string

func LoadWorkspaceConfig() (WorkspaceConfigFile, error) {
	generic_error := "Could neither load nor create a .gr-workspace file: "
	must_reload := false

	data_bin, err := os.ReadFile(GrWorkspaceFileLocation)

	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) { //To see if we got ANY OTHER KIND of error than „not exists” (that's different from what directly checking fs.ErrExist does) - these errors are irrecoverable, so we crash
			return nil, errors.New(generic_error + err.Error())
		}

		file, err := os.Create(GrWorkspaceFileLocation)
		if err != nil {
			return nil, errors.New(generic_error + err.Error())
		}
		defer file.Close()
		if err = PopulateWorkspaceConfigFile(*file); err != nil {
			return nil, errors.New("Couldn't write to a just-created .gr-workspace file: " + err.Error())
		}
		must_reload = true
	}

	if must_reload {
		if data_bin, err = os.ReadFile(GrWorkspaceFileLocation); err != nil {
			return nil, errors.New("Couldn't load a just-created .gr-workspace file" + err.Error())
		}
	}
	data := string(data_bin)
	lines := strings.Split(data, "\n")

	if len(lines) < 3 {
		return nil, errors.New("your .gr-workspace is malformatted: it has " + fmt.Sprint(len(lines)) + " lines (less than the minimum of 3)")
	}

	if lines[0] != "FMTv1" {
		return nil, errors.New("your .gr-workspace uses a different format version than v1, which is what this version of Grinch supports")
	}

	return lines, nil
}

func CheckAndAddKnownMrpack(name string, wcf WorkspaceConfigFile) (bool, error) {
	for i, l := range wcf {
		if i < 3 {
			continue //We only care about 4th line (3rd index) and beyond
		}
		if l == name {
			return true, nil
		}
	}

	err := AppendToWorkspaceConfig(name)
	return false, err
}

func PopulateWorkspaceConfigFile(file os.File) error {
	_, err := file.WriteString("FMTv1\n# TODO: Have a proper comment here\n")
	return err
}

func AppendToWorkspaceConfig(line string) error {
	f, err := os.OpenFile(GrWorkspaceFileLocation, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("\n" + line)
	return err
}

func FindNewMrpack(cmd_args []string) (string, error) {
	if len(cmd_args) > 0 {
		return cmd_args[0], AppendToWorkspaceConfig(cmd_args[0])
	}

	wcf, err := LoadWorkspaceConfig()
	if err != nil {
		return "", err //This is a „true error”
	}

	files, err := os.ReadDir(".")
	if err != nil {
		return "", err //This is a „true error”
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".mrpack") {
			continue
		}
		known, err := CheckAndAddKnownMrpack(f.Name(), wcf)
		if !known {
			return f.Name(), err //Even if err != nil, the search was sucesful, so we should always check for that first (and optionally warn that appending failed, if err != nil)
		}
	}

	return "", nil //Not finding anything isn't an error, per say.
}
