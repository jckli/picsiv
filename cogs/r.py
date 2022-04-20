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

    @commands.slash_command(name="streetmoe", description="Gets a random image from r/streetmoe")
    async def streetmoe(self, ctx, timeperiod: Option(str, "Pick a time period", autocomplete=get_timeperiods) = None):
        rs = r_utils.randomString(length=6)
        timeperiodstring = f"&t={timeperiod}" if timeperiod is not None else ""
        try:
            img = r_utils.requestimg(f"https://streetmoe.jackli.dev/api?_={rs}{timeperiodstring}")
        except:
            errorEmbed = discord.Embed(title="Error", url="Could not get image from API. Please try again.", color=0xff524f)
            await ctx.respond(embed=errorEmbed)
            return
        embed = discord.Embed(title="r/streetmoe", url="https://www.reddit.com/r/streetmoe/", color=0x0096fa)
        embed.set_image(url=img["imglink"])
        embed.set_footer(text="Powered by https://streetmoe.jackli.dev/")
        await ctx.respond(embed=embed)

    @commands.slash_command(name="animehoodies", description="Gets a random image from r/animehoodies")
    async def animehoodies(self, ctx, timeperiod: Option(str, "Pick a time period", autocomplete=get_timeperiods) = None):
        rs = r_utils.randomString(length=6)
        timeperiodstring = f"&t={timeperiod}" if timeperiod is not None else ""
        try:
            img = r_utils.requestimg(f"https://animehoodies.jackli.dev/api?_={rs}{timeperiodstring}")
        except:
            errorEmbed = discord.Embed(title="Error", url="Could not get image from API. Please try again.", color=0xff524f)
            await ctx.respond(embed=errorEmbed)
            return
        embed = discord.Embed(title="r/animehoodies", url="https://www.reddit.com/r/animehoodies/", color=0x0096fa)
        embed.set_image(url=img["imglink"])
        embed.set_footer(text="Powered by https://animehoodies.jackli.dev/")
        await ctx.respond(embed=embed)

    @commands.slash_command(name="animewallpaper", description="Gets a random image from r/animewallpaper")
    async def animewallpaper(self, ctx, timeperiod: Option(str, "Pick a time period", autocomplete=get_timeperiods) = None):
        rs = r_utils.randomString(length=6)
        timeperiodstring = f"&t={timeperiod}" if timeperiod is not None else ""
        try:
            img = r_utils.requestimg(f"https://aniwp.jackli.dev/api?_={rs}{timeperiodstring}")
        except:
            errorEmbed = discord.Embed(title="Error", url="Could not get image from API. Please try again.", color=0xff524f)
            await ctx.respond(embed=errorEmbed)
            return
        embed = discord.Embed(title="r/animehoodies", url="https://www.reddit.com/r/animewallpaper/", color=0x0096fa)
        embed.set_image(url=img["imglink"])
        embed.set_footer(text="Powered by https://aniwp.jackli.dev/")
        await ctx.respond(embed=embed)

    @commands.slash_command(name="moescape", description="Gets a random image from r/moescape")
    async def moescape(self, ctx, timeperiod: Option(str, "Pick a time period", autocomplete=get_timeperiods) = None):
        rs = r_utils.randomString(length=6)
        timeperiodstring = f"&t={timeperiod}" if timeperiod is not None else ""
        try:
            img = r_utils.requestimg(f"https://moescape.jackli.dev/api?_={rs}{timeperiodstring}")
        except:
            errorEmbed = discord.Embed(title="Error", url="Could not get image from API. Please try again.", color=0xff524f)
            await ctx.respond(embed=errorEmbed)
            return
        embed = discord.Embed(title="r/moescape", url="https://www.reddit.com/r/moescape/", color=0x0096fa)
        embed.set_image(url=img["imglink"])
        embed.set_footer(text="Powered by https://moescape.jackli.dev/")
        await ctx.respond(embed=embed)

    @commands.slash_command(name="wholesomeyuri", description="Gets a random image from r/wholesomeyuri")
    async def wholesomeyuri(self, ctx, timeperiod: Option(str, "Pick a time period", autocomplete=get_timeperiods) = None):
        rs = r_utils.randomString(length=6)
        timeperiodstring = f"&t={timeperiod}" if timeperiod is not None else ""
        try:
            img = r_utils.requestimg(f"https://wsyuri.jackli.dev/api?_={rs}{timeperiodstring}")
        except:
            errorEmbed = discord.Embed(title="Error", url="Could not get image from API. Please try again.", color=0xff524f)
            await ctx.respond(embed=errorEmbed)
            return
        embed = discord.Embed(title="r/wholesomeyuri", url="https://www.reddit.com/r/wholesomeyuri/", color=0x0096fa)
        embed.set_image(url=img["imglink"])
        embed.set_footer(text="Powered by https://wsyuri.jackli.dev/")
        await ctx.respond(embed=embed)
    
def setup(bot):
    bot.add_cog(Reddit(bot))