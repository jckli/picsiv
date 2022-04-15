import string
import random
import requests
import json

def requestimg(link):
    response = requests.get(link)
    return json.loads(response.text)

def randomString(length):
    res = ''.join(random.choices(string.ascii_lowercase + string.digits, k = length))
    return str(res)