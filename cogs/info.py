import discord
from discord.ext import commands
import time
from datetime import datetime
from datetime import timedelta

startTime = time.time()

class InfoButtons(discord.ui.View):
    def __init__(self, support_server, github):
        super().__init__()
        self.add_item(discord.ui.Button(label="Support Server", url=support_server))
        self.add_item(discord.ui.Button(label="GitHub", url=github))

class Information(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.slash_command(name="help", description="Displays all commands for Picsiv")
    async def help(self, ctx):
        embed = discord.Embed(title="Picsiv Commands", color=0x0096fa)
        embed.add_field(name="help", value="Shows this message.", inline=False)
        embed.add_field(name="picsiv", value="Displays basic information about Picsiv.", inline=False)
        embed.add_field(name="streetmoe", value="Gets a random image from r/streetmoe.", inline=False)
        embed.add_field(name="animehoodies", value="Gets a random image from r/animehoodies.", inline=False)
        embed.add_field(name="animewallpaper", value="Gets a random image from r/animewallpaper.", inline=False)
        embed.add_field(name="moescape", value="Gets a random image from r/moescape.", inline=False)
        await ctx.respond(embed=embed)

    @commands.slash_command(name="picsiv", description="Displays basic information about Picsiv")
    async def picsiv(self, ctx):
        activeServers = self.bot.guilds
        botUsers = 0
        for i in activeServers:
            botUsers += i.member_count
        currentTime = time.time()
        differenceUptime = int(round(currentTime - startTime))
        uptime = str(timedelta(seconds = differenceUptime))
        botinfo = discord.Embed(
            title="Picsiv",
            color=0x0096fa,
            timestamp=datetime.now(),
            description=f"Thanks for using Picsiv bot! Any questions can be brought up in the support server. This bot is also open-source! All code can be found on GitHub (Please leave a star ⭐ if you enjoy the bot).\n\n**Server Count:** {len(self.bot.guilds)}\n**Bot Users:** {botUsers}\n**Bot Uptime:** {uptime}"
        )
        botinfo.set_author(name="Picsiv", icon_url=self.bot.user.avatar.url)
        await ctx.respond(embed=botinfo, view=InfoButtons("https://discord.gg/UcYspqftTF", "https://github.com/jckli/picsiv"))

def setup(bot):
    bot.add_cog(Information(bot))