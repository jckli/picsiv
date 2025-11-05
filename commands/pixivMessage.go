package commands

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"

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
	var c utils.PixivCache
	if found {
		err := json.Unmarshal([]byte(resp), &c)
		if err != nil {
			return err
		}
	} else {
		illustResp, err := utils.RequestHibiApiIllust(id)
		if err != nil {
			b.Logger.Error("Failed to request Hibi API (illust): " + err.Error())
			sendErrorReply(b, fmt.Sprintf("Could not contact the API for Pixiv ID: %s\nRequester: %s", id, e.Message.Author.ID.String()))
			return err
		}
		illust, ok := utils.ParseHibiApiIllust(illustResp)
		if !ok {
			b.Logger.Error("Failed to parse Hibi API response for ID: " + id)
			sendErrorReply(b, fmt.Sprintf("Could not parse Hibi API for Pixiv ID: %s\nRequester: %s", id, e.Message.Author.ID.String()))
			return fmt.Errorf("Failed to parse illust.")
		}

		cache := utils.PixivCache{
			Title:   illustResp.Title,
			Caption: illust.Caption,
			Author: struct {
				Name     string `json:"name"`
				Account  string `json:"account"`
				ImageUrl string `json:"image_url"`
			}{
				Name:     illustResp.User.Name,
				Account:  illustResp.User.Account,
				ImageUrl: illustResp.User.ProfileImageUrls.Medium,
			},
			Urls:           illust.Urls,
			TotalView:      illustResp.TotalView,
			TotalBookmarks: illustResp.TotalBookmarks,
		}

		jsonByte, err := json.Marshal(cache)
		if err != nil {
			return err
		}

		jsonString := string(jsonByte)

		b.Cache.Add(id, jsonString)

		c = cache
	}

	pageInt, _ := strconv.Atoi(page)

	embed := discord.NewEmbedBuilder().
		SetAuthorName(fmt.Sprintf("%s (@%s)", c.Author.Name, c.Author.Account)).
		SetAuthorIcon(utils.ConvertPixivImage(c.Author.ImageUrl)).
		SetTitle(c.Title).
		SetDescription(c.Caption).
		SetColor(0x0096fa).
		SetImage(c.Urls[pageInt-1]).
		AddField("👀", strconv.Itoa(c.TotalView), true).
		AddField("🔖", strconv.Itoa(c.TotalBookmarks), true).
		Build()

	maxPage := len(c.Urls)
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
		if !isValidURL(urlRaw) || strings.Contains(messageContent, "<"+urlRaw+">") {
			return
		}

		id := regexp.MustCompile(`artworks/(\d+)`).FindStringSubmatch(urlRaw)
		if len(id) < 2 {
			return
		}

		illustResp, err := utils.RequestHibiApiIllust(id[1])
		if err != nil {
			b.Logger.Error("Failed to request Hibi API (illust): " + err.Error())
			sendErrorReply(b, fmt.Sprintf("Could not contact the API for Pixiv ID: %s\nRequester: %s", id[1], e.Message.Author.ID.String()))
			return
		}

		illust, ok := utils.ParseHibiApiIllust(illustResp)
		if !ok {
			b.Logger.Error("Failed to parse Hibi API response for ID: " + id[1])
			sendErrorReply(b, fmt.Sprintf("Could not parse Hibi API for Pixiv ID: %s\nRequester: %s", id[1], e.Message.Author.ID.String()))
			return
		}

		if illust.Nsfw {
			channel, ok := e.Channel()
			if !ok {
				return
			}
			nsfw := channel.NSFW()
			if channel.Type() == 11 {
				gChannel, ok := e.Client().Caches().GuildMessageChannel(*channel.ParentID())
				if !ok {
					return
				}
				nsfw = gChannel.NSFW()

			}
			if !nsfw {
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
			if err != nil {
				b.Logger.Error("Failed to request Hibi API (ugoira): " + err.Error())
				sendErrorReply(b, fmt.Sprintf("Could not contact the Ugoira API for Pixiv ID: %s\nRequester: %s", id[1], e.Message.Author.ID.String()))
				return
			}

			ugoira, err := utils.ParseHibiApiUgoira(ugoiraResp)
			if err != nil {
				b.Logger.Error("Failed to parse Hibi API ugoira response: " + err.Error())
				sendErrorReply(b, fmt.Sprintf("Could not parse Ugoira data from the API for Pixiv ID: %s\nRequester: %s", id[1], e.Message.Author.ID.String()))
				return
			}

			file := discord.NewFile("ugoira.gif", "", ugoira)
			embed := discord.NewEmbedBuilder().
				SetAuthorName(fmt.Sprintf("%s (@%s)", illustResp.User.Name, illustResp.User.Account)).
				SetAuthorIcon(utils.ConvertPixivImage(illustResp.User.ProfileImageUrls.Medium)).
				SetTitle(illustResp.Title).
				SetDescription(illust.Caption).
				SetColor(0x0096fa).
				SetImage("attachment://ugoira.gif").
				AddField("👀", strconv.Itoa(illustResp.TotalView), true).
				AddField("🔖", strconv.Itoa(illustResp.TotalBookmarks), true).
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
					SetAuthorName(fmt.Sprintf("%s (@%s)", illustResp.User.Name, illustResp.User.Account)).
					SetAuthorIcon(utils.ConvertPixivImage(illustResp.User.ProfileImageUrls.Medium)).
					SetTitle(illustResp.Title).
					SetDescription(illust.Caption).
					SetImage(illust.Urls[0]).
					SetColor(0x0096fa).
					AddField("👀", strconv.Itoa(illustResp.TotalView), true).
					AddField("🔖", strconv.Itoa(illustResp.TotalBookmarks), true).
					Build()
				components := pixivComponents(id[1], "-1", "2", "1", strconv.Itoa(len(illust.Urls)))

				e.Client().Rest().CreateMessage(e.ChannelID, discord.MessageCreate{
					Embeds: []discord.Embed{embed},
					MessageReference: &discord.MessageReference{
						MessageID: &e.Message.ID,
					},
					Components: components,
					AllowedMentions: &discord.AllowedMentions{
						RepliedUser: false,
					},
				})
				return
			} else {
				embed := discord.NewEmbedBuilder().
					SetAuthorName(fmt.Sprintf("%s (@%s)", illustResp.User.Name, illustResp.User.Account)).
					SetAuthorIcon(utils.ConvertPixivImage(illustResp.User.ProfileImageUrls.Medium)).
					SetTitle(illustResp.Title).
					SetDescription(illust.Caption).
					SetColor(0x0096fa).
					SetImage(illust.Urls[0]).
					AddField("👀", strconv.Itoa(illustResp.TotalView), true).
					AddField("🔖", strconv.Itoa(illustResp.TotalBookmarks), true).
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

func sendErrorReply(b *dbot.Bot, message string) {
	id := snowflake.GetEnv("DEV_ERROR_CHANNEL_ID")
	embed := discord.NewEmbedBuilder().
		SetTitle("Error").
		SetDescription(message).
		SetColor(0xff524f).
		Build()

	b.Client.Rest().CreateMessage(id, discord.MessageCreate{
		Embeds: []discord.Embed{embed},
	})
}
