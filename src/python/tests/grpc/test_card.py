from typing import Callable
from pinochle.grpc import Card, CardSuit


def test_can_create_instance(card: Card):
    assert card


def test_has_no_default_suit(card: Card):
    assert card.suit == CardSuit.Unspecified  # type: ignore


def test_has_no_default_symbol(card: Card):
    assert not card.symbol


def test_create_accepts_suit_parameter(card_factory: Callable):
    card: Card = card_factory(suit=CardSuit.Hearts)  # type: ignore

    assert card.suit == CardSuit.Hearts  # type: ignore


def test_create_accepts_symbol_parameter(card_factory: Callable):
    card: Card = card_factory(symbol="1")

    assert card.symbol == "1"
