import pytest

from pinochle.grpc import Board, Card, Deck, Game, Meld, Player, PinochleClient
from typing import Callable


@pytest.fixture(scope="module")
def grpc_add_to_server():
    from pinochle.grpc import add_PinochleServiceServicer_to_server

    return add_PinochleServiceServicer_to_server


@pytest.fixture(scope="module")
def grpc_servicer():
    from pinochle.grpc import PinochleServiceServicer

    return PinochleServiceServicer()


@pytest.fixture(scope="module")
def grpc_stub_cls(grpc_channel):
    from pinochle.grpc import PinochleClient

    return PinochleClient


@pytest.fixture(scope="module")
def client(grpc_stub) -> PinochleClient:
    return grpc_stub


@pytest.fixture
def board() -> Board:
    return Board()


@pytest.fixture
def card() -> Card:
    return Card()


@pytest.fixture
def card_factory() -> Callable[..., Card]:
    return Card


@pytest.fixture
def deck() -> Deck:
    return Deck()


@pytest.fixture
def game() -> Game:
    return Game()


@pytest.fixture
def meld() -> Meld:
    return Meld()


@pytest.fixture
def player() -> Player:
    return Player()
