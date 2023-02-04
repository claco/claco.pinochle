package commands

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/alexflint/go-arg"
	"github.com/claco/claco.pinochle/pb"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type GamesArgs struct {
	AddressArgs

	CreateArgs *CreateArgs `arg:"subcommand:create" help:"create a new game"`
	ListArgs   *ListArgs   `arg:"subcommand:list" help:"list all games"`
	ResumeArgs *ResumeArgs `arg:"subcommand:resume" help:"resume an existing game"`
	ShowArgs   *ShowArgs   `arg:"subcommand:show" help:"display game details"`
	StartArgs  *StartArgs  `arg:"subcommand:start" help:"start an existing game"`
}

type CreateArgs struct {
	Name string `arg:"positional,required:name" help:"name for new game"`
}

type ListArgs struct {
}

type ResumeArgs struct {
	KeyArgs
}

type ShowArgs struct {
	KeyArgs
}

type StartArgs struct {
	KeyArgs
}

func (gargs *GamesArgs) Execute(parser *arg.Parser) (int, error) {
	if gargs.CreateArgs != nil {
		return gargs.Create(gargs.CreateArgs)
	} else if gargs.ListArgs != nil {
		return gargs.List(gargs.ListArgs)
	} else if gargs.ResumeArgs != nil {
		return gargs.Resume(gargs.ResumeArgs)
	} else if gargs.ShowArgs != nil {
		return gargs.Show(gargs.ShowArgs)
	} else if gargs.StartArgs != nil {
		return gargs.Start(gargs.StartArgs)
	}

	usage := bytes.NewBufferString(" ")
	parser.WriteUsage(usage)

	log.Warn(usage.String())

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

func (gargs *GamesArgs) Resume(args *ResumeArgs) (int, error) {
	request := &pb.ResumeGameRequest{}

	id, err := uuid.Parse(args.Id)
	if err == nil {
		request.Id = id.String()
	} else {
		request.Slug = args.Id
	}

	client, err := gargs.NewClient()
	if err != nil {
		return 1, err
	} else {
		log.WithField("request", request).Debugf("making request: %s", request.ProtoReflect().Descriptor().FullName())
	}

	response, err := client.ResumeGame(context.Background(), request)
	if err != nil {
		return 1, err
	} else {
		log.WithField("response", response).Debugf("received response: %s", response.ProtoReflect().Descriptor().FullName())

		log.Info(response.Game)
	}

	return 0, nil
}

func (gargs *GamesArgs) Show(args *ShowArgs) (int, error) {
	request := &pb.GetGameRequest{}

	id, err := uuid.Parse(args.Id)
	if err == nil {
		request.Id = id.String()
	} else {
		request.Slug = args.Id
	}

	client, err := gargs.NewClient()
	if err != nil {
		return 1, err
	} else {
		log.WithField("request", request).Debugf("making request: %s", request.ProtoReflect().Descriptor().FullName())
	}

	response, err := client.GetGame(context.Background(), request)
	if err != nil {
		return 1, err
	} else {
		log.WithField("response", response).Debugf("received response: %s", response.ProtoReflect().Descriptor().FullName())

		log.Info(response.Game)
	}

	return 0, nil
}

func (gargs *GamesArgs) Start(args *StartArgs) (int, error) {
	request := &pb.StartGameRequest{}

	id, err := uuid.Parse(args.Id)
	if err == nil {
		request.Id = id.String()
	} else {
		request.Slug = args.Id
	}

	client, err := gargs.NewClient()
	if err != nil {
		return 1, err
	} else {
		log.WithField("request", request).Debugf("making request: %s", request.ProtoReflect().Descriptor().FullName())
	}

	response, err := client.StartGame(context.Background(), request)
	if err != nil {
		return 1, err
	} else {
		log.WithField("response", response).Debugf("received response: %s", response.ProtoReflect().Descriptor().FullName())

		log.Info(response.Game)
	}

	return 0, nil
}
