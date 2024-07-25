package commands

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"

	"github.com/jckli/picsiv/dbot"
	"github.com/jckli/picsiv/utils"
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

func pixivComponents(
	id, prevPage, nextPage, curPage, maxPage string,
) []discord.ContainerComponent {
	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.NewDangerButton("", "/pixiv/"+id+"/page/"+prevPage).
				WithEmoji(discord.ComponentEmoji{Name: "◀"}).
				WithDisabled(prevPage == "-1"),
			discord.NewSecondaryButton(fmt.Sprintf("%s/%s", curPage, maxPage), "page-counter").
				WithDisabled(true),
			discord.NewSuccessButton("", "/pixiv/"+id+"/page/"+nextPage).
				WithEmoji(discord.ComponentEmoji{Name: "▶"}).
				WithDisabled(nextPage == "-1"),
		},
	}
}

func PixivButtonHandler(e *handler.ComponentEvent, b *dbot.Bot) error {
	id := e.Variables["id"]
	page := e.Variables["page"]

	resp, found := b.Cache.Get(id)
	var urls []string
	if found {
		urls = strings.Split(resp, ",")
	} else {
		illustResp, err := utils.RequestHibiApiIllust(id)
		if err != nil {
			return err
		}
		illust, ok := utils.ParseHibiApiIllust(illustResp)
		if !ok {
			return fmt.Errorf("Failed to parse illust.")
		}

		urlString := strings.Join(illust.Urls, ",")
		b.Cache.Add(id, urlString)

		urls = illust.Urls
	}

	pageInt, _ := strconv.Atoi(page)

	embed := discord.NewEmbedBuilder().
		SetTitle("Full Pixiv Images").
		SetColor(0x0096fa).
		SetImage(urls[pageInt-1]).
		Build()

	maxPage := len(urls)
	prevPage := strconv.Itoa(pageInt - 1)
	nextPage := strconv.Itoa(pageInt + 1)
	maxPageStr := strconv.Itoa(maxPage)

	if pageInt == 1 {
		prevPage = "-1"
	}
	if pageInt == maxPage {
		nextPage = "-1"
	}

	components := pixivComponents(id, prevPage, nextPage, page, maxPageStr)

	e.UpdateMessage(discord.MessageUpdate{
		Embeds:     &[]discord.Embed{embed},
		Components: &components,
	})

	return nil
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

		if illust.Nsfw {
			channel, ok := e.Channel()
			if !ok {
				return
			}

			if !channel.NSFW() {
				embed := discord.NewEmbedBuilder().
					SetTitle("Error").
					SetDescription("This image is NSFW. Please resend the link in a NSFW channel to view this image.").
					SetColor(0xff524f).
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
				embed := discord.NewEmbedBuilder().
					SetTitle("Full Pixiv Images").
					SetImage(illust.Urls[0]).
					SetColor(0x0096fa).
					Build()
				components := pixivComponents(id[1], "-1", "2", "1", strconv.Itoa(len(illust.Urls)))

				e.Client().Rest().CreateMessage(e.ChannelID, discord.MessageCreate{
					Embeds:     []discord.Embed{embed},
					Components: components,
				})
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
