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
	"github.com/disgoorg/snowflake/v2"
	"log/slog"

	lru "github.com/hashicorp/golang-lru/v2/expirable"
)

type Config struct {
	Token       string
	DevMode     bool
	DevServerID snowflake.ID
}

type Bot struct {
	Client  bot.Client
	Logger  *slog.Logger
	Version string
	Cache   *lru.LRU[string, string]
	Config  Config
}

func New(version string) *Bot {
	devServerID, _ := strconv.Atoi(os.Getenv("DEV_SERVER_ID"))

	logger := slog.Default()
	logger.Info("Starting bot version: " + version)

	return &Bot{
		Client:  nil,
		Logger:  logger,
		Version: version,
		Cache:   nil,
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

func (b *Bot) InitializeCache() *lru.LRU[string, string] {
	c := lru.NewLRU[string, string](300, nil, time.Minute*30)
	b.Logger.Info("Cache initialized.")

	return c
}
