package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/hbollon/go-edlib"
	"github.com/newtoallofthis123/lyrics_go/db"
	"github.com/newtoallofthis123/lyrics_go/files"
)

const FARSIDE_LINK = "https://farside.link/dumb"

func getInstanceUrl() (string, error) {
	for {
		client := http.Client{}

		resp, err := client.Get(FARSIDE_LINK)
		if err != nil {
			time.Sleep(5 * time.Millisecond)
			continue
		}

		defer resp.Body.Close()

		// test by sending a request to the instance
		_, err = client.Get(fmt.Sprintf("%s/search?q=%s", resp.Request.URL.String(), "test"))
		if err != nil {
			time.Sleep(5 * time.Millisecond) // Wait before retrying
			continue
		}

		return resp.Request.URL.String(), nil
	}
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

func FindExisting(entries []db.Entry, query string) (db.Entry, bool) {
	searchable, options := SearchableEntires(entries, false)

	similar, err := edlib.FuzzySearchThreshold(query, options, 0.2, edlib.Levenshtein)
	if err != nil {
		panic(err)
	}

	if len(similar) == 0 {
		return db.Entry{}, false
	}

	entry := searchable[similar]
	if entry == "" {
		return db.Entry{}, false
	}

	for _, v := range entries {
		if v.Url == entry {
			return v, true
		}
	}

	return db.Entry{}, false
}

func ConvertHTMLToMarkdown(html string) (string, error) {
	converter := md.NewConverter("", true, nil)
	return converter.ConvertString(html)
}

func GetInstance() (string, error) {
	instance := files.JoinHomeDir("INSTANCE")
	if _, err := os.Stat(instance); os.IsNotExist(err) {
		os.Create(instance)

		url, err := getInstanceUrl()
		if err != nil {
			return "", err
		}

		f, err := os.OpenFile(instance, os.O_WRONLY, 0644)
		if err != nil {
			return "", err
		}

		_, err = f.WriteString(fmt.Sprintf("%s|%d\n", url, 10))
		if err != nil {
			return "", err
		}
	}

	f, err := os.Open(instance)
	if err != nil {
		return "", err
	}

	var url string
	var count int
	_, err = fmt.Fscanf(f, "%s|%d", &url, &count)
	if err != nil {
		return "", err
	}

	f.Truncate(0)
	f.Seek(0, 0)

	if count <= 1 {
		url, err = getInstanceUrl()
		if err != nil {
			return "", err
		}
		count = 10
	}

	_, err = f.WriteString(fmt.Sprintf("%s|%d\n", url, count-1))
	if err != nil {
		return "", err
	}

	if err := f.Close(); err != nil {
		return "", err
	}

	return path.Join(instance, "lyrics.db"), nil
}
