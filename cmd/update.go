package cmd

import (
	"github.com/iliyankg/colab-shield/cli/client"
	"github.com/iliyankg/colab-shield/cli/config"
	"github.com/iliyankg/colab-shield/cli/gitutils"
	"github.com/iliyankg/colab-shield/protos"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update claimed files with changes",
	Long:  `Update claimed files with changes`,
	Run: func(cmd *cobra.Command, args []string) {
		files, err := getFilesOfInterestFromStaged()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get git staged files")
		}

		hashes, err := gitutils.GetGitBlobHashes(log.Logger, files)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get git hashes")
		}

		headHashes, err := gitutils.GetGitBlobHEADHashes(log.Logger, files)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get git HEAD hashes")
		}

		payload, err := newUpdateFilesRequest(files, hashes, headHashes)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to map files to hash")
		}

		ctx, cancel := buildContext(config.ProjectId(), gitUser)
		defer cancel()
		conn, client := client.NewColabShieldClient(config.ServerHost(), config.ServerPortGrpc())
		defer conn.Close()

		response, err := client.Update(ctx, payload)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to update files")
		}

		if response.Status != protos.Status_OK {
			log.Fatal().Msg("Failed to update files")
		}
	},
}

func getFilesOfInterestFromStaged() ([]string, error) {
	allStagedfiles, err := gitutils.GetGitStagedFiles(log.Logger)
	if err != nil {
		return nil, err
	}

	return filterToFilesOfInterest(allStagedfiles)
}

func newUpdateFilesRequest(files []string, hashes []string, headHashes []string) (*protos.UpdateFilesRequest, error) {
	if len(files) != len(hashes) || len(files) != len(headHashes) {
		return nil, ErrFileToHashMissmatch
	}

	updateFileInfos := make([]*protos.UpdateFileInfo, 0, len(files))
	for i, file := range files {
		updateFileInfos = append(updateFileInfos, &protos.UpdateFileInfo{
			FileId:   file,
			FileHash: hashes[i],
			OldHash:  headHashes[i],
		})
	}

	return &protos.UpdateFilesRequest{
		BranchName: gitBranch,
		Files:      updateFileInfos,
	}, nil
}
