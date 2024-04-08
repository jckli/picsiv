package commands

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"time"
)

var startTime = time.Now()

var pingCommand = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "Pong!",
}

func PingHandler(e *handler.CommandEvent) error {
	var ping string
	if e.Client().HasGateway() {
		ping = e.Client().Gateway().Latency().String()
	}

	embed := discord.NewEmbedBuilder().
		SetTitle("Pong! 🏓").
		SetDescription("My ping is " + ping).
		SetColor(0x0096fa).
		SetTimestamp(e.CreatedAt()).
		Build()

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().SetEmbeds(embed).Build(),
	)
}

var infoCommand = discord.SlashCommandCreate{
	Name:        "picsiv",
	Description: "Display basic information about Picsiv",
}

func InfoHandler(e *handler.CommandEvent) error {
	var (
		guildCount  int
		memberCount int
	)
	e.Client().Caches().GuildsForEach(func(guild discord.Guild) {

		guildCount++
		memberCount += guild.MemberCount
	})

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

	botUser, _ := e.Client().Caches().SelfUser()

	description := fmt.Sprintf(
		"Thanks for using Picsiv bot! Any questions can be brought up in the support server. This bot is also open-source! All code can be found on GitHub (Please leave a star ⭐ if you enjoy the bot).\n\nPrivacy Policy: https://picsiv.hayasaka.moe/privacy\n\n**Server Count:** %d\n**User Count:** %d\n**Bot Uptime**: %s",
		guildCount,
		memberCount,
		uptimeStr,
	)

	embed := discord.NewEmbedBuilder().
		SetTitle("Picsiv").
		SetAuthor("Picsiv", "", *botUser.AvatarURL()).
		SetColor(0x0096fa).
		SetDescription(description).
		SetTimestamp(e.CreatedAt()).
		Build()

	var actionRow discord.ActionRowComponent
	actionRow = actionRow.AddComponents(
		discord.NewLinkButton("Support Server", "https://discord.gg/Fr2BhuCkET"),
	)
	actionRow = actionRow.AddComponents(
		discord.NewLinkButton("GitHub", "https://github.com/jckli/picsiv"),
	)

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetEmbeds(embed).
			SetContainerComponents(actionRow).
			Build(),
	)
}
