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

//var pc = "DESKTOP-OBVOF9A"
//var laptop = "DESKTOP-EEEPFL6"
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
	Conn, err = sql.Open("sqlserver", "sqlserver://user:user@DESKTOP-EEEPFL6?database=DigitalLibrary")
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

// SelectBook find name of book or author
func SelectBook(w http.ResponseWriter, search string) (SearchStruct, error) {
	tsql := fmt.Sprintf(`Select NameBook,NameAuthor 
	FROM Books,BookAuthor,Authors, Genres 
	WHERE Books.IDBook = BookAuthor.IDBook
	AND BookAuthor.IDBook = Authors.IDAuthor
	AND Genres.IDGenre = Books.IDGenre
	AND (NameBook
		 LIKE '%%%s%%' OR 
		 Authors.NameAuthor LIKE  '%%%s%%'
		  OR Genres.NameGenre LIKE '%%%s%%' );`, search, search, search)
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

// TODO: Fill dropdown with existing genres
func GetGenres(w http.ResponseWriter) []string {

	ctx := context.Background()
	err := Conn.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Connected!")

	rows, err := Conn.QueryContext(ctx, "GetTopGenres")
	if err != nil {
		log.Printf("Error in get top genres querycontext!")
		log.Fatal(err.Error())
	}
	var TopGenresSlice []string

	TopGenresSlice = nil

	// tsql := fmt.Sprintf("SELECT TOP 5 NameGenre FROM Genres;")
	// rows, err := Conn.Query(tsql)
	// if err != nil {
	// 	log.Fatal("Error select top 5 genres row: " + err.Error())
	// }
	// defer rows.Close()

	for rows.Next() {
		var NameGenre string
		// Get values from row.
		err := rows.Scan(&NameGenre)
		if err != nil {
			log.Fatal("Error scan namegenre row: " + err.Error())
		}

		TopGenresSlice = append(TopGenresSlice, NameGenre)
		//fmt.Fprintf(w, "ID: %d, Login: %s, Name: %s,Surname %s,Age: %d,\n", IDReader, Login, Name, Surname, Age)

	}
	return TopGenresSlice
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

func GetAllGenres(w http.ResponseWriter) SearchStruct {

	ctx := context.Background()
	err := Conn.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Connected!")

	rows, err := Conn.QueryContext(ctx, "GetAllGenres")
	if err != nil {
		log.Printf("Error in get top genres querycontext!")
		log.Fatal(err.Error())
	}

	// tsql := fmt.Sprintf("SELECT TOP 5 NameGenre FROM Genres;")
	// rows, err := Conn.Query(tsql)
	// if err != nil {
	// 	log.Fatal("Error select top 5 genres row: " + err.Error())
	// }
	// defer rows.Close()

	searchResult := SearchStruct{}

	for rows.Next() {
		var NameGenreScan, DescribeGenreScan string
		// Get values from row.
		err := rows.Scan(&NameGenreScan, &DescribeGenreScan)
		if err != nil {
			log.Fatal("Error scan namegenre describegenre row: " + err.Error())
		}

		currentGenre := Genre{
			NameGenre:     NameGenreScan,
			DescribeGenre: DescribeGenreScan,
		}

		searchResult.Genres = append(searchResult.Genres, currentGenre)
		//fmt.Fprintf(w, "ID: %d, Login: %s, Name: %s,Surname %s,Age: %d,\n", IDReader, Login, Name, Surname, Age)

	}
	return searchResult
}

func AboutGenre(w http.ResponseWriter, genre string) (SearchStruct, error) {

	ctx := context.Background()
	err := Conn.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Connected!")

	rows, err := Conn.QueryContext(ctx, "GetGenre", sql.Named("NameGenre", sql.Out{Dest: &genre}))
	if err != nil {
		log.Printf("Error in get genre  querycontext!")
		log.Fatal(err.Error())
	}
	var count int

	searchResults := SearchStruct{}

	for rows.Next() {
		var NameGenreScan, DescribeGenreScan, NameBookScan string
		// Get values from row.
		err := rows.Scan(&NameGenreScan, &DescribeGenreScan, &NameBookScan)
		if err != nil {
			log.Fatal("Error scan row: " + err.Error())
		}

		instanceGenre := Genre{
			NameGenre:     NameGenreScan,
			DescribeGenre: DescribeGenreScan,
		}

		searchResults, err = SelectBook(w, NameGenreScan)

		searchResults.Genres = append(searchResults.Genres, instanceGenre)

		//fmt.Fprintf(w, "NameBook: %s, NameAuthor: %s\n", NameBook, NameAuthor)
		count++
	}

	return searchResults, nil
}

func AboutAuthor(w http.ResponseWriter, author string) (SearchStruct, error) {

	ctx := context.Background()
	err := Conn.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Connected!")

	rows, err := Conn.QueryContext(ctx, "GetAuthor", sql.Named("NameAuthor", sql.Out{Dest: &author}))
	if err != nil {
		log.Printf("Error in get author  querycontext!")
		log.Fatal(err.Error())
	}
	var count int

	searchResults := SearchStruct{}

	for rows.Next() {
		var NameAuthorScan, DescribeAuthorScan string
		// Get values from row.
		err := rows.Scan(&NameAuthorScan, &DescribeAuthorScan)
		if err != nil {
			log.Fatal("Error scan row: " + err.Error())
		}

		instanceAuthor := Author{
			NameAuthor:     NameAuthorScan,
			DescribeAuthor: DescribeAuthorScan,
		}

		searchResults, err = SelectBook(w, NameAuthorScan)
		if err != nil {
			log.Printf("Error in about author  querycontext!")
			log.Fatal(err.Error())
		}

		searchResults.Authors = append(searchResults.Authors, instanceAuthor)

		//fmt.Fprintf(w, "NameBook: %s, NameAuthor: %s\n", NameBook, NameAuthor)
		count++
	}

	return searchResults, nil
}

func AboutBook(w http.ResponseWriter, book string) (SearchStruct, error) {

	ctx := context.Background()
	err := Conn.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Connected!")

	rows, err := Conn.QueryContext(ctx, "GetBook", sql.Named("NameBook", sql.Out{Dest: &book}))
	if err != nil {
		log.Printf("Error in get book  querycontext!")
		log.Fatal(err.Error())
	}
	var count int

	searchResults := SearchStruct{}

	for rows.Next() {
		var NameBookScan, DescribeBookScan, NameAuthorScan, NameGenreScan string
		// Get values from row.
		err := rows.Scan(&NameBookScan, &DescribeBookScan, &NameAuthorScan, &NameGenreScan)
		if err != nil {
			log.Fatal("Error scan row: " + err.Error())
		}

		instanceBook := Book{
			NameBook:     NameBookScan,
			DescribeBook: DescribeBookScan,
			NameAuthor:   NameAuthorScan,
			NameGenre:    NameGenreScan,
		}

		searchResults.Books = append(searchResults.Books, instanceBook)

		//fmt.Fprintf(w, "NameBook: %s, NameAuthor: %s\n", NameBook, NameAuthor)
		count++
	}

	return searchResults, nil
}
