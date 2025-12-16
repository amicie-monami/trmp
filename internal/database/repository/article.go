package repository

import (
	"database/sql"
	"log"
	"trmp/internal/model"
)

type ArticleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// GetAll возвращает список всех статей (карточки)
func (r *ArticleRepository) GetAll() ([]model.ArticleCard, error) {
	query := `
		SELECT id, cover_url, title, tags, description 
		FROM articles 
		ORDER BY id DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error querying articles: %v", err)
		return nil, err
	}
	defer rows.Close()

	var articles []model.ArticleCard
	for rows.Next() {
		var article model.ArticleCard
		var tagsString string

		err := rows.Scan(
			&article.ID,
			&article.CoverURL,
			&article.Title,
			&tagsString,
			&article.Description,
		)
		if err != nil {
			log.Printf("Error scanning article row: %v", err)
			continue
		}

		// Преобразуем строку тегов в массив
		article.Tags = model.ParseTags(tagsString)
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	log.Printf("Found %d articles in database", len(articles))
	return articles, nil
}

// GetByID возвращает полную статью по ID
func (r *ArticleRepository) GetByID(id int) (*model.Article, error) {
	query := `
		SELECT id, cover_url, title, tags, description, content
		FROM articles 
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)

	article := &model.Article{}
	var tagsString string

	err := row.Scan(
		&article.ID,
		&article.CoverURL,
		&article.Title,
		&tagsString,
		&article.Description,
		&article.Content,
	)

	if err == sql.ErrNoRows {
		log.Printf("Article not found with ID: %d", id)
		return nil, nil
	}

	if err != nil {
		log.Printf("Error scanning article row: %v", err)
		return nil, err
	}

	// Преобразуем строку тегов в массив
	article.Tags = model.ParseTags(tagsString)

	return article, nil
}
