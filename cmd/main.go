package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hbollon/go-edlib"
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

	queries, err := dbConn.GetAllQueries()
	if err != nil {
		panic(err)
	}

	if len(queries) > 0 {
		var texts []string
		for _, v := range queries {
			texts = append(texts, v.Title)
		}

		similar, err := edlib.FuzzySearch(query, texts, edlib.Levenshtein)
		if err == nil {
			query = similar
		}
	}

	s := cli.GetSpinner()
	var result db.Entry

	var instance string

	s.Message("Searching...")
	s.Start()

	var searchResults []scraper.SearchResult

	existingQueries, err := dbConn.GetByQuery(query)

	fmt.Println(existingQueries)

	if err == nil && len(existingQueries) > 0 {
		searchResults = make([]scraper.SearchResult, len(existingQueries))
		for i, v := range existingQueries {
			searchResults[i] = EntryToResult(v)
		}
	} else {
		fmt.Println(chalk.Yellow, "Cache Miss, Scraping...", chalk.Reset)

		s.Message("Acquiring Instance...")
		s.Start()
		instance, err = utils.GetInstance()
		if err != nil {
			panic(err)
		}
		s.StopMessage(fmt.Sprintf("Acquired Instance: %s", instance))
		s.Stop()

		request := scraper.GetLyrics{
			Instance: instance,
			Query:    query,
		}
		searchResults, err = request.ParseSearchPage()
		if err != nil {
			s.StopFail()
			panic(err)
		}
		if len(searchResults) == 0 {
			s.StopFail()
			panic("No Results Found")
		}
	}

	s.StopMessage(fmt.Sprintf("Found %d Results", len(searchResults)))
	s.Stop()

	var options []string
	for _, v := range searchResults {
		err := dbConn.InsertQuery(query, ResultToEntry(v))
		if err != nil {
			panic(err)
		}
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

	fmt.Println(result)
}
