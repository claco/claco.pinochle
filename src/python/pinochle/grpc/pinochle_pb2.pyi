from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

Clubs: CardSuit
Completed: GameStatus
DESCRIPTOR: _descriptor.FileDescriptor
Diamonds: CardSuit
Hearts: CardSuit
New: GameStatus
Playing: GameStatus
Spades: CardSuit
Unspecified: CardSuit

class Board(_message.Message):
    __slots__ = ["discards", "melds", "stock"]
    DISCARDS_FIELD_NUMBER: _ClassVar[int]
    MELDS_FIELD_NUMBER: _ClassVar[int]
    STOCK_FIELD_NUMBER: _ClassVar[int]
    discards: _containers.RepeatedCompositeFieldContainer[Card]
    melds: _containers.RepeatedCompositeFieldContainer[Meld]
    stock: _containers.RepeatedCompositeFieldContainer[Card]
    def __init__(self, stock: _Optional[_Iterable[_Union[Card, _Mapping]]] = ..., discards: _Optional[_Iterable[_Union[Card, _Mapping]]] = ..., melds: _Optional[_Iterable[_Union[Meld, _Mapping]]] = ...) -> None: ...

class Card(_message.Message):
    __slots__ = ["suit", "symbol"]
    SUIT_FIELD_NUMBER: _ClassVar[int]
    SYMBOL_FIELD_NUMBER: _ClassVar[int]
    suit: CardSuit
    symbol: str
    def __init__(self, suit: _Optional[_Union[CardSuit, str]] = ..., symbol: _Optional[str] = ...) -> None: ...

class CreateGameRequest(_message.Message):
    __slots__ = ["name"]
    NAME_FIELD_NUMBER: _ClassVar[int]
    name: str
    def __init__(self, name: _Optional[str] = ...) -> None: ...

class CreateGameResponse(_message.Message):
    __slots__ = ["game"]
    GAME_FIELD_NUMBER: _ClassVar[int]
    game: Game
    def __init__(self, game: _Optional[_Union[Game, _Mapping]] = ...) -> None: ...

class Deck(_message.Message):
    __slots__ = ["cards"]
    CARDS_FIELD_NUMBER: _ClassVar[int]
    cards: _containers.RepeatedCompositeFieldContainer[Card]
    def __init__(self, cards: _Optional[_Iterable[_Union[Card, _Mapping]]] = ...) -> None: ...

class Game(_message.Message):
    __slots__ = ["board", "id", "name", "players", "slug", "status"]
    BOARD_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    PLAYERS_FIELD_NUMBER: _ClassVar[int]
    SLUG_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    board: Board
    id: int
    name: str
    players: _containers.RepeatedCompositeFieldContainer[Player]
    slug: str
    status: GameStatus
    def __init__(self, id: _Optional[int] = ..., slug: _Optional[str] = ..., name: _Optional[str] = ..., board: _Optional[_Union[Board, _Mapping]] = ..., players: _Optional[_Iterable[_Union[Player, _Mapping]]] = ..., status: _Optional[_Union[GameStatus, str]] = ...) -> None: ...

class GetGameRequest(_message.Message):
    __slots__ = ["slug"]
    SLUG_FIELD_NUMBER: _ClassVar[int]
    slug: str
    def __init__(self, slug: _Optional[str] = ...) -> None: ...

class GetGameResponse(_message.Message):
    __slots__ = ["game"]
    GAME_FIELD_NUMBER: _ClassVar[int]
    game: Game
    def __init__(self, game: _Optional[_Union[Game, _Mapping]] = ...) -> None: ...

class ListGamesRequest(_message.Message):
    __slots__ = []
    def __init__(self) -> None: ...

class Meld(_message.Message):
    __slots__ = []
    def __init__(self) -> None: ...

class Player(_message.Message):
    __slots__ = []
    def __init__(self) -> None: ...

class StartGameRequest(_message.Message):
    __slots__ = ["slug"]
    SLUG_FIELD_NUMBER: _ClassVar[int]
    slug: str
    def __init__(self, slug: _Optional[str] = ...) -> None: ...

class StartGameResponse(_message.Message):
    __slots__ = ["game"]
    GAME_FIELD_NUMBER: _ClassVar[int]
    game: Game
    def __init__(self, game: _Optional[_Union[Game, _Mapping]] = ...) -> None: ...

class CardSuit(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []

class GameStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
