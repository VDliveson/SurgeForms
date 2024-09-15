import coloredlogs, logging

logging.basicConfig(level=logging.ERROR)
LOGGER = logging.getLogger()
coloredlogs.install(fmt="%(asctime)s - %(message)s", level="INFO", logger=LOGGER)
