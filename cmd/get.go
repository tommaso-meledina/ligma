/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tom/ligma/internal/cache"
	"github.com/tom/ligma/internal/config"
	"github.com/tom/ligma/internal/spdx"
)

var getSimulateIO bool

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:           "get",
	Short:          "Output license text by SPDX ID",
	Long:           `Fetch and print the full license text for the given SPDX license ID.`,
	Args:           cobra.ExactArgs(1),
	SilenceUsage:   true,
	SilenceErrors:  true,
	RunE:           runGet,
}

func runGet(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if getSimulateIO {
		return fmt.Errorf("simulated I/O error: %w", ErrIOOrNetwork)
	}
	dir, err := config.LigmaDir()
	if err != nil {
		return err
	}
	cacheDir := filepath.Join(dir, "_cache")
	template := spdx.DefaultDetailsURLTemplate
	if cfg.SPDXGetURLTemplate != "" {
		template = cfg.SPDXGetURLTemplate
	}
	id := cfg.Resolve(args[0])
	text, err := cache.FetchDetails(cmd.Context(), cacheDir, cache.TTL(cfg.CacheTTL), template, id)
	if err != nil {
		if errors.Is(err, spdx.ErrNotFound) {
			return fmt.Errorf("license not found: %s: %w", id, ErrNotFound)
		}
		return fmt.Errorf("fetch license %s: %v: %w", id, err, ErrIOOrNetwork)
	}
	useJSON, _ := cmd.Flags().GetBool("json")
	if useJSON {
		out := struct {
			ID          string `json:"id"`
			LicenseText string `json:"licenseText"`
		}{ID: id, LicenseText: text}
		b, err := json.Marshal(out)
		if err != nil {
			return fmt.Errorf("%w: failed to encode JSON: %v", ErrIOOrNetwork, err)
		}
		_, _ = os.Stdout.Write(b)
		return nil
	}
	_, _ = os.Stdout.WriteString(text)
	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().BoolP("json", "j", false, "output as JSON")
	getCmd.Flags().BoolVar(&getSimulateIO, "simulate-io-error", false, "simulate I/O or network error (dev)")
	_ = getCmd.Flags().MarkHidden("simulate-io-error")
}
