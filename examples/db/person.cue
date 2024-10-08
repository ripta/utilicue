package db

#Title: string

#Identity: {
	First:  string
	Middle: string
	Last:   string
	Nick?:  string
}

#Person: {
	Title: #Title
	Name:  #Identity
	Age:   int
}

#People: [...#Person]
