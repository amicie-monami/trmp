package repository

import (
	"database/sql"
	"log"
	"sort"
	"strings"
	"trmp/internal/model"
)

type SearchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) *SearchRepository {
	return &SearchRepository{db: db}
}

func (r *SearchRepository) SearchAll(query string, tags []string, userID int) ([]model.ArticleCard, []model.WriterCard, error) {
	// Ищем статьи
	articles, err := r.SearchArticles(query, tags, userID)
	if err != nil {
		return nil, nil, err
	}

	// Ищем писателей
	writers, err := r.SearchWriters(query, tags, userID)
	if err != nil {
		// Если ошибка при поиске писателей, возвращаем только статьи
		log.Printf("Error searching writers: %v", err)
		return articles, []model.WriterCard{}, nil
	}

	return articles, writers, nil
}

// SearchArticles поиск статей по тексту и тегам
func (r *SearchRepository) SearchArticles(query string, tags []string, userID int) ([]model.ArticleCard, error) {
	var conditions []string
	var args []interface{}

	// Базовый запрос
	sqlQuery := `
		SELECT a.id, a.cover_url, a.title, a.tags, a.description,
		       CASE WHEN fa.user_id IS NOT NULL THEN 1 ELSE 0 END as is_favorite
		FROM articles a
		LEFT JOIN favorite_articles fa ON a.id = fa.article_id AND fa.user_id = ?
	`

	args = append(args, userID)

	// Поиск по тексту (заголовок, описание, контент, теги)
	if query != "" {
		searchTerm := "%" + strings.ToLower(query) + "%"
		conditions = append(conditions, `
			(LOWER(a.title) LIKE ? OR 
			 LOWER(a.description) LIKE ? OR 
			 LOWER(a.content) LIKE ? OR
			 LOWER(a.tags) LIKE ?)
		`)
		args = append(args, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	// Фильтрация по тегам (AND логика - должны быть ВСЕ указанные теги)
	if len(tags) > 0 {
		for _, tag := range tags {
			conditions = append(conditions, "LOWER(a.tags) LIKE ?")
			args = append(args, "%"+strings.ToLower(tag)+"%")
		}
	}

	// Добавляем WHERE если есть условия
	if len(conditions) > 0 {
		sqlQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	sqlQuery += " ORDER BY a.id DESC"

	log.Printf("Search articles query: %s", sqlQuery)

	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.ArticleCard
	for rows.Next() {
		var article model.ArticleCard
		var tagsString string
		var favorite int

		err := rows.Scan(
			&article.ID,
			&article.CoverURL,
			&article.Title,
			&tagsString,
			&article.Description,
			&favorite,
		)
		if err != nil {
			log.Printf("Error scanning article: %v", err)
			continue
		}

		article.Tags = model.ParseTags(tagsString)
		article.IsFavorite = (favorite == 1)
		articles = append(articles, article)
	}

	return articles, nil
}

// SearchWriters поиск писателей по тексту и тегам
func (r *SearchRepository) SearchWriters(query string, tags []string, userID int) ([]model.WriterCard, error) {
	var conditions []string
	var args []interface{}

	// Базовый запрос
	sqlQuery := `
		SELECT w.id, w.name, w.portrait_url, w.tags,
		       CASE WHEN fw.user_id IS NOT NULL THEN 1 ELSE 0 END as is_favorite
		FROM writers w
		LEFT JOIN favorite_writers fw ON w.id = fw.writer_id AND fw.user_id = ?
	`

	args = append(args, userID)

	// Поиск по тексту (имя, страна, род деятельности)
	if query != "" {
		searchTerm := "%" + strings.ToLower(query) + "%"
		conditions = append(conditions, `
			(LOWER(w.name) LIKE ? OR 
			 LOWER(w.country) LIKE ? OR
			 LOWER(w.occupation) LIKE ?)
		`)
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	// Фильтрация по тегам
	if len(tags) > 0 {
		for _, tag := range tags {
			conditions = append(conditions, "LOWER(w.tags) LIKE ?")
			args = append(args, "%"+strings.ToLower(tag)+"%")
		}
	}

	// Добавляем WHERE если есть условия
	if len(conditions) > 0 {
		sqlQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	sqlQuery += " ORDER BY w.name"

	log.Printf("Search writers query: %s", sqlQuery)

	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var writers []model.WriterCard
	for rows.Next() {
		var writer model.WriterCard
		var tagsString string
		var favorite int

		err := rows.Scan(
			&writer.ID,
			&writer.Name,
			&writer.PortraitURL,
			&tagsString,
			&favorite,
		)
		if err != nil {
			continue
		}

		writer.Tags = model.ParseTags(tagsString)
		writer.IsFavorite = (favorite == 1)
		writers = append(writers, writer)
	}

	return writers, nil
}

// GetAllTags возвращает все уникальные теги (и из статей и из писателей)
func (r *SearchRepository) GetAllTags() ([]string, error) {
	// Получаем теги из статей
	articleTagsQuery := `SELECT DISTINCT tags FROM articles WHERE tags IS NOT NULL AND tags != ''`
	writerTagsQuery := `SELECT DISTINCT tags FROM writers WHERE tags IS NOT NULL AND tags != ''`

	// Объединяем все теги
	tagSet := make(map[string]bool)

	// Теги из статей
	rows, err := r.db.Query(articleTagsQuery)
	if err == nil {
		for rows.Next() {
			var tagsString string
			rows.Scan(&tagsString)
			tags := model.ParseTags(tagsString)
			for _, tag := range tags {
				tagSet[strings.ToLower(tag)] = true
			}
		}
		rows.Close()
	}

	// Теги из писателей
	rows, err = r.db.Query(writerTagsQuery)
	if err == nil {
		for rows.Next() {
			var tagsString string
			rows.Scan(&tagsString)
			tags := model.ParseTags(tagsString)
			for _, tag := range tags {
				tagSet[strings.ToLower(tag)] = true
			}
		}
		rows.Close()
	}

	// Преобразуем в массив и сортируем
	var tags []string
	for tag := range tagSet {
		tags = append(tags, tag)
	}

	// Сортируем по алфавиту
	sort.Strings(tags)

	return tags, nil
}
