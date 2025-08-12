package main

import (
	"database/sql"
	"time"
)

type Database interface {
	Init() error
	Write(time time.Time, float65 float64) error
	Reset() error
}

type SQLite struct {
	db *sql.DB
}

func (s SQLite) Init() error {
	return s.CreateTable()
}

func (s SQLite) CreateTable() error {
	const query = `
    CREATE TABLE IF NOT EXISTS power(
        timestamp TIMESTAMP PRIMARY KEY,
        load REAL
    );
    `
	_, err := s.db.Exec(query)
	return err
}

func (s SQLite) Write(time time.Time, load float64) error {
	const query = `INSERT INTO power VALUES (?, ?)`

	_, err := s.db.Exec(query, time, load)
	return err
}

func (s SQLite) Reset() error {
	const query = `DROP TABLE IF EXISTS power;`

	if _, err := s.db.Exec(query); err == nil {
		return s.CreateTable()
	} else {
		return err
	}

}
