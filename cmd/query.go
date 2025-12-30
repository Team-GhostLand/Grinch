package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/Team-GhostLand/Grinch/util"

	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:     "query <type> [modpack]",
	Aliases: []string{"q", "qry", "qr", "get", "g"},
	Short:   "Prints some value from the selected/default modpack's definition in grinch.kdl or modrinth.index.json",
	Long: `Prints some value from a modpack's definition. Supported query types:
- version: prints the modpack version as specified in modrinth.index.json

If no modpack is specified (from among those in grinch.kdl), the default modpack (either from grinch.kdl or the workspace config) will be used.`,

	Run: func(cmd *cobra.Command, args []string) {

		//COMMAND PARSING:
		genericArgErr := "Wrong arguments supplied"

		if len(args) < 1 {
			util.Hndl(errors.New("no query type specified"), genericArgErr, false)
		}

		if len(args) > 2 {
			util.Hndl(errors.New(fmt.Sprint(len(args))+" args is more than the maximum of 2\n[TIP: If you want to input a space without it being interpreted as an argument separator, prefix it with \\ or surround the entire piece of text containing space(s) with quotation marks]"), genericArgErr, false)
		}

		//DATA PARSING:
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
		if args[0] == "version" {
			fmt.Println(mri.Ver)
		} else {
			util.Hndl(errors.New("unknown query type "+args[0]), genericArgErr, false) //TODO: Add more query types, after re-doing some stuff in grinch.kdl
		}
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
}