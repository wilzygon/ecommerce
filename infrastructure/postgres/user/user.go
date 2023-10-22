package user

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wilzygon/ecommerce/model"
)

var (
	psqlInsert = `INSERT INTO users(id, email, password, details, created_at) 
	VALUES($1, $2, $3, $4, $5)`

	psqlGetAll = `SELECT id, email, password, details, created_at FROM users`
)

type User struct {
	db *pgxpool.Pool
}

// New retorna un nuevo User storage
func New(db *pgxpool.Pool) User {
	return User{db: db}
}

func (u User) Create(m *model.User) error {
	_, err := u.db.Exec(
		context.Background(),
		psqlInsert,
		m.ID,
		m.Email,
		m.Password,
		m.IsAdmin,
		m.Details,
		m.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u User) GetByEmail(email string) (model.User, error) {
	query := psqlGetAll + "WHERE email = $1"
	row := u.db.QueryRow(
		context.Background(),
		query,
		email,
	)

	return u.scanRow(row)
}

func (u User) GetAll() (model.Users, error) {
	rows, err := u.db.Query(
		context.Background(),
		psqlGetAll,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ms := model.Users{}
	for rows.Next() {
		m, err := u.scanRow(rows)
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	return ms, nil

}

func (u User) scanRow(s pgx.Row) (model.User, error) {
	m := model.User{}
	updatedAtNull := sql.NullInt64{}

	err := s.Scan(
		&m.ID,
		&m.Email,
		&m.Password,
		&m.IsAdmin,
		&m.Details,
		&m.CreatedAt,
		&updatedAtNull,
	)
	if err != nil {
		return m, err
	}

	m.UpdatedAt = updatedAtNull.Int64

	return m, nil
}