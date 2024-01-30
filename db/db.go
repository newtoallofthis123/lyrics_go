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

type Query struct {
	Id     int
	Query  string
	Artist string
	Title  string
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

func (sqlite3 *DbConnection) InsertQuery(query string, q Query) error {
	sql, args, err := sq.Insert("queries").Columns("query", "artist", "title", "url").Values(query, q.Artist, q.Title, q.Url).ToSql()
	if err != nil {
		return err
	}

	_, err = sqlite3.db.Exec(sql, args...)
	return err
}

func (sqlite3 *DbConnection) GetByQuery(query string) ([]Query, error) {
	var queries []Query

	sql, args, err := sq.Select("id", "artist", "title", "url").From("queries").Where(sq.Like{"query": query}).ToSql()
	if err != nil {
		return queries, err
	}

	rows, err := sqlite3.db.Query(sql, args...)
	if err != nil {
		return queries, err
	}

	for rows.Next() {
		var q Query
		err = rows.Scan(&q.Id, &q.Artist, &q.Title, &q.Url)
		if err != nil {
			return queries, err
		}

		queries = append(queries, q)
	}

	return queries, nil
}

func (sqlite3 *DbConnection) GetAllQueries() ([]Query, error) {
	var queries []Query

	sql, args, err := sq.Select("id", "artist", "title", "url").From("queries").ToSql()
	if err != nil {
		return queries, err
	}

	rows, err := sqlite3.db.Query(sql, args...)
	if err != nil {
		return queries, err
	}

	for rows.Next() {
		var q Query
		err = rows.Scan(&q.Id, &q.Artist, &q.Title, &q.Url)
		if err != nil {
			return queries, err
		}

		queries = append(queries, q)
	}

	return queries, nil
}

func (sqlite3 *DbConnection) GetAllOnlyQueries() ([]string, error) {
	var queries []string

	sql, args, err := sq.Select("query").From("queries").ToSql()
	if err != nil {
		return queries, err
	}

	rows, err := sqlite3.db.Query(sql, args...)
	if err != nil {
		return queries, err
	}

	for rows.Next() {
		var q string
		err = rows.Scan(&q)
		if err != nil {
			return queries, err
		}

		queries = append(queries, q)
	}

	return queries, nil
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

func (sqlite3 *DbConnection) GetAllEntries() ([]Entry, error) {
	var entries []Entry

	sql, args, err := sq.Select("id", "artist", "title", "lyrics", "url").From("lyrics").ToSql()
	if err != nil {
		return entries, err
	}

	rows, err := sqlite3.db.Query(sql, args...)
	if err != nil {
		return entries, err
	}

	for rows.Next() {
		var e Entry
		err = rows.Scan(&e.Id, &e.Artist, &e.Title, &e.Lyrics, &e.Url)
		if err != nil {
			return entries, err
		}

		entries = append(entries, e)
	}

	return entries, nil
}
