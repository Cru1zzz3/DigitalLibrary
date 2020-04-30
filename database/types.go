package database

//SearchStruct which contain collection of Books
type SearchStruct struct {
	Books []Book
}

//Book contain NameBook and NameAuthor
type Book struct {
	NameBook   string
	NameAuthor string
}
