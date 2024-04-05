package utils

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"github.com/andybons/gogif"
	"image"
	"image/gif"
	"image/jpeg"
	"os"
	"strings"

	"github.com/valyala/fasthttp"
)

type HibiApiIllustResponse struct {
	Illust struct {
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
		Tools          []string    `json:"tools"`
		CreateDate     string      `json:"create_date"`
		PageCount      int         `json:"page_count"`
		Width          int         `json:"width"`
		Height         int         `json:"height"`
		SanityLevel    int         `json:"sanity_level"`
		XRestrict      int         `json:"x_restrict"`
		Series         interface{} `json:"series"`
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
	} `json:"illust"`
}
type ParsedHibiApiIllust struct {
	Nsfw   bool
	Urls   []string
	Ugoira bool
}

type HibiApiUgoiraResponse struct {
	UgoiraMetadata struct {
		ZipUrls struct {
			Medium string `json:"medium"`
		} `json:"zip_urls"`
		Frames []struct {
			File  string `json:"file"`
			Delay int    `json:"delay"`
		} `json:"frames"`
	} `json:"ugoira_metadata"`
}

func RequestHibiApiIllust(id string) (*HibiApiIllustResponse, error) {
	url := os.Getenv("HIBIAPI_URL") + "/api/pixiv/illust?id=" + id
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
	url := os.Getenv("HIBIAPI_URL") + "/api/pixiv/ugoira_metadata?id=" + id
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
	ugoira := illustResp.Illust.Type == "ugoira"
	nsfw := illustResp.Illust.SanityLevel >= 5
	urls := []string{}
	if illustResp.Illust.MetaSinglePage.OriginalImageUrl != "" {
		rawImageUrl := illustResp.Illust.MetaSinglePage.OriginalImageUrl
		path := strings.Split(rawImageUrl, "https://i.pximg.net/")[1]
		mirrorUrl := "https://pximg.jackli.dev/" + path

		urls = append(urls, mirrorUrl)
	} else {
		for _, page := range illustResp.Illust.MetaPages {
			rawImageUrl := page.ImageUrls.Original
			path := strings.Split(rawImageUrl, "https://i.pximg.net/")[1]
			mirrorUrl := "https://pximg.jackli.dev/" + path

			urls = append(urls, mirrorUrl)
		}
	}

	return &ParsedHibiApiIllust{
		Nsfw:   nsfw,
		Urls:   urls,
		Ugoira: ugoira,
	}, true
}

func ParseHibiApiUgoira(ugoiraResp *HibiApiUgoiraResponse) (*bytes.Buffer, error) {
	rawZipUrl := ugoiraResp.UgoiraMetadata.ZipUrls.Medium
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

	frameCount := len(ugoiraResp.UgoiraMetadata.Frames)
	totalDelay := 0
	for _, frame := range ugoiraResp.UgoiraMetadata.Frames {
		totalDelay += frame.Delay
	}
	dur := totalDelay / frameCount

	g := &gif.GIF{}
	for _, f := range zipReader.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		img, err := jpeg.Decode(rc)
		if err != nil {
			return nil, err
		}
		b := img.Bounds()
		palettedImage := image.NewPaletted(b, nil)
		quantizer := gogif.MedianCutQuantizer{NumColor: 64}
		quantizer.Quantize(palettedImage, b, img, image.Point{})
		g.Image = append(g.Image, palettedImage)
		g.Delay = append(g.Delay, dur)
		rc.Close()
	}

	gifBuffer := &bytes.Buffer{}
	err = gif.EncodeAll(gifBuffer, g)
	if err != nil {
		return nil, err
	}

	return gifBuffer, nil

}
