from bs4 import BeautifulSoup as bs
import requests
import json
import io
import zipfile
import imageio.v3 as imageio
import urllib3

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
    response = requests.get(url=f"https://hibiapi.hayasaka.moe/api/pixiv/illust?id={id}")
    return response.text

def request_hibiapi_ugoria(id):
    response = requests.get(url=f"https://hibiapi.hayasaka.moe/api/pixiv/ugoira_metadata?id={id}")
    return response.text

def parse_hibiapi(data):
    data = json.loads(data)
    illust = data["illust"]

    # check if ugoria
    ugoria = False
    if illust["type"] == "ugoira":
        ugoria = True

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

    result = {"nsfw": nsfw, "links": links, "ugoria": ugoria}
    return result

def parse_ugoria(id):
    ugoria_data = request_hibiapi_ugoria(id)
    ugoria_json = json.loads(ugoria_data)
    ugoria = ugoria_json["ugoira_metadata"]

    raw_zip_link = ugoria["zip_urls"]["medium"]
    path = raw_zip_link.split("https://i.pximg.net/")[1]
    mirpath = "https://pximg.jackli.dev/" + path

    http = urllib3.PoolManager()

    zip_resp = http.request("GET", mirpath)
    zip_file = zipfile.ZipFile(io.BytesIO(zip_resp.data))

    frame_count = len(ugoria["frames"])
    total_delay = sum(frame["delay"] for frame in ugoria["frames"])
    dur = total_delay / frame_count

    gif = io.BytesIO()
    images = []
    for name in zip_file.namelist():
        images.append(imageio.imread(zip_file.open(name)))
    imageio.imwrite(gif, images, extension=".gif", duration=dur, loop=0)

    return gif
