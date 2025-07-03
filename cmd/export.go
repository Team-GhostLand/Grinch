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

	cp "github.com/otiai10/copy"
)

var exportCmd = &cobra.Command{
	Use:     "export [modpack]",
	Aliases: []string{"e", "exp", "xprt", "x"},
	Short:   "Exports your Grinch project as an Mrpack",
	//Long: `TODO`,
	Run: func(cmd *cobra.Command, args []string) {

		//--PARSING ARGS/FLAGS--
		mode_flag_parese_error := "Couldn't determine export mode"
		em_quick, err := cmd.Flags().GetBool("quick")
		if err != nil {
			util.Hndl(err, mode_flag_parese_error, false)
		}
		em_dev, err := cmd.Flags().GetBool("dev")
		if err != nil {
			util.Hndl(err, mode_flag_parese_error, false)
		}
		em_slim, err := cmd.Flags().GetBool("slim")
		if err != nil {
			util.Hndl(err, mode_flag_parese_error, false)
		}
		em_tweakable, err := cmd.Flags().GetBool("tweakable")
		if err != nil {
			util.Hndl(err, mode_flag_parese_error, false)
		}

		//There's probably a better way to do it, using Go's funky little string marker thingys (like when parsing KDL), but I don't know Go well enough to know how to use them
		em := trans.Default
		if em_quick {
			em = trans.Quick
		}
		if em_dev {
			em = trans.Dev
		}
		if em_slim {
			em = trans.Slim
		}
		if em_tweakable {
			em = trans.Tweakable
		}

		to, err := cmd.Flags().GetString("to")
		if err != nil {
			log.Println("WARN: Couldn't parse the --to flag - will act like it wasn't there")
		}

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
		err = trans.DoExportJsonTransforms(em)
		util.Hndl(err, "Couldn't execute the JSON transforms necessary by your export mode", true)

		file_transform_error := "Couldn't execute the file transforms necessary by your export mode"
		if em == trans.Default {
			err = trans.ResolveServerRemovals()
			util.Hndl(err, file_transform_error, true)
		}
		if em == trans.Dev {
			err = trans.SwapServerGitToDev()
			util.Hndl(err, file_transform_error, true)
		}

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

	exportCmd.Flags().Bool("default", false, "Default mode: Makes your modpack a fully standards-compliant Mrpack, that works flawlessly in both standards-compliant client launchers (that's important for later) and server-side Modrinth loaders. Replaces the non-standard (but much easier for us to operate on) server-overrides/REMOVALS.txt file with a proper client-overrides/ directory. Anything else (like what mods are optional, etc.) should alreadly be configured to be standards-compliant, if you used grinch import on an earlier-grinch-exported project (*if not, please refer to the docs to find out how to „Grinchify” an Mrpack). You don't need to set this flag explicitly - in fact, the code doesn't even check for its presence (apart from making sure it wasn't combined with other mode flags), due to the fact that this is the default behaviour. Cannot be combined with other mode flags.")
	exportCmd.Flags().BoolP("quick", "q", false, "Quick mode: Much quicker, as it doesn't do any file/JSON transforms whatsoever. Sufficient if you're exporting it for client-side only or if your server's Modrinth loader can understand the non-standard server-overrides/REMOVALS.txt file. Cannot be combined with other mode flags.")
	exportCmd.Flags().BoolP("dev", "d", false, "Dev mode: Transforms JSON and files in such a way that makes your pack suitable for development using Modrinth's official launcher. Replaces server-overrides/ with overrides/.SERVERSIDE (becasue otherwise Modrinth launcher would discard them, due to it being a client-side launcher) and marks any client-unsupported mod as client-required (same reason). These steps will later be reversed when running grinch import (though the latter requires you to have prefixes configured correctly, which Grinch won't help you with - thankfully, you only need to do this once and they'll be persisted across future imports/exports). Cannot be combined with other mode flags.")
	exportCmd.Flags().BoolP("slim", "s", false, "Slim mode: Marks every client-optional mod as client-unsupported. This is done to mimick the standards-compliant behaviour of „letting users choose whether they want to install optional mods” within the Modrinth launcher (which - hilariously - isn't compliant with Modrinth's own standard), by letting the users pick whether they want a „full” modpack experience or a „slim” one (if you provide both files, of course). Since this option only targets clients (or, one specific client, really), we don't do any extra transforms to ensure proper server support (ie. server-overrides/REMOVALS.txt stays). Cannot be combined with other mode flags.")
	exportCmd.Flags().BoolP("tweakable", "t", false, "Tweakable mode: Marks every client-optional mod as disabled. This is done to mimick the standards-compliant behaviour of „letting users choose whether they want to install optional mods” within the Modrinth launcher (which - hilariously - isn't compliant with Modrinth's own standard), by letting the users manually re-enable mods that they want. Since this option only targets clients (or, one specific client, really), we don't do any extra transforms to ensure proper server support (ie. server-overrides/REMOVALS.txt stays). Cannot be combined with other mode flags.")
	exportCmd.Flags().StringP("to", "T", "", "Renames the output file. Cannot be used in --dev mode, as that runs the risk of accidentially importing your Mrpack over a wrong Grinch project later down the line.")

	exportCmd.MarkFlagsMutuallyExclusive("to", "dev")
	exportCmd.MarkFlagsMutuallyExclusive("quick", "dev", "slim", "tweakable", "default")
}
