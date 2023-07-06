package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type FileRepository struct {
	db *sql.DB
}

type File struct {
	ID          int64     `json:"id"`
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
}

func (r FileRepository) Create(user *User, uuid, name, path, contentType string, size int64) error {
	_, err := r.db.Exec(`INSERT INTO
		file(user_id, uuid, name, path, content_type, size)
		VALUES(?, ?, ?, ?, ?, ?)`, user.ID, uuid, name, path, contentType, size)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (m FileRepository) GetAll(user *User, filters Filters) ([]*File, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, uuid, name, path, size, content_type, created_at
		FROM file
		WHERE user_id = $1
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{user.ID, filters.limit(), filters.offset()}

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	files := []*File{}

	for rows.Next() {
		var file File

		err := rows.Scan(
			&totalRecords,
			&file.ID,
			&file.UUID,
			&file.Name,
			&file.Path,
			&file.Size,
			&file.ContentType,
			&file.CreatedAt,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		files = append(files, &file)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return files, metadata, nil
}
