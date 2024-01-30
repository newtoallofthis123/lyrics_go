package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/newtoallofthis123/lyrics_go/cli"
	"github.com/newtoallofthis123/lyrics_go/db"
	"github.com/newtoallofthis123/lyrics_go/scraper"
	"github.com/newtoallofthis123/lyrics_go/utils"
	"github.com/ttacon/chalk"
)

func main() {
	args := os.Args[1:]

	var query string

	if len(args) == 0 {
		query = cli.GetQuery()
	} else {
		query = strings.Join(args, " ")
	}

	dbConn, err := db.GetDbConnection()
	if err != nil {
		panic(err)
	}

	entries, err := dbConn.GetAllEntries()
	if err != nil {
		panic(err)
	}

	s := cli.GetSpinner()
	var results []db.Entry

	existing, ok := utils.FindExisting(entries, query)
	if ok {
		results = append(results, existing...)

		fmt.Println(results)
	} else {
		fmt.Println(chalk.Yellow, "Cache Miss, Scraping...", chalk.Reset)

		s.Message("Acquiring Instance...")
		s.Start()
		instance, err := utils.GetInstance()
		if err != nil {
			panic(err)
		}
		s.StopMessage(fmt.Sprintf("Acquired Instance: %s", instance))
		s.Stop()

		s.Message("Searching...")
		request := scraper.GetLyrics{
			Instance: instance,
			Query:    query,
		}

		s.Start()
		searchResults, err := request.ParseSearchPage()
		if err != nil {
			s.StopFail()
			panic(err)
		}
		if len(searchResults) == 0 {
			s.StopFail()
			panic("No Results Found")
		}

		s.StopMessage("Search Complete")
		s.Stop()

		var options []string
		for _, v := range searchResults {
			options = append(options, fmt.Sprintf("%s - %s", v.Title, v.Artist))
		}

		userChoice := cli.GetOptions(options)

		var choice scraper.SearchResult
		for _, v := range searchResults {
			if userChoice == fmt.Sprintf("%s - %s", v.Title, v.Artist) {
				choice = v
			}
		}

		s.Message("Getting Lyrics...")
		s.Start()
		lyrics, err := choice.GetSongLyrics()
		if err != nil {
			s.StopFail()
			panic(err)
		}
		s.StopMessage("Lyrics Acquired")
		s.Stop()

		if lyrics.Lyrics == "" {
			s.StopFail()
			panic("No Lyrics Found, Try Again")
		}

		s.Message("Inserting Into Cache...")
		s.Start()
		err = dbConn.InsertEntry(LyricsToEntry(lyrics))
		if err != nil {
			s.StopFail()
			panic(err)
		}
		s.StopMessage("Inserted Into Cache")
		s.Stop()
	}
}
