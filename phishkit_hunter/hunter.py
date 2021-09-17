from .core import Core
from .deobfuscate import Deobfuscate
from .downloader import Downloader
from .generate import Generate
from .parser import Parser


class Hunter(Core):

    downloader = Downloader()
    deobfuscate = Deobfuscate()
    generate = Generate()

    def inspect(self, value: str):
        value = self.deobfuscate.url(value)
        self.__logger.info(f"Deobfuscated url: {value}")
        return_list = set()
        for item in self.generate.subdirectories(value):
            parsed = Parser()
            parsed.parse_links(item)
            for parsed_files in parsed.parsed_files:
                if parsed_files:
                    download = self.downloader.download(parsed_files)
                    if download:
                        return_list.add(download)
        parsed = Parser()
        parsed.parse_links(value)
        for item in parsed.parsed_files:
            if item:
                for generated_files in self.generate.subdirectories(item):
                    if generated_files:
                        download = self.downloader.download(generated_files)
                        if download:
                            return_list.add(download)
        return { 'files': list(return_list) }
