import logging, os, sys, typer

from typing import Optional, Sequence
from pinochle.cli.commands import game, service

app = typer.Typer(no_args_is_help=True, add_completion=False)
app.add_typer(game.app, name="game", help="Pinochle Game Commands")
app.add_typer(service.app, name="service", help="Pinochle Service Commands")

logger = logging.getLogger(__package__)


def main(argv: Optional[Sequence[str]] = []) -> int:
    (log_format, log_level) = "%(message)s", logging.INFO

    exitCode = 0

    if not argv:
        argv = sys.argv[1:]

    if "--debug" in argv or os.getenv("LOG_LEVEL", "").strip().upper() == "DEBUG":
        log_format = "%(asctime)s %(name)s %(levelname)s [%(filename)s:%(lineno)s:%(funcName)s] %(message)s"
        log_level = logging.DEBUG

        if "--debug" in sys.argv:
            sys.argv.remove("--debug")

    logging.basicConfig(encoding="utf-8", datefmt="%Y-%m-%dT%H:%M:%S%z", format=log_format, level=log_level)

    try:
        exitCode = app()
    except SystemExit as ex:
        exitCode = ex.code
    except BaseException as ex:
        logger.error(ex)

        exitCode = 2

    return int(exitCode)  # type: ignore
