package models

type Post struct {
	PostId          uint64   `json:"post_id"`
	Title           string   `json:"title"`
	Content         string   `json:"content"`
	Author          string   `json:"author"`
	PublicationDate string   `json:"publication_date"`
	Tags            []string `json:"tags"`
}
