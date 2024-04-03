package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/picsiv/src/dbot"
)

var CommandList = []discord.ApplicationCommandCreate{}

func CommandHandlers(b *dbot.Bot) *handler.Mux {
	h := handler.New()

	h.Command("/ping", PingHandler)

	return h
}