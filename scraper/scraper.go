package scraper

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/newtoallofthis123/lyrics_go/utils"
)

type GetLyrics struct {
	Instance string
	Query    string
}

type SearchResult struct {
	Artist string
	Title  string
	Url    string
}

func (lr *GetLyrics) ParseSearchPage() ([]SearchResult, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("%s/search?q=%s", lr.Instance, utils.ConvertToQuery(lr.Query)))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var lyrics []SearchResult

	doc.Find("#search-results").Each(func(i int, s *goquery.Selection) {
		items := s.Find("#search-item")

		items.Each(func(i int, si *goquery.Selection) {
			artist := si.Find("div span").Text()
			title := si.Find("div h2").Text()
			url, _ := si.Attr("href")

			lyrics = append(lyrics, SearchResult{
				Artist: artist,
				Title:  title,
				Url:    url,
			})
		})
	})

	return lyrics, nil
}

type Lyrics struct {
	Artist string
	Title  string
	Url    string
	Lyrics string
}

func (sr *SearchResult) GetSongLyrics() (Lyrics, error) {
	client := http.Client{
		Timeout: 20 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("%s%s", utils.FARSIDE_LINK, sr.Url))
	if err != nil {
		return Lyrics{}, err
	}

	if resp.StatusCode != 200 {
		return Lyrics{}, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Lyrics{}, err
	}

	defer resp.Body.Close()

	var lyrics string

	doc.Find("#lyrics").Each(func(i int, s *goquery.Selection) {
		lyrics, err = s.Html()
		if err != nil {
			panic(err)
		}
	})

	return Lyrics{
		Artist: sr.Artist,
		Title:  sr.Title,
		Url:    sr.Url,
		Lyrics: lyrics,
	}, nil
}
