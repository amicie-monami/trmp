package model

type Article struct {
	ID          int
	Title       string
	Description string
	CoverURL    string
	IsFavorite  bool
	Content     string
}
