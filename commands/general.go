package commands

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/picsiv/dbot"
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

var (
	helpCommand = discord.SlashCommandCreate{
		Name:        "help",
		Description: "Displays all commands",
	}

	cachedHelpEmbed discord.Embed
	helpOnce        sync.Once
)

func HelpHandler(e *handler.CommandEvent, b *dbot.Bot) error {
	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}

	helpOnce.Do(func() {
		cachedHelpEmbed = buildHelpEmbed(b)
	})

	_, err := e.UpdateInteractionResponse(discord.MessageUpdate{
		Embeds: &[]discord.Embed{cachedHelpEmbed},
	})
	return err
}

func buildHelpEmbed(b *dbot.Bot) discord.Embed {
	var sb strings.Builder

	for _, cmd := range CommandList {
		if slashCmd, ok := cmd.(discord.SlashCommandCreate); ok {
			sb.WriteString(fmt.Sprintf("**/%s** - %s\n", slashCmd.Name, slashCmd.Description))
			for _, opt := range slashCmd.Options {
				if sub, ok := opt.(discord.ApplicationCommandOptionSubCommand); ok {
					sb.WriteString(fmt.Sprintf("> `/%s %s` - %s\n", slashCmd.Name, sub.Name, sub.Description))
				}
				if group, ok := opt.(discord.ApplicationCommandOptionSubCommandGroup); ok {
					for _, sub := range group.Options {
						sb.WriteString(fmt.Sprintf("> `/%s %s %s` - %s\n", slashCmd.Name, group.Name, sub.Name, sub.Description))
					}
				}
			}
			sb.WriteString("\n")
		}
	}

	botIcon := ""
	if self, ok := b.Client.Caches.SelfUser(); ok {
		botIcon = self.EffectiveAvatarURL()
	}

	return discord.NewEmbed().
		WithAuthor("Picsiv", "", botIcon).
		WithDescription("Picsiv will automatically respond to all `pixiv.net` links with the full image! There is no setup required.\n\n" + sb.String()).
		WithColor(0x0096fa).
		WithFooterText(fmt.Sprintf("Version: %s", b.Version))
}
