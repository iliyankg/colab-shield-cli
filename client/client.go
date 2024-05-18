package client

import (
	"fmt"

	"github.com/iliyankg/colab-shield/protos"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewColabShieldClient(host string, port int) (*grpc.ClientConn, protos.ColabShieldClient) {
	// FIXME: Fix the insecure.NewCredentials() call
	conn, err := grpc.Dial(fmt.Sprintf("%v:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect")
	}

	return conn, protos.NewColabShieldClient(conn)
}
