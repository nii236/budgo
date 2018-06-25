package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type conn struct {
	*sqlx.DB
}

func connect(dbuser, dbpass, dbname, dbport string) (*conn, error) {
	fmt.Println("dbuser", dbuser)
	fmt.Println("dbport", dbport)
	fmt.Println("dbname", dbname)
	fmt.Println("dbpass", dbpass)
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s port=%s dbname=%s password=%s sslmode=disable", dbuser, dbport, dbname, dbpass))
	if err != nil {
		log.Fatalln(err)
	}

	return &conn{db}, nil
}

// Record is a single transaction
type Record struct {
	Description string
	Category    string
	Cents       int
}

func (c *conn) saveRecord(name, category string, cents int) error {
	_, err := c.Exec("INSERT INTO records (description, category, cents) VALUES ($1, $2, $3)", name, category, cents)
	if err != nil {
		return err
	}
	return nil
}

func (c *conn) records() ([]*Record, error) {
	result := []*Record{}
	err := c.Select(&result, "SELECT description, category, cents FROM records ORDER BY created_at")
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *conn) total() (int, error) {
	result := 0
	err := c.Select(&result, "SELECT SUM(cents) FROM records ORDER BY created_at")
	if err != nil {
		return 0, err
	}
	return result, nil
}
