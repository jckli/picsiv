package dbot

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/paginator"
	"github.com/disgoorg/snowflake/v2"
	"log/slog"
)

type Config struct {
	Token       string
	DevMode     bool
	DevServerID snowflake.ID
}

type Bot struct {
	Client    bot.Client
	Logger    *slog.Logger
	Version   string
	Paginator *paginator.Manager
	Config    Config
}

func New(version string) *Bot {
	devServerID, _ := strconv.Atoi(os.Getenv("DEV_SERVER_ID"))

	logger := slog.Default()
	logger.Info("Starting bot version: " + version)

	return &Bot{
		Client:  nil,
		Logger:  logger,
		Version: version,
		Paginator: paginator.New(
			paginator.WithTimeout(15*time.Minute),
			paginator.WithEmbedColor(0x0096fa),
			paginator.WithButtonsConfig(
				paginator.ButtonsConfig{
					First: nil,
					Back: &paginator.ComponentOptions{
						Emoji: paginator.DefaultConfig().ButtonsConfig.Back.Emoji,
						Label: paginator.DefaultConfig().ButtonsConfig.Back.Label,
						Style: discord.ButtonStyleDanger,
					},
					Stop: nil,
					Next: &paginator.ComponentOptions{
						Emoji: paginator.DefaultConfig().ButtonsConfig.Next.Emoji,
						Label: paginator.DefaultConfig().ButtonsConfig.Next.Label,
						Style: discord.ButtonStyleSuccess,
					},
					Last: nil,
				},
			),
		),
		Config: Config{
			Token:       os.Getenv("TOKEN"),
			DevMode:     os.Getenv("DEV_MODE") == "true",
			DevServerID: snowflake.ID(devServerID),
		},
	}
}

func (b *Bot) Setup(listeners ...bot.EventListener) bot.Client {
	client, err := disgo.New(
		b.Config.Token,
		bot.WithLogger(b.Logger),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		),
		bot.WithEventListeners(listeners...),
		bot.WithEventListeners(b.Paginator),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds),
			cache.WithCaches(cache.FlagMessages),
			cache.WithCaches(cache.FlagChannels),
		),
	)
	if err != nil {
		b.Logger.Error("Error while building DisGo client: ", err)
	}

	return client
}

func (b *Bot) ReadyEvent(_ *events.Ready) {
	err := b.Client.SetPresence(
		context.TODO(),
		gateway.WithWatchingActivity("for pixiv links"),
		gateway.WithOnlineStatus(discord.OnlineStatusOnline),
	)
	if err != nil {
		b.Logger.Error("Error while setting presence: ", err)
	}

	b.Logger.Info("Bot presence set successfully.")
}
