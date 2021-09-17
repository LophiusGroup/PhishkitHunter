import os
import requests
from .utils.logger import LoggingBase


class Core(metaclass=LoggingBase):

    download_location = os.environ.get('DOWNLOAD_LOCATION', os.path.dirname(os.path.abspath(__file__)))
    network = requests.Session()
