package model

type ArticleCard struct {
	ID          int      `json:"id"`
	CoverURL    string   `json:"cover_url"`
	Title       string   `json:"title"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	IsFavorite  bool     `json:"is_favorite"`
}

type Article struct {
	ID          int      `json:"id"`
	CoverURL    string   `json:"cover_url"`
	Title       string   `json:"title"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	IsFavorite  bool     `json:"is_favorite"`
}
