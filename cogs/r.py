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
    async def streetmoe(self, ctx, 
        timeperiod: Option(str, "Pick a time period", autocomplete=get_timeperiods) = None, 
        nsfw: Option(bool, "NSFW?") = None
    ):
        await ctx.defer()
        channelnsfw = ctx.channel.is_nsfw()
        if nsfw and not channelnsfw:
            nsfwEmbed = discord.Embed(title="NSFW", color=0x0096fa, description="This channel is not NSFW. Please switch to an NSFW channel to use the NSFW true option.")
            await ctx.respond(embed=nsfwEmbed)
            return
        elif channelnsfw and nsfw is None:
            nsfw = True 
        elif not channelnsfw and nsfw is None:
            nsfw = False
        rs = r_utils.randomString(length=6)
        timeperiodstring = f"&t={timeperiod}" if timeperiod is not None else ""
        nsfwstring = f"&nsfw={nsfw}"
        try:
            img = r_utils.requestimg(f"https://streetmoe.jackli.dev/api?_={rs}{timeperiodstring}{nsfwstring}")
        except:
            errorEmbed = discord.Embed(title="Error", url="Could not get image from API. Please try again.", color=0xff524f)
            await ctx.respond(embed=errorEmbed)
            return
        embed = discord.Embed(title="r/streetmoe", url="https://www.reddit.com/r/streetmoe/", color=0x0096fa)
        embed.set_image(url=img["imglink"])
        embed.set_footer(text="Powered by https://streetmoe.jackli.dev/")
        await ctx.respond(embed=embed)

    @commands.slash_command(name="animehoodies", description="Gets a random image from r/animehoodies")
    async def animehoodies(self, ctx, 
        timeperiod: Option(str, "Pick a time period", autocomplete=get_timeperiods) = None,
        nsfw: Option(bool, "NSFW?") = None
    ):
        await ctx.defer()
        channelnsfw = ctx.channel.is_nsfw()
        if nsfw and not channelnsfw:
            nsfwEmbed = discord.Embed(title="NSFW", color=0x0096fa, description="This channel is not NSFW. Please switch to an NSFW channel to use the NSFW true option.")
            await ctx.respond(embed=nsfwEmbed)
            return
        elif channelnsfw and nsfw is None:
            nsfw = True 
        elif not channelnsfw and nsfw is None:
            nsfw = False
        rs = r_utils.randomString(length=6)
        timeperiodstring = f"&t={timeperiod}" if timeperiod is not None else ""
        nsfwstring = f"&nsfw={nsfw}"
        try:
            img = r_utils.requestimg(f"https://animehoodies.jackli.dev/api?_={rs}{timeperiodstring}{nsfwstring}")
        except:
            errorEmbed = discord.Embed(title="Error", url="Could not get image from API. Please try again.", color=0xff524f)
            await ctx.respond(embed=errorEmbed)
            return
        embed = discord.Embed(title="r/animehoodies", url="https://www.reddit.com/r/animehoodies/", color=0x0096fa)
        embed.set_image(url=img["imglink"])
        embed.set_footer(text="Powered by https://animehoodies.jackli.dev/")
        await ctx.respond(embed=embed)

    @commands.slash_command(name="animewallpaper", description="Gets a random image from r/animewallpaper")
    async def animewallpaper(self, ctx, 
        timeperiod: Option(str, "Pick a time period", autocomplete=get_timeperiods) = None,
        nsfw: Option(bool, "NSFW?") = None
    ):
        await ctx.defer()
        channelnsfw = ctx.channel.is_nsfw()
        if nsfw and not channelnsfw:
            nsfwEmbed = discord.Embed(title="NSFW", color=0x0096fa, description="This channel is not NSFW. Please switch to an NSFW channel to use the NSFW true option.")
            await ctx.respond(embed=nsfwEmbed)
            return
        elif channelnsfw and nsfw is None:
            nsfw = True 
        elif not channelnsfw and nsfw is None:
            nsfw = False
        rs = r_utils.randomString(length=6)
        timeperiodstring = f"&t={timeperiod}" if timeperiod is not None else ""
        nsfwstring = f"&nsfw={nsfw}"
        try:
            img = r_utils.requestimg(f"https://aniwp.jackli.dev/api?_={rs}{timeperiodstring}{nsfwstring}")
        except:
            errorEmbed = discord.Embed(title="Error", url="Could not get image from API. Please try again.", color=0xff524f)
            await ctx.respond(embed=errorEmbed)
            return
        embed = discord.Embed(title="r/animewallpaper", url="https://www.reddit.com/r/animewallpaper/", color=0x0096fa)
        embed.set_image(url=img["imglink"])
        embed.set_footer(text="Powered by https://aniwp.jackli.dev/")
        await ctx.respond(embed=embed)

    @commands.slash_command(name="moescape", description="Gets a random image from r/moescape")
    async def moescape(self, ctx, 
        timeperiod: Option(str, "Pick a time period", autocomplete=get_timeperiods) = None,
        nsfw: Option(bool, "NSFW?") = None
    ):
        await ctx.defer()
        channelnsfw = ctx.channel.is_nsfw()
        if nsfw and not channelnsfw:
            nsfwEmbed = discord.Embed(title="NSFW", color=0x0096fa, description="This channel is not NSFW. Please switch to an NSFW channel to use the NSFW true option.")
            await ctx.respond(embed=nsfwEmbed)
            return
        elif channelnsfw and nsfw is None:
            nsfw = True 
        elif not channelnsfw and nsfw is None:
            nsfw = False
        rs = r_utils.randomString(length=6)
        timeperiodstring = f"&t={timeperiod}" if timeperiod is not None else ""
        nsfwstring = f"&nsfw={nsfw}"
        try:
            img = r_utils.requestimg(f"https://moescape.jackli.dev/api?_={rs}{timeperiodstring}{nsfwstring}")
        except:
            errorEmbed = discord.Embed(title="Error", url="Could not get image from API. Please try again.", color=0xff524f)
            await ctx.respond(embed=errorEmbed)
            return
        embed = discord.Embed(title="r/moescape", url="https://www.reddit.com/r/moescape/", color=0x0096fa)
        embed.set_image(url=img["imglink"])
        embed.set_footer(text="Powered by https://moescape.jackli.dev/")
        await ctx.respond(embed=embed)

    @commands.slash_command(name="wholesomeyuri", description="Gets a random image from r/wholesomeyuri")
    async def wholesomeyuri(self, ctx, 
        timeperiod: Option(str, "Pick a time period", autocomplete=get_timeperiods) = None,
        nsfw: Option(bool, "NSFW?") = None
    ):
        await ctx.defer()
        channelnsfw = ctx.channel.is_nsfw()
        if nsfw and not channelnsfw:
            nsfwEmbed = discord.Embed(title="NSFW", color=0x0096fa, description="This channel is not NSFW. Please switch to an NSFW channel to use the NSFW true option.")
            await ctx.respond(embed=nsfwEmbed)
            return
        elif channelnsfw and nsfw is None:
            nsfw = True 
        elif not channelnsfw and nsfw is None:
            nsfw = False
        rs = r_utils.randomString(length=6)
        timeperiodstring = f"&t={timeperiod}" if timeperiod is not None else ""
        nsfwstring = f"&nsfw={nsfw}"
        try:
            img = r_utils.requestimg(f"https://wsyuri.jackli.dev/api?_={rs}{timeperiodstring}{nsfwstring}")
        except:
            errorEmbed = discord.Embed(title="Error", url="Could not get image from API. Please try again.", color=0xff524f)
            await ctx.respond(embed=errorEmbed)
            return
        embed = discord.Embed(title="r/wholesomeyuri", url="https://www.reddit.com/r/wholesomeyuri/", color=0x0096fa)
        embed.set_image(url=img["imglink"])
        embed.set_footer(text="Powered by https://wsyuri.jackli.dev/")
        await ctx.respond(embed=embed)
    
    @commands.slash_command(name="awwnime", description="Gets a random image from r/awwnime")
    async def awwnime(self, ctx, 
        timeperiod: Option(str, "Pick a time period", autocomplete=get_timeperiods) = None,
        nsfw: Option(bool, "NSFW?") = None
    ):
        await ctx.defer()
        channelnsfw = ctx.channel.is_nsfw()
        if nsfw and not channelnsfw:
            nsfwEmbed = discord.Embed(title="NSFW", color=0x0096fa, description="This channel is not NSFW. Please switch to an NSFW channel to use the NSFW true option.")
            await ctx.respond(embed=nsfwEmbed)
            return
        elif channelnsfw and nsfw is None:
            nsfw = True 
        elif not channelnsfw and nsfw is None:
            nsfw = False
        rs = r_utils.randomString(length=6)
        timeperiodstring = f"&t={timeperiod}" if timeperiod is not None else ""
        nsfwstring = f"&nsfw={nsfw}"
        try:
            img = r_utils.requestimg(f"https://awwnime.jackli.dev/api?_={rs}{timeperiodstring}{nsfwstring}")
        except:
            errorEmbed = discord.Embed(title="Error", url="Could not get image from API. Please try again.", color=0xff524f)
            await ctx.respond(embed=errorEmbed)
            return
        embed = discord.Embed(title="r/awwnime", url="https://www.reddit.com/r/awwnime/", color=0x0096fa)
        embed.set_image(url=img["imglink"])
        embed.set_footer(text="Powered by https://awwnime.jackli.dev")
        await ctx.respond(embed=embed)

    @commands.slash_command(name="animeirl", description="Gets a random image from r/anime_irl")
    async def animeirl(self, ctx, 
        timeperiod: Option(str, "Pick a time period", autocomplete=get_timeperiods) = None,
        nsfw: Option(bool, "NSFW?") = None
    ):
        await ctx.defer()
        channelnsfw = ctx.channel.is_nsfw()
        if nsfw and not channelnsfw:
            nsfwEmbed = discord.Embed(title="NSFW", color=0x0096fa, description="This channel is not NSFW. Please switch to an NSFW channel to use the NSFW true option.")
            await ctx.respond(embed=nsfwEmbed)
            return
        elif channelnsfw and nsfw is None:
            nsfw = True 
        elif not channelnsfw and nsfw is None:
            nsfw = False
        rs = r_utils.randomString(length=6)
        timeperiodstring = f"&t={timeperiod}" if timeperiod is not None else ""
        nsfwstring = f"&nsfw={nsfw}"
        try:
            img = r_utils.requestimg(f"https://animeirl.jackli.dev/api?_={rs}{timeperiodstring}{nsfwstring}")
        except:
            errorEmbed = discord.Embed(title="Error", url="Could not get image from API. Please try again.", color=0xff524f)
            await ctx.respond(embed=errorEmbed)
            return
        embed = discord.Embed(title="r/anime_irl", url="https://www.reddit.com/r/anime_irl/", color=0x0096fa)
        embed.set_image(url=img["imglink"])
        embed.set_footer(text="Powered by https://animeirl.jackli.dev")
        await ctx.respond(embed=embed)
    
def setup(bot):
    bot.add_cog(Reddit(bot))