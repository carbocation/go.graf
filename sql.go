/*
Demonstrates usage of pq

Expects:
Postgres
Host: Localhost
Port: 5432
DB: postgres
User: (your current logged-in user)
Password: (none)
Sslmode: disabled

Schema exists: golang
Table exists: yourtable 

It then pulls down some records from golang.yourtable and prints them out.1
*/
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	// Open a connection
	db, err := sql.Open("postgres", "dbname=postgres sslmode=disable")
	if err != nil {
		fmt.Printf("Connection Error: %s", err)
	}
	defer db.Close()

	// Prepare a statement
	stmt, err := db.Prepare("SELECT * FROM golang.yourtable WHERE rank > 1")
	if err != nil {
		fmt.Printf("Statement Preparation Error: %s", err)
	}

	// Query from that prepared statement
	rows, err := stmt.Query()
	if err != nil {
		fmt.Printf("Query Error: %v", err)
	}

	// Show the *Rows ptr
	fmt.Printf("Row pointer: %#v \n", rows)
	
	cols, err := rows.Columns()
	if err != nil {
		fmt.Printf("Column error: %s", err)
	}
	
	fmt.Printf("Columns: %s \n", cols)
	
	// Iterate over the rows
	for rows.Next() {
		var rank int
		var username, password string
		err = rows.Scan(&rank, &username, &password)
		fmt.Printf("Record: %#i, %s, %s \n", rank, username, password)
	}

}
