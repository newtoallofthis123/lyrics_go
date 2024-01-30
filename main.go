package main

import (
	"fmt"

	"github.com/newtoallofthis123/lyrics_go/db"
	"github.com/newtoallofthis123/lyrics_go/scraper"
	"github.com/newtoallofthis123/lyrics_go/utils"
)

func main() {
	fmt.Println("Hello World!")

	instance, err := utils.GetInstance()
	if err != nil {
		panic(err)
	}

	query := scraper.GetLyrics{
		Instance: instance,
		Query:    "Emily James Arthur",
	}

	searchResults, err := query.ParseSearchPage()
	if err != nil {
		panic(err)
	}

	fmt.Println(searchResults)

	dbConn, err := db.GetDbConnection()
	if err != nil {
		panic(err)
	}

	toConsider := searchResults[0]

	var result scraper.Lyrics

	entry, err := dbConn.GetEntryByUrl(toConsider.Url)
	if err != nil {
		lyrics, err := searchResults[0].GetSongLyrics()
		if err != nil {
			panic(err)
		}

		err = dbConn.InsertEntry(LyricsToEntry(lyrics))
		if err != nil {
			panic(err)
		}

		result = lyrics
	} else {
		result = EntryToLyrics(entry)
	}

	fmt.Println(result)
}

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
