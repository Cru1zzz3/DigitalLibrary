package database

//SearchStruct which contain collection of Books
type SearchStruct struct {
	Books   []Book
	Authors []Author
	Genres  []Genre
}

//Book contain NameBook and NameAuthor
type Book struct {
	NameBook     string
	DescribeBook string
	NameAuthor   string
	NameGenre    string
}

type Author struct {
	NameAuthor     string
	DescribeAuthor string
	Books          []Book
	//IDUser         string
}

type Genre struct {
	NameGenre     string
	DescribeGenre string
	Books         []Book
}
