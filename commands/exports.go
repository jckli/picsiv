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
	pixivCommand,
	sugoiArtCommand,
}

func CommandHandlers(b *dbot.Bot) *handler.Mux {
	h := handler.New()

	h.Command("/help", func(e *handler.CommandEvent) error {
		return HelpHandler(e, b)
	})
	h.Command("/ping", PingHandler)
	h.Command("/picsiv", InfoHandler)

	h.Route("/pixiv", func(h handler.Router) {
		h.Command("/random", func(e *handler.CommandEvent) error {
			return PixivRandomHandler(e, b)
		})
		h.Command("/illust", func(e *handler.CommandEvent) error {
			return PixivIllustHandler(e, b)
		})
		h.Command("/user", func(e *handler.CommandEvent) error {
			return PixivUserHandler(e, b)
		})
		h.Command("/ranking", func(e *handler.CommandEvent) error {
			return PixivRankingHandler(e, b)
		})
		h.Autocomplete("/random", PixivAutocompleteHandler)
		h.Autocomplete("/ranking", PixivAutocompleteHandler)

		h.Component("/{id}/page/{page}", func(e *handler.ComponentEvent) error {
			return PixivButtonHandler(e, b)
		})
		h.Component("/user/{userId}/idx/{index}", func(e *handler.ComponentEvent) error {
			return PixivUserButtonHandler(e, b)
		})
		h.Component("/user/{userId}/viewillust/{illustId}", func(e *handler.ComponentEvent) error {
			return PixivUserViewIllustButtonHandler(e, b)
		})
		h.Component("/ranking/{mode}/{date}/idx/{index}", func(e *handler.ComponentEvent) error {
			return PixivRankingButtonHandler(e, b)
		})
	})

	h.Route("/reddit", func(h handler.Router) {
		h.Command("/", RedditHandler)
		h.Autocomplete("/", RedditAutocompleteHandler)
	})

	h.Route("/sugoiart", func(h handler.Router) {
		h.Command("/", SugoiArtHandler)
		h.Autocomplete("/", SugoiArtAutocompleteHandler)
	})

	return h
}
