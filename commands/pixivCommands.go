package commands

import (
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"

	"github.com/jckli/picsiv/dbot"
	"github.com/jckli/picsiv/utils"
)

var pixivCommand = discord.SlashCommandCreate{
	Name:        "pixiv",
	Description: "Interact with Pixiv",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "random",
			Description: "Get a random Pixiv post",
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
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "illust",
			Description: "Display a Pixiv illust by ID",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "illust",
					Description: "The Pixiv illust to display (can be URL or ID)",
					Required:    true,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "user",
			Description: "Display a Pixiv user profile and browse their works",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "user",
					Description: "The Pixiv user to display (can be URL or ID)",
					Required:    true,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "ranking",
			Description: "Display top ranked Pixiv illustrations",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:         "mode",
					Description:  "The ranking mode (e.g. day, week, month, rookie)",
					Required:     false,
					Autocomplete: true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "date",
					Description: "Date to fetch ranking from, in YYYY-MM-DD format",
					Required:    false,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "search",
			Description: "Search Pixiv illustrations by tag or keyword",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "query",
					Description: "The tag or keyword to search for",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:         "type",
					Description:  "The search type (tag, exact_tag, title_and_caption)",
					Required:     false,
					Autocomplete: true,
				},
				discord.ApplicationCommandOptionString{
					Name:         "sort",
					Description:  "The sort order (date_desc or popular_desc)",
					Required:     false,
					Autocomplete: true,
				},
				discord.ApplicationCommandOptionBool{
					Name:        "nsfw",
					Description: "Allow NSFW search results",
					Required:    false,
				},
			},
		},
	},
}

func pixivTypeAutocompleteHandler(e *handler.AutocompleteEvent) error {
	choices := []discord.AutocompleteChoice{
		discord.AutocompleteChoiceString{Name: "Tag (Partial)", Value: "tag"},
		discord.AutocompleteChoiceString{Name: "Exact Tag", Value: "exact_tag"},
		discord.AutocompleteChoiceString{Name: "Title & Caption", Value: "title_and_caption"},
	}
	return e.AutocompleteResult(choices)
}

func PixivAutocompleteHandler(e *handler.AutocompleteEvent) error {
	sortOption, sOk := e.Data.Option("sort")
	modeOption, mOk := e.Data.Option("mode")
	typeOption, tOk := e.Data.Option("type")
	if (sOk && sortOption.Focused) || (mOk && modeOption.Focused) {
		return pixivSortAutocompleteHandler(e)
	}
	if tOk && typeOption.Focused {
		return pixivTypeAutocompleteHandler(e)
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

func PixivRandomHandler(e *handler.CommandEvent, b *dbot.Bot) error {
	err := e.DeferCreateMessage(false)
	if err != nil {
		return err
	}
	data := e.SlashCommandInteractionData()
	sort := data.String("sort")
	date := data.String("date")
	nsfw := data.Bool("nsfw")

	if date != "" {
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			embed := discord.NewEmbed().
				WithTitle("Error").
				WithDescription("Invalid date format. Please use the format YYYY-MM-DD").
				WithColor(0xff524f)
			_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
				Embeds: &[]discord.Embed{embed},
			})
			return err
		}
	}

	channelNsfw := false
	if channel, ok := e.Channel().MessageChannel.(discord.GuildMessageChannel); ok {
		channelNsfw = channel.NSFW()
	}

	if nsfw && !channelNsfw {
		embed := discord.NewEmbed().
			WithTitle("Error").
			WithDescription("This image is NSFW. Please resend the link in a NSFW channel to view this image.").
			WithColor(0xff524f)
		_, err := e.UpdateInteractionResponse(discord.MessageUpdate{
			Embeds: &[]discord.Embed{embed},
		})
		return err
	}

	if strings.Contains(sort, "r18") {
		nsfw = true
	}

	resp, err := utils.RequestPximgApi(sort, date, nsfw)
	if err != nil || resp == nil || resp.Status != 200 {
		return errorHandler(e)
	}

	id := ParseIllustID(resp.Data.Illust)
	if id == "" {
		return errorHandler(e)
	}

	embed, file, components, buildErr := BuildPixivPost(id, channelNsfw, b)
	if buildErr != nil {
		embed = discord.NewEmbed().
			WithTitle("Error").
			WithDescription(buildErr.Error()).
			WithColor(0xff524f)
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Embeds: &[]discord.Embed{embed},
		})
		return err
	}

	msgUpdate := discord.MessageUpdate{
		Embeds: &[]discord.Embed{embed},
	}
	if file != nil {
		msgUpdate.Files = []*discord.File{file}
	}
	if len(components) > 0 {
		msgUpdate.Components = &components
	}

	_, err = e.UpdateInteractionResponse(msgUpdate)
	return err
}

