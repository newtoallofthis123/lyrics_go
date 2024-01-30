package files

import (
	"os"
	"path"
)

func getHomeDir() string {
	return os.ExpandEnv("$HOME")
}

func GetLyricsDir() string {
	return path.Join(getHomeDir(), ".lyrics")
}

func GetDbPath() string {
	if _, err := os.Stat(GetLyricsDir()); os.IsNotExist(err) {
		os.Mkdir(GetLyricsDir(), 0755)
	}
	return path.Join(GetLyricsDir(), "lyrics.db")
}

func JoinHomeDir(base string) string {
	return path.Join(getHomeDir(), base)
}
