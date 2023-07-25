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

type Timestamp int64

func (p *Timestamp) Scan(value interface{}) error {
	t := value.(time.Time).UnixMilli()
	*p = Timestamp(t)
	return nil
}

type Contact struct {
	Id           int64      `json:"-"`
	UserId       int64      `json:"-"`
	Uuid         string     `json:"uuid"`
	EmailAddress *string    `json:"email_address"`
	FirstName    *string    `json:"firstname"`
	LastName     *string    `json:"lastname"`
	CreatedAt    Timestamp  `json:"created_at"`
	ModifiedAt   *Timestamp `json:"modified_at"`
	TimelineId   int64      `json:"-"`
	HistoryId    int64      `json:"-"`
	LastStmt     int        `json:"-"`
}

type contactHistory struct {
	History          int64 `json:"last_history_id"`
	ContactsInserted []*Contact
	ContactsUpdated  []*Contact
	ContactsTrashed  []*Contact
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

	query := `
		INSERT INTO contact (user_id, uuid, email_address, firstname, lastname)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING * ;`

	contact.Uuid = uuid.NewString()

	args := []interface{}{user.ID, contact.Uuid, contact.EmailAddress, contact.FirstName, contact.LastName}

	err := r.db.QueryRowContext(ctx, query, args...).Scan(contact.Scan()...)
	if err != nil {
		return nil, err
	}

	// contact.CreatedAtUnix = contact.CreatedAt.UnixMilli()
	// contact.ModifiedAtUnix = contact.ModifiedAt.UnixMilli()

	return contact, nil
}

func (r *ContactsRepository) GetAll(user *User) ([]*Contact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT *
		FROM contact
		WHERE user_id = $1 AND
		last_stmt < 2
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

func (r *ContactsRepository) GetHistory(user *User, history *History) (*contactHistory, error) {
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
		FROM contact
		WHERE user_id = $1 AND
		      last_stmt = 0 AND
			  history_id > $2
		ORDER BY created_at DESC;`

	args := []interface{}{user.ID, history.LastHistoryId}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	contactHistory := &contactHistory{
		ContactsInserted: []*Contact{},
		ContactsUpdated:  []*Contact{},
		ContactsTrashed:  []*Contact{},
	}

	for rows.Next() {
		var contact Contact

		err := rows.Scan(contact.Scan()...)

		if err != nil {
			return nil, err
		}

		contactHistory.ContactsInserted = append(contactHistory.ContactsInserted, &contact)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// updated rows
	query = `
		SELECT *
		FROM contact
		WHERE user_id = $1 AND
		      last_stmt = 1 AND
			  history_id > $2
		ORDER BY created_at DESC;`

	args = []interface{}{user.ID, history.LastHistoryId}

	rows, err = tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var contact Contact

		err := rows.Scan(contact.Scan()...)

		if err != nil {
			return nil, err
		}

		contactHistory.ContactsUpdated = append(contactHistory.ContactsUpdated, &contact)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// trashed rows
	query = `
		SELECT *
		FROM contact
		WHERE user_id = $1 AND
		      last_stmt = 2 AND
			  history_id > $2
		ORDER BY created_at DESC;`

	args = []interface{}{user.ID, history.LastHistoryId}

	rows, err = tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var contact Contact

		err := rows.Scan(contact.Scan()...)

		if err != nil {
			return nil, err
		}

		contactHistory.ContactsTrashed = append(contactHistory.ContactsTrashed, &contact)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// history
	query = `
	SELECT last_history_id FROM contact_history_seq;`

	err = tx.QueryRowContext(ctx, query).Scan(&contactHistory.History)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return contactHistory, nil
}
