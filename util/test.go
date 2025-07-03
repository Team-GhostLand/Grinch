package util

import (
	"fmt"
	"log"
)

func TestConfigLoad() {
	pcf, err := LoadProjectConfig()
	Hndl(err, "Couldn't load project config", false)
	fmt.Printf("%#v\n", pcf)
}

func TestWkspcLoad1() {
	s, err := LoadWorkspaceConfig()
	Hndl(err, "Couldn't load workspace config", false)
	log.Println(s)
}

func TestSelection() {
	pcf, err := LoadProjectConfig()
	Hndl(err, "Couldn't load project config", false)
	s, err := SelectModpack(pcf)
	Hndl(err, "Couldn't select modpack", false)
	log.Println(s.Path)
}

func TestJsonTransforms() {
	mi, err := GetMrIndexJson()
	Hndl(err, "Couldn't open JSON", false)
	log.Println("PRE:\n", mi)
	DoClientJsonTransforms(&mi, MssRequired, MssUnsupported, true)
	log.Println("POST:\n", mi)
	err = SetMrIndexJson(mi)
	Hndl(err, "Couldn't save JSON", false)
}
