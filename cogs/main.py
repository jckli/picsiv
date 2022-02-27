import discord
from discord.ext import commands
import validators
from core import pixiv

class Commands(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.slash_command(name="help")
    async def help(self, ctx):
        embed = discord.Embed(title="Help", color=0x0096fa)

    @commands.Cog.listener()
    async def on_message(self, message):
        if validators.url(message.content) is True:
            if "pixiv.net" and "artworks" in message.content:
                pixivid = message.content.split("/artworks/")[1]
                json = pixiv.request_hibiapi(pixivid)
                links = pixiv.parse_hibiapi(json)
                mirlinks = []
                for i in links:
                    path = i.split("https://i.pximg.net/")[1]
                    mirimg = "https://pximg.jackli.dev/" + path
                    mirlinks.append(mirimg)
                for link in mirlinks:
                    embed = discord.Embed(title="Full pixiv Image", color=0x0096fa)
                    embed.set_image(url=link)
                    await message.channel.send(embed=embed, reference=message, mention_author=False)

def setup(bot):
    bot.add_cog(Commands(bot))