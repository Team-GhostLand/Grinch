package util

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type WorkspaceConfigFile = []string

func LoadWorkspaceConfig() (WorkspaceConfigFile, error) {
	filename := ".gr-workspace"
	generic_error := "Could neither load nor create a .gr-workspace file: "
	must_reload := false

	data_bin, err := os.ReadFile(filename)

	if err != nil {
		if !os.IsNotExist(err) { //To see if we got ANY OTHER KIND of error than „not exists” (that's different from what os.IsExist does) - these errors are irrecoverable, so we crash
			return nil, errors.New(generic_error + err.Error())
		}

		file, err := os.Create(".gr-workspace")
		if err != nil {
			return nil, errors.New(generic_error + err.Error())
		}
		if err = PopulateWorkspaceFile(*file); err != nil {
			return nil, errors.New("Couldn't write to a just-created .gr-workspace file: " + err.Error())
		}
		must_reload = true
		file.Close()
	}

	if must_reload {
		if data_bin, err = os.ReadFile(filename); err != nil {
			return nil, errors.New("Couldn't load a just-created .gr-workspace file" + err.Error())
		}
	}
	data := string(data_bin)
	lines := strings.Split(data, "\n")

	if len(lines) < 3 {
		return nil, errors.New("your .gr-workspace is malformatted: it has " + fmt.Sprint(len(lines)) + " lines (less than minimum of 3)")
	}

	if lines[0] != "FMTv1" {
		return nil, errors.New("your .gr-workspace uses a different format version than v1, which is what this version of Grinch supports")
	}

	return lines, nil
}

/*func CheckAndAddKnownMrpack() (bool, error) {

}*/

func PopulateWorkspaceFile(file os.File) error {
	_, err := file.WriteString("FMTv1\n# yap yap yap\n")
	return err
}
