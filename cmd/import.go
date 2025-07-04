package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Team-GhostLand/Grinch/trans"
	"github.com/Team-GhostLand/Grinch/util"

	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:     "import [name]",
	Aliases: []string{"i", "imp"},
	Short:   "Imports an mrpack, that you previously exported in --dev mode (and presumably made some changes), back as a Grinch project. If you want to use it to start a new Grinch project - you can, but some preparations are needed beforehand. Docs coming soon.",
	//Long: `TODO`,
	Run: func(cmd *cobra.Command, args []string) {

		//--PARSING ARGS--
		if len(args) > 1 {
			util.Hndl(errors.New(fmt.Sprint(len(args))+" is more than the maximum of 1\n[TIP: If you want to input a space without it being interpreted as an argument separator, prefix it with \\ or surround the entire piece of text containing space(s) with quotation marks]"), "Too many arguments", false)
		}

		//--CONFIGS--
		pcf, err := util.LoadProjectConfig()
		util.Hndl(err, "Couldn't load config", false)

		var mp *util.ModpackDefinition
		mp, err = util.SelectModpack(pcf)
		util.Hndl(err, "Couldn't select modpack", false)

		folder_path := filepath.FromSlash(mp.Path)
		if folder_path == "" {
			util.Hndl(errors.New("chosen modpack "+mp.Name+" has no path associated with it"), "Couldn't select modpack", false)
		}

		mrpack_path, err := util.FindNewMrpack(args)
		if mrpack_path == "" {
			if err == nil {
				err = errors.New("no undiscovered (ie. unlisted in .gr-workspace) mrpacks present in your working directory")
			}
			util.Hndl(err, "Couldn't find any mrpack to operate on", false)
		}
		if err != nil {
			log.Println("WARN: Found an mrpack to import (" + mrpack_path + "), but couldn't mark it as known - it may get accidentially picked up by another grinch import in the future")
		}

		//--UNZIP--
		_, err = util.IsSafelyCreateable(util.Tempdir)
		if err != nil {
			util.Hndl(err, "Cannot safely create a .temp directory", false)
		}

		err = util.Unzip(mrpack_path, util.Tempdir)
		if err != nil {
			util.Hndl(err, "Cannot unzip "+mrpack_path+" to "+util.Tempdir, true) //Although there is no need to cleanup if unzip fails COMPLETELY, it might also fail partially (and leave a half-full .temp folder behind) - hence the true here
		}

		//--TRANSFORMS--
		err = trans.DoImportJsonTransforms()
		if err != nil {
			util.Hndl(err, "Couldn't execute the JSON transforms necessary for import", true)
		}
		err = trans.SwapServerDevToGit()
		if err != nil {
			util.Hndl(err, "Couldn't execute the file transforms necessary for import", true)
		}
		//TODO: Constraints

		//--TURN TEMP INTO THE INTENDED DIR--
		err = os.RemoveAll(folder_path)
		if err != nil {
			util.Hndl(err, "Couldn't remove your old "+folder_path, true)
		}
		err = os.Rename(util.Tempdir, folder_path)
		if err != nil {
			util.Hndl(err, "Couldn't turn "+util.Tempdir+" into "+folder_path, false) //We don't want to remove the .temp if this happens, so that the users can rename it themselves - hence the false here
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
