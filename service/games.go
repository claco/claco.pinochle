package service

import (
	"context"
	"github.com/claco/claco.pinochle/pb"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

func (svc *PinichleService) CreateGame(ctx context.Context, request *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	board := &pb.Board{Id: uuid.NewString()}
	game := &pb.Game{Id: uuid.NewString(), Board: board, Name: request.Name, Slug: svc.Slugify(request.Name)}
	response := &pb.CreateGameResponse{}

	log.WithField("game", game).Debugf("creating game: %s/%s", game.Id, game.Slug)

	batch := &pgx.Batch{}
	batch.Queue("INSERT INTO boards (id) VALUES ($1)", board.Id)
	batch.Queue("INSERT INTO games (id, slug, name, board_id, status) VALUES ($1, $2, $3, $4, $5)", game.Id, game.Slug, game.Name, game.Board.Id, game.Status)

	results := svc.db.SendBatch(ctx, batch)
	_, err := results.Exec()

	defer results.Close()

	if err != nil {
		return response, err
	} else {
		response.Game = game
	}

	return response, nil
}

func (svc *PinichleService) GetGame(ctx context.Context, request *pb.GetGameRequest) (*pb.GetGameResponse, error) {
	key := svc.GetKey(request)
	response := &pb.GetGameResponse{}

	log.WithField("id", key).Debugf("getting game: %s", key)

	row := svc.db.QueryRow(context.Background(), "SELECT id, name, slug, board_id, status from games where id=? or slug=?", key)

	var id, name, slug, board_id string
	var status uint8
	var game = &pb.Game{}

	log.WithField("game", game).Debugf("sending game: %s/%s", game.Id, game.Slug)

	err := row.Scan(&id, &name, &slug, &board_id, &status)

	if err != nil {
		return response, err
	} else {
		game.Id = id
		game.Name = name
		game.Slug = slug
		game.Board = &pb.Board{Id: board_id}
		game.Status = pb.GameStatus(status)

		response.Game = game
	}

	return response, nil
}

func (svc *PinichleService) ListGames(request *pb.ListGamesRequest, stream pb.PinochleService_ListGamesServer) error {
	rows, err := svc.db.Query(context.Background(), "SELECT id, name, slug, board_id, status from games")

	if err != nil {
		return err
	} else {
		defer rows.Close()
	}

	for rows.Next() {
		var id, name, slug, board_id string
		var status uint8
		var game = &pb.Game{}

		rows.Scan(&id, &name, &slug, &board_id, &status)
		game.Id = id
		game.Name = name
		game.Slug = slug
		game.Board = &pb.Board{Id: board_id}
		game.Status = pb.GameStatus(status)

		log.WithField("game", game).Debugf("sending game: %s/%s", game.Id, game.Slug)

		if err := stream.Send(game); err != nil {
			return err
		}
	}

	return nil
}
