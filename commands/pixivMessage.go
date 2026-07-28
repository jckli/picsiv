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
) []discord.LayoutComponent {
	return []discord.LayoutComponent{
		discord.NewActionRow(
			discord.NewDangerButton("", "/pixiv/"+id+"/page/"+prevPage).
				WithEmoji(discord.ComponentEmoji{Name: "◀"}).
				WithDisabled(prevPage == "-1"),
			discord.NewSecondaryButton(fmt.Sprintf("%s/%s", curPage, maxPage), "page-counter").
				WithDisabled(true),
			discord.NewSuccessButton("", "/pixiv/"+id+"/page/"+nextPage).
				WithEmoji(discord.ComponentEmoji{Name: "▶"}).
				WithDisabled(nextPage == "-1"),
		),
	}
}

func PixivButtonHandler(e *handler.ComponentEvent, b *dbot.Bot) error {
	id := e.Vars["id"]
	page := e.Vars["page"]

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

	embed := discord.NewEmbed().
		WithAuthorName(fmt.Sprintf("%s (@%s)", c.Author.Name, c.Author.Account)).
		WithAuthorIcon(utils.ConvertPixivImage(c.Author.ImageUrl)).
		WithTitle(c.Title).
		WithDescription(c.Caption).
		WithColor(0x0096fa).
		WithImage(c.Urls[pageInt-1]).
		AddField("👀", strconv.Itoa(c.TotalView), true).
		AddField("🔖", strconv.Itoa(c.TotalBookmarks), true)

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

	return e.UpdateMessage(discord.MessageUpdate{
		Embeds:     &[]discord.Embed{embed},
		Components: &components,
	})
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
			if thread, ok := channel.(discord.GuildThread); ok && thread.ParentID() != nil {
				gChannel, ok := e.Client().Caches.GuildMessageChannel(*thread.ParentID())
				if ok {
					nsfw = gChannel.NSFW()
				}
			}
			if !nsfw {
				embed := discord.NewEmbed().
					WithTitle("Error").
					WithDescription("This image is NSFW. Please resend the link in a NSFW channel to view this image.").
					WithColor(0xff524f)
				_, _ = e.Client().Rest.CreateMessage(e.ChannelID, discord.NewMessageCreate().
					WithEmbeds(embed).
					WithMessageReferenceByID(e.Message.ID).
					WithAllowedMentions(&discord.AllowedMentions{
						RepliedUser: false,
					}),
				)
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
			embed := discord.NewEmbed().
				WithAuthorName(fmt.Sprintf("%s (@%s)", illustResp.User.Name, illustResp.User.Account)).
				WithAuthorIcon(utils.ConvertPixivImage(illustResp.User.ProfileImageUrls.Medium)).
				WithTitle(illustResp.Title).
				WithDescription(illust.Caption).
				WithColor(0x0096fa).
				WithImage("attachment://ugoira.gif").
				AddField("👀", strconv.Itoa(illustResp.TotalView), true).
				AddField("🔖", strconv.Itoa(illustResp.TotalBookmarks), true)
			_, _ = e.Client().Rest.CreateMessage(e.ChannelID, discord.NewMessageCreate().
				WithEmbeds(embed).
				WithFiles(file).
				WithMessageReferenceByID(e.Message.ID).
				WithAllowedMentions(&discord.AllowedMentions{
					RepliedUser: false,
				}),
			)
			return
		} else {
			if len(illust.Urls) > 1 {
				embed := discord.NewEmbed().
					WithAuthorName(fmt.Sprintf("%s (@%s)", illustResp.User.Name, illustResp.User.Account)).
					WithAuthorIcon(utils.ConvertPixivImage(illustResp.User.ProfileImageUrls.Medium)).
					WithTitle(illustResp.Title).
					WithDescription(illust.Caption).
					WithImage(illust.Urls[0]).
					WithColor(0x0096fa).
					AddField("👀", strconv.Itoa(illustResp.TotalView), true).
					AddField("🔖", strconv.Itoa(illustResp.TotalBookmarks), true)
				components := pixivComponents(id[1], "-1", "2", "1", strconv.Itoa(len(illust.Urls)))

				_, _ = e.Client().Rest.CreateMessage(e.ChannelID, discord.NewMessageCreate().
					WithEmbeds(embed).
					WithComponents(components...).
					WithMessageReferenceByID(e.Message.ID).
					WithAllowedMentions(&discord.AllowedMentions{
						RepliedUser: false,
					}),
				)
				return
			} else {
				embed := discord.NewEmbed().
					WithAuthorName(fmt.Sprintf("%s (@%s)", illustResp.User.Name, illustResp.User.Account)).
					WithAuthorIcon(utils.ConvertPixivImage(illustResp.User.ProfileImageUrls.Medium)).
					WithTitle(illustResp.Title).
					WithDescription(illust.Caption).
					WithColor(0x0096fa).
					WithImage(illust.Urls[0]).
					AddField("👀", strconv.Itoa(illustResp.TotalView), true).
					AddField("🔖", strconv.Itoa(illustResp.TotalBookmarks), true)
				_, _ = e.Client().Rest.CreateMessage(e.ChannelID, discord.NewMessageCreate().
					WithEmbeds(embed).
					WithMessageReferenceByID(e.Message.ID).
					WithAllowedMentions(&discord.AllowedMentions{
						RepliedUser: false,
					}),
				)
				return
			}
		}
	}
}

func sendErrorReply(b *dbot.Bot, message string) {
	id := snowflake.GetEnv("DEV_ERROR_CHANNEL_ID")
	if id == 0 {
		return
	}
	embed := discord.NewEmbed().
		WithTitle("Error").
		WithDescription(message).
		WithColor(0xff524f)

	_, _ = b.Client.Rest.CreateMessage(id, discord.NewMessageCreate().WithEmbeds(embed))
}
