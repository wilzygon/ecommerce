package model

import (
	"github.com/google/uuid"
)

// Proveedor model of table proveedores
type Proveedor struct {
	ID           uuid.UUID `json:"id"`
	CodProveedor int64     `json:"cod_proveedor"`
	Nombre       string    `json:"nombre"`
	RucCi        string    `json:"ruc_ci"`
	Telefono     string    `json:"telefono"`
	Direccion    string    `json:"direccion"`
	Email        string    `json:"email"`
	CreatedAt    int64     `json:"created_at"`
	UpdatedAt    int64     `json:"updated_at"`
}

func (p Proveedor) HasID() bool {
	return p.ID != uuid.Nil
}

// Proveedores slice of Proveedor
type Proveedores []Proveedor

func (p Proveedores) IsEmpty() bool { return len(p) == 0 }
