package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/picsiv/utils"
)

var sugoiArtCommand = discord.SlashCommandCreate{
	Name:        "sugoiart",
	Description: "Gets a random post from my sugoiart API",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:         "orientation",
			Description:  "The orientation of the image",
			Required:     false,
			Autocomplete: true,
		},
	},
}

func SugoiArtAutocompleteHandler(e *handler.AutocompleteEvent) error {
	orientationOption, oOk := e.Data.Option("orientation")
	if oOk && orientationOption.Focused {
		return orientationAutocompleteHandler(e)
	}

	return e.AutocompleteResult(nil)
}

func orientationAutocompleteHandler(e *handler.AutocompleteEvent) error {
	orientations := []string{
		"landscape",
		"portrait",
		"square",
	}

	choices := make([]discord.AutocompleteChoice, 0, 6)

	str := e.Data.String("orientation")

	if str == "" {
		for _, orientation := range orientations {
			choices = append(choices, discord.AutocompleteChoiceString{
				Name:  orientation,
				Value: orientation,
			})
		}
	} else {
		potentialOrientations := fuzzySearch(orientations, str)
		for _, orientation := range potentialOrientations {
			choices = append(choices, discord.AutocompleteChoiceString{
				Name:  orientation,
				Value: orientation,
			})
		}
	}

	return e.AutocompleteResult(choices)
}

func SugoiArtHandler(e *handler.CommandEvent) error {
	e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil)
	data := e.SlashCommandInteractionData()
	orientation := data.String("orientation")

	resp, err := utils.RequestSugoiArt(orientation)
	if err != nil {
		return errorHandler(e)
	}

	if resp.Status != 200 {
		return errorHandler(e)
	}

	embed := discord.NewEmbedBuilder().
		SetTitle("SugoiArt").
		SetURL(resp.Url).
		SetColor(0x0096fa).
		SetImage(resp.Url).
		SetFooterText("Powered by https://art.hayasaka.moe").
		Build()

	_, err = e.UpdateInteractionResponse(
		discord.MessageUpdate{
			Embeds: &[]discord.Embed{embed},
		},
	)
	return err

}
