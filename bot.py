import discord
from discord.ext import commands
import dotenv
import os

dotenv.load_dotenv()

bot = commands.Bot()
bot.remove_command("help")

for file in os.listdir("./cogs"):
    if file.endswith(".py"):
        name = file[:-3]
        bot.load_extension(f"cogs.{name}")

@bot.event
async def on_ready():
    print(f"Bot is online")
    await bot.change_presence(activity=discord.Activity(type=discord.ActivityType.watching, name="pixiv"))

try:
    bot.run(os.environ.get("TOKEN"))
except Exception as err:
    print(f"Error: {err}")