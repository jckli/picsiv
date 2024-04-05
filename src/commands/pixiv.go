package commands

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"

	"github.com/jckli/picsiv/src/dbot"
	"github.com/jckli/picsiv/src/utils"

	"github.com/disgoorg/paginator"
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

func OnMessageCreate(e *events.MessageCreate, b *dbot.Bot) {
	if e.Message.Author.Bot || e.Message.Author.System {
		return
	}

	messageContent := e.Message.Content

	if strings.Contains(messageContent, "pixiv.net") &&
		strings.Contains(messageContent, "artworks") {
		urlRaw := regexp.MustCompile(`https?://[^\s]+\d`).FindString(messageContent)
		if !isValidURL(urlRaw) {
			return
		}
		id := regexp.MustCompile(`artworks/(\d+)`).FindStringSubmatch(urlRaw)
		if len(id) < 2 {
			return
		}

		illustResp, err := utils.RequestHibiApiIllust(id[1])
		if err != nil || illustResp.Error != nil {
			return
		}
		illust, ok := utils.ParseHibiApiIllust(illustResp)
		if !ok {
			return
		}

		if illust.Ugoira {
			ugoiraResp, err := utils.RequestHibiApiUgoria(id[1])
			if err != nil || ugoiraResp == nil {
				return
			}

			ugoira, err := utils.ParseHibiApiUgoira(ugoiraResp)
			file := discord.NewFile("ugoira.gif", "", ugoira)
			embed := discord.NewEmbedBuilder().
				SetTitle("Full Pixiv Ugoira").
				SetColor(0x0096fa).
				SetImage("attachment://ugoira.gif").
				Build()
			e.Client().Rest().CreateMessage(e.ChannelID, discord.MessageCreate{
				Embeds: []discord.Embed{embed},
				Files:  []*discord.File{file},
				MessageReference: &discord.MessageReference{
					MessageID: &e.Message.ID,
				},
				AllowedMentions: &discord.AllowedMentions{
					RepliedUser: false,
				},
			})
			return
		} else {
			if len(illust.Urls) > 1 {
				fmt.Println("Creating paginator")
				_, err := b.Paginator.CreateMessage(e.Client(), e.ChannelID, paginator.Pages{
					ID: e.MessageID.String(),
					PageFunc: func(page int, embed *discord.EmbedBuilder) {
						embed.SetTitle("Full Pixiv Images").
							SetImage(illust.Urls[page]).
							SetFooterText(fmt.Sprintf("%d/%d", page+1, len(illust.Urls)))
					},
					Pages:      len(illust.Urls),
					ExpireMode: paginator.ExpireModeAfterLastUsage,
				}, false)
				fmt.Println("Paginator created")
				fmt.Println(err)
				if err != nil {
					return
				}
				return
			} else {
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
