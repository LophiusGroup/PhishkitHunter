from bs4 import BeautifulSoup

from .core import Core


class Parser(Core):

    extensions = [
        'zip',
        'exe',
        'msi',
        'mp4',
        'ps1',
        'txt',
        'log',
        'apk',
        'dll',
        'bin'
    ]
    url_list = []
    url_file_list = []
    searched_url_list = []

    @property
    def parsed_urls(self):
        return_url_list = []
        if self.url_list:
            for url in self.url_list:
                return_url_list.append(url)
        return return_url_list

    @property
    def parsed_files(self):
        return_url_list = []
        if self.url_file_list:
            for url in self.url_file_list:
                return_url_list.append(url)
        return return_url_list

    def parse_links(self, links):
        url_list = []
        if not isinstance(links, list):
            links = [links]
        for link in links:
            if link not in self.searched_url_list:
                if link.endswith('/'):
                    content = ''
                    self.searched_url_list.append(link)
                    response = self.network.get(link)
                    if response:
                        try:
                            content = response.content.decode('utf-8')
                        except:
                            continue
                    soup = BeautifulSoup(content, 'html.parser')
                    index_of = self.parse_indexof(soup, link)
                    if index_of:
                        for item in index_of:
                            if item not in self.searched_url_list:
                                url_list.append(item)
                    else:
                        for item in self.parse_href_links(soup):
                            new_url = '{}/{}'.format(link, item)
                            if new_url not in self.searched_url_list:
                                url_list.append(new_url)
                else:
                    for extension in self.extensions:
                        if link.endswith(extension):
                            if link not in self.searched_url_list:
                                self.url_file_list.append(link)
        if url_list:
            for item in url_list:
                self.url_list.append(item)
            self.parse_links(self.url_list)

    def parse_indexof(self, soup, url):
        return_list = []
        try:
            if 'index of' in (soup.title.string).lower():
                for link in self.parse_href_links(soup):
                    new_link_path = '{}/{}'.format(url.rstrip('/'), link)
                    if new_link_path.endswith('//'):
                        new_link_path = new_link_path[:-1]
                    return_list.append(new_link_path)
                return return_list
        except:
            pass
    
    def parse_href_links(self, soup):
        return_list = []
        for link in soup.find_all('a'):
            return_list.append(link.get('href'))
        return return_list
