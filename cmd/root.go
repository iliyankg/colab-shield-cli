package cmd

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	gitUser   string
	gitBranch string
	files     []string
	rootCmd   = &cobra.Command{
		Use:   "colab-shield",
		Short: "A CLI tool for colaborative work with hard to merge files.",
		Long:  `A CLI tool for colaborative work with hard to merge files. It does this by providing an interface to a backend server which tracks file changes and versions.`,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&gitUser, "gitUser", "u", "", "git user name")
	rootCmd.MarkFlagRequired("gitUser")

	rootCmd.PersistentFlags().StringVarP(&gitBranch, "gitBranch", "b", "", "git branch")
	rootCmd.MarkFlagRequired("gitBranch")

	claimFilesCmd.Flags().StringArrayVarP(&files, "file", "f", []string{}, "files to lock")
}

// Execute executes the root command.
func Execute() error {
	// Ensure context is root of a git repository
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		log.Error().Msg("Make sure you are in the root of a git repository!")
		log.Fatal().Msg(".git folder does not exist")
	}

	return rootCmd.Execute()
}
