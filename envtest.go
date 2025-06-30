package main

import (
	"fmt"

	"github.com/sblinch/kdl-go"
)

func Main() {

	data := `format-version 1

			modpacks {
    			modpack "Something" path="./modpack";
    			modpack "Something else" path="./modpack2" {
					property1;
					prop2;
			    }
    			modpack "Something even different";
			}
		`

	type ConfigFile struct {
		Version int `kdl:"format-version"`
		MPs     struct {
			MP []struct {
				Name string                 `kdl:",arg"`
				Path string                 `kdl:"path"`
				Kids map[string]interface{} `kdl:",children"`
			} `kdl:"modpack,multiple"`
		} `kdl:"modpacks"`
	}

	var cf ConfigFile
	if err := kdl.Unmarshal([]byte(data), &cf); err == nil {
		fmt.Printf("%#v\n", cf.MPs.MP[2])
	} else {
		fmt.Printf("Error: \n%#v\n", err)
	}
}
