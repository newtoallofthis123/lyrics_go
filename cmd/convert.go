package main

import (
	"github.com/newtoallofthis123/lyrics_go/db"
	"github.com/newtoallofthis123/lyrics_go/scraper"
)

func EntryToLyrics(e db.Entry) scraper.Lyrics {
	return scraper.Lyrics{
		Artist: e.Artist,
		Title:  e.Title,
		Url:    e.Url,
		Lyrics: e.Lyrics,
	}
}

func LyricsToEntry(l scraper.Lyrics) db.Entry {
	return db.Entry{
		Artist: l.Artist,
		Title:  l.Title,
		Url:    l.Url,
		Lyrics: l.Lyrics,
	}
}

func ResultToEntry(r scraper.SearchResult) db.Entry {
	return db.Entry{
		Artist: r.Artist,
		Title:  r.Title,
		Url:    r.Url,
	}
}

func EntryToResult(e db.Entry) scraper.SearchResult {
	return scraper.SearchResult{
		Artist: e.Artist,
		Title:  e.Title,
		Url:    e.Url,
	}
}

func ResultToQuery(r scraper.SearchResult) db.Query {
	return db.Query{
		Query:  r.Title,
		Artist: r.Artist,
		Title:  r.Title,
		Url:    r.Url,
	}
}

func QueryToResult(q db.Query) scraper.SearchResult {
	return scraper.SearchResult{
		Artist: q.Artist,
		Title:  q.Title,
		Url:    q.Url,
	}
}
