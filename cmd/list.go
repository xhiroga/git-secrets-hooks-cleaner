/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
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
		fmt.Println(verbose)
		logger := getLogger(verbose)
		defer logger.Sync() // flushes buffer, if any
		sugar := logger.Sugar()

		baseDir, _ := cmd.Flags().GetString("base-dir")
		sugar.Debugf("Base dir: %s", baseDir)

		if baseDir == "" {
			currentDir, _ := os.Getwd()
			sugar.Debugf("Current dir: %s", currentDir)
			baseDir = currentDir
		}

		matches, _ := filepath.Glob(baseDir + "/**/.git/hooks/*")
		sugar.Debugf("Found %d hooks", len(matches))
		sugar.Debug(matches)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().String("base-dir", "", "Base directory to search for hooks")
	listCmd.Flags().Bool("verbose", false, "Verbose output")
}
