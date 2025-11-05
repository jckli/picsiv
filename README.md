# picsiv

~~woo! another speedrun project! this project took me about 1 hour and 30 minutes to make.~~ this project has been updated to include some new nifty features. so its taken much longer than that to make lol.

## invite link

https://jackli.dev/picsiv

## features

picsiv is a simple and nifty bot that makes browsing anime art much easier. it can do a variety of features, like grabbing cool art from subreddits as well as simply sending a full image of a pixiv link with maybe more features to come.

originally the bot was made to simply send the full quality image of a pixiv link posted in chat. because pixiv original embeds only send half of the image, as well as because their image cdn forbids people from sending the raw image link, i made this bot lol. this feature utilizes my pixiv mirror (created using cloudflare workers), https://pximg.jackli.dev/, and [hibiapi](https://github.com/mixmoe/HibiAPI)

theres also features that use reddit's api, and a bunch of my other cloudflare workers which can be found [here](https://github.com/jckli/art-workers). these commands grab images from a variety of subreddits and send the cool art into chat.

## screenshots

### original pixiv embed (that discord provides originally)

<img src="https://cdn.hayasaka.moe/aea5d9y2mkl2.jpg" />

_we dont like this i want to see the full beauty of this anime girl_

### pixiv bot!!!

<img src="https://cdn.hayasaka.moe/cknu8ebhssnt.jpg" />

**WE LIKE THIS WOOOOOOO!!!!!**

### reddit commands

<img src="https://cdn.hayasaka.moe/4k3t08qxpknu.jpg" />

## setup

ideally just use the one that i host at https://jackli.dev/picsiv, but if you want to host your own, follow these steps:

1. create discord bot from the [Discord Developer Portal](https://discord.com/developers/applications/) website and save the token
2. host a hibiapi instance, or find a public one that you can use
3. clone this repository
4. go into the folder and create a .env file with the following contents:

```
VERSION=your_version_here
TOKEN=your_token_here
DEV_SERVER_ID=your_dev_server_id_here
DEV_MODE=dev_mode_true_or_false_here
PIXIV_API_URL=your_api_url_here
DEV_ERROR_CHANNEL_ID=your_dev_channel_id_here
```

4. open terminal and run `go run main.go`
5. invite the bot to your server and use it
