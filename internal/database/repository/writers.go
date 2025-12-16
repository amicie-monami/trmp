package repository

import (
	"database/sql"
	"log"
	"sort"
	"strings"
	"trmp/internal/model"
)

type WriterRepository struct {
	db *sql.DB
}

func NewWriterRepository(db *sql.DB) *WriterRepository {
	return &WriterRepository{db: db}
}

// GetAll возвращает список всех писателей (карточки)
func (r *WriterRepository) GetAll() ([]model.WriterCard, error) {
	query := `
		SELECT id, name, portrait_url, is_favorite, tags 
		FROM writers 
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error querying writers: %v", err)
		return nil, err
	}
	defer rows.Close()

	var writers []model.WriterCard
	for rows.Next() {
		var writer model.WriterCard
		var tagsString string

		err := rows.Scan(
			&writer.ID,
			&writer.Name,
			&writer.PortraitURL,
			&writer.IsFavorite,
			&tagsString,
		)
		if err != nil {
			log.Printf("Error scanning writer row: %v", err)
			continue
		}

		// Преобразуем строку тегов в массив
		writer.Tags = model.ParseTags(tagsString)
		writers = append(writers, writer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	log.Printf("Found %d writers", len(writers))
	return writers, nil
}

// GetByID возвращает полную биографию писателя по ID
func (r *WriterRepository) GetByID(id int) (*model.WriterBiography, error) {
	query := `
		SELECT id, name, portrait_url, tags, lifespan, country, occupation, is_favorite, content
		FROM writers 
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)

	writer := &model.WriterBiography{}
	var tagsString string

	err := row.Scan(
		&writer.ID,
		&writer.Name,
		&writer.PortraitURL,
		&tagsString,
		&writer.Lifespan,
		&writer.Country,
		&writer.Occupation,
		&writer.IsFavorite,
		&writer.Content,
	)

	if err == sql.ErrNoRows {
		log.Printf("Writer not found with ID: %d", id)
		return nil, nil
	}

	if err != nil {
		log.Printf("Error scanning writer row: %v", err)
		return nil, err
	}

	// Преобразуем строку тегов в массив
	writer.Tags = model.ParseTags(tagsString)

	return writer, nil
}

// Search возвращает писателей по поисковому запросу (по имени ИЛИ тегам)
func (r *WriterRepository) Search(queryStr string) ([]model.WriterCard, error) {
	if queryStr == "" {
		return r.GetAll()
	}

	query := `
		SELECT id, name, portrait_url, is_favorite, tags 
		FROM writers 
		WHERE name LIKE ? OR tags LIKE ?
		ORDER BY name
	`

	searchPattern := "%" + queryStr + "%"
	rows, err := r.db.Query(query, searchPattern, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var writers []model.WriterCard
	for rows.Next() {
		var writer model.WriterCard
		var tagsString string

		err := rows.Scan(
			&writer.ID,
			&writer.Name,
			&writer.PortraitURL,
			&writer.IsFavorite,
			&tagsString,
		)
		if err != nil {
			continue
		}

		writer.Tags = model.ParseTags(tagsString)
		writers = append(writers, writer)
	}

	return writers, nil
}

// SearchByName поиск писателей только по имени
func (r *WriterRepository) SearchByName(name string) ([]model.WriterCard, error) {
	query := `
		SELECT id, name, portrait_url, is_favorite, tags 
		FROM writers 
		WHERE name LIKE ? 
		ORDER BY name
	`

	rows, err := r.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var writers []model.WriterCard
	for rows.Next() {
		var writer model.WriterCard
		var tagsString string

		err := rows.Scan(
			&writer.ID,
			&writer.Name,
			&writer.PortraitURL,
			&writer.IsFavorite,
			&tagsString,
		)
		if err != nil {
			continue
		}

		writer.Tags = model.ParseTags(tagsString)
		writers = append(writers, writer)
	}

	return writers, nil
}

// SearchByTags поиск писателей по тегам
func (r *WriterRepository) SearchByTags(tags []string) ([]model.WriterCard, error) {
	if len(tags) == 0 {
		return r.GetAll()
	}

	// Создаем условие для каждого тега (AND логика)
	var conditions []string
	var args []interface{}

	for _, tag := range tags {
		conditions = append(conditions, "tags LIKE ?")
		args = append(args, "%"+strings.TrimSpace(tag)+"%")
	}

	query := `
		SELECT id, name, portrait_url, is_favorite, tags 
		FROM writers 
		WHERE ` + strings.Join(conditions, " AND ") + `
		ORDER BY name
	`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var writers []model.WriterCard
	for rows.Next() {
		var writer model.WriterCard
		var tagsString string

		err := rows.Scan(
			&writer.ID,
			&writer.Name,
			&writer.PortraitURL,
			&writer.IsFavorite,
			&tagsString,
		)
		if err != nil {
			continue
		}

		writer.Tags = model.ParseTags(tagsString)
		writers = append(writers, writer)
	}

	return writers, nil
}

// GetAllTags возвращает все уникальные теги
func (r *WriterRepository) GetAllTags() ([]string, error) {
	query := `
		SELECT DISTINCT tags FROM writers
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tagSet := make(map[string]bool)

	for rows.Next() {
		var tagsString string
		err := rows.Scan(&tagsString)
		if err != nil {
			continue
		}

		tags := model.ParseTags(tagsString)
		for _, tag := range tags {
			tagSet[tag] = true
		}
	}

	var tags []string
	for tag := range tagSet {
		tags = append(tags, tag)
	}

	// Сортируем
	strings := sort.StringSlice(tags)
	strings.Sort()

	return tags, nil
}

// ToggleFavorite переключает статус избранного
func (r *WriterRepository) ToggleFavorite(id int) error {
	query := `
		UPDATE writers 
		SET is_favorite = NOT is_favorite 
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)
	return err
}

