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

var commitMsgTemplate = `#!/usr/bin/env bash
git secrets --commit_msg_hook -- "$@"
`
var preCommitTemplate = `#!/usr/bin/env bash
git secrets --pre_commit_hook -- "$@"
`
var prepareCommitMsg = `#!/usr/bin/env bash
git secrets --prepare_commit_msg_hook -- "$@"
`

func isFromTemplate(path string) (bool, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	content := string(bytes)
	switch filepath.Base(path) {
	case "commit-msg":
		return content == commitMsgTemplate, nil
	case "pre-commit":
		return content == preCommitTemplate, nil
	case "prepare-commit-msg":
		return content == prepareCommitMsg, nil
	}
	return false, nil
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
				fromTemplate, err := isFromTemplate(path)
				if err != nil {
					return err
				}
				sugar.Debugf("From template: %v", fromTemplate)
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			sugar.Fatalf("Error walking directory: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().String("base-dir", "", "Base directory to search for hooks")
	listCmd.Flags().Bool("verbose", false, "Verbose output")
}
