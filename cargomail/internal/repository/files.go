package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type FilesRepository struct {
	db *sql.DB
}

type File struct {
	ID          int64     `json:"-"`
	UUID        string    `json:"uuid"`
	Hash        string    `json:"-"`
	Name        string    `json:"name"`
	Path        string    `json:"-"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
}

func (r FilesRepository) Create(user *User, uuid string, checksum []byte, name string, path string, contentType string, size int64) (time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `INSERT INTO
		file (user_id, uuid, hash, name, path, content_type, size)
		VALUES(?, ?, ?, ?, ?, ?, ?) RETURNING created_at;`, user.ID, uuid, checksum, name, path, contentType, size)

	err := row.Scan(&user.CreatedAt)
	if row.Err() == sql.ErrNoRows {
		return time.Time{}, nil
	}
	if err != nil {
		return time.Time{}, err
	}

	return user.CreatedAt, nil
}

func (r FilesRepository) GetAll(user *User, filters Filters) ([]*File, Metadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var query string
	if filters.Page == 0 || filters.PageSize == 0 {
		query = fmt.Sprintf(`
		SELECT count(*) OVER(), id, uuid, hash, name, path, size, content_type, created_at
		FROM file
		WHERE user_id = $1
		ORDER BY %s %s, id ASC`, filters.sortColumn(), filters.sortDirection())
	} else {
		query = fmt.Sprintf(`
		SELECT count(*) OVER(), id, uuid, hash, name, path, size, content_type, created_at
		FROM file
		WHERE user_id = $1
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())
	}

	args := []interface{}{user.ID, filters.limit(), filters.offset()}

	rows, err := r.db.QueryContext(ctx, query, args...)
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
			&file.Hash,
			&file.Name,
			&file.Path,
			&file.Size,
			&file.ContentType,
			&file.CreatedAt,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		file.Hash = fmt.Sprintf("%x", file.Hash)

		files = append(files, &file)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return files, metadata, nil
}

func (r FilesRepository) DeleteByUuidList(user *User, uuidList []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if len(uuidList) > 0 {
		uuids := fmt.Sprintf("%v", uuidList)
		uuids = uuids[1 : len(uuids)-1]
		uuids = strings.ReplaceAll(uuids, " ", `","`)

		_, err := r.db.ExecContext(ctx, `DELETE FROM
			file
			WHERE user_id = $1 AND
				  uuid IN ("`+uuids+`");`, user.ID)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (r FilesRepository) GetOriginalFileName(user *User, uuid string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT name, size, content_type, created_at
		FROM file
		WHERE user_id = $1 AND
		      uuid = $2`

	args := []interface{}{user.ID, uuid}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return "", err
	}

	defer rows.Close()

	files := []*File{}

	for rows.Next() {
		var file File

		err := rows.Scan(
			&file.Name,
			&file.Size,
			&file.ContentType,
			&file.CreatedAt,
		)

		if err != nil {
			return "", err
		}

		file.Hash = fmt.Sprintf("%x", file.Hash)

		files = append(files, &file)
	}

	if err = rows.Err(); err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "not found", ErrFileNameNotFound
	}

	return files[0].Name, nil
}
