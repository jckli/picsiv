package commands

import (
	"encoding/json"
	"fmt"
	"math/rand"
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
	if maxPage == "1" || maxPage == "" {
		return nil
	}
	row := discord.NewActionRow().AddComponents(
		discord.NewDangerButton("", "/pixiv/"+id+"/page/"+prevPage).
			WithEmoji(discord.ComponentEmoji{Name: "◀"}).
			WithDisabled(prevPage == "-1"),
		discord.NewSecondaryButton(fmt.Sprintf("%s/%s", curPage, maxPage), "page-counter").
			WithDisabled(true),
		discord.NewSuccessButton("", "/pixiv/"+id+"/page/"+nextPage).
			WithEmoji(discord.ComponentEmoji{Name: "▶"}).
			WithDisabled(nextPage == "-1"),
	)
	return []discord.LayoutComponent{row}
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

		if len(illust.Urls) > 0 {
			utils.PrefetchImage(illust.Urls[0])
		}

		cache := utils.PixivCache{
			Title:   illustResp.Title,
			Caption: illust.Caption,
			Author: struct {
				ID       int64  `json:"id"`
				Name     string `json:"name"`
				Account  string `json:"account"`
				ImageUrl string `json:"image_url"`
			}{
				ID:       illustResp.User.ID,
				Name:     illustResp.User.Name,
				Account:  illustResp.User.Account,
				ImageUrl: illustResp.User.ProfileImageUrls.Medium,
			},
			Urls:           illust.Urls,
			OriginalUrls:   illust.OriginalUrls,
			TotalView:      illustResp.TotalView,
			TotalBookmarks: illustResp.TotalBookmarks,
			Tags:           utils.FormatTags(illustResp.Tags),
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

	if pageInt-1 >= 0 && pageInt-1 < len(c.Urls) {
		utils.PrefetchImage(c.Urls[pageInt-1])
	}

	footerText := c.Tags
	if footerText == "" {
		footerText = "Pixiv"
	}

	embed := discord.NewEmbed().
		WithAuthorName(fmt.Sprintf("%s (@%s)", c.Author.Name, c.Author.Account)).
		WithAuthorURL(fmt.Sprintf("https://www.pixiv.net/users/%d", c.Author.ID)).
		WithAuthorIcon(utils.ConvertPixivImage(c.Author.ImageUrl)).
		WithTitle(c.Title).
		WithURL(fmt.Sprintf("https://www.pixiv.net/artworks/%s", id)).
		WithDescription(c.Caption).
		WithColor(0x0096fa).
		WithImage(c.Urls[pageInt-1]).
		WithFooterText(footerText).
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

	msgUpdate := discord.MessageUpdate{
		Embeds: &[]discord.Embed{embed},
	}
	if len(components) > 0 {
		msgUpdate.Components = &components
	}

	return e.UpdateMessage(msgUpdate)
}

func ParseIllustID(input string) string {
	if match := regexp.MustCompile(`artworks/(\d+)`).FindStringSubmatch(input); len(match) > 1 {
		return match[1]
	}
	if match := regexp.MustCompile(`illust_id=(\d+)`).FindStringSubmatch(input); len(match) > 1 {
		return match[1]
	}
	if match := regexp.MustCompile(`/(\d+)_p\d+`).FindStringSubmatch(input); len(match) > 1 {
		return match[1]
	}
	input = strings.TrimSpace(input)
	if regexp.MustCompile(`^\d+$`).MatchString(input) {
		return input
	}
	return ""
}

func BuildPixivPost(id string, isNSFWChannel bool, b *dbot.Bot) (discord.Embed, *discord.File, []discord.LayoutComponent, error) {
	illustResp, err := utils.RequestHibiApiIllust(id)
	if err != nil {
		b.Logger.Error("Failed to request Hibi API (illust): " + err.Error())
		return discord.Embed{}, nil, nil, fmt.Errorf("Could not contact the API for Pixiv ID: %s", id)
	}

	illust, ok := utils.ParseHibiApiIllust(illustResp)
	if !ok {
		b.Logger.Error("Failed to parse Hibi API response for ID: " + id)
		return discord.Embed{}, nil, nil, fmt.Errorf("Could not parse Hibi API for Pixiv ID: %s", id)
	}

	if len(illust.Urls) > 0 {
		utils.PrefetchImage(illust.Urls[0])
	}

	cache := utils.PixivCache{
		Title:   illustResp.Title,
		Caption: illust.Caption,
		Author: struct {
			ID       int64  `json:"id"`
			Name     string `json:"name"`
			Account  string `json:"account"`
			ImageUrl string `json:"image_url"`
		}{
			ID:       illustResp.User.ID,
			Name:     illustResp.User.Name,
			Account:  illustResp.User.Account,
			ImageUrl: illustResp.User.ProfileImageUrls.Medium,
		},
		Urls:           illust.Urls,
		OriginalUrls:   illust.OriginalUrls,
		TotalView:      illustResp.TotalView,
		TotalBookmarks: illustResp.TotalBookmarks,
		Tags:           utils.FormatTags(illustResp.Tags),
	}
	if jsonByte, err := json.Marshal(cache); err == nil {
		b.Cache.Add(id, string(jsonByte))
	}

	if illust.Nsfw && !isNSFWChannel {
		embed := discord.NewEmbed().
			WithTitle("Error").
			WithDescription("This image is NSFW. Please resend the link in a NSFW channel to view this image.").
			WithColor(0xff524f)
		return embed, nil, nil, nil
	}

	footerText := cache.Tags
	if footerText == "" {
		footerText = "Pixiv"
	}

	if illust.Ugoira {
		ugoiraResp, err := utils.RequestHibiApiUgoria(id)
		if err != nil {
			b.Logger.Error("Failed to request Hibi API (ugoira): " + err.Error())
			return discord.Embed{}, nil, nil, fmt.Errorf("Could not contact the Ugoira API for Pixiv ID: %s", id)
		}

		ugoira, err := utils.ParseHibiApiUgoira(ugoiraResp)
		if err != nil {
			b.Logger.Error("Failed to parse Hibi API ugoira response: " + err.Error())
			return discord.Embed{}, nil, nil, fmt.Errorf("Could not parse Ugoira data from the API for Pixiv ID: %s", id)
		}

		file := discord.NewFile("ugoira.gif", "", ugoira)
		embed := discord.NewEmbed().
			WithAuthorName(fmt.Sprintf("%s (@%s)", illustResp.User.Name, illustResp.User.Account)).
			WithAuthorURL(fmt.Sprintf("https://www.pixiv.net/users/%d", illustResp.User.ID)).
			WithAuthorIcon(utils.ConvertPixivImage(illustResp.User.ProfileImageUrls.Medium)).
			WithTitle(illustResp.Title).
			WithURL(fmt.Sprintf("https://www.pixiv.net/artworks/%s", id)).
			WithDescription(illust.Caption).
			WithColor(0x0096fa).
			WithImage("attachment://ugoira.gif").
			WithFooterText(footerText).
			AddField("👀", strconv.Itoa(illustResp.TotalView), true).
			AddField("🔖", strconv.Itoa(illustResp.TotalBookmarks), true)

		return embed, file, nil, nil
	}

	embed := discord.NewEmbed().
		WithAuthorName(fmt.Sprintf("%s (@%s)", illustResp.User.Name, illustResp.User.Account)).
		WithAuthorURL(fmt.Sprintf("https://www.pixiv.net/users/%d", illustResp.User.ID)).
		WithAuthorIcon(utils.ConvertPixivImage(illustResp.User.ProfileImageUrls.Medium)).
		WithTitle(illustResp.Title).
		WithURL(fmt.Sprintf("https://www.pixiv.net/artworks/%s", id)).
		WithDescription(illust.Caption).
		WithColor(0x0096fa).
		WithImage(illust.Urls[0]).
		WithFooterText(footerText).
		AddField("👀", strconv.Itoa(illustResp.TotalView), true).
		AddField("🔖", strconv.Itoa(illustResp.TotalBookmarks), true)

	var components []discord.LayoutComponent
	if len(illust.Urls) > 1 {
		components = pixivComponents(id, "-1", "2", "1", strconv.Itoa(len(illust.Urls)))
	}

	return embed, nil, components, nil
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

		id := ParseIllustID(urlRaw)
		if id == "" {
			return
		}

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

		embed, file, components, buildErr := BuildPixivPost(id, nsfw, b)
		if buildErr != nil {
			sendErrorReply(b, fmt.Sprintf("%s\nRequester: %s", buildErr.Error(), e.Message.Author.ID.String()))
			return
		}

		msgCreate := discord.NewMessageCreate().
			WithEmbeds(embed).
			WithMessageReferenceByID(e.Message.ID).
			WithAllowedMentions(&discord.AllowedMentions{
				RepliedUser: false,
			})
		if file != nil {
			msgCreate = msgCreate.WithFiles(file)
		}
		if len(components) > 0 {
			msgCreate = msgCreate.WithComponents(components...)
		}

		_, _ = e.Client().Rest.CreateMessage(e.ChannelID, msgCreate)
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

func ParseUserID(input string) string {
	if match := regexp.MustCompile(`users/(\d+)`).FindStringSubmatch(input); len(match) > 1 {
		return match[1]
	}
	if match := regexp.MustCompile(`member\.php\?id=(\d+)`).FindStringSubmatch(input); len(match) > 1 {
		return match[1]
	}
	input = strings.TrimSpace(input)
	if regexp.MustCompile(`^\d+$`).MatchString(input) {
		return input
	}
	return ""
}

func pixivUserComponents(userID string, index int, totalIllusts int, currentIllustID int64, pageCount int) []discord.LayoutComponent {
	prevIdx := strconv.Itoa(index - 1)
	nextIdx := strconv.Itoa(index + 1)

	indicatorText := fmt.Sprintf("Illust %d/%d", index+1, totalIllusts)
	if index == 0 {
		indicatorText = fmt.Sprintf("Illust %d/%d (Latest)", index+1, totalIllusts)
	}

	btnPrev := discord.NewDangerButton("◀ Prev", fmt.Sprintf("/pixiv/user/%s/idx/%s", userID, prevIdx)).
		WithDisabled(index <= 0)
	btnIndicator := discord.NewSecondaryButton(indicatorText, "user-indicator").
		WithDisabled(true)
	btnNext := discord.NewSuccessButton("Next ▶", fmt.Sprintf("/pixiv/user/%s/idx/%s", userID, nextIdx)).
		WithDisabled(index >= totalIllusts-1)

	btnRandom := discord.NewSecondaryButton("🎲 Random", fmt.Sprintf("/pixiv/user/%s/idx/random", userID))

	row1Btns := []discord.InteractiveComponent{btnPrev, btnIndicator, btnNext}
	if pageCount > 1 {
		btnViewPages := discord.NewPrimaryButton(fmt.Sprintf("🖼️ View All %d Images", pageCount), fmt.Sprintf("/pixiv/user/%s/viewillust/%d", userID, currentIllustID))
		row1Btns = append(row1Btns, btnViewPages)
	}
	row1Btns = append(row1Btns, btnRandom)

	row1 := discord.NewActionRow().AddComponents(row1Btns...)
	return []discord.LayoutComponent{row1}
}

func BuildPixivUserPost(userID string, index int, isNSFWChannel bool, b *dbot.Bot) (discord.Embed, []discord.LayoutComponent, error) {
	userResp, err := utils.RequestHibiApiUser(userID)
	if err != nil || userResp == nil {
		return discord.Embed{}, nil, fmt.Errorf("Could not contact the API for Pixiv User ID: %s", userID)
	}

	illusts, err := utils.RequestHibiApiUserIllusts(userID)
	if err != nil {
		return discord.Embed{}, nil, fmt.Errorf("Could not fetch illustrations for Pixiv User ID: %s", userID)
	}

	var filtered []utils.HibiApiIllustResponse
	for _, ill := range illusts {
		if !isNSFWChannel && (ill.SanityLevel >= 5 || ill.XRestrict >= 1) {
			continue
		}
		filtered = append(filtered, ill)
	}

	if len(filtered) == 0 {
		if !isNSFWChannel {
			return discord.Embed{}, nil, fmt.Errorf("All recent artworks by this artist are NSFW. Please run this command in a NSFW channel to view them.")
		}
		return discord.Embed{}, nil, fmt.Errorf("This artist has no illustrations.")
	}

	if index < 0 {
		index = 0
	}
	if index >= len(filtered) {
		index = len(filtered) - 1
	}

	currentIllust := filtered[index]
	illust, ok := utils.ParseHibiApiIllust(&currentIllust)
	if !ok || len(illust.Urls) == 0 {
		return discord.Embed{}, nil, fmt.Errorf("Could not parse illustration data for ID: %d", currentIllust.ID)
	}

	utils.PrefetchImage(illust.Urls[0])

	bio := userResp.User.Comment
	if bio == "" {
		bio = "No bio provided."
	}
	cleanedBio := utils.ConvertMarkdown(bio)

	pubDate := currentIllust.CreateDate
	if len(pubDate) >= 10 {
		pubDate = pubDate[:10]
	}

	embed := discord.NewEmbed().
		WithAuthorName(fmt.Sprintf("%s (@%s)", userResp.User.Name, userResp.User.Account)).
		WithAuthorURL(fmt.Sprintf("https://www.pixiv.net/users/%d", userResp.User.ID)).
		WithAuthorIcon(utils.ConvertPixivImage(userResp.User.ProfileImageUrls.Medium)).
		WithDescription(cleanedBio).
		WithColor(0x0096fa).
		AddField("🎨 Illusts", strconv.Itoa(userResp.Profile.TotalIllusts), true).
		AddField("📚 Manga", strconv.Itoa(userResp.Profile.TotalManga), true).
		AddField("👥 Following", strconv.Itoa(userResp.Profile.TotalFollowUsers), true).
		AddField("Current Illustration", fmt.Sprintf("[%s](https://www.pixiv.net/artworks/%d)", currentIllust.Title, currentIllust.ID), false).
		WithImage(illust.Urls[0]).
		WithFooterText(fmt.Sprintf("👀 %s • 🔖 %s • Published %s", strconv.Itoa(currentIllust.TotalView), strconv.Itoa(currentIllust.TotalBookmarks), pubDate))

	components := pixivUserComponents(userID, index, len(filtered), currentIllust.ID, len(illust.Urls))
	return embed, components, nil
}

func PixivUserButtonHandler(e *handler.ComponentEvent, b *dbot.Bot) error {
	userID := e.Vars["userId"]
	idxStr := e.Vars["index"]

	nsfw := false
	if channel, ok := e.Channel().MessageChannel.(discord.GuildMessageChannel); ok {
		nsfw = channel.NSFW()
	}

	index := 0
	if idxStr == "random" {
		illusts, err := utils.RequestHibiApiUserIllusts(userID)
		if err == nil && len(illusts) > 0 {
			var validIndices []int
			for i, ill := range illusts {
				if nsfw || (ill.SanityLevel < 5 && ill.XRestrict < 1) {
					validIndices = append(validIndices, i)
				}
			}
			if len(validIndices) > 0 {
				index = validIndices[rand.Intn(len(validIndices))]
			}
		}
	} else {
		index, _ = strconv.Atoi(idxStr)
	}

	embed, components, err := BuildPixivUserPost(userID, index, nsfw, b)
	if err != nil {
		return err
	}

	msgUpdate := discord.MessageUpdate{
		Embeds: &[]discord.Embed{embed},
	}
	if len(components) > 0 {
		msgUpdate.Components = &components
	}

	return e.UpdateMessage(msgUpdate)
}

func PixivUserViewIllustButtonHandler(e *handler.ComponentEvent, b *dbot.Bot) error {
	illustID := e.Vars["illustId"]

	nsfw := false
	if channel, ok := e.Channel().MessageChannel.(discord.GuildMessageChannel); ok {
		nsfw = channel.NSFW()
	}

	embed, file, components, buildErr := BuildPixivPost(illustID, nsfw, b)
	if buildErr != nil {
		return e.CreateMessage(discord.NewMessageCreate().
			WithContent(buildErr.Error()).
			WithEphemeral(true))
	}

	msgCreate := discord.NewMessageCreate().
		WithEmbeds(embed)
	if file != nil {
		msgCreate = msgCreate.WithFiles(file)
	}
	if len(components) > 0 {
		msgCreate = msgCreate.WithComponents(components...)
	}

	return e.CreateMessage(msgCreate)
}

