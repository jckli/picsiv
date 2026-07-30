package utils

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"os"
	"strconv"
	"strings"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/andybons/gogif"
	"github.com/valyala/fasthttp"
)

type HibiApiIllustResponse struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	ImageUrls struct {
		SquareMedium string `json:"square_medium"`
		Medium       string `json:"medium"`
		Large        string `json:"large"`
	} `json:"image_urls"`
	Caption  string `json:"caption"`
	Restrict int    `json:"restrict"`
	User     struct {
		ID               int64  `json:"id"`
		Name             string `json:"name"`
		Account          string `json:"account"`
		ProfileImageUrls struct {
			Medium string `json:"medium"`
		} `json:"profile_image_urls"`
		IsFollowed bool `json:"is_followed"`
	} `json:"user"`
	Tags []struct {
		Name           string `json:"name"`
		TranslatedName string `json:"translated_name"`
	} `json:"tags"`
	Tools          []string `json:"tools"`
	CreateDate     string   `json:"create_date"`
	PageCount      int      `json:"page_count"`
	Width          int      `json:"width"`
	Height         int      `json:"height"`
	SanityLevel    int      `json:"sanity_level"`
	XRestrict      int      `json:"x_restrict"`
	Series         any      `json:"series"`
	MetaSinglePage struct {
		OriginalImageUrl string `json:"original_image_url"`
	} `json:"meta_single_page"`
	MetaPages []struct {
		ImageUrls struct {
			SquareMedium string `json:"square_medium"`
			Medium       string `json:"medium"`
			Large        string `json:"large"`
			Original     string `json:"original"`
		} `json:"image_urls"`
	} `json:"meta_pages"`
	TotalView            int  `json:"total_view"`
	TotalBookmarks       int  `json:"total_bookmarks"`
	IsBookmarked         bool `json:"is_bookmarked"`
	Visible              bool `json:"visible"`
	IsMuted              bool `json:"is_muted"`
	TotalComments        int  `json:"total_comments"`
	IllustAIType         int  `json:"illust_ai_type"`
	IllustBookStyle      int  `json:"illust_book_style"`
	CommentAccessControl int  `json:"comment_access_control"`
}

type ParsedHibiApiIllust struct {
	Nsfw         bool
	Urls         []string
	OriginalUrls []string
	Ugoira       bool
	Caption      string
}

type HibiApiUgoiraResponse struct {
	ZipUrls struct {
		Medium string `json:"medium"`
	} `json:"zip_urls"`
	Frames []struct {
		File  string `json:"file"`
		Delay int    `json:"delay"`
	} `json:"frames"`
}

type PximgApiResponse struct {
	Status int `json:"status"`
	Data   struct {
		Illust string `json:"illust"`
		Nsfw   bool   `json:"nsfw"`
	} `json:"data"`
}

type PixivCache struct {
	Title   string `json:"title"`
	Caption string `json:"caption"`
	Author  struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		Account  string `json:"account"`
		ImageUrl string `json:"image_url"`
	} `json:"author"`
	TotalView      int      `json:"total_view"`
	TotalBookmarks int      `json:"total_bookmarks"`
	Urls           []string `json:"urls"`
	OriginalUrls   []string `json:"original_urls"`
	Tags           string   `json:"tags"`
}

func FormatTags(tags []struct {
	Name           string `json:"name"`
	TranslatedName string `json:"translated_name"`
}) string {
	if len(tags) == 0 {
		return ""
	}
	var result []string
	maxTags := 8
	if len(tags) < maxTags {
		maxTags = len(tags)
	}
	for i := 0; i < maxTags; i++ {
		tagName := tags[i].Name
		tagName = strings.ReplaceAll(tagName, " ", "_")
		result = append(result, "#"+tagName)
	}
	formatted := strings.Join(result, " ")
	if len(formatted) > 180 {
		if idx := strings.LastIndex(formatted[:180], " "); idx != -1 {
			formatted = formatted[:idx]
		} else {
			formatted = formatted[:180]
		}
	}
	return formatted
}

var markdownConverter = func() *converter.Converter {
	c := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
		),
	)
	c.Register.RendererFor("a", converter.TagTypeInline, renderChildrenOnly, converter.PriorityEarly)
	return c
}()

func GetPublicApiUrl() string {
	publicUrl := os.Getenv("PIXIV_PUBLIC_API_URL")
	if publicUrl != "" {
		return strings.TrimSuffix(publicUrl, "/")
	}
	apiUrl := os.Getenv("PIXIV_API_URL")
	if strings.HasPrefix(apiUrl, "http://python-monoapi-service") {
		return "https://pymapi.hayasaka.moe"
	}
	return strings.TrimSuffix(apiUrl, "/")
}

func PrefetchImage(url string) {
	if url == "" {
		return
	}
	go func() {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		req.Header.SetMethod("GET")
		req.SetRequestURI(url)
		_ = client.Do(req, resp)
	}()
}

