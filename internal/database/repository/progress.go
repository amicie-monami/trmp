package repository

import (
	"database/sql"
	"strconv"
	"trmp/internal/model"
)

type ProgressRepository struct {
	db *sql.DB
}

func NewProgressRepository(db *sql.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

func (r *ProgressRepository) GetUserProgress(userID int) (*model.ProgressResponse, error) {
	query := `
		SELECT item_type, item_id, progress 
		FROM reading_progress 
		WHERE user_id = ?
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	response := &model.ProgressResponse{
		Writers:  make(map[string]float64),
		Articles: make(map[string]float64),
	}

	for rows.Next() {
		var itemType string
		var itemID int
		var progress float64

		err := rows.Scan(&itemType, &itemID, &progress)
		if err != nil {
			continue
		}

		itemIDStr := strconv.Itoa(itemID)

		if itemType == "writer" {
			response.Writers[itemIDStr] = progress
		} else if itemType == "article" {
			response.Articles[itemIDStr] = progress
		}
	}

	return response, nil
}

func (r *ProgressRepository) UpdateProgress(userID int, itemType string, itemID int, progress float64) error {
	query := `
		INSERT INTO reading_progress (user_id, item_type, item_id, progress)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(user_id, item_type, item_id) 
		DO UPDATE SET progress = excluded.progress, updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.Exec(query, userID, itemType, itemID, progress)
	return err
}

func (r *ProgressRepository) BulkUpdateProgress(userID int, writers, articles map[string]float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for idStr, progress := range writers {
		itemID, err := strconv.Atoi(idStr)
		if err != nil || itemID <= 0 {
			continue
		}

		_, err = tx.Exec(`
			INSERT INTO reading_progress (user_id, item_type, item_id, progress)
			VALUES (?, 'writer', ?, ?)
			ON CONFLICT(user_id, item_type, item_id) 
			DO UPDATE SET progress = excluded.progress, updated_at = CURRENT_TIMESTAMP
		`, userID, itemID, progress)

		if err != nil {
			return err
		}
	}

	for idStr, progress := range articles {
		itemID, err := strconv.Atoi(idStr)
		if err != nil || itemID <= 0 {
			continue
		}

		_, err = tx.Exec(`
			INSERT INTO reading_progress (user_id, item_type, item_id, progress)
			VALUES (?, 'article', ?, ?)
			ON CONFLICT(user_id, item_type, item_id) 
			DO UPDATE SET progress = excluded.progress, updated_at = CURRENT_TIMESTAMP
		`, userID, itemID, progress)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ProgressRepository) GetItemProgress(userID int, itemType string, itemID int) (float64, error) {
	query := `
		SELECT progress FROM reading_progress 
		WHERE user_id = ? AND item_type = ? AND item_id = ?
	`

	var progress float64
	err := r.db.QueryRow(query, userID, itemType, itemID).Scan(&progress)
	if err == sql.ErrNoRows {
		return 0.0, nil
	}

	return progress, err
}
