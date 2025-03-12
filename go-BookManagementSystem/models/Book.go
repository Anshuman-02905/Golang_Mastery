package models

type Book struct {
	UID    int    `json:"id"`
	NAME   string `json:"name"`
	AUTHOR string `json:"author"`
	GENRE  string `json:"genre"`
}
