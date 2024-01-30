package utils

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/hbollon/go-edlib"
	"github.com/newtoallofthis123/lyrics_go/db"
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

func SearchableEntires(e []db.Entry, readable bool) (map[string]string, []string) {
	var entries map[string]string = make(map[string]string, len(e))
	var texts []string = make([]string, len(e))

	for _, v := range e {
		if !readable {
			entries[fmt.Sprintf("%s %s", v.Title, v.Artist)] = v.Url
		} else {
			entries[fmt.Sprintf("%s By %s", v.Artist, v.Title)] = v.Url
		}
	}

	for k := range entries {
		texts = append(texts, k)
	}

	return entries, texts
}

func FindExisting(entries []db.Entry, query string) ([]db.Entry, bool) {
	searchable, options := SearchableEntires(entries, false)

	similar, err := edlib.FuzzySearchSetThreshold(query, options, 3, 0.2, edlib.Levenshtein)
	if err != nil {
		panic(err)
	}

	if len(similar) == 0 {
		return nil, false
	}

	var results []db.Entry

	for _, v := range similar {
		entry := searchable[v]
		for _, e := range entries {
			if e.Url == entry {
				results = append(results, e)
			}
		}
	}

	if len(results) == 0 {
		return nil, false
	}

	return results, true
}
