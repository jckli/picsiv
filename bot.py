import discord
from discord.ext import commands
import dotenv
import os

dotenv.load_dotenv()

bot = commands.Bot(intents=discord.Intents(message_content=True, messages=True, guilds=True))
bot.remove_command("help")

for file in os.listdir("./cogs"):
    if file.endswith(".py"):
        name = file[:-3]
        bot.load_extension(f"cogs.{name}")

@bot.event
async def on_ready():
    print(f"Bot is online")
    await bot.change_presence(activity=discord.Activity(type=discord.ActivityType.watching, name="for pixiv links"))

@bot.event
async def on_message(message):
    if message.author == bot.user:
        return

    if ("twitter.com/" in message.content or "x.com/" in message.content) and not message.content.startswith("twitter.com") and not message.content.startswith("x.com"):
        urls = message.content.split()
        modified_urls = []

        for url in urls:
            if "twitter.com" in url or "x.com" in url:
                # Replace twitter.com or x.com with vxtwitter.com
                modified_url = url.replace("twitter.com", "vxtwitter.com").replace("x.com", "vxtwitter.com")
                modified_urls.append(modified_url)
            else:
                modified_urls.append(url)

        modified_content = ' '.join(modified_urls)

        await message.reply(modified_content, allowed_mentions=discord.AllowedMentions.none())

    await bot.process_commands(message)


try:
    bot.run(os.environ.get("TOKEN"))
except Exception as err:
    print(f"Error: {err}")
