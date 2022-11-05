/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func getLogger(verbose bool) *zap.Logger {
	if verbose {
		logger, _ := zap.NewDevelopment()
		return logger
	}
	logger, _ := zap.NewProduction()
	return logger
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		logger := getLogger(verbose)
		defer logger.Sync() // flushes buffer, if any
		sugar := logger.Sugar()
		sugar.Debugf("Verbose: %v", verbose)

		baseDir, _ := cmd.Flags().GetString("base-dir")
		sugar.Debugf("Base dir: %s", baseDir)

		if baseDir == "" {
			currentDir, _ := os.Getwd()
			sugar.Debugf("Current dir: %s", currentDir)
			baseDir = currentDir
		}

		// filepath.Glob doesn't support **, so we need to walk the directory
		// https://stackoverflow.com/a/26809999/7869792

		var files []string
		err := filepath.WalkDir(baseDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if filepath.Base(filepath.Dir(path)) == "hooks" {
				sugar.Debugf("Found file: %s", path)
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			sugar.Fatalf("Error walking directory: %s", err)
		}

		// matches, _ := filepath.Glob(baseDir + "/**/.git/hooks/*")
		// sugar.Debugf("Found %d hooks", len(matches))
		// sugar.Debug(matches)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().String("base-dir", "", "Base directory to search for hooks")
	listCmd.Flags().Bool("verbose", false, "Verbose output")
}
