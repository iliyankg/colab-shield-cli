package cmd

import (
	"github.com/iliyankg/colab-shield/cli/client"
	"github.com/iliyankg/colab-shield/cli/config"
	"github.com/iliyankg/colab-shield/protos"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	claimFilesCmd.MarkFlagRequired("file")

	rootCmd.AddCommand(releaseCmd)
}

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Release file(s) previously claimed.",
	Long:  `Release file(s) previously claimed.`,
	Run: func(cmd *cobra.Command, args []string) {
		filteredFiles, err := filterToFilesOfInterest(files)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to filter to files of interest")
		}

		payload := newReleaseFilesRequest(filteredFiles)

		ctx, cancel := buildContext(config.ProjectId(), gitUser)
		defer cancel()
		conn, client := client.NewColabShieldClient(config.ServerHost(), config.ServerPortGrpc())
		defer conn.Close()

		response, err := client.Release(ctx, payload)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to release files")
		}

		for _, file := range response.RejectedFiles {
			log.Info().Msgf("Rejected - %s - %s", file.FileId, file.RejectReason.String())
		}

		if response.Status != protos.Status_OK {
			log.Fatal().Msg("Failed to release files")
		}
	},
}

// newReleaseFilesRequest creates a new ReleaseFilesRequest from the given files
func newReleaseFilesRequest(filesToRelease []string) *protos.ReleaseFilesRequest {
	return &protos.ReleaseFilesRequest{
		BranchName: gitBranch,
		FileIds:    filesToRelease,
	}
}
