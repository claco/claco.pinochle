import os

from pinochle.grpc import PinochleClient
from pinochle.grpc.pinochle_pb2_grpc import grpc
from typer import Argument, Option


class Client:
    def __init__(self, service_address: str | None = None) -> None:
        if not service_address:
            service_address = os.environ.get("PINOCHLE_SERVICE_ADDRESS", "localhost:50051")

        self.channel = grpc.insecure_channel(target=service_address)
        self.service = PinochleClient(channel=self.channel)


def ServiceAddressOption(default=None, **kwargs):
    return Option(default, help="Pinochle service address", metavar="host:port")
