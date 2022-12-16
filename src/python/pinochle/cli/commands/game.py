import logging, typer

from pinochle.grpc import PinochleClient, Game
from pinochle.grpc.pinochle_pb2 import CreateGameRequest, CreateGameResponse, GetGameRequest, GetGameResponse
from pinochle.grpc.pinochle_pb2 import StartGameRequest, StartGameResponse
from pinochle.grpc.pinochle_pb2 import ListGamesRequest
from pinochle.cli.commands import Client

app = typer.Typer(no_args_is_help=True, name="game", help="Pinochle Game Commands", add_completion=False)
client = Client().service
logger = logging.getLogger(__package__)


@app.command(help="Create a new game", no_args_is_help=False)
def create_game(name: str | None = None):
    logger.info(f"Creating game...")

    request = CreateGameRequest(name=name)
    response: CreateGameResponse = client.CreateGame(request)

    logger.debug(f"request={repr(request)}, response={repr(response)}")

    if response:
        logger.info(f"Successfully created game")
    else:
        logger.error(f"Error creating game")


@app.command(help="Get existing game details", no_args_is_help=False)
def get_game(slug: str):
    logger.info(f"Retrieving game...")

    request = GetGameRequest(slug=slug)
    response: GetGameResponse = client.GetGame(request)

    logger.debug(f"request={repr(request)}, response={repr(response)}")

    if response:
        logger.info(f"Successfully retrieved game")
    else:
        logger.error(f"Error retrieving game")

    logger.info(response)


@app.command(help="List existing games", no_args_is_help=False)
def list_games():
    logger.info(f"Retrieving games...\n")

    request = ListGamesRequest()

    logger.debug(f"request={repr(request)}")

    for game in client.ListGames(request):
        logger.debug(f"game={repr(game)}")

        # if game:
        #     logger.info(f"Successfully retrieved games")
        # else:
        #     logger.error(f"Error retrieving game")

        logger.info(game)


@app.command(help="Start existing game", no_args_is_help=False)
def start_game(slug: str):
    logger.info(f"Starting game...\n")

    request = StartGameRequest(slug=slug)
    response: StartGameResponse = client.StartGame(request)

    logger.debug(f"request={repr(request)}, response={repr(response)}")

    if response:
        logger.info(f"Successfully started game")
    else:
        logger.error(f"Error starting game")

    logger.info(response)
