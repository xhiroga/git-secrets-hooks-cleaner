/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
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

type UseGitSecrets int

const (
	No UseGitSecrets = iota
	Modified
	TemplateCompatible
)

var commitMsgTemplate = `#!/usr/bin/env bash
git secrets --commit_msg_hook -- "$@"
`
var preCommitTemplate = `#!/usr/bin/env bash
git secrets --pre_commit_hook -- "$@"
`
var prepareCommitMsg = `#!/usr/bin/env bash
git secrets --prepare_commit_msg_hook -- "$@"
`

func isTemplateCompatible(path string, content string) (bool, error) {
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
func doesUseGitSecrets(path string) (UseGitSecrets, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return No, err
	}
	content := string(bytes)
	compatible, err := isTemplateCompatible(path, content)
	if err != nil {
		return No, err
	}
	if compatible {
		return TemplateCompatible, nil
	}
	if strings.Contains(content, "git secrets") {
		return Modified, nil
	}
	return No, nil
}

type hook struct {
	path          string
	useGitSecrets UseGitSecrets
}

func getHooks(baseDir string, sugar *zap.SugaredLogger) ([]hook, error) {
	var hooks []hook
	// filepath.Glob doesn't support **, so we need to walk the directory
	// https://stackoverflow.com/a/26809999/7869792
	err := filepath.WalkDir(baseDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Base(filepath.Dir(filepath.Dir(path))) == ".git" && filepath.Base(filepath.Dir(path)) == "hooks" {
			sugar.Debugf("Found file: %s", path)
			useGitSecrets, err := doesUseGitSecrets(path)
			if err != nil {
				return err
			}
			sugar.Debugf("Use git-secrets: %v", useGitSecrets)
			hooks = append(hooks, hook{path: path, useGitSecrets: useGitSecrets})
		}
		return nil
	})
	if err != nil {
		sugar.Fatalf("Error walking directory: %s", err)
		return nil, err
	}
	return hooks, nil
}

func showHooks(hooks []hook) {
	fmt.Println("Path\tUse git-secrets")
	for _, hook := range hooks {
		fmt.Printf("%s\t", hook.path)
		switch hook.useGitSecrets {
		case No:
			color.Green("No")
		case Modified:
			color.Yellow("Modified")
		case TemplateCompatible:
			color.Red("Template compatible")
		}
	}
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

		baseDir, _ := cmd.Flags().GetString("base")
		sugar.Debugf("Base dir: %s", baseDir)

		if baseDir == "" {
			currentDir, _ := os.Getwd()
			sugar.Debugf("Current dir: %s", currentDir)
			baseDir = currentDir
		}

		hooks, err := getHooks(baseDir, sugar)
		if err != nil {
			sugar.Fatalf("Error getting hooks: %s", err)
		}
		showHooks(hooks)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().String("base", "", "Base directory to search for hooks")
	listCmd.Flags().Bool("verbose", false, "Verbose output")
}
