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

var importCmd = &cobra.Command{
	Use:     "import [name]",
	Aliases: []string{"i", "imp"},
	Short:   "Imports an mrpack, that you previously exported in dev-mode, back as a Grinch project.",
	Long:    "Imports an mrpack, that you previously exported in --dev mode (and presumably made some changes), back as a Grinch project. If you want to use it to start an entirely new Grinch project based on an existing Mrpack - you can, but some preparations are needed beforehand (such as setting up prefixes, making sure to use our format for override layering, or (in case it's the 1st project in the repo) setting up project config). Please refer to the README on how to do said preparations.",
	Run: func(cmd *cobra.Command, args []string) {

		//--PARSING ARGS--
		if len(args) > 1 {
			util.Hndl(errors.New(fmt.Sprint(len(args))+" is more than the maximum of 1\n[TIP: If you want to input a space without it being interpreted as an argument separator, prefix it with \\ or surround the entire piece of text containing space(s) with quotation marks]"), "Too many arguments", false)
		}

		//--CONFIGS--
		pcf, err := util.LoadProjectConfig(util.GrProjectFileLocation)
		util.Hndl(err, "Couldn't load project config", false)
		wcf, err := util.LoadWorkspaceConfig(util.GrWorkspaceFileLocation)
		util.Hndl(err, "Couldn't load workspace config", false)

		var mp *util.ModpackDefinition
		mp, err = util.SelectModpack(pcf, wcf)
		util.Hndl(err, "Couldn't select modpack", false)

		folder_path := filepath.FromSlash(mp.Path)
		if folder_path == "" {
			util.Hndl(errors.New("chosen modpack "+mp.Name+" has no path associated with it"), "Couldn't select modpack", false)
		}

		mrpack_path, err := util.FindNewMrpack(args, ".", util.GrWorkspaceFileLocation, wcf) //Because this is an inherently machine-tied operation (.gr-workspace should be .gitignored - and if you manually inserted a path, when you're probably already doing it using your OS's path conventions), there's no need to do filepath.FromSlash()
		if mrpack_path == "" {
			if err == nil {
				err = errors.New("no undiscovered (ie. unlisted in .gr-workspace) mrpacks present in your working directory")
			}
			util.Hndl(err, "Couldn't find any mrpack to operate on", false)
		}
		if err != nil {
			log.Println("WARN: Found an mrpack to import (" + mrpack_path + "), but couldn't mark it as known - it may get accidentally picked up by another grinch import in the future")
		}

		//--UNZIP--
		_, err = util.IsSafelyCreatable(util.Tempdir)
		util.Hndl(err, "Cannot safely create a .temp directory", false)

		err = util.Unzip(mrpack_path, util.Tempdir)
		util.Hndl(err, "Cannot unzip "+mrpack_path+" to "+util.Tempdir, true) //Although there is no need to cleanup if unzip fails COMPLETELY, it might also fail partially (and leave a half-full .temp folder behind) - hence the true here

		//--TRANSFORMS--
		file_transform_error := "Couldn't execute the JSON transforms necessary for import" //JSON STUFF:
		mi, err := util.GetMrIndexJson(util.MrIndexFileLocation)
		util.Hndl(err, file_transform_error, true)

		util.DoPrefixSideSupportJsonTransforms(&mi, trans.ImportTransformPredicates, "GR_", false)
		util.DoPrefixSideSupportJsonTransforms(&mi, trans.ImportTransformPredicates, "GRd_", true)
		trans.ApplyJsonParamsOnImport(&mi, mp.Params)
		trans.SortMrIndexOnImport(&mi)

		err = util.SetMrIndexJson(mi, util.MrIndexFileLocation)
		util.Hndl(err, file_transform_error, true)

		err = trans.SwapServerDevToGit()
		util.Hndl(err, "Couldn't execute the file transforms necessary for import", true)

		//--BACKUP OLD PROJECT AND REPLACE IT WITH NEW ONE--
		_, err = util.IsSafelyCreatable(util.Backup)
		util.Hndl(err, "Cannot safely backup your previous "+folder_path+" as .old", true)

		err = os.Rename(folder_path, util.Backup)
		util.Hndl(err, "Couldn't backup your previous "+folder_path+" as .old", true)

		defer os.RemoveAll(util.Backup) //Defer won't run if the app crashed via a util.Hndl call, therefore it's safe to blindly defer the removal of backups, without checking whether said backups might actually come in handy.

		err = os.Rename(util.Tempdir, folder_path)
		util.Hndl(err, "Couldn't turn "+util.Tempdir+" into "+folder_path+", but otherwise completed the import - please rename it manually", false)

	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}