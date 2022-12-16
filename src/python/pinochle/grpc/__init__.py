import os, sys

# Fix protoc generated `import` statement issues
# See https://github.com/protocolbuffers/protobuf/issues/7061
current_python_path = os.path.dirname(os.path.realpath(__file__))
sys.path.append(current_python_path)

from pinochle.grpc.pinochle_pb2 import Board, Card, CardSuit, Game, Meld, Player
from pinochle.grpc.pinochle_pb2_grpc import PinochleServiceStub as PinochleClient
from pinochle.grpc.pinochle_pb2_grpc import add_PinochleServiceServicer_to_server
from pinochle.grpc.servicers import PinochleServiceServicer
from pinochle.grpc.services import PinochleService
