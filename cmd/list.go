package cmd

import (
	"github.com/iliyankg/colab-shield/cli/client"
	"github.com/iliyankg/colab-shield/cli/config"
	"github.com/iliyankg/colab-shield/protos"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	pathToList string
	cursor     uint64
	pageSize   int64
)

func init() {
	listCmd.Flags().StringVarP(&pathToList, "path", "p", "", "Path to list files")
	listCmd.Flags().Uint64VarP(&cursor, "cursor", "c", 0, "Cursor for pagination")
	listCmd.Flags().Int64VarP(&pageSize, "page-size", "s", 10, "Page size for pagination")

	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all at specified path.",
	Long: `Lists all files the backend knows about at the specified path. Provides pagination.
	Empty path lists all files.
	Cursor of 0 means start from the beginning.
	PageSize is a "suggestion" and number of entries may be less or more than this. (Redis Quirk)`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("List called")

		payload := &protos.ListFilesRequest{
			FolderPath: pathToList,
			Cursor:     cursor,
			PageSize:   pageSize,
		}

		ctx, cancel := buildContext(config.ProjectId(), gitUser)
		defer cancel()
		conn, client := client.NewColabShieldClient(config.ServerHost(), config.ServerPortGrpc())
		defer conn.Close()

		response, err := client.ListFiles(ctx, payload)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to list files")
		}

		log.Info().Msgf("Next cursor: %d", response.NextCursor)
		log.Info().Msgf("File Infos:")
		for _, file := range response.Files {
			log.Info().Interface("fileInfo", file).Send()
		}
	},
}
