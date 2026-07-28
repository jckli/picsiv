package commands

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var startTime = time.Now()

var pingCommand = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "Pong!",
}

func PingHandler(e *handler.CommandEvent) error {
	var ping string
	if e.Client().HasGateway() {
		ping = e.Client().Gateway.Latency().String()
	}

	embed := discord.NewEmbed().
		WithTitle("Pong! 🏓").
		WithDescription("My ping is " + ping).
		WithColor(0x0096fa).
		WithTimestamp(e.CreatedAt())

	return e.CreateMessage(
		discord.NewMessageCreate().WithEmbeds(embed),
	)
}

var infoCommand = discord.SlashCommandCreate{
	Name:        "picsiv",
	Description: "Displays basic information about Picsiv",
}

func InfoHandler(e *handler.CommandEvent) error {
	var (
		guildCount  int
		memberCount int
	)
	for guild := range e.Client().Caches.Guilds() {
		guildCount++
		memberCount += guild.MemberCount
	}

	uptime := time.Since(startTime)
	days := uptime / (24 * time.Hour)
	uptime %= 24 * time.Hour
	hours := uptime / time.Hour
	uptime %= time.Hour
	minutes := uptime / time.Minute
	uptime %= time.Minute
	seconds := uptime / time.Second

	var uptimeStr string
	if days > 0 {
		uptimeStr = fmt.Sprintf("%d days, ", days)
	}
	uptimeStr += fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	botUser, _ := e.Client().Caches.SelfUser()

	description := fmt.Sprintf(
		"Thanks for using Picsiv bot! Any questions can be brought up in the support server. This bot is also open-source! All code can be found on GitHub (Please leave a star ⭐ if you enjoy the bot).\n\nPrivacy Policy: https://picsiv.hayasaka.moe/privacy\n\n**Server Count:** %d\n**User Count:** %d\n**Bot Uptime**: %s",
		guildCount,
		memberCount,
		uptimeStr,
	)

	embed := discord.NewEmbed().
		WithTitle("Picsiv").
		WithAuthor("Picsiv", "", botUser.EffectiveAvatarURL()).
		WithColor(0x0096fa).
		WithDescription(description).
		WithTimestamp(e.CreatedAt())

	actionRow := discord.NewActionRow(
		discord.NewLinkButton("Support Server", "https://discord.gg/Fr2BhuCkET"),
		discord.NewLinkButton("GitHub", "https://github.com/jckli/picsiv"),
	)

	return e.CreateMessage(
		discord.NewMessageCreate().
			WithEmbeds(embed).
			WithComponents(actionRow),
	)
}

var helpCommand = discord.SlashCommandCreate{
	Name:        "help",
	Description: "Displays all commands",
}

func HelpHandler(e *handler.CommandEvent) error {
	botUser, _ := e.Client().Caches.SelfUser()

	description := fmt.Sprintf(
		"**ping**: Pong! Shows the current ping of Picsiv.\n**picsiv**: Displays basic information about Picsiv.\n**help**: Displays all commands.\n**reddit**: Gets a random post from an art subreddit.",
	)

	embed := discord.NewEmbed().
		WithTitle("Picsiv Commands").
		WithAuthor("Picsiv", "", botUser.EffectiveAvatarURL()).
		WithColor(0x0096fa).
		WithDescription("Picsiv will automatically respond to all `pixiv.net` links with the full image! There is no setup required.").
		AddField("Commands", description, false)

	return e.CreateMessage(
		discord.NewMessageCreate().WithEmbeds(embed),
	)
}
