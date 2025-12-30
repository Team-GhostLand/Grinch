package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/Team-GhostLand/Grinch/util"

	"github.com/spf13/cobra"
)

var vqCmd = &cobra.Command{
	Use:   "vq [modpack]",
	Short: "Version Query - prints selected/default modpack's version",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) > 1 {
			util.Hndl(errors.New(fmt.Sprint(len(args))+" is more than the maximum of 1\n[TIP: If you want to input a space without it being interpreted as an argument separator, prefix it with \\ or surround the entire piece of text containing space(s) with quotation marks]"), "Too many arguments", false)
		}

		pcf, err := util.LoadProjectConfig(util.GrProjectFileLocation)
		util.Hndl(err, "Couldn't load project config", false)
		wcf, err := util.LoadWorkspaceConfig(util.GrWorkspaceFileLocation)
		util.Hndl(err, "Couldn't load workspace config", false)

		var mp *util.ModpackDefinition
		if len(args) == 0 {
			mp, err = util.SelectModpack(pcf, wcf)
			util.Hndl(err, "Couldn't select modpack", false)
		} else {
			mp, err = util.FindModpackByName(pcf, args[0])
			util.Hndl(err, "Couldn't find modpack "+args[0], false)
		}

		path := filepath.FromSlash(mp.Path)
		if path == "" {
			util.Hndl(errors.New("chosen modpack "+mp.Name+" has no path associated with it"), "Couldn't select modpack", false)
		}

		mri, err := util.GetMrIndexJson(filepath.FromSlash(path + "/modrinth.index.json"))
		util.Hndl(err, "Couldn't read Modrinth Index:", false)

		//END RESULT:
		fmt.Println(mri.Ver) //TODO: Add more query types?
	},
}

func init() {
	rootCmd.AddCommand(vqCmd)
}