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
