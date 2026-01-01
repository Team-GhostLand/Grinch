package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sblinch/kdl-go"
)

type ProjectConfigFile struct {
	Version int    `kdl:"format-version"`
	Default string `kdl:"default-modpack"`
	MPs     struct {
		MP []ModpackDefinition `kdl:"modpack,multiple"`
	} `kdl:"modpacks"`
}

type ModpackDefinition struct {
	Name   string        `kdl:",arg"`
	Path   string        `kdl:"path"`
	Params PackDefParams `kdl:",children"`
}

type PackDefParams struct {
	Names       NameSet `kdl:"names"`
	Description string  `kdl:"description"`
}

type NameSet struct {
	Default   string `kdl:"default"`
	Quick     string `kdl:"quick"`
	Dev       string `kdl:"dev"`
	Slim      string `kdl:"slim"`
	Tweakable string `kdl:"tweakable"`
	Git       string `kdl:"git"`
}

func LoadProjectConfig(path string) (ProjectConfigFile, error) {
	var pcf ProjectConfigFile
	data, err := os.ReadFile(filepath.FromSlash(path))
	if err != nil {
		return pcf, err
	}
	err = kdl.Unmarshal([]byte(data), &pcf)
	if err != nil {
		return pcf, err
	}
	if pcf.Version != 4 {
		return pcf, errors.New("this version of grinch uses config version 3, but yours is written in version " + fmt.Sprint(pcf.Version))
	}
	if len(pcf.MPs.MP) == 0 {
		return pcf, errors.New("there aren't any modpacks defined in your config file")
	}
	return pcf, nil
}

func SelectModpack(pcf ProjectConfigFile, wcf WorkspaceConfigFile) (*ModpackDefinition, error) {

	if wcf[2] == "" && pcf.Default == "" && len(pcf.MPs.MP) > 1 {
		return nil, errors.New("you don't have any modpacks selected in neither workspace nor project settings, but you have more than one defined, so we cannot auto-select")
	} else if wcf[2] == "" && pcf.Default == "" && len(pcf.MPs.MP) == 1 {
		return &pcf.MPs.MP[0], nil
	} else if wcf[2] != "" { //We don't care whether default-modpack is set, since workspace takes precedence
		return FindModpackByName(pcf, wcf[2])
	} else if wcf[2] == "" && pcf.Default != "" {
		return FindModpackByName(pcf, pcf.Default)
	}

	return nil, errors.New("Something went HORRIBLY wrong when selecting modpacks - PLEASE do report this error to us!!!") //This code is unreachable. If it ever IS reached, then we got a len(pcf.MPs.MP) < 1, which shouldn't happen under normal circumstances (LoadConfig() throws if there are no modpacks), hence the ominous error message (we even broke convention of „no capitals” and „no exclamation marks” just to signify how important is it).
}

func FindModpackByName(pcf ProjectConfigFile, name string) (*ModpackDefinition, error) {
	if name == "" {
		return nil, errors.New("you're trying to search for a modpack without a name")
	}

	for _, pack := range pcf.MPs.MP {
		if name == pack.Name {
			return &pack, nil
		}
	}

	return nil, errors.New("modpack " + name + " isn't defined anywhere in grinch.kdl")
}

func GetFileExportName(mp *ModpackDefinition, mode ExportMode, nameOverride string) string {
	ext := "mrpack"
	if nameOverride != "" {
		return EnsureExtension(nameOverride, ext)
	} else {
		return EnsureExtension(GetIngameExportName(mp, mode), ext)
	}
}

func GetIngameExportName(mp *ModpackDefinition, mode ExportMode) string {
	switch mode {
	case EmDefault:
		return mp.Params.Names.Default
	case EmQuick:
		return mp.Params.Names.Quick
	case EmDev:
		return mp.Params.Names.Dev
	case EmSlim:
		return mp.Params.Names.Slim
	case EmTweakable:
		return mp.Params.Names.Tweakable
	default:
		return mp.Params.Names.Git
	}
}