package util

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type MrIndexModSideSupport string

const (
	MssRequired    MrIndexModSideSupport = "required"
	MssOptional    MrIndexModSideSupport = "optional"
	MssUnsupported MrIndexModSideSupport = "unsupported"
)

type MrIndexModSideSupportPair struct {
	Client MrIndexModSideSupport `json:"client"`
	Server MrIndexModSideSupport `json:"server"`
}

type MrIndexModInstance struct {
	Path    string                    `json:"path"`
	Hashes  map[string]string         `json:"hashes"`
	Side    MrIndexModSideSupportPair `json:"env"`
	Sources []string                  `json:"downloads"`
	Size    int                       `json:"fileSize"`
}

type MrIndex struct {
	Game string               `json:"game"`
	Fmt  int                  `json:"formatVersion"`
	Ver  string               `json:"versionId"`
	Name string               `json:"name"`
	Desc string               `json:"summary"`
	Mods []MrIndexModInstance `json:"files"`
	Deps map[string]string    `json:"dependencies"`
}

func GetMrIndexJson(path string) (MrIndex, error) {
	var mi MrIndex
	data, err := os.ReadFile(filepath.FromSlash(path))
	if err != nil {
		return mi, err
	}
	err = json.Unmarshal(data, &mi)
	return mi, err
}

func SetMrIndexJson(mi MrIndex, path string) error {
	data, err := json.MarshalIndent(mi, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.FromSlash(path), data, ReasonableFilePerms)
}

func DoClientsideSupportJsonTransforms(mi *MrIndex, from, to MrIndexModSideSupport, disable bool) {
	for i, m := range mi.Mods {
		if m.Side.Client == from {
			mi.Mods[i].Side.Client = to
			if disable {
				mi.Mods[i].Path = EnsureExtension(m.Path, DisabledExtension)
			}
		}
	}
}

func DoServersideSupportJsonTransforms(mi *MrIndex, from, to MrIndexModSideSupport, disable bool) {
	for i, m := range mi.Mods {
		if m.Side.Server == from {
			mi.Mods[i].Side.Server = to
			if disable {
				mi.Mods[i].Path = EnsureExtension(m.Path, DisabledExtension)
			}
		}
	}
}

func DoPrefixSideSupportJsonTransforms(mi *MrIndex, predicates map[string]MrIndexModSideSupportPair, pfx string, disable bool) {
	for i, m := range mi.Mods {
		if strings.HasPrefix(IsolateEndPathElement(m.Path), pfx) { //Step 1: Find a mod that has Grinch's prefix

			for prd, s := range predicates { //Step 2: If its prefix is one of those „transform supported sides” ones - do said side transforms
				if strings.HasPrefix(IsolateEndPathElement(m.Path), (pfx + prd)) {
					mi.Mods[i].Side = s
					break
				}
			}

			if disable { //Step 3: Mark as either enabled or disabled, depending on what was passed in
				mi.Mods[i].Path = EnsureExtension(m.Path, DisabledExtension)
			} else {
				mi.Mods[i].Path = strings.TrimSuffix(m.Path, ("." + DisabledExtension))
			}
		}
	}
}

//TODO: Make disable bool's behavior consistent across all funcs above (right now, the upper 2 leave disabled status alone if disable==false, while the lower will make it enabled in such case). Instead, make an enum to hold all 3 possible options ("leave alone", "disable", "enable")