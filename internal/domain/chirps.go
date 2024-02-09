package domain

type Chirp struct {
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
	ID       int    `json:"id"`
}
