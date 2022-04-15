import discord
from discord.ext import commands
from discord.commands import Option
from core import r_utils

timeperiods = ["hour", "day", "week", "month", "year", "all"]

class Reddit(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    async def get_timeperiods(ctx: discord.AutocompleteContext):
        return [tp for tp in timeperiods if tp.startswith(ctx.value.lower())]

    @commands.slash_command(name="streetmoe")
    async def streetmoe(self, ctx, timeperiod: Option(str, "Pick a time period:", autocomplete=get_timeperiods) = None):
        rs = r_utils.randomString(length=6)
        img = ""
        if timeperiod is None:
            img = r_utils.requestimg(f"https://streetmoe.jackli.dev/api?_={rs}")
        else:
            img = r_utils.requestimg(f"https://streetmoe.jackli.dev/api?_={rs}&t={timeperiod}")
        embed = discord.Embed(title="r/streetmoe", url="https://www.reddit.com/r/streetmoe/", color=0x0096fa)
        embed.set_image(url=img["imglink"])
        embed.set_footer(text="Powered by https://streetmoe.jackli.dev/")
        await ctx.respond(embed=embed)

def setup(bot):
    bot.add_cog(Reddit(bot))