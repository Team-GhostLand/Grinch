package util

import (
	"fmt"
	"log"
)

func Hndl(err error, with string) {
	if err != nil {
		log.Fatal("ERROR: ", with, ": ", err)
	}
}

func TestConfigLoad() {
	if cf, err := LoadProjectConfig(); err == nil {
		fmt.Printf("%#v\n", cf)
	} else {
		Hndl(err, "Couldn't load project config")
	}
}

func TestWkspcLoad1() {
	if s, err := LoadWorkspaceConfig(); err == nil {
		log.Println(s)
	} else {
		Hndl(err, "Couldn't load workspace config")
	}
}

func TestSelection() {
	if pcf, err := LoadProjectConfig(); err == nil {
		if s, err := SelectModpack(pcf); err == nil {
			log.Println(s.Path)
		} else {
			Hndl(err, "Couldn't select modpack")
		}
	} else {
		Hndl(err, "Couldn't load project config")
	}
}
