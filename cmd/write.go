/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tom/ligma/internal/config"
	"github.com/tom/ligma/internal/spdx"
)

// writeFetchDetails is the details fetcher; swappable for tests.
var writeFetchDetails = spdx.FetchLicenseDetails

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Write license text to a file by SPDX ID",
	Long:  `Fetch the license for the given SPDX ID and write it to LICENSE (one arg) or to the given path (two args). With no arguments, writes the configured favorite to LICENSE in the current directory. Overwrites if the file exists.`,
	Args:  cobra.RangeArgs(0, 2),
	SilenceUsage: true,
	SilenceErrors: true,
	RunE:  runWrite,
}

func runWrite(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	template := spdx.DefaultDetailsURLTemplate
	if cfg.SPDXGetURLTemplate != "" {
		template = cfg.SPDXGetURLTemplate
	}

	var id, path string
	if len(args) == 0 {
		if cfg.Favorite == nil || *cfg.Favorite == "" {
			return fmt.Errorf("favorite license is not set; run with <id> or set favorite in ~/.ligma/config.json")
		}
		id = cfg.Resolve(*cfg.Favorite)
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get working directory: %w", err)
		}
		path = filepath.Join(cwd, "LICENSE")
	} else {
		id = cfg.Resolve(args[0])
		if len(args) == 2 {
			path = args[1]
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("get working directory: %w", err)
			}
			path = filepath.Join(cwd, "LICENSE")
		}
	}

	text, err := writeFetchDetails(cmd.Context(), template, id)
	if err != nil {
		if errors.Is(err, spdx.ErrNotFound) {
			return fmt.Errorf("license not found: %s: %w", id, ErrNotFound)
		}
		return fmt.Errorf("fetch license %s: %v: %w", id, err, ErrIOOrNetwork)
	}

	if err := os.WriteFile(path, []byte(text), 0644); err != nil {
		return fmt.Errorf("write %s: %v: %w", path, err, ErrIOOrNetwork)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(writeCmd)
}
