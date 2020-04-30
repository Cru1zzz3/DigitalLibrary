package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"net/http"

	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

var server = "DESKTOP-OBVOF9A"
var port = "1433"
var user = "user"
var password = "user"
var database = "DigitalLibrary"

// Conn uses for shared connection between functions/queries
var Conn *sql.DB

// ConnectToDB connects to MSSQL database and returns connection to database
func ConnectToDB() (*sql.DB, error) {
	//connString := fmt.Sprintf("server=%s,port=%s,user id=%s, Iassword=%s,database=%s", server, port, user, IDPassword, database)
	//conn, err := sql.Open("sqlserver", connString)

	var err error
	Conn, err = sql.Open("sqlserver", "sqlserver://user:user@DESKTOP-OBVOF9A?database=DigitalLibrary")
	if err != nil {
		log.Println("Error creating connection pool: " + err.Error())
		return nil, err
	}

	ctx := context.Background()
	err = Conn.PingContext(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	log.Printf("Connected!")
	return Conn, nil

}

//RegisterUser adds to DB new user
func RegisterUser(Nickname string, Email string, Hash string) error {

	ctx := context.Background()
	var err error

	// Check if database is alive.
	if Conn == nil {
		err = errors.New("Register user: db is null")
		return err
	}

	err = Conn.PingContext(ctx)
	if err != nil {
		return err
	}

	tsql := "INSERT INTO Users (Nickname,Email,Hash) Values (@Nickname,@Email,@Hash)"
	prepared, err := Conn.Prepare(tsql)
	if err != nil {
		log.Println("Error prepare new user row: "+err.Error(), prepared)
		return err
	}
	defer prepared.Close()

	row := prepared.QueryRowContext(
		ctx,
		sql.Named("Nickname", Nickname),
		sql.Named("Email", Email),
		sql.Named("Hash", Hash),
	)

	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Successfully registred user with nickname", Nickname)
		} else {
			panic(err)
		}
	}
	return nil
}

//LoginUser ...
func LoginUser(Nickname string) (bool, string, error) {
	ctx := context.Background()
	var err error

	// Check if database is alive.
	if Conn == nil {
		err = errors.New("Register user: db is null")
		return false, "", err
	}

	err = Conn.PingContext(ctx)
	if err != nil {
		return false, "", err
	}

	var nicknameFromDB, hashFromDB string // FQ = from query

	tsql := fmt.Sprintf("SELECT Nickname, Hash FROM Users where Nickname = '%s';", Nickname)
	err = Conn.QueryRow(tsql, 1).Scan(&nicknameFromDB, &hashFromDB)
	if err != nil {
		log.Println("User is not registered!: ", err)
		return false, "", nil
	}

	return true, hashFromDB, nil
}

// InsertReader uses autoincrement IDReader field
func InsertReader(Login string, Name string, Surname string, Age int) error {

	ctx := context.Background()
	var err error

	if Conn == nil {
		err = errors.New("Insert reader: db is null")
		return err
	}

	// Check if database is alive.
	err = Conn.PingContext(ctx)
	if err != nil {
		return err
	}

	//
	tsql := "INSERT INTO Readers (Login,Name,Surname,Age) Values (@Login,@Name,@Surname,@Age)"
	prepared, err := Conn.Prepare(tsql)
	if err != nil {
		log.Println("Error insert row: "+err.Error(), prepared)
		return err
	}
	defer prepared.Close()

	row := prepared.QueryRowContext(
		ctx,
		sql.Named("Login", Login),
		sql.Named("Name", Name),
		sql.Named("Surname", Surname),
		sql.Named("Age", Age))

	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Successfully inserted new Reader!")
		} else {
			panic(err)
		}
	}
	return nil
}

// SelectReader executes select query on Readers table
func SelectReader(w http.ResponseWriter) error {
	tsql := fmt.Sprintf("SELECT * FROM Readers;")
	rows, err := Conn.Query(tsql)
	if err != nil {
		log.Fatal("Error select row: " + err.Error())
	}
	defer rows.Close()
	var count int

	for rows.Next() {
		var IDReader, Age int
		var Login, Name, Surname string
		// Get values from row.
		err := rows.Scan(&IDReader, &Login, &Name, &Surname, &Age)
		if err != nil {
			log.Fatal("Error scan row: " + err.Error())
		}

		fmt.Fprintf(w, "ID: %d, Login: %s, Name: %s,Surname %s,Age: %d,\n", IDReader, Login, Name, Surname, Age)
		count++
	}

	return nil
}

// SelectBook TODO: Finish search via searchbar
func SelectBook(w http.ResponseWriter, search string) (SearchStruct, error) {
	tsql := fmt.Sprintf(`Select NameBook,NameAuthor 
	FROM Books,BookAuthor,Authors 
	WHERE Books.IDBook = BookAuthor.IDBook
	AND BookAuthor.IDBook = Authors.IDAuthor
	AND (NameBook LIKE '%%%s%%' OR Authors.NameAuthor LIKE  '%%%s%%');`, search, search)
	rows, err := Conn.Query(tsql)
	if err != nil {
		log.Fatal("Error select row: " + err.Error())
	}
	defer rows.Close()

	var count int

	searchResults := SearchStruct{}

	for rows.Next() {
		var NameBookScan, NameAuthorScan string
		// Get values from row.
		err := rows.Scan(&NameBookScan, &NameAuthorScan)
		if err != nil {
			log.Fatal("Error scan row: " + err.Error())
		}

		instanceBook := Book{
			NameBook:   NameBookScan,
			NameAuthor: NameAuthorScan,
		}

		searchResults.Books = append(searchResults.Books, instanceBook)

		//fmt.Fprintf(w, "NameBook: %s, NameAuthor: %s\n", NameBook, NameAuthor)
		count++
	}

	return searchResults, nil
}

// db, err := db.Query("select * from dbo.Formulars")
// if err != nil {
// 	log.Fatal("Prepare failed: ", err.Error())
// }
// defer stmt.Close()

// row := stmt.QueryRow()
// var IDBook int
// err = row.Scan(&IDBook)
// if err != nil {
// 	log.Fatal("Scan failed:", err.Error())
// }
// fmt.Printf("IDBook:%d\n", IDBook)
