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

	return os.WriteFile(filepath.FromSlash(path), data, ReasonablePerms)
}

func DoClientsideSupportJsonTransforms(mi *MrIndex, from, to MrIndexModSideSupport, disable bool) {
	for i, m := range mi.Mods {
		if m.Side.Client == from {
			mi.Mods[i].Side.Client = to
			if disable {
				mi.Mods[i].Path = EnsureExtension(m.Path, "disabled")
			}
		}
	}
}

func DoServersideSupportJsonTransforms(mi *MrIndex, from, to MrIndexModSideSupport, disable bool) {
	for i, m := range mi.Mods {
		if m.Side.Server == from {
			mi.Mods[i].Side.Server = to
			if disable {
				mi.Mods[i].Path = EnsureExtension(m.Path, "disabled")
			}
		}
	}
}

func DoPrefixSideSupportJsonTransforms(mi *MrIndex, predicates map[string]MrIndexModSideSupportPair, p string) {
	for i, m := range mi.Mods {
		if strings.HasPrefix(IsolateEndPathElement(m.Path), p) {
			for p, s := range predicates {
				if strings.HasPrefix(IsolateEndPathElement(m.Path), p) {
					mi.Mods[i].Side = s
					break
				}
			}
		}
	}
}
