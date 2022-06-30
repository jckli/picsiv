import discord
from discord.ext import commands
from discord.commands import Option
from core import r_utils

orientations = ["landscape", "portrait", "square"]

class ArtApis(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    async def get_orientations(ctx: discord.AutocompleteContext):
        return [ot for ot in orientations if ot.startswith(ctx.value.lower())]

    @commands.slash_command(name="sugoiart", description="Gets a random image from the sugoiart API")
    async def sugoiart(self, ctx, orientation: Option(str, "Choose an orientation of the image fetched", autocomplete=get_orientations) = None):
        await ctx.defer()
        rs = r_utils.randomString(length=6)
        orientationstring = f"&orien={orientation}" if orientation is not None else ""
        try:
            img = r_utils.requestimg(f"https://art.hayasaka.moe/api/art/random?_={rs}{orientationstring}")
        except:
            errorEmbed = discord.Embed(title="Error", url="Could not get image from API. Please try again.", color=0xff524f)
            await ctx.respond(embed=errorEmbed)
            return
        embed = discord.Embed(title="sugoiart", url="https://art.hayasaka.moe/", color=0x0096fa)
        embed.set_image(url=img["url"])
        embed.set_footer(text="Powered by sugoiart (https://art.hayasaka.moe/)")
        await ctx.respond(embed=embed)
        return

def setup(bot):
    bot.add_cog(ArtApis(bot))