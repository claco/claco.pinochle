from pinochle.grpc import PinochleClient
from pinochle.grpc.pinochle_pb2 import CreateGameRequest, CreateGameResponse, GetGameRequest, GetGameResponse


def test_can_create_a_new_game(client: PinochleClient):
    request = CreateGameRequest()
    response: CreateGameResponse = client.CreateGame(request)

    assert response


def test_can_get_an_existing_game(client: PinochleClient):
    request = GetGameRequest()
    response: GetGameResponse = client.GetGame(request)

    assert response
