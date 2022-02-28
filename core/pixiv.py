from bs4 import BeautifulSoup as bs
import requests
import json

def request_pixiv(url):
    headers = {
        "Host": "www.pixiv.net",
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:99.0) Gecko/20100101 Firefox/99.0",
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
        "Accept-Language": "en-US,en;q=0.5",
        "Accept-Encoding": "gzip, deflate, br",
        "DNT": "1",
        "Connection": "keep-alive",
        "Upgrade-Insecure-Requests": "1",
        "Sec-Fetch-Dest": "document",
        "Sec-Fetch-Mode": "navigate",
        "Sec-Fetch-Site": "none",
        "Sec-Fetch-User": "?1"
    }
    response = requests.get(url=url, headers=headers)
    return response.text

def parse_pixiv(html):
    soup = bs(html, "html.parser")
    print(soup)
    imgbase = soup.find("div", {"class": "sc-1qpw8k9-0 gTFqQV"})
    print(imgbase)
    return imgbase["src"]

def request_hibiapi(id):
    response = requests.get(url=f"https://api.zettai.moe/api/pixiv/illust?id={id}")
    return response.text

def parse_hibiapi(data):
    data = json.loads(data)
    illust = data["illust"]

    # get images
    links = []
    if illust["meta_single_page"] != {}:
        links.append(illust["meta_single_page"]["original_image_url"])
    elif illust["meta_pages"] != []:
        for i in illust["meta_pages"]:
            links.append(i["image_urls"]["original"])
    else:
        return None

    # get level
    if illust["sanity_level"] >= 5:
        nsfw = True
    else:
        nsfw = False

    result = {"nsfw": nsfw, "links": links}
    return result