# picsiv

woo! another speedrun project! this project took me about 1 hour and 30 minutes to make. 

<details>
    <summary>proof</summary>
    <img src="https://cdn.hayasaka.moe/5dtw2vdd1jzg.jpg" />
</details>

## invite link

https://discord.com/oauth2/authorize?client_id=947361674703302738&scope=applications.commands%20bot&permissions=3072

*i forgot to update this readme even after it was hosted lol*

## features

it simply sends the full quality image of the pixiv link. because pixiv original embeds only send half of the image, as well as because their image cdn forbids people from sending the raw image link, i made this bot lol.

this utilizes my pixiv mirror (created using cloudflare workers), https://pximg.jackli.dev/, and [hibiapi](https://github.com/mixmoe/HibiAPI)

theres no commands, just send a pixiv link and it will automatically send the full image.

## screenshots

### original pixiv embed (that discord provides originally)

<img src="https://cdn.hayasaka.moe/aea5d9y2mkl2.jpg" />

*we dont like this i want to see the full beauty of this anime girl*

### pixiv bot!!!

<img src="https://cdn.hayasaka.moe/cknu8ebhssnt.jpg" />

**WE LIKE THIS WOOOOOOO!!!!!**

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