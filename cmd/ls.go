/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tom/ligma/internal/cache"
	"github.com/tom/ligma/internal/config"
	"github.com/tom/ligma/internal/spdx"
)

// popularIDs is the static set of "popular" SPDX IDs for --popular. Order preserved in output.
var popularIDs = []string{"MIT", "Apache-2.0", "GPL-2.0", "BSD-3-Clause", "ISC"}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:          "ls",
	Short:        "List all available SPDX licenses",
	Long:         `Fetch and list all available SPDX license identifiers from the official SPDX license list.`,
	SilenceUsage: true,
	SilenceErrors: true,
	RunE:         runLs,
}

// lsListURLOverride, when non-empty, is used instead of spdx.DefaultListURL (for tests).
var lsListURLOverride string

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().BoolP("json", "j", false, "output as JSON")
	lsCmd.Flags().String("filter", "", "case-insensitive filter on license ID or name")
	lsCmd.Flags().Bool("popular", false, "restrict to a popular set (MIT, Apache-2.0, GPL-2.0, BSD-3-Clause, ISC)")
}

func runLs(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	dir, err := config.LigmaDir()
	if err != nil {
		return err
	}
	cacheDir := filepath.Join(dir, "_cache")
	url := spdx.DefaultListURL
	if lsListURLOverride != "" {
		url = lsListURLOverride
	} else if cfg.SPDXListURL != "" {
		url = cfg.SPDXListURL
	}
	ctx := context.Background()
	list, err := cache.FetchList(ctx, cacheDir, cache.TTL(cfg.CacheTTL), url)
	if err != nil {
		return fmt.Errorf("%w: failed to fetch license list: %v", ErrIOOrNetwork, err)
	}
	// Apply --popular first, then --filter (both before output formatting).
	if ok, _ := cmd.Flags().GetBool("popular"); ok {
		set := make(map[string]bool)
		for _, id := range popularIDs {
			set[id] = true
		}
		n := 0
		for _, l := range list {
			if set[l.LicenseID] {
				list[n] = l
				n++
			}
		}
		list = list[:n]
	}
	if term, _ := cmd.Flags().GetString("filter"); term != "" {
		term = strings.ToLower(term)
		n := 0
		for _, l := range list {
			if strings.Contains(strings.ToLower(l.LicenseID), term) || strings.Contains(strings.ToLower(l.Name), term) {
				list[n] = l
				n++
			}
		}
		list = list[:n]
	}
	useJSON, _ := cmd.Flags().GetBool("json")
	if useJSON {
		enc := json.NewEncoder(os.Stdout)
		if err := enc.Encode(list); err != nil {
			return fmt.Errorf("%w: failed to encode JSON: %v", ErrIOOrNetwork, err)
		}
		return nil
	}
	for _, l := range list {
		fmt.Println(l.LicenseID)
	}
	return nil
}
