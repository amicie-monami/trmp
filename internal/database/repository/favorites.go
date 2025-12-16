package repository

import (
	"database/sql"
)

type FavoritesRepository struct {
	db *sql.DB
}

func NewFavoritesRepository(db *sql.DB) *FavoritesRepository {
	return &FavoritesRepository{db: db}
}

//
// writers
//

func (r *FavoritesRepository) IsWriterFavorite(userID, writerID int) (bool, error) {
	query := `
		SELECT COUNT(*) FROM favorite_writers 
		WHERE user_id = ? AND writer_id = ?
	`

	var count int
	err := r.db.QueryRow(query, userID, writerID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *FavoritesRepository) AddWriterToFavorite(userID, writerID int) error {
	query := `
		INSERT OR IGNORE INTO favorite_writers (user_id, writer_id) 
		VALUES (?, ?)
	`

	_, err := r.db.Exec(query, userID, writerID)
	return err
}

func (r *FavoritesRepository) RemoveWriterFromFavorite(userID, writerID int) error {
	query := `
		DELETE FROM favorite_writers 
		WHERE user_id = ? AND writer_id = ?
	`

	_, err := r.db.Exec(query, userID, writerID)
	return err
}

func (r *FavoritesRepository) GetFavoriteWriters(userID int) ([]int, error) {
	query := `
		SELECT writer_id FROM favorite_writers 
		WHERE user_id = ? 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var writerIDs []int
	for rows.Next() {
		var writerID int
		err := rows.Scan(&writerID)
		if err != nil {
			continue
		}
		writerIDs = append(writerIDs, writerID)
	}

	return writerIDs, nil
}

func (r *FavoritesRepository) ToggleWriterFavorite(userID, writerID int) (bool, error) {
	isFavorite, err := r.IsWriterFavorite(userID, writerID)
	if err != nil {
		return false, err
	}

	if isFavorite {
		err = r.RemoveWriterFromFavorite(userID, writerID)
		return false, err
	} else {
		err = r.AddWriterToFavorite(userID, writerID)
		return true, err
	}
}

//
// articles
//

func (r *FavoritesRepository) GetFavoriteArticles(userID int) ([]int, error) {
	query := `
		SELECT article_id FROM favorite_articles 
		WHERE user_id = ? 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articleIDs []int
	for rows.Next() {
		var articleID int
		err := rows.Scan(&articleID)
		if err != nil {
			continue
		}
		articleIDs = append(articleIDs, articleID)
	}

	return articleIDs, nil
}

// IsArticleFavorite проверяет, добавлена ли статья в избранное
func (r *FavoritesRepository) IsArticleFavorite(userID, articleID int) (bool, error) {
	query := `SELECT COUNT(*) FROM favorite_articles WHERE user_id = ? AND article_id = ?`
	var count int
	err := r.db.QueryRow(query, userID, articleID).Scan(&count)
	return count > 0, err
}

// AddArticleToFavorite добавляет статью в избранное
func (r *FavoritesRepository) AddArticleToFavorite(userID, articleID int) error {
	query := `INSERT OR IGNORE INTO favorite_articles (user_id, article_id) VALUES (?, ?)`
	_, err := r.db.Exec(query, userID, articleID)
	return err
}

// RemoveArticleFromFavorite удаляет статью из избранного
func (r *FavoritesRepository) RemoveArticleFromFavorite(userID, articleID int) error {
	query := `DELETE FROM favorite_articles WHERE user_id = ? AND article_id = ?`
	_, err := r.db.Exec(query, userID, articleID)
	return err
}

// ToggleArticleFavorite переключает статус избранного для статьи
func (r *FavoritesRepository) ToggleArticleFavorite(userID, articleID int) (bool, error) {
	isFavorite, err := r.IsArticleFavorite(userID, articleID)
	if err != nil {
		return false, err
	}

	if isFavorite {
		err = r.RemoveArticleFromFavorite(userID, articleID)
		return false, err
	} else {
		err = r.AddArticleToFavorite(userID, articleID)
		return true, err
	}
}
