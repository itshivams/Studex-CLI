package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "studex-cli",
	Short: "Studex CLI",
	Long:  `A command line interface for Studex platform to search users, login, and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		StartInteractiveMode()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Studex CLI",
	Long:  `All software has versions. This is Studex CLI's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("studex-cli version 1.0.0")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Print version")
	
	rootCmd.AddCommand(versionCmd)
	
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		v, _ := cmd.Flags().GetBool("version")
		if v {
			fmt.Println("studex-cli version 1.0.0")
			return
		}
		StartInteractiveMode()
	}
}
