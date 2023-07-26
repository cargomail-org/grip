package repository

import (
	"context"
	"database/sql"
	"log"
	"reflect"
	"time"
)

type FilesRepository struct {
	db *sql.DB
}

type File struct {
	Id          string     `json:"id"`
	UserId      int64      `json:"-"`
	Checksum    string     `json:"-"`
	Name        string     `json:"name"`
	Path        string     `json:"-"`
	Size        int64      `json:"file_size"`
	ContentType string     `json:"content_type"`
	CreatedAt   Timestamp  `json:"created_at"`
	ModifiedAt  *Timestamp `json:"modified_at"`
	TimelineId  int64      `json:"-"`
	HistoryId   int64      `json:"-"`
	LastStmt    int        `json:"-"`
}

type fileHistory struct {
	History       int64 `json:"last_history_id"`
	FilesInserted []*File
	FilesTrashed  []*File
}

func (f *File) Scan() []interface{} {
	s := reflect.ValueOf(f).Elem()
	numCols := s.NumField()
	columns := make([]interface{}, numCols)
	for i := 0; i < numCols; i++ {
		field := s.Field(i)
		columns[i] = field.Addr().Interface()
	}
	return columns
}

func (r FilesRepository) Create(user *User, file *File) (*File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO
			file (user_id, checksum, name, path, content_type, size)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING * ;`

	args := []interface{}{user.Id, file.Checksum, file.Name, file.Path, file.ContentType, file.Size}

	err := r.db.QueryRowContext(ctx, query, args...).Scan(file.Scan()...)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (r FilesRepository) GetAll(user *User) ([]*File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT *
			FROM file
			WHERE user_id = $1 AND
			last_stmt < 2
			ORDER BY created_at DESC;`

	args := []interface{}{user.Id}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files := []*File{}

	for rows.Next() {
		var file File

		err := rows.Scan(file.Scan()...)
		if err != nil {
			return nil, err
		}

		files = append(files, &file)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func (r *FilesRepository) GetHistory(user *User, history *History) (*fileHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// inserted rows
	query := `
		SELECT *
			FROM file
			WHERE user_id = $1 AND
				last_stmt = 0 AND
				history_id > $2
			ORDER BY created_at DESC;`

	args := []interface{}{user.Id, history.LastHistoryId}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	fileHistory := &fileHistory{
		FilesInserted: []*File{},
		FilesTrashed:  []*File{},
	}

	for rows.Next() {
		var file File

		err := rows.Scan(file.Scan()...)

		if err != nil {
			return nil, err
		}

		fileHistory.FilesInserted = append(fileHistory.FilesInserted, &file)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// trashed rows
	query = `
		SELECT *
			FROM file
			WHERE user_id = $1 AND
				last_stmt = 2 AND
				history_id > $2
			ORDER BY created_at DESC;`

	args = []interface{}{user.Id, history.LastHistoryId}

	rows, err = tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var file File

		err := rows.Scan(file.Scan()...)

		if err != nil {
			return nil, err
		}

		fileHistory.FilesTrashed = append(fileHistory.FilesTrashed, &file)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// history
	query = `
	SELECT last_history_id
	   FROM file_history_seq
	   WHERE user_id = $1 ;`

	args = []interface{}{user.Id}

	err = tx.QueryRowContext(ctx, query, args...).Scan(&fileHistory.History)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return fileHistory, nil
}

func (r *FilesRepository) TrashByIdList(user *User, idList string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if len(idList) > 0 {
		query := `
		UPDATE file
			SET last_stmt = 2
			WHERE user_id = $1 AND
			id IN (SELECT value FROM json_each($2));`

		args := []interface{}{user.Id, idList}

		_, err := r.db.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r FilesRepository) DeleteByIdList(user *User, idList string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if len(idList) > 0 {
		query := `
		DELETE
			FROM file
			WHERE user_id = $1 AND
			id IN (SELECT value FROM json_each($2));`

		args := []interface{}{user.Id, idList}

		_, err := r.db.ExecContext(ctx, query, args...)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (r FilesRepository) GetOriginalFileName(user *User, id string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT *
			FROM file
			WHERE user_id = $1 AND
				id = $2;`

	file := &File{}

	args := []interface{}{user.Id, id}

	err := r.db.QueryRowContext(ctx, query, args...).Scan(file.Scan()...)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return file.Name, nil
}