func PixivIllustHandler(e *handler.CommandEvent, b *dbot.Bot) error {
	err := e.DeferCreateMessage(false)
	if err != nil {
		return err
	}

	input := e.SlashCommandInteractionData().String("illust")
	id := ParseIllustID(input)
	if id == "" {
		embed := discord.NewEmbed().
			WithTitle("Error").
			WithDescription("Invalid Pixiv URL or Illust ID.").
			WithColor(0xff524f)
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Embeds: &[]discord.Embed{embed},
		})
		return err
	}

	nsfw := false
	if channel, ok := e.Channel().MessageChannel.(discord.GuildMessageChannel); ok {
		nsfw = channel.NSFW()
	}

	embed, file, components, buildErr := BuildPixivPost(id, nsfw, b)
	if buildErr != nil {
		embed = discord.NewEmbed().
			WithTitle("Error").
			WithDescription(buildErr.Error()).
			WithColor(0xff524f)
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Embeds: &[]discord.Embed{embed},
		})
		return err
	}

	msgUpdate := discord.MessageUpdate{
		Embeds: &[]discord.Embed{embed},
	}
	if file != nil {
		msgUpdate.Files = []*discord.File{file}
	}
	if len(components) > 0 {
		msgUpdate.Components = &components
	}

	_, err = e.UpdateInteractionResponse(msgUpdate)
	return err
}

func PixivUserHandler(e *handler.CommandEvent, b *dbot.Bot) error {
	err := e.DeferCreateMessage(false)
	if err != nil {
		return err
	}

	input := e.SlashCommandInteractionData().String("user")
	id := ParseUserID(input)
	if id == "" {
		embed := discord.NewEmbed().
			WithTitle("Error").
			WithDescription("Invalid Pixiv User ID or URL.").
			WithColor(0xff524f)
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Embeds: &[]discord.Embed{embed},
		})
		return err
	}

	nsfw := false
	if channel, ok := e.Channel().MessageChannel.(discord.GuildMessageChannel); ok {
		nsfw = channel.NSFW()
	}

	embed, components, buildErr := BuildPixivUserPost(id, 0, nsfw, b)
	if buildErr != nil {
		embed = discord.NewEmbed().
			WithTitle("Error").
			WithDescription(buildErr.Error()).
			WithColor(0xff524f)
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Embeds: &[]discord.Embed{embed},
		})
		return err
	}

	msgUpdate := discord.MessageUpdate{
		Embeds: &[]discord.Embed{embed},
	}
	if len(components) > 0 {
		msgUpdate.Components = &components
	}

	_, err = e.UpdateInteractionResponse(msgUpdate)
	return err
}

func PixivRankingHandler(e *handler.CommandEvent, b *dbot.Bot) error {
	err := e.DeferCreateMessage(false)
	if err != nil {
		return err
	}

	data := e.SlashCommandInteractionData()
	mode := data.String("mode")
	date := data.String("date")

	if date != "" {
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			embed := discord.NewEmbed().
				WithTitle("Error").
				WithDescription("Invalid date format. Please use the format YYYY-MM-DD").
				WithColor(0xff524f)
			_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
				Embeds: &[]discord.Embed{embed},
			})
			return err
		}
	}

	nsfw := false
	if channel, ok := e.Channel().MessageChannel.(discord.GuildMessageChannel); ok {
		nsfw = channel.NSFW()
	}

	embed, components, buildErr := BuildPixivRankingPost(mode, date, 0, nsfw, b)
	if buildErr != nil {
		embed = discord.NewEmbed().
			WithTitle("Error").
			WithDescription(buildErr.Error()).
			WithColor(0xff524f)
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Embeds: &[]discord.Embed{embed},
		})
		return err
	}

	msgUpdate := discord.MessageUpdate{
		Embeds: &[]discord.Embed{embed},
	}
	if len(components) > 0 {
		msgUpdate.Components = &components
	}

	_, err = e.UpdateInteractionResponse(msgUpdate)
	return err
}

func PixivSearchHandler(e *handler.CommandEvent, b *dbot.Bot) error {
	err := e.DeferCreateMessage(false)
	if err != nil {
		return err
	}

	data := e.SlashCommandInteractionData()
	query := data.String("query")
	typeOpt := data.String("type")
	sortOpt := data.String("sort")
	nsfwOpt, _ := data.OptBool("nsfw")

	nsfw := false
	if channel, ok := e.Channel().MessageChannel.(discord.GuildMessageChannel); ok {
		nsfw = channel.NSFW()
	}

	embed, components, buildErr := BuildPixivSearchPost(query, typeOpt, sortOpt, 0, nsfwOpt, nsfw, b)
	if buildErr != nil {
		embed = discord.NewEmbed().
			WithTitle("Error").
			WithDescription(buildErr.Error()).
			WithColor(0xff524f)
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Embeds: &[]discord.Embed{embed},
		})
		return err
	}

	msgUpdate := discord.MessageUpdate{
		Embeds: &[]discord.Embed{embed},
	}
	if len(components) > 0 {
		msgUpdate.Components = &components
	}

	_, err = e.UpdateInteractionResponse(msgUpdate)
	return err
}


