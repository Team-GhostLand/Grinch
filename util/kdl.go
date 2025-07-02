package util

import (
	"errors"
	"os"

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
	Name    string    `kdl:",arg"`
	Path    string    `kdl:"path"`
	NameOut string    `kdl:"preferred-name"`
	Filters FilterSet `kdl:",children"`
}

type FilterSet struct {
	Allow      []string `kdl:"allow"`
	Expect     []string `kdl:"expect"`
	Disallowed []string `kdl:"disallowed"`
}

func LoadProjectConfig() (ProjectConfigFile, error) {
	var pcf ProjectConfigFile
	data, err := os.ReadFile("grinch.kdl")
	if err != nil {
		return pcf, err
	}
	err = kdl.Unmarshal([]byte(data), &pcf)
	if err != nil {
		return pcf, err
	}
	if len(pcf.MPs.MP) == 0 {
		return pcf, errors.New("there aren't any modpacks defined in your config file")
	}
	return pcf, nil
}

func SelectModpack(pcf ProjectConfigFile) (*ModpackDefinition, error) {
	wcf, err := LoadWorkspaceConfig()
	if err != nil {
		return nil, err
	}

	if wcf[2] == "" && pcf.Default == "" && len(pcf.MPs.MP) > 1 {
		return nil, errors.New("you don't have any modpacks selected in neither workspace nor projects settings, but you have more than one defined, so we cannot auto-select")
	} else if wcf[2] == "" && pcf.Default == "" && len(pcf.MPs.MP) == 1 {
		return &pcf.MPs.MP[0], nil
	} else if wcf[2] != "" { //We don't care whether default-modpack is set, since workspace takes precednece
		return FindModpackByName(pcf, wcf[2])
	} else if wcf[2] == "" && pcf.Default != "" {
		return FindModpackByName(pcf, pcf.Default)
	}

	return nil, errors.New("Something went HORRIBLY wrong when selecting modpacks - PLEASE do report this error to us!!!") //This code should never be reached
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

	return nil, errors.New("modpack " + name + " not found")
}
