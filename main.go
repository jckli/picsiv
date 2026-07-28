package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jckli/picsiv/commands"
	"github.com/jckli/picsiv/dbot"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	picsiv := dbot.New(os.Getenv("VERSION"))

	h := commands.CommandHandlers(picsiv)

	client := picsiv.Setup(
		bot.WithEventListeners(h),
		bot.WithEventListenerFunc(picsiv.ReadyEvent),
		bot.WithEventListenerFunc(picsiv.OnGuildUpdate),
		bot.WithEventListenerFunc(func(e *events.MessageCreate) {
			commands.OnMessageCreate(e, picsiv)
		}),
	)

	picsiv.Client = client

	cache := picsiv.InitializeCache()
	picsiv.Cache = cache
	ticker := time.NewTicker(time.Minute * 30)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			picsiv.Cache = picsiv.InitializeCache()
		}
	}()

	var guildIDs []snowflake.ID
	if picsiv.Config.DevMode {
		picsiv.Logger.Info(
			fmt.Sprintf(
				"Running in dev mode. Syncing commands to server ID: %s",
				picsiv.Config.DevServerID,
			),
		)
		guildIDs = []snowflake.ID{picsiv.Config.DevServerID}
	} else {
		picsiv.Logger.Info(
			"Running in global mode. Syncing commands globally.",
		)
	}

	err := handler.SyncCommands(client, commands.CommandList, guildIDs)
	if err != nil {
		picsiv.Logger.Error(fmt.Sprintf("Failed to sync commands: %s", err.Error()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.OpenGateway(ctx)
	if err != nil {
		picsiv.Logger.Error("Error while connecting: " + err.Error())
	}
	defer client.Close(context.TODO())

	picsiv.Logger.Info("Bot is now running.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
	picsiv.Logger.Info("Shutting down...")
}
