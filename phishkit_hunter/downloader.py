import os
from urllib.parse import urlparse
from .core import Core


class Downloader(Core):

    file_name = None

    def download(self, url):
        if not os.path.exists(self.download_location):
            os.makedirs(self.download_location)
        response = None
        parsed_url = urlparse(url)
        url_directory = parsed_url.netloc.replace('.', '_')
        if not os.path.exists(os.path.join(self.download_location, url_directory)):
            os.makedirs(os.path.join(self.download_location, url_directory))
        file_name = parsed_url.path.split('/')[-1]
        full_path = '{}/{}/{}'.format(self.download_location, url_directory, file_name)
        response = self.network.get(url, stream=True, verify=False)
        if response.status_code == 200 and int(response.headers.get('Content-Length', 0)) > 100:
            with open(full_path, 'wb') as download:
                download.write(response.content)
            return os.path.abspath(full_path)
        return None
