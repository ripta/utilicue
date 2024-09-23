package db

#Identity: {
	First:  string
	Middle: string
	Last:   string
	Nick:   string
}

#Person: {
	Name: #Identity
	Age:  int
}

#People: [...#Person]
