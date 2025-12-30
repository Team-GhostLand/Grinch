package util

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type WorkspaceConfigFile = []string

func LoadWorkspaceConfig(path string) (WorkspaceConfigFile, error) {
	generic_error := "Could neither load nor create a .gr-workspace file: "
	must_reload := false

	data_bin, err := os.ReadFile(filepath.FromSlash(path))

	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) { //To see if we got ANY OTHER KIND of error than „not exists” (that's different from what directly checking fs.ErrExist does) - these errors are irrecoverable, so we crash
			return nil, errors.New(generic_error + err.Error())
		}

		file, err := os.Create(filepath.FromSlash(path))
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
		if data_bin, err = os.ReadFile(filepath.FromSlash(path)); err != nil {
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

func CheckAndAddKnownMrpack(name, path string, wcf WorkspaceConfigFile) (bool, error) {
	for i, l := range wcf {
		if i < 3 {
			continue //We only care about 4th line (3rd index) and beyond
		}
		if l == name {
			return true, nil
		}
	}

	err := AppendToWorkspaceConfig(name, path)
	return false, err
}

func PopulateWorkspaceConfigFile(file os.File) error {
	_, err := file.WriteString("FMTv1\n# PLEASE READ THIS COMMENT TO THE VERY END, AS YOU MUST BE VERY CAREFUL WHEN EDITING THIS FILE! This format is VERY particular about line numbers. The line above contains the format version - DO NOT change that. Do not move it, either. FMTv? must be EXACTLY on the 1st line. The line below, ie. the 3rd line (again - don't move them), specifies what modpack are we working on right now. If it has any content, it overrides whatever was set in grinch.kdl - otherwise (if it's blank), it gets ignored and grinch.kdl is followed. Please note, that the line can either be blank or have content, but it nevertheless MUST exist. If the file has less than 3 lines, it fails to parse. Finally, the next lines (4th and beyond) are for storing discovered MRPACKS (ie. those that were either exported, or already imported) - so that you don't need to specify a filename whenever using grinch import. ONE MORE IMPORTANT THING: This file is supposed to stay local to your PC. Do not commit it to any Git repos. It's best to .gitignore it, alongside all *.mrpack files.\n")
	return err
}

func AppendToWorkspaceConfig(line, path string) error {
	f, err := os.OpenFile(filepath.FromSlash(path), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("\n" + line)
	return err
}

func FindNewMrpack(cmd_args []string, dir_path, wcf_path string, wcf WorkspaceConfigFile) (string, error) {
	if len(cmd_args) > 0 {
		return cmd_args[0], AppendToWorkspaceConfig(cmd_args[0], wcf_path)
	}

	files, err := os.ReadDir(filepath.FromSlash(dir_path))
	if err != nil {
		return "", err //This is a „true error”
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".mrpack") {
			continue
		}
		known, err := CheckAndAddKnownMrpack(f.Name(), wcf_path, wcf)
		if !known {
			return f.Name(), err //Even if err != nil, the search was successful, so we should always check for that first (and optionally warn that appending failed, if err != nil)
		}
	}

	return "", nil //Not finding anything isn't an error, per say.
}