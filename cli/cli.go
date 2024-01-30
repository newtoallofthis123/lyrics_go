package cli

import (
	"time"

	"github.com/charmbracelet/huh"
	"github.com/theckman/yacspin"
)

func GetSpinner() *yacspin.Spinner {
	cfg := yacspin.Config{
		Frequency:         100 * time.Millisecond,
		CharSet:           yacspin.CharSets[9],
		SuffixAutoColon:   true,
		StopCharacter:     "✓ ",
		StopFailCharacter: "✗ ",
		StopFailMessage:   "Failed!",
		StopColors:        []string{"fgGreen"},
	}

	spinner, err := yacspin.New(cfg)
	if err != nil {
		panic(err)
	}

	return spinner
}

func GetQuery() string {
	var query string

	form := huh.NewInput().Title("Enter The Song Name?").Value(&query).Placeholder("Get Ready to Sing!").WithTheme(huh.ThemeBase16())

	err := form.Run()
	if err != nil {
		panic(err)
	}

	return query
}

func GetOptions(options []string) string {
	var result string

	form := huh.NewSelect[string]().
		Options(huh.NewOptions(options...)...).
		Title("Select a song").Value(&result)

	err := form.Run()
	if err != nil {
		panic(err)
	}

	return result
}
