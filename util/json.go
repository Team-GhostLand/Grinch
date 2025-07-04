package util

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type MrIndexModSideSupport string

const (
	MssRequired    MrIndexModSideSupport = "required"
	MssOptional    MrIndexModSideSupport = "optional"
	MssUnsupported MrIndexModSideSupport = "unsupported"
)

type MrIndexModInstance struct {
	Path   string            `json:"path"`
	Hashes map[string]string `json:"hashes"`
	Side   struct {
		Client MrIndexModSideSupport `json:"client"`
		Server MrIndexModSideSupport `json:"server"`
	} `json:"env"`
	Sources []string `json:"downloads"`
	Size    int      `json:"fileSize"`
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

func GetMrIndexJson() (MrIndex, error) {
	var mi MrIndex
	data, err := os.ReadFile(filepath.FromSlash(MrIndexFileLocation))
	if err != nil {
		return mi, err
	}
	err = json.Unmarshal(data, &mi)
	return mi, err
}

func SetMrIndexJson(mi MrIndex) error {
	//data, err := os.ReadFile("testindex.json")
	data, err := json.MarshalIndent(mi, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.FromSlash(MrIndexFileLocation), data, ReasonablePerms)
}

func DoClientJsonTransforms(mi *MrIndex, from, to MrIndexModSideSupport, disable bool) {
	for i, m := range mi.Mods {
		if m.Side.Client == from {
			mi.Mods[i].Side.Client = to
			if disable {
				mi.Mods[i].Path = EnsureExtension(m.Path, "disabled")
			}
		}
	}
}
