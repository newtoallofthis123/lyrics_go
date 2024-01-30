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
	"github.com/theckman/yacspin"
	"github.com/ttacon/chalk"
)

func main() {
	fmt.Println(chalk.Blue, "Lyrics Go", chalk.Reset)

	db, err := db.GetDbConnection()
	if err != nil {
		panic(err)
	}

	spinner := cli.GetSpinner()

	cmd := os.Args[1]
	options := os.Args[2:]
	switch cmd {
	case "search":
		search(spinner, db, options)
	default:
		fmt.Println("Unknown command")
	}
}

func search(spinner *yacspin.Spinner, db *db.DbConnection, options []string) {
	var query string

	if len(options) < 1 {
		query = cli.GetQuery()
	} else {
		query = strings.Join(options, " ")
	}

	queries, _ := db.GetAllOnlyQueries()

	spinner.Message("Accessing Instance")
	spinner.Start()
	instance, err := utils.GetInstance()
	if err != nil {
		spinner.StopFailMessage("Failed to access instance")
		spinner.StopFail()
		panic(err)
	}
	spinner.Stop()

	similar, err := edlib.FuzzySearchThreshold(query, queries, 0.2, edlib.Levenshtein)
	if err == nil && len(similar) > 0 {
		query = similar
	} else {
		req := scraper.GetLyrics{
			Instance: instance,
			Query:    query,
		}

		spinner.Message("Fetching Search Results")
		spinner.Start()

		searchResults, err := req.ParseSearchPage()
		if err != nil {
			spinner.StopFailMessage("Failed to fetch search results")
			spinner.StopFail()
			panic(err)
		}
		spinner.Stop()
		if len(searchResults) == 0 {
			fmt.Println(chalk.Red, "No results found", chalk.Reset)
			return
		}

		var options []string

		for i, v := range searchResults {
			options = append(options, fmt.Sprintf("%d. %s By %s", i+1, v.Title, v.Artist))

			db.InsertQuery(query, ResultToQuery(v))
		}

		// choice := cli.GetOptions(options)
	}
}
