import discord
from discord.ext import commands

class Information(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.slash_command(name="help")
    async def help(self, ctx):
        embed = discord.Embed(title="Picsiv Commands", color=0x0096fa)
        embed.add_field(name="help", value="Shows this message.", inline=False)
        embed.add_field(name="streetmoe", value="Get a random image from r/streetmoe", inline=False)
        embed.add_field(name="animehoodies", value="Get a random image from r/animehoodies", inline=False)
        embed.add_field(name="animewallpaper", value="Get a random image from r/animewallpaper", inline=False)
        embed.add_field(name="moescape", value="Get a random image from r/moescape", inline=False)
        await ctx.respond(embed=embed)

def setup(bot):
    bot.add_cog(Information(bot))