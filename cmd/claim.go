package cmd

import (
	"github.com/iliyankg/colab-shield/cli/client"
	"github.com/iliyankg/colab-shield/cli/config"
	"github.com/iliyankg/colab-shield/cli/gitutils"
	"github.com/iliyankg/colab-shield/protos"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	claimMode int32
	softClaim bool
)

func init() {
	claimFilesCmd.MarkFlagRequired("file")

	claimFilesCmd.Flags().Int32VarP(&claimMode, "mode", "m", int32(protos.ClaimMode_EXCLUSIVE), "claim mode")
	claimFilesCmd.Flags().BoolVarP(&softClaim, "soft-claim", "s", false, "Soft claim only exposes any files that may get rejected if any. Nothing is saved to the DB")

	rootCmd.AddCommand(claimFilesCmd)
}

var claimFilesCmd = &cobra.Command{
	Use:   "claim",
	Short: "Claim file(s) for editing",
	Long:  `Claim file(s) for editing`,
	Run: func(cmd *cobra.Command, args []string) {
		if !validateClaimMode(claimMode) {
			log.Fatal().Msg("Invalid claim mode. Must be 1 (EXCLUSIVE) or 2 (SHARED).")
		}

		filteredFiles, err := filterToFilesOfInterest(files)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to filter to files of interest")
		}

		hashes, err := gitutils.GetGitBlobHEADHashes(log.Logger, filteredFiles)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get git hashes")
		}

		log.Info().Msgf("Git hash for files %s: %s", filteredFiles, hashes)

		payload, err := newClaimFilesRequest(filteredFiles, hashes, protos.ClaimMode(claimMode), softClaim)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to map files to hash")
		}

		ctx, cancel := buildContext(config.ProjectId(), gitUser)
		defer cancel()
		conn, client := client.NewColabShieldClient(config.ServerHost(), config.ServerPortGrpc())
		defer conn.Close()

		response, err := client.Claim(ctx, payload)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to lock files")
		}

		for _, file := range response.RejectedFiles {
			log.Warn().Msgf("Rejected - %s - %s", file.FileId, file.RejectReason.String())
		}

		if response.Status != protos.Status_OK {
			log.Fatal().Msg("status not OK")
		}
	},
}

// newClaimFilesRequest creates a new ClaimFilesRequest from the given files and hashes
// Same claim mode is applied to all files.
// TODO: Per file claim mode support?
func newClaimFilesRequest(files []string, hashes []string, claimMode protos.ClaimMode, softClaim bool) (*protos.ClaimFilesRequest, error) {
	if len(files) != len(hashes) {
		return nil, ErrFileToHashMissmatch
	}

	claimFileInfos := make([]*protos.ClaimFileInfo, 0, len(files))
	for i, file := range files {
		claimFileInfos = append(claimFileInfos, &protos.ClaimFileInfo{
			FileId:    file,
			FileHash:  hashes[i],
			ClaimMode: claimMode,
		})
	}

	return &protos.ClaimFilesRequest{
		BranchName: gitBranch,
		Files:      claimFileInfos,
		SoftClaim:  softClaim,
	}, nil
}

// validateClaimMode validates the incoming integer is valid.
// FIXME: This ties is to protobuf and not ideal.
func validateClaimMode(mode int32) bool {
	switch mode {
	case int32(protos.ClaimMode_EXCLUSIVE):
	case int32(protos.ClaimMode_SHARED):
		return true
	default:
		return false
	}

	return false
}
