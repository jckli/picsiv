# picsiv

~~woo! another speedrun project! this project took me about 1 hour and 30 minutes to make.~~ this project has been updated to include some new nifty features. so its taken much longer than that to make lol.

## invite link

https://jackli.dev/picsiv

## features

picsiv is a simple and nifty bot that makes browsing anime art much easier. it can do a variety of features, like grabbing cool art from subreddits as well as simply sending a full image of a pixiv link with more features to come.

originally the bot was made to simply send the full quality image of a pixiv link posted in chat. because pixiv original embeds only send half of the image, as well as because their image cdn forbids people from sending the raw image link, i made this bot lol. this feature utilizes my pixiv mirror (created using cloudflare workers), https://pximg.jackli.dev/, and [hibiapi](https://github.com/mixmoe/HibiAPI)

theres also features that use reddit's api, and a bunch of my other cloudflare workers which can be found [here](https://github.com/jckli/art-workers). these commands grab images from a variety of subreddits and send the cool art into chat.

## screenshots

### original pixiv embed (that discord provides originally)

<img src="https://cdn.hayasaka.moe/aea5d9y2mkl2.jpg" />

*we dont like this i want to see the full beauty of this anime girl*

### pixiv bot!!!

<img src="https://cdn.hayasaka.moe/cknu8ebhssnt.jpg" />

**WE LIKE THIS WOOOOOOO!!!!!**

### reddit commands

<img src="https://cdn.hayasaka.moe/4k3t08qxpknu.jpg" />

## setup

1. create discord bot from the [Discord Developer Portal](https://discord.com/developers/applications/) website and save the token
2. clone this repository
3. go into the folder and create a .env file with the following contents:
```
TOKEN=your_token_here
```
*your_token_here should be replaced with your own token from step 1*

4. open terminal and run `pip install -r requirements.txt`
5. run the bot with `python -u bot.py`