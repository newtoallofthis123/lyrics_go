package utils

import (
	"net/http"
	"net/url"
)

const FARSIDE_LINK = "https://farside.link/dumb"

func GetInstance() (string, error) {
	client := http.Client{}

	resp, err := client.Get(FARSIDE_LINK)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	return resp.Request.URL.String(), nil
}

func ConvertToQuery(query string) string {
	return url.QueryEscape(query)
}
