package repository

import (
	"context"
	"database/sql"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type ContactsRepository struct {
	db *sql.DB
}

type Contact struct {
	ID           int64     `json:"-"`
	UserId       int64     `json:"-"`
	Uuid         string    `json:"uuid"`
	EmailAddress string    `json:"email_address"`
	FirstName    string    `json:"firstname"`
	LastName     string    `json:"lastname"`
	TimelineId   int64     `json:"timeline_id"`
	HistoryId    int64     `json:"history_id"`
	LastStmt     int       `json:"last_stmt"`
	CreatedAt    time.Time `json:"created_at"`
}

func (c *Contact) Scan() []interface{} {
	s := reflect.ValueOf(c).Elem()
	numCols := s.NumField()
	columns := make([]interface{}, numCols)
	for i := 0; i < numCols; i++ {
		field := s.Field(i)
		columns[i] = field.Addr().Interface()
	}
	return columns
}

func (r *ContactsRepository) Create(user *User, contact *Contact) (*Contact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO contact (user_id, uuid, email_address, firstname, lastname)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;`

	contact.Uuid = uuid.NewString()

	args := []interface{}{user.ID, contact.Uuid, contact.EmailAddress, contact.FirstName, contact.LastName}

	err = tx.QueryRowContext(ctx, query, args...).Scan(&contact.ID)
	if err != nil {
		return &Contact{}, err
	}

	query = `
		SELECT *
		FROM contact
		WHERE id = $1 AND user_id = $2
		ORDER BY created_at DESC;`

	args = []interface{}{contact.ID, user.ID}

	err = tx.QueryRowContext(ctx, query, args...).Scan(contact.Scan()...)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return contact, nil
}

func (r *ContactsRepository) GetAll(user *User) ([]*Contact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT *
		FROM contact
		WHERE user_id = $1
		ORDER BY created_at DESC;`

	args := []interface{}{user.ID}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	contacts := []*Contact{}

	for rows.Next() {
		var contact Contact

		err := rows.Scan(contact.Scan()...)

		if err != nil {
			return nil, err
		}

		contacts = append(contacts, &contact)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}
