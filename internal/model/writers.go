package model

import "strings"

type WriterCard struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	PortraitURL string   `json:"portrait_url"`
	IsFavorite  bool     `json:"is_favorite"`
	Tags        []string `json:"tags"`
}

type WriterBiography struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	PortraitURL string   `json:"portrait_url"`
	Tags        []string `json:"tags"`
	Lifespan    string   `json:"lifespan"`
	Country     string   `json:"country"`
	Occupation  string   `json:"occupation"`
	IsFavorite  bool     `json:"is_favorite"`
	Content     string   `json:"content"`
}

func ParseTags(tagsString string) []string {
	if tagsString == "" {
		return []string{}
	}

	var tags []string
	rawTags := strings.Split(tagsString, ",")
	for _, tag := range rawTags {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			tags = append(tags, trimmed)
		}
	}
	return tags
}

func TagsToString(tags []string) string {
	return strings.Join(tags, ", ")
}
