package user

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wilzygon/ecommerce/infrastructure/postgres"
	"github.com/wilzygon/ecommerce/model"
)

const table = "users"

var fields = []string{
	"id",
	"email",
	"password",
	"is_admin",
	"details",
	"created_at",
	"updated_at",
}
var (
	//psqlInsert = `INSERT INTO users(id, email, password, details, created_at) VALUES($1, $2, $3, $4, $5)`
	psqlInsert = postgres.BuildSQLInsert(table, fields)

	//psqlGetAll = `SELECT id, email, password, details, created_at FROM users`
	psqlGetAll = postgres.BuildSQLSelect(table, fields)
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
		postgres.Int64ToNull(m.UpdatedAt),
	)
	if err != nil {
		return err
	}

	return nil
}

func (u User) GetByEmail(email string) (model.User, error) {
	query := psqlGetAll + " WHERE email = $1"
	row := u.db.QueryRow(
		context.Background(),
		query,
		email,
	)

	return u.scanRow(row, true)
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
	for rows.Next() { //Iteración sobre todos los registros que estamos obteniendo de la BD
		m, err := u.scanRow(rows, false) //Indico que no quiero el password
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	return ms, nil

}

// scanRow toma cada uno de los registros de la base de datos y lo convierte en una estructura usuario
func (u User) scanRow(s pgx.Row, WithPassword bool) (model.User, error) {
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

	if !WithPassword {
		m.Password = ""
	}

	return m, nil
}
