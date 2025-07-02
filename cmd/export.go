/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Team-GhostLand/Grinch/trans"
	"github.com/Team-GhostLand/Grinch/util"
	"github.com/spf13/cobra"

	cp "github.com/otiai10/copy"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:     "export [modpack]",
	Aliases: []string{"e", "exp", "xprt", "x"},
	Short:   "A brief description of your command",
	Long: `EXPORT: A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		//--PARSING ARGS/FLAGS--
		em := trans.Quick //TODO: Get this from args
		to := ""          //TODO: Get this from args

		if len(args) > 1 {
			util.Hndl(errors.New(fmt.Sprint(len(args))+" is more than the maximum of 1\n[TIP: If you want to input a space without it being interpreted as an argument separator, prefix it with \\ or surround the entire piece of text containing space(s) with quotation marks]"), "Too many arguments", false)
		}

		//--CONFIGS--
		pcf, err := util.LoadProjectConfig()
		util.Hndl(err, "Couldn't load config", false)

		var mp *util.ModpackDefinition
		if len(args) == 0 {
			mp, err = util.SelectModpack(pcf)
			util.Hndl(err, "Couldn't select modpack", false)
		} else {
			mp, err = util.FindModpackByName(pcf, args[0])
			util.Hndl(err, "Couldn't find modpack "+args[0], false)
		}

		path := filepath.FromSlash(mp.Path)
		if path == "" {
			util.Hndl(errors.New("chosen modpack "+mp.Name+" has no path associated with it"), "Couldn't select modpack", false)
		}

		exn := util.GetExportName(mp, to)

		//--TEMPDIR--
		_, err = util.IsSafelyCreateable(util.Tempdir)
		if err != nil {
			util.Hndl(err, "Cannot safely create a .temp directory", false)
		}

		err = cp.Copy(path, util.Tempdir)
		if err != nil {
			util.Hndl(err, "Cannot copy "+path+" to "+util.Tempdir, false)
		}

		defer os.RemoveAll(util.Tempdir)

		//--TRANSFORMS--
		if em == trans.Default {
			trans.ResolveServerRemovals()
		} else if em == trans.Dev {
			trans.DoExportJsonTransforms(em)
			trans.SwapServerGitToDev()
		} else if em != trans.Quick { //ie. slim or tweakable - where (and also in dev, but that was done above) we have to do some transforms
			trans.DoExportJsonTransforms(em)
		} //else: it's --quick mode, so we do nothing

		//--ZIP THAT BODYBAG UP--
		_, err = util.IsSafelyCreateable(exn)
		if err != nil {
			util.Hndl(err, "Cannot safely create "+exn, true)
		}
		err = util.MakeZipFile(util.Tempdir, exn)
		if err != nil {
			util.Hndl(err, "Cannot ZIP "+util.Tempdir+" to "+exn, true)
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	exportCmd.Flags().BoolP("quick", "q", false, "faster")
	exportCmd.Flags().BoolP("dev", "d", false, "dev mode")
	exportCmd.Flags().BoolP("slim", "s", false, "SE")
	exportCmd.Flags().BoolP("tweakable", "t", false, "SE+")
	exportCmd.Flags().StringP("to", "T", "", "output file")
}
