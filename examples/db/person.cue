package db

// Title is a record of a person's title.
#Title: string

// Identity is a record of a person's name.
#Identity: {
	First:  string
	Middle: string
	Last:   string
	Nick?:  string
}

// Person is a record of a person.
#Person: {
	Title: #Title
	Name:  #Identity
	Age:   int
}

// People is a list of person records.
#People: [...#Person]
