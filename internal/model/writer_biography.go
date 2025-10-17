package model

type WriterBiography struct {
	ID          int64
	Name        string
	PortraitURL string
	Lifespan    string
	Country     string
	Occuptation string
	IsFavorite  bool
	Content     string
}
