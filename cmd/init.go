package cmd

import (
	"context"
	"time"

	"github.com/iliyankg/colab-shield/cli/client"
	"github.com/iliyankg/colab-shield/cli/config"
	"github.com/iliyankg/colab-shield/protos"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes colabshield in the current directory.",
	Long:  `Initializes colabshield in the current directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		conn, client := client.NewColabShieldClient(config.ServerHost(), config.ServerPortGrpc())
		defer conn.Close()

		payload := &protos.InitProjectRequest{
			ProjectId: gitBranch,
		}

		response, err := client.InitProject(ctx, payload)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to init project")
		}

		if response.Status != protos.Status_OK {
			log.Fatal().Msg("status not OK")
		}
	},
}
