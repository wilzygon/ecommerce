package proveedor

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wilzygon/ecommerce/infrastructure/postgres"
	"github.com/wilzygon/ecommerce/model"
)

const table = "proveedores"

var fields = []string{
	"id",
	"cod_proveedor",
	"nombre",
	"ruc_ci",
	"telefono",
	"direccion",
	"email",
	"created_at",
	"updated_at",
}

var (
	psqlInsert = postgres.BuildSQLInsert(table, fields)
	psqlUpdate = postgres.BuildSQLUpdateByID(table, fields)
	psqlDelete = postgres.BuildSQLDelete(table)
	psqlGetAll = postgres.BuildSQLSelect(table, fields)
)

// Proveedor struct that implement the interface domain.proveedor.Storage
type Proveedor struct {
	db *pgxpool.Pool
}

// New returns a new Proveedor storage
func New(db *pgxpool.Pool) Proveedor {
	return Proveedor{db}
}

// Create creates a model.Proveedor
func (p Proveedor) Create(m *model.Proveedor) error {
	_, err := p.db.Exec(
		context.Background(),
		psqlInsert,
		m.ID,
		m.CodProveedor,
		m.Nombre,
		m.RucCi,
		m.Telefono,
		m.Direccion,
		m.Email,
		m.CreatedAt,
		postgres.Int64ToNull(m.UpdatedAt),
	)
	if err != nil {
		return err
	}

	return nil
}

// Update this method updates a model.Proveedor by id
func (p Proveedor) Update(m *model.Proveedor) error {
	_, err := p.db.Exec(
		context.Background(),
		psqlUpdate,
		m.CodProveedor,
		m.Nombre,
		m.RucCi,
		m.Telefono,
		m.Direccion,
		m.Email,
		m.UpdatedAt,
		m.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a model.Proveedor by id
func (p Proveedor) Delete(ID uuid.UUID) error {
	_, err := p.db.Exec(
		context.Background(),
		psqlDelete,
		ID,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetByID gets an ordered model.Proveedor with filters
func (p Proveedor) GetByID(ID uuid.UUID) (model.Proveedor, error) {
	query := psqlGetAll + " WHERE id = $1"
	row := p.db.QueryRow(
		context.Background(),
		query,
		ID,
	)

	return p.scanRow(row)
}

// GetAll gets all model.Proveedores with Fields
func (p Proveedor) GetAll() (model.Proveedores, error) {
	rows, err := p.db.Query(
		context.Background(),
		psqlGetAll,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//ms := model.Products{}
	var ms model.Proveedores
	for rows.Next() {
		m, err := p.scanRow(rows)
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	return ms, nil
}

func (p Proveedor) scanRow(s pgx.Row) (model.Proveedor, error) {
	m := model.Proveedor{}

	updatedAtNull := sql.NullInt64{}

	err := s.Scan(
		&m.ID,
		&m.CodProveedor,
		&m.Nombre,
		&m.RucCi,
		&m.Telefono,
		&m.Direccion,
		&m.Email,
		&m.CreatedAt,
		&updatedAtNull,
	)
	if err != nil {
		return m, err
	}

	m.UpdatedAt = updatedAtNull.Int64

	return m, nil
}
