/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ligma",
	Short: "List, view, and write SPDX licenses.",
	Long:  `A small CLI to list, view, and write SPDX licenses. Data comes from the official SPDX website.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute runs the root command and returns the exit code (0–3). main calls os.Exit(Execute()).
// All process exit is via os.Exit in main only; subcommands and internal/ must return errors.
func Execute() int {
	err := rootCmd.Execute()
	if err == nil {
		return 0
	}
	fmt.Fprintln(os.Stderr, err.Error())
	return exitCodeFrom(err)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.licensegen.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
