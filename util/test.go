package util

import (
	"fmt"
	"log"
)

func TestConfigLoad() {
	pcf, err := LoadProjectConfig(GrProjectFileLocation)
	Hndl(err, "Couldn't load project config", false)
	fmt.Printf("%#v\n", pcf)
}

func TestWkspcLoad1() {
	wcf, err := LoadWorkspaceConfig(GrWorkspaceFileLocation)
	Hndl(err, "Couldn't load workspace config", false)
	log.Println(wcf)
}

func TestSelection() {
	pcf, err := LoadProjectConfig(GrProjectFileLocation)
	Hndl(err, "Couldn't load project config", false)
	wcf, err := LoadWorkspaceConfig(GrWorkspaceFileLocation)
	Hndl(err, "Couldn't load workspace config", false)
	mp, err := SelectModpack(pcf, wcf)
	Hndl(err, "Couldn't select modpack", false)
	log.Println(mp.Path)
}

func TestJsonTransforms() {
	mi, err := GetMrIndexJson(MrIndexFileLocation)
	Hndl(err, "Couldn't open JSON", false)
	log.Println("PRE:\n", mi)
	DoClientsideSupportJsonTransforms(&mi, MssRequired, MssUnsupported, true)
	log.Println("POST:\n", mi)
	err = SetMrIndexJson(mi, MrIndexFileLocation)
	Hndl(err, "Couldn't save JSON", false)
}
