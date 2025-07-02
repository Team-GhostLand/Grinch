/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Team-GhostLand/Grinch/util"
	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:     "export",
	Aliases: []string{"e", "exp", "xprt", "x"},
	Short:   "A brief description of your command",
	Long: `EXPORT: A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		util.Hndl(err, "Couldn't open working directory")
		fmt.Println("export called from " + dir)

		cf, err := util.LoadProjectConfig()
		util.Hndl(err, "Couldn't load config")

		util.SelectModpack(cf)
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
	exportCmd.Flags().StringP("clone", "c", "", "clones")
	exportCmd.Flags().BoolP("quick", "q", false, "faster")
	exportCmd.Flags().BoolP("dev", "d", false, "dev mode")
	exportCmd.Flags().BoolP("slim", "s", false, "SE")
	exportCmd.Flags().BoolP("tweakable", "t", false, "SE+")
	exportCmd.Flags().BoolP("remove", "r", false, "cleanup")
}
