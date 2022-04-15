import discord
from discord.ext import commands

class Information(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.slash_command(name="help")
    async def help(self, ctx):
        embed = discord.Embed(title="Help", color=0x0096fa, description="There are no commands! Simply send a pixiv link and I will send the full image.")
        await ctx.respond(embed=embed)

def setup(bot):
    bot.add_cog(Information(bot))