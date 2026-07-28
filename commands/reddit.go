package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/picsiv/utils"
)

var redditCommand = discord.SlashCommandCreate{
	Name:        "reddit",
	Description: "Gets a random post from an art subreddit",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:         "subreddit",
			Description:  "The subreddit to get a post from",
			Required:     true,
			Autocomplete: true,
		},
		discord.ApplicationCommandOptionString{
			Name:         "timeperiod",
			Description:  "The time period to get a post from",
			Required:     false,
			Autocomplete: true,
		},
		discord.ApplicationCommandOptionBool{
			Name:        "nsfw",
			Description: "Allow NSFW posts",
			Required:    false,
		},
	},
}

func RedditAutocompleteHandler(e *handler.AutocompleteEvent) error {
	subredditOption, srOk := e.Data.Option("subreddit")
	timeperiodOption, tpOk := e.Data.Option("timeperiod")
	if srOk && subredditOption.Focused {
		return subredditAutocompleteHandler(e)
	}
	if tpOk && timeperiodOption.Focused {
		return timeperiodAutocompleteHandler(e)
	}
	return e.AutocompleteResult(nil)
}

func subredditAutocompleteHandler(e *handler.AutocompleteEvent) error {
	subreddits := []string{
		"streetmoe",
		"animehoodies",
		"animewallpaper",
		"moescape",
		"wholesomeyuri",
		"awwnime",
		"animeirl",
		"saltsanime",
		"megane",
		"imaginarymaids",
		"cutelittlefangs",
		"pouts",
	}

	choices := make([]discord.AutocompleteChoice, 0, 25)

	str := e.Data.String("subreddit")

	if str == "" {
		for _, subreddit := range subreddits {
			choices = append(choices, discord.AutocompleteChoiceString{
				Name:  subreddit,
				Value: subreddit,
			})
		}
	} else {
		potentialSubreddits := fuzzySearch(subreddits, str)
		for _, subreddit := range potentialSubreddits {
			choices = append(choices, discord.AutocompleteChoiceString{
				Name:  subreddit,
				Value: subreddit,
			})
		}
	}

	return e.AutocompleteResult(choices)
}

func timeperiodAutocompleteHandler(e *handler.AutocompleteEvent) error {
	timeperiods := []string{
		"hour",
		"day",
		"week",
		"month",
		"year",
		"all",
	}

	choices := make([]discord.AutocompleteChoice, 0, 6)

	str := e.Data.String("timeperiod")

	if str == "" {
		for _, timeperiod := range timeperiods {
			choices = append(choices, discord.AutocompleteChoiceString{
				Name:  timeperiod,
				Value: timeperiod,
			})
		}
	} else {
		potentialTimeperiods := fuzzySearch(timeperiods, str)
		for _, timeperiod := range potentialTimeperiods {
			choices = append(choices, discord.AutocompleteChoiceString{
				Name:  timeperiod,
				Value: timeperiod,
			})
		}
	}

	return e.AutocompleteResult(choices)
}

func RedditHandler(e *handler.CommandEvent) error {
	err := e.DeferCreateMessage(false)
	if err != nil {
		return err
	}
	data := e.SlashCommandInteractionData()
	subreddit := data.String("subreddit")
	timeperiod := data.String("timeperiod")
	nsfw := data.Bool("nsfw")

	if channel, ok := e.Channel().MessageChannel.(discord.GuildMessageChannel); ok && nsfw &&
		!channel.NSFW() {
		embed := discord.NewEmbed().
			WithTitle("Error").
			WithDescription("This image is NSFW. Please resend the link in a NSFW channel to view this image.").
			WithColor(0xff524f)
		_, err := e.UpdateInteractionResponse(
			discord.MessageUpdate{
				Embeds: &[]discord.Embed{embed},
			},
		)
		return err
	}

	resp, err := utils.RequestReddit(subreddit, timeperiod, nsfw)
	if err != nil {
		return errorHandler(e)
	}

	if resp.Status != 200 {
		return errorHandler(e)
	}

	embed := discord.NewEmbed().
		WithTitle("r/" + subreddit).
		WithURL(resp.Data.Illust).
		WithColor(0x0096fa).
		WithImage(resp.Data.Illust).
		WithFooterText("Powered by https://reddit.jackli.dev/" + subreddit)

	_, err = e.UpdateInteractionResponse(
		discord.MessageUpdate{
			Embeds: &[]discord.Embed{embed},
		},
	)
	return err
}

func errorHandler(e *handler.CommandEvent) error {
	embed := discord.NewEmbed().
		WithTitle("Error").
		WithDescription("Could not get image from API. Please try again later.").
		WithColor(0xff524f)

	_, err := e.UpdateInteractionResponse(
		discord.MessageUpdate{
			Embeds: &[]discord.Embed{embed},
		},
	)
	return err
}
