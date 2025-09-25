package utils

import (
	"encoding/json"
)

type SugoiArtResponse struct {
	Status      int    `json:"status"`
	Url         string `json:"url"`
	Sha         string `json:"sha"`
	Orientation string `json:"orientation"`
}

func RequestSugoiArt(orientation string) (*SugoiArtResponse, error) {
	url := "https://api.art.hayasaka.moe/v1/art/random"
	if orientation != "" {
		url += "?o=" + orientation
	}
	resp, err := getRequest(url)
	if err != nil {
		return nil, err
	}

	respBody := SugoiArtResponse{}
	if err := json.Unmarshal(resp, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}
