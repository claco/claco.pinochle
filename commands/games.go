package commands

import (
	"context"
	"errors"
	"github.com/claco/claco.pinochle/pb"
	log "github.com/sirupsen/logrus"
	"io"
)

type GamesArgs struct {
	AddressArgs

	CreateArgs *CreateArgs `arg:"subcommand:create" help:"create a new game"`
	ListArgs   *ListArgs   `arg:"subcommand:list" help:"list all games"`
}

type CreateArgs struct {
	Name string `arg:"positional:name" help:"name for new game"`
}

type ListArgs struct{}

func (gargs *GamesArgs) Execute() (int, error) {
	if gargs.CreateArgs != nil {
		return gargs.Create(gargs.CreateArgs)
	} else if gargs.ListArgs != nil {
		return gargs.List(gargs.ListArgs)
	}

	return 127, errors.New("command not implemented: games")
}

func (gargs *GamesArgs) Create(args *CreateArgs) (int, error) {
	request := &pb.CreateGameRequest{Name: args.Name}

	client, err := gargs.NewClient()
	if err != nil {
		return 1, err
	} else {
		log.WithField("request", request).Debugf("making request: %s", request.ProtoReflect().Descriptor().FullName())
	}

	response, err := client.CreateGame(context.Background(), request)
	if err != nil {
		return 1, err
	} else {
		log.WithField("response", response).Debugf("received response: %s", response.ProtoReflect().Descriptor().FullName())
	}

	return 0, nil
}

func (gargs *GamesArgs) List(args *ListArgs) (int, error) {
	request := &pb.ListGamesRequest{}

	client, err := gargs.NewClient()
	if err != nil {
		return 1, err
	} else {
		log.WithField("request", request).Debugf("making request: %s", request.ProtoReflect().Descriptor().FullName())
	}

	stream, err := client.ListGames(context.Background(), request)
	if err != nil {
		return 1, err
	} else {
		for {
			game, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return 1, err
			}

			log.Info(game)
		}
	}

	return 0, nil
}