// GetAllWithFavorites возвращает список всех писателей с информацией об избранном для конкретного пользователя
func (r *WriterRepository) GetAllWithFavorites(userID int) ([]model.WriterCard, error) {
	query := `
		SELECT w.id, w.name, w.portrait_url, w.tags,
		       CASE WHEN fw.id IS NOT NULL THEN 1 ELSE 0 END as is_favorite
		FROM writers w
		LEFT JOIN favorite_writers fw ON w.id = fw.writer_id AND fw.user_id = ?
		ORDER BY w.name
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		log.Printf("Error querying writers with favorites: %v", err)
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
			log.Printf("Error scanning writer row: %v", err)
			continue
		}

		writer.Tags = model.ParseTags(tagsString)
		writer.IsFavorite = (favorite == 1)
		writers = append(writers, writer)
	}

	return writers, nil
}

// GetByIDWithFavorite возвращает писателя с информацией об избранном для конкретного пользователя
func (r *WriterRepository) GetByIDWithFavorite(id, userID int) (*model.WriterBiography, error) {
	query := `
		SELECT w.id, w.name, w.portrait_url, w.tags, w.lifespan, 
		       w.country, w.occupation, w.content,
		       CASE WHEN fw.id IS NOT NULL THEN 1 ELSE 0 END as is_favorite
		FROM writers w
		LEFT JOIN favorite_writers fw ON w.id = fw.writer_id AND fw.user_id = ?
		WHERE w.id = ?
	`

	row := r.db.QueryRow(query, userID, id)

	writer := &model.WriterBiography{}
	var tagsString string
	var favorite int

	err := row.Scan(
		&writer.ID,
		&writer.Name,
		&writer.PortraitURL,
		&tagsString,
		&writer.Lifespan,
		&writer.Country,
		&writer.Occupation,
		&writer.Content,
		&favorite,
	)

	if err == sql.ErrNoRows {
		log.Printf("Writer not found with ID: %d", id)
		return nil, nil
	}

	if err != nil {
		log.Printf("Error scanning writer row: %v", err)
		return nil, err
	}

	writer.Tags = model.ParseTags(tagsString)
	writer.IsFavorite = (favorite == 1)

	return writer, nil
}
