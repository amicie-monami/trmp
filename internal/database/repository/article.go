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

// GetAllWithFavorites возвращает статьи с информацией об избранном
func (r *ArticleRepository) GetAllWithFavorites(userID int) ([]model.ArticleCard, error) {
	query := `
		SELECT a.id, a.cover_url, a.title, a.tags, a.description,
		       CASE WHEN fa.user_id IS NOT NULL THEN 1 ELSE 0 END as is_favorite
		FROM articles a
		LEFT JOIN favorite_articles fa ON a.id = fa.article_id AND fa.user_id = ?
		ORDER BY a.id DESC
	`

	rows, err := r.db.Query(query, userID)
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
			continue
		}

		article.Tags = model.ParseTags(tagsString)
		article.IsFavorite = (favorite == 1)
		articles = append(articles, article)
	}

	return articles, nil
}

func (r *ArticleRepository) GetByIDWithFavorite(id, userID int) (*model.Article, error) {
	query := `
		SELECT a.id, a.cover_url, a.title, a.tags, a.description, a.content,
		       CASE WHEN fa.user_id IS NOT NULL THEN 1 ELSE 0 END as is_favorite
		FROM articles a
		LEFT JOIN favorite_articles fa ON a.id = fa.article_id AND fa.user_id = ?
		WHERE a.id = ?
	`

	row := r.db.QueryRow(query, userID, id)

	article := &model.Article{}
	var tagsString string
	var favorite int

	err := row.Scan(
		&article.ID,
		&article.CoverURL,
		&article.Title,
		&tagsString,
		&article.Description,
		&article.Content,
		&favorite, // Сканируем is_favorite
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	article.Tags = model.ParseTags(tagsString)
	article.IsFavorite = (favorite == 1) // Устанавливаем поле
	return article, nil
}
