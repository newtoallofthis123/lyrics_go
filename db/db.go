package db

import (
	sq "github.com/Masterminds/squirrel"
)

type Entry struct {
	Id     int
	Artist string
	Title  string
	Lyrics string
	Url    string
}

func (sqlite3 *DbConnection) InsertEntry(e Entry) error {
	sql, args, err := sq.Insert("lyrics").Columns("artist", "title", "lyrics", "url").Values(e.Artist, e.Title, e.Lyrics, e.Url).ToSql()
	if err != nil {
		return err
	}

	_, err = sqlite3.db.Exec(sql, args...)
	return err
}

func (sqlite3 *DbConnection) GetEntryByUrl(url string) (Entry, error) {
	var e Entry

	sql, args, err := sq.Select("id", "artist", "title", "lyrics", "url").From("lyrics").Where(sq.Eq{"url": url}).ToSql()
	if err != nil {
		return e, err
	}

	err = sqlite3.db.QueryRow(sql, args...).Scan(&e.Id, &e.Artist, &e.Title, &e.Lyrics, &e.Url)
	if err != nil {
		return e, err
	}

	return e, nil
}
