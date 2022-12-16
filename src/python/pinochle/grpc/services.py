import debugpy, logging, os

from concurrent import futures
from pinochle.grpc.pinochle_pb2_grpc import grpc, add_PinochleServiceServicer_to_server
from pinochle.grpc.servicers import PinochleServiceServicer

logger = logging.getLogger(__package__)


class PinochleService:
    def run(self, service_address: str | None = None) -> int:
        if not service_address:
            service_address = os.environ.get("PINOCHLE_SERVICE_ADDRESS", "localhost:50051")

        interceptors = []
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10), interceptors=interceptors)
        server.add_insecure_port(address=service_address)

        servicer = PinochleServiceServicer()

        add_PinochleServiceServicer_to_server(servicer, server)

        logger.info(f"Running Deployments Server (Python): Listening on {repr(service_address)}")
        logger.info("  ↪ Press CTRL-C to stop the service")

        if logger.getEffectiveLevel() == logging.DEBUG:
            logger.info(f"Running Remote Debugger (Python): Listening on 'localhost:5678'")

            debugpy.listen(5678)
            debugpy.breakpoint()

        try:
            server.start()
            server.wait_for_termination()
        except KeyboardInterrupt:
            logger.info("Server shutdown successfully")
        except BaseException as ex:
            logger.error(f"Server shutdown unexpectedly: {repr(ex)}")

            return 1

        return 0