func RequestPximgApi(mode, date string, nsfw bool) (*PximgApiResponse, error) {
	rs := generateRandomString(10)
	if mode != "" {
		mode = "&mode=" + mode
	}
	if date != "" {
		date = "&date=" + date
	}

	url := "https://pximg.jackli.dev/api" + "?_=" + rs + mode + date + "&nsfw=" + strconv.FormatBool(
		nsfw,
	)
	resp, err := getRequest(url)
	if err != nil {
		return nil, err
	}

	respBody := PximgApiResponse{}
	if err := json.Unmarshal(resp, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}

func ConvertPixivImage(original string) string {
	if strings.Contains(original, "https://i.pximg.net/") {
		path := strings.Split(original, "https://i.pximg.net/")[1]
		return "https://pximg.jackli.dev/" + path
	}
	return original
}

func RequestHibiApiIllust(id string) (*HibiApiIllustResponse, error) {
	url := os.Getenv("PIXIV_API_URL") + "/v1/pixiv/illust/details/" + id
	resp, err := getRequest(url)
	if err != nil {
		return nil, err
	}

	respBody := HibiApiIllustResponse{}
	if err := json.Unmarshal(resp, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}

func RequestHibiApiUgoria(id string) (*HibiApiUgoiraResponse, error) {
	url := os.Getenv("PIXIV_API_URL") + "/v1/pixiv/ugoira_metadata/" + id
	resp, err := getRequest(url)
	if err != nil {
		return nil, err
	}

	respBody := HibiApiUgoiraResponse{}
	if err := json.Unmarshal(resp, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}

func ParseHibiApiIllust(illustResp *HibiApiIllustResponse) (*ParsedHibiApiIllust, bool) {
	if illustResp == nil {
		return nil, false
	}
	mainUrl := "https://pximg.jackli.dev/"
	ugoira := illustResp.Type == "ugoira"
	nsfw := illustResp.SanityLevel >= 5
	urls := []string{}
	originalUrls := []string{}

	if len(illustResp.MetaPages) > 0 {
		for _, page := range illustResp.MetaPages {
			rawOrigUrl := page.ImageUrls.Original
			if rawOrigUrl == "" {
				rawOrigUrl = page.ImageUrls.Large
			}
			if rawOrigUrl == "" {
				rawOrigUrl = page.ImageUrls.Medium
			}
			origPath := strings.Split(rawOrigUrl, "https://i.pximg.net/")[1]
			urls = append(urls, mainUrl+origPath)
			originalUrls = append(originalUrls, mainUrl+origPath)
		}
	} else {
		rawOrigUrl := illustResp.MetaSinglePage.OriginalImageUrl
		if rawOrigUrl == "" {
			rawOrigUrl = illustResp.ImageUrls.Large
		}
		if rawOrigUrl == "" {
			rawOrigUrl = illustResp.ImageUrls.Medium
		}
		if rawOrigUrl != "" {
			origPath := strings.Split(rawOrigUrl, "https://i.pximg.net/")[1]
			urls = append(urls, mainUrl+origPath)
			originalUrls = append(originalUrls, mainUrl+origPath)
		}
	}

	cleanedCaption, err := markdownConverter.ConvertString(illustResp.Caption)
	if err != nil {
		cleanedCaption = illustResp.Caption
	}

	return &ParsedHibiApiIllust{
		Nsfw:         nsfw,
		Urls:         urls,
		OriginalUrls: originalUrls,
		Ugoira:       ugoira,
		Caption:      cleanedCaption,
	}, true
}

func ParseHibiApiUgoira(ugoiraResp *HibiApiUgoiraResponse) (*bytes.Buffer, error) {
	rawZipUrl := ugoiraResp.ZipUrls.Medium
	path := strings.Split(rawZipUrl, "https://i.pximg.net/")[1]
	mirrorUrl := "https://pximg.jackli.dev/" + path

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(mirrorUrl)
	if err := client.Do(req, resp); err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(resp.Body()), int64(len(resp.Body())))
	if err != nil {
		return nil, err
	}

	if len(zipReader.File) == 0 {
		return nil, fmt.Errorf("zip file contains no frames")
	}

	f0, err := zipReader.File[0].Open()
	if err != nil {
		return nil, err
	}
	img0, err := jpeg.Decode(f0)
	if err != nil {
		f0.Close()
		return nil, err
	}
	f0.Close()

	b0 := img0.Bounds()
	quantizer := gogif.MedianCutQuantizer{NumColor: 256}
	globalPaletted := image.NewPaletted(b0, nil)
	quantizer.Quantize(globalPaletted, b0, img0, image.Point{})

	globalPalette := globalPaletted.Palette

	g := &gif.GIF{}
	for i, f := range zipReader.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}

		img, err := jpeg.Decode(rc)
		if err != nil {
			rc.Close()
			return nil, err
		}

		b := img.Bounds()
		palettedImage := image.NewPaletted(b, globalPalette)
		draw.FloydSteinberg.Draw(palettedImage, b, img, image.Point{})
		g.Image = append(g.Image, palettedImage)
		g.Delay = append(g.Delay, ugoiraResp.Frames[i].Delay/10)
		rc.Close()
	}

	gifBuffer := &bytes.Buffer{}
	err = gif.EncodeAll(gifBuffer, g)
	if err != nil {
		return nil, err
	}

	return gifBuffer, nil
}

func renderChildrenOnly(ctx converter.Context, w converter.Writer, node *html.Node) converter.RenderStatus {
	href := dom.GetAttributeOr(node, "href", "")
	w.WriteString(href)
	return converter.RenderSuccess
}
