package commands

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"

	"github.com/jckli/picsiv/src/utils"
)

func isValidURL(toTest string) bool {
	if toTest == "" {
		return false
	}
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func OnMessageCreate(e *events.MessageCreate) {
	if e.Message.Author.Bot || e.Message.Author.System {
		return
	}

	messageContent := e.Message.Content

	if strings.Contains(messageContent, "pixiv.net") &&
		strings.Contains(messageContent, "artworks") {
		fmt.Println("pixiv link detected")
		urlRaw := regexp.MustCompile(`https?://[^\s]+\d`).FindString(messageContent)
		if !isValidURL(urlRaw) {
			return
		}
		id := regexp.MustCompile(`artworks/(\d+)`).FindStringSubmatch(urlRaw)
		if len(id) < 2 {
			return
		}
		fmt.Println("id: ", id[1])

		illustResp, err := utils.RequestHibiApiIllust(id[1])
		if err != nil || illustResp == nil {
			return
		}
		illust, ok := utils.ParseHibiApiIllust(illustResp)
		if !ok {
			return
		}
		fmt.Println("illust: ", illust)
		if illust.Ugoira {
			fmt.Println("ugoira detected")
			ugoiraResp, err := utils.RequestHibiApiUgoria(id[1])
			if err != nil || ugoiraResp == nil {
				return
			}

			ugoira, err := utils.ParseHibiApiUgoira(ugoiraResp)
			_ = discord.NewFile("ugoira.zip", "ugoria gif", ugoira)
			return
		} else {
			fmt.Println("image detected")
			if len(illust.Urls) > 1 {
				return
			} else {
				fmt.Println("sending image")
				embed := discord.NewEmbedBuilder().
					SetTitle("Full Pixiv Image").
					SetColor(0x0096fa).
					SetImage(illust.Urls[0]).
					Build()
				e.Client().Rest().CreateMessage(e.ChannelID, discord.MessageCreate{
					Embeds: []discord.Embed{embed},
					MessageReference: &discord.MessageReference{
						MessageID: &e.Message.ID,
					},
					AllowedMentions: &discord.AllowedMentions{
						RepliedUser: false,
					},
				})
				return
			}
		}
	}
}
