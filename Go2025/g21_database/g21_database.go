// g21_database.go
// Learning go, Example of Matering Go, §5 Packahes and functions, example of packahe in GitHub
//
// 2025-06-18	PV		First version

package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("g21_database")

	arguments := os.Args
	if len(arguments) != 6 {
		fmt.Println("Please provide: hostname port username password db")
		return
	}

	host := arguments[1]
	p := arguments[2]
	user := arguments[3]
	pass := arguments[4]
	database := arguments[5]

	// Port number SHOULD BE an integer
	port, err := strconv.Atoi(p)
	if err != nil {
		fmt.Println("Not a valid port number:", err)
		return
	}

	// connection string
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, database)

	// open PostgreSQL database
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("Open():", err)
		return
	}
	defer db.Close()

	fmt.Println("Database opened successfully")

	// Get all databases
	rows, err := db.Query(`SELECT "datname" FROM "pg_database" WHERE datistemplate = false`)
	if err != nil {
		fmt.Println("Query", err)
		return
	}

	// In order to execute a SELECT query, you need to create it first. As the presented
	// SELECT query contains no parameters, which means that it does not change based on
	// variables, you can pass it to the Query() function and execute it. The live outcome of
	// the SELECT query is kept in the rows variable, which is a cursor. You do not get all the
	// results from the database, as a query might return millions of records, but you get
	// them one by one—this is the point of using a cursor.
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println("Scan", err)
			return
		}
		fmt.Println("*", name)
	}
	defer rows.Close()

	// The previous code shows how to process the results of a SELECT query, which can
	// be from nothing to lots of rows. As the rows variable is a cursor, you advance from
	// row to row by calling Next(). After that, you need to assign the values returned from
	// the SELECT query into Go variables, in order to use them. This happens with a call
	// to Scan(), which requires pointer parameters. If the SELECT query returns multiple
	// values, you need to put multiple parameters in Scan(). Lastly, you must call Close()
	// with defer for the rows variable in order to close the statement and free various
	// types of used resources.

	// Get all tables from __current__ database
	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name`
	rows, err = db.Query(query)

	if err != nil {
		fmt.Println("Query", err)
		return
	}
	
	// We are going to execute another SELECT query in the current database, as provided
	// by the user. The definition of the SELECT query is kept in the query variable for
	// simplicity and for creating easy to read code. The contents of the query variable are
	// passed to the db.Query() method.
	// This is how you process the rows that are returned from SELECT
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println("Scan", err)
			return
		}
		fmt.Println("+T", name)
	}
	defer rows.Close()
}
