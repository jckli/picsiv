package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/picsiv/dbot"
)

var CommandList = []discord.ApplicationCommandCreate{
	helpCommand,
	pingCommand,
	infoCommand,
	redditCommand,
}

func CommandHandlers(b *dbot.Bot) *handler.Mux {
	h := handler.New()

	h.Command("/help", HelpHandler)
	h.Command("/ping", PingHandler)
	h.Command("/picsiv", InfoHandler)
	h.Route("/reddit", func(h handler.Router) {
		h.Command("/", RedditHandler)
		h.Autocomplete("/", RedditAutocompleteHandler)
	})

	h.Route("/pixiv", func(h handler.Router) {
		// h.Command("/", PixivHandler)

		h.Component("/{id}/page/{page}", func(e *handler.ComponentEvent) error {
			return PixivButtonHandler(e, b)
		})
	})

	return h
}
