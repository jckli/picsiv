import discord
from discord.ext import commands, pages
import validators
import re
import io
from core import pixiv

class Commands(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.Cog.listener()
    async def on_message(self, message):
        if "pixiv.net" and "artworks" in message.content:
            urlRaw = re.search("(?P<url>https?://[^\s]+\d)", message.content)
            if urlRaw == None:
                return
            else:
                url = urlRaw.group("url")
            if validators.url(url) is True:
                ctx = await self.bot.get_context(message)
                pixivid = url.split("/artworks/")[1]
                json = pixiv.request_hibiapi(pixivid)
                data = pixiv.parse_hibiapi(json)
                channel = await discord.utils.get_or_fetch(message.guild, "channel", message.channel.id)
                if data["nsfw"] is True:
                    if channel.nsfw is False:
                        embed = discord.Embed(title="NSFW", color=0x0096fa, description="This image is NSFW. Please resend the link in a NSFW channel to view this image.")
                        await ctx.send(embed=embed, reference=message, mention_author=False)
                        return
                if data["ugoria"] is True:
                    ugoria = pixiv.parse_ugoria(pixivid)
                    ugoria.seek(0)
                    file = discord.File(fp=ugoria, filename="ugoria.gif")
                    embed = discord.Embed(title="Full pixiv Ugoria", color=0x0096fa)
                    embed.set_image(url="attachment://ugoria.gif")
                    await message.channel.send(file=file, embed=embed, reference=message, mention_author=False)
                    return
                mirlinks = []
                for i in data["links"]:
                    path = i.split("https://i.pximg.net/")[1]
                    mirimg = "https://pximg.jackli.dev/" + path
                    mirlinks.append(mirimg)
                if len(data["links"]) > 1:
                    imgpages = []
                    for count, link in enumerate(mirlinks):
                        imgpages.append(discord.Embed(title=f"Full pixiv Image", color=0x0096fa))
                        imgpages[count].set_image(url=link)
                    paginator = pages.Paginator(pages=imgpages, use_default_buttons=False, timeout=None, author_check=False)
                    paginator.add_button(
                        pages.PaginatorButton("prev", label="<", style=discord.ButtonStyle.red)
                    )
                    paginator.add_button(
                        pages.PaginatorButton(
                            "page_indicator", style=discord.ButtonStyle.gray, disabled=True
                        )
                    )
                    paginator.add_button(
                        pages.PaginatorButton("next", label=">", style=discord.ButtonStyle.green)
                    )
                    await paginator.send(ctx, reference=message, mention_author=False)
                else:
                    embed = discord.Embed(title="Full pixiv Image", color=0x0096fa)
                    embed.set_image(url=mirlinks[0])
                    await message.channel.send(embed=embed, reference=message, mention_author=False)

def setup(bot):
    bot.add_cog(Commands(bot))