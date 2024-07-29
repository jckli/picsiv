package commands

import (
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"

	"github.com/jckli/picsiv/utils"
)

var pixivCommand = discord.SlashCommandCreate{
	Name:        "pixiv",
	Description: "Gets a random post from Pixiv",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:         "sort",
			Description:  "The sort method to get a post from",
			Required:     false,
			Autocomplete: true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "date",
			Description: "Date to fetch posts from, in the format YYYY-MM-DD",
			Required:    false,
		},
		discord.ApplicationCommandOptionBool{
			Name:        "nsfw",
			Description: "Allow NSFW posts",
			Required:    false,
		},
	},
}

func PixivAutocompleteHandler(e *handler.AutocompleteEvent) error {
	sortOption, sOk := e.Data.Option("sort")
	if sOk && sortOption.Focused {
		return pixivSortAutocompleteHandler(e)
	}
	return e.AutocompleteResult(nil)
}

func pixivSortAutocompleteHandler(e *handler.AutocompleteEvent) error {
	sortOptions := []string{
		"day",
		"week",
		"month",
		"day_male",
		"day_female",
		"week_original",
		"week_rookie",
		"day_ai",
		"day_manga",
		"week_manga",
		"month_manga",
		"week_rookie_manga",
	}

	if channel, ok := e.Channel().MessageChannel.(discord.GuildMessageChannel); ok &&
		channel.NSFW() {
		r18Options := []string{
			"day_r18",
			"day_male_r18",
			"day_female_r18",
			"week_r18",
			"week_r18g",
			"day_r18_ai",
			"day_r18_manga",
			"week_r18_manga",
		}
		sortOptions = append(sortOptions, r18Options...)
	}

	choices := make([]discord.AutocompleteChoice, 0, 30)

	str := e.Data.String("sort")

	if str == "" {
		for _, sort := range sortOptions {
			choices = append(choices, discord.AutocompleteChoiceString{
				Name:  sort,
				Value: sort,
			})
		}
	} else {
		potentialSorts := fuzzySearch(sortOptions, str)
		for _, sort := range potentialSorts {
			choices = append(choices, discord.AutocompleteChoiceString{
				Name:  sort,
				Value: sort,
			})
		}
	}

	return e.AutocompleteResult(choices)
}

func PixivRandomHandler(e *handler.CommandEvent) error {
	e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil)
	data := e.SlashCommandInteractionData()
	sort := data.String("sort")
	date := data.String("date")
	nsfw := data.Bool("nsfw")

	if date != "" {
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			embed := discord.NewEmbedBuilder().
				SetTitle("Error").
				SetDescription("Invalid date format. Please use the format YYYY-MM-DD").
				SetColor(0xff524f).
				Build()
			_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
				Embeds: &[]discord.Embed{embed},
			},
			)
			return err
		}
	}

	if channel, ok := e.Channel().MessageChannel.(discord.GuildMessageChannel); ok && nsfw &&
		!channel.NSFW() {
		embed := discord.NewEmbedBuilder().
			SetTitle("Error").
			SetDescription("This image is NSFW. Please resend the link in a NSFW channel to view this image.").
			SetColor(0xff524f).
			Build()
		_, err := e.UpdateInteractionResponse(discord.MessageUpdate{
			Embeds: &[]discord.Embed{embed},
		},
		)
		return err
	}

	if strings.Contains(sort, "r18") {
		nsfw = true
	}

	resp, err := utils.RequestPximgApi(sort, date, nsfw)
	if err != nil {
		return errorHandler(e)
	}

	if resp.Status != 200 {
		return errorHandler(e)
	}

	mirrorImg := utils.ConvertPixivImage(resp.Data.Illust)

	embed := discord.NewEmbedBuilder().
		SetTitle("Random Pixiv Post").
		SetURL(mirrorImg).
		SetColor(0x0096fa).
		SetImage(mirrorImg).
		SetFooterText("Powered by https://pximg.jackli.dev").
		Build()

	_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
		Embeds: &[]discord.Embed{embed},
	},
	)

	return err
}
