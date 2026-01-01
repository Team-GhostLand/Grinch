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
	Long: `Prints some value from a modpack's definition.

Supported query types:
- path: prints the path to where the modpack's source is located (as specified in grinch.kdl)
- name: prints the modpack's name when exporting in Default Mode (as specified in grinch.kdl), with templates resolved
- name_quick: prints the modpack's name when exporting in Quick Mode (as specified in grinch.kdl), with templates resolved
- name_dev: prints the modpack's name when exporting in Dev Mode (as specified in grinch.kdl), with templates resolved
- name_slim: prints the modpack's name when exporting in Slim Mode (as specified in grinch.kdl), with templates resolved
- name_tweakable: prints the modpack's name when exporting in Tweakable Mode (as specified in grinch.kdl), with templates resolved
- rawname: prints the modpack's name when exporting in Default Mode (as specified in grinch.kdl), without template resolution
- rawname_quick: prints the modpack's name when exporting in Quick Mode (as specified in grinch.kdl), without template resolution
- rawname_dev: prints the modpack's name when exporting in Dev Mode (as specified in grinch.kdl), without template resolution
- rawname_slim: prints the modpack's name when exporting in Slim Mode (as specified in grinch.kdl), without template resolution
- rawname_tweakable: prints the modpack's name when exporting in Tweakable Mode (as specified in grinch.kdl), without template resolution
- version: prints the modpack version as specified in modrinth.index.json
- description: prints the modpack's description (as specified in grinch.kdl)
- gitname: prints the modpack's gitname (aka „The one that gets pushed into modrinth.index.json during import, so that it ends up on Git later”) from grinch.kdl
- mcversion: prints the Minecraft version as specified in modrinth.index.json
- name_mr: prints the modpack's name from modrinth.index.json (which SHOULD be the same as gitname in grinch.kdl (if Grinch is operating normally), so this query type is mostly for automated tests of Grinch itself)
- description_mr: prints the modpack's description (as specified in modrinth.index.json - which SHOULD be the same as that in grinch.kdl (if Grinch is operating normally), so this query type is mostly for automated tests of Grinch itself)

The modpack to query (from among those in grinch.kdl) can be specified as the second argument.
If no modpack is specified, the default modpack (either from grinch.kdl or the workspace config) will be used.`,

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
		if len(args) == 1 {
			mp, err = util.SelectModpack(pcf, wcf)
			util.Hndl(err, "Couldn't select modpack", false)
		} else {
			mp, err = util.FindModpackByName(pcf, args[1])
			util.Hndl(err, "Couldn't find modpack "+args[1], false)
		}

		path := filepath.FromSlash(mp.Path)
		if path == "" {
			util.Hndl(errors.New("chosen modpack "+mp.Name+" has no path associated with it"), "Couldn't select modpack", false)
		}

		mri, err := util.GetMrIndexJson(filepath.FromSlash(path + "/modrinth.index.json"))
		util.Hndl(err, "Couldn't read Modrinth Index:", false)

		//END RESULT:
		switch args[0] {
		case "path":
			fmt.Println(mp.Path)
		case "name":
			fmt.Println(util.ResolveTemplateString(&mri, mp.Params.Names.Default))
		case "name_quick":
			fmt.Println(util.ResolveTemplateString(&mri, mp.Params.Names.Quick))
		case "name_dev":
			fmt.Println(util.ResolveTemplateString(&mri, mp.Params.Names.Dev))
		case "name_slim":
			fmt.Println(util.ResolveTemplateString(&mri, mp.Params.Names.Slim))
		case "name_tweakable":
			fmt.Println(util.ResolveTemplateString(&mri, mp.Params.Names.Tweakable))
		case "rawname":
			fmt.Println(mp.Params.Names.Default)
		case "rawname_quick":
			fmt.Println(mp.Params.Names.Quick)
		case "rawname_dev":
			fmt.Println(mp.Params.Names.Dev)
		case "rawname_slim":
			fmt.Println(mp.Params.Names.Slim)
		case "rawname_tweakable":
			fmt.Println(mp.Params.Names.Tweakable)
		case "version":
			fmt.Println(mri.Ver)
		case "description":
			fmt.Println(mp.Params.Description)
		case "gitname":
			fmt.Println(mp.Params.Names.Git)
		case "mcversion":
			fmt.Println(mri.Deps["minecraft"])
		case "name_mr":
			fmt.Println(mri.Name)
		case "description_mr":
			fmt.Println(mri.Desc)
		default:
			util.Hndl(errors.New("unknown query type "+args[0]), genericArgErr, false) //TODO: Add more query types, after re-doing some stuff in grinch.kdl
		}
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
}