package utils

import (
	"crypto/rand"
	"encoding/json"
	"strconv"
)

type RedditResponse struct {
	Status int `json:"status"`
	Data   struct {
		Illust string `json:"illust"`
		Nsfw   bool   `json:"nsfw"`
	} `json:"data"`
}

func RequestReddit(subreddit, timeperiod string, nsfw bool) (*RedditResponse, error) {
	rs := generateRandomString(10)
	if timeperiod != "" {
		timeperiod = "&timeperiod=" + timeperiod
	}
	url := "https://reddit.jackli.dev/" + subreddit + "/api?_=" + rs + timeperiod + "&nsfw=" + strconv.FormatBool(
		nsfw,
	)
	resp, err := getRequest(url)
	if err != nil {
		return nil, err
	}

	respBody := RedditResponse{}
	if err := json.Unmarshal(resp, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return string(b)
}
