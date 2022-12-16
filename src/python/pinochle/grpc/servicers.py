import json, logging, random, secrets, sqlite3

from google.protobuf import json_format
from grpc import ServicerContext
from pinochle.grpc.pinochle_pb2_grpc import PinochleServiceServicer as PinochleServicerBase
from pinochle.grpc.pinochle_pb2 import CreateGameRequest, CreateGameResponse, GetGameRequest, GetGameResponse
from pinochle.grpc.pinochle_pb2 import StartGameResponse
from pinochle.grpc.pinochle_pb2 import Game, GameStatus, Board, Card, CardSuit

logger = logging.getLogger(__package__)


class Deck:
    def __init__(self, shuffle=False) -> None:
        self.cards = []
        self.suit_order = [CardSuit.Spades, CardSuit.Diamonds, CardSuit.Clubs, CardSuit.Hearts]  # type: ignore
        # self.suit_symbols = {CardSuit.Spades: "♠️", CardSuit.Diamonds: "♦️", CardSuit.Clubs: "♣️", CardSuit.Hearts: "♥️"}  # type: ignore
        self.suit_symbols = {CardSuit.Spades: "S", CardSuit.Diamonds: "D", CardSuit.Clubs: "C", CardSuit.Hearts: "H"}  # type: ignore
        self.symbol_order = ["A"] + list(range(2, 10)) + ["J", "Q", "K"]

        for suit in self.suit_order:
            for symbol in self.symbol_order:
                card = Card(suit=suit, symbol=str(symbol))
                self.cards.append(card)

        if shuffle:
            random.shuffle(self.cards)

    def __repr__(self) -> str:
        cards = []

        for card in self.cards:
            cards.append(f"{self.suit_symbols[card.suit]}{card.symbol}")

        return ",".join(cards)


class PinochleServiceServicer(PinochleServicerBase):
    def __init__(self) -> None:
        db = self.db = sqlite3.connect("pinochle.db", isolation_level=None, check_same_thread=False)
        db.row_factory = sqlite3.Row

        db.executescript(
            """
            CREATE TABLE IF NOT EXISTS games (
                id      INTEGER PRIMARY KEY AUTOINCREMENT,
                slug    VARCHAR(10) NOT NULL,
                name    VARCHAR(25) NOT NULL,
                status  TINYINT DEFAULT 0
            );

            CREATE TABLE IF NOT EXISTS boards (
                id      INTEGER PRIMARY KEY AUTOINCREMENT,
                game_id INTEGER UNIQUE NOT NULL,
                stock   VARCHAR
            );
            """
        )

    def CreateGame(self, request: CreateGameRequest, context: ServicerContext) -> CreateGameResponse:
        logger.debug(f"request={repr(request)}, context={repr(context)}")

        slug = secrets.token_urlsafe(10)
        name = request.name or "Untitled Game"
        logger.debug(f"id={repr(slug)}, name={repr(name)}")

        record = dict(
            self.db.execute("INSERT INTO games (slug, name) VALUES (?, ?) RETURNING *", (slug, name)).fetchone()
        )
        logger.info(record)
        game = json_format.ParseDict(record, message=Game())

        return CreateGameResponse(game=game)

    def GetGame(self, request: GetGameRequest, context: ServicerContext) -> GetGameResponse:
        logger.debug(f"request={repr(request)}, context={repr(context)}")

        slug = request.slug
        record = dict(self.db.execute("SELECT id,slug,name,status FROM games WHERE slug=?", (slug,)).fetchone())
        game = json_format.ParseDict(record, message=Game())
        logger.debug(f"record={repr(record)}")

        if game.status:
            record = dict(self.db.execute("SELECT stock FROM boards WHERE game_id=?", (game.id,)).fetchone())
            logger.debug(f"record={repr(record)}")

            deck = Deck()  # ditch this
            suits = dict([(value, key) for key, value in deck.suit_symbols.items()])

            for card_spec in record["stock"].split(","):
                card_spec = card_spec.strip()
                suit, symbol = card_spec[:1], card_spec[1:]
                card = Card(suit=suits[suit], symbol=symbol)
                game.board.stock.append(card)

            # json_format.ParseDict(record, message=game.board)

        return GetGameResponse(game=game)

    def ListGames(self, request, context):
        logger.debug(f"request={repr(request)}, context={repr(context)}")

        for record in self.db.execute("SELECT id,slug,name,status FROM games;"):
            record = dict(record)

            logger.debug(f"record={repr(record)}")

            yield json_format.ParseDict(record, message=Game())

    def StartGame(self, request, context):
        logger.debug(f"request={repr(request)}, context={repr(context)}")

        game = self.GetGame(request, context).game

        if not game.status:
            deck = Deck(shuffle=True)

            logging.info(repr(deck))

            with self.db:
                status = GameStatus.Playing  # type: ignore
                board_record = self.db.execute(
                    "INSERT OR REPLACE INTO boards (game_id, stock) VALUES (?, ?)", (game.id, repr(deck))
                )
                self.db.execute("UPDATE games SET status = ? where id = ?", (status, game.id))
                game.status = status
                game.board.stock.extend(deck.cards)

        return StartGameResponse(game=game)
