package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type ResourcesRepository struct {
	db *sql.DB
}

type Resource struct {
	ID          int64     `json:"id"`
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
}

func (r ResourcesRepository) Create(user *User, uuid, name, path, contentType string, size int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `INSERT INTO
		resources (user_id, uuid, name, path, content_type, size)
		VALUES(?, ?, ?, ?, ?, ?)`, user.ID, uuid, name, path, contentType, size)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r ResourcesRepository) GetAll(user *User, filters Filters) ([]*Resource, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, uuid, name, path, size, content_type, created_at
		FROM resources
		WHERE user_id = $1
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{user.ID, filters.limit(), filters.offset()}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	resources := []*Resource{}

	for rows.Next() {
		var resource Resource

		err := rows.Scan(
			&totalRecords,
			&resource.ID,
			&resource.UUID,
			&resource.Name,
			&resource.Path,
			&resource.Size,
			&resource.ContentType,
			&resource.CreatedAt,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		resources = append(resources, &resource)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return resources, metadata, nil
}
