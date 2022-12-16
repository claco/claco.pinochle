import logging, typer

from pinochle.cli.commands import ServiceAddressOption
from pinochle.grpc import PinochleService

app = typer.Typer(no_args_is_help=True, name="service", help="Pinochle Service Commands", add_completion=False)

logger = logging.getLogger(__package__)


@app.command(help="Run the Pinochle service")
def run(service_address: str = ServiceAddressOption()) -> int:
    return PinochleService().run(service_address=service_address)
