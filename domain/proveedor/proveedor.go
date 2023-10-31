package proveedor

import (
	"github.com/google/uuid"
	"github.com/wilzygon/ecommerce/model"
)

type UseCase interface {
	Create(m *model.Proveedor) error
	Update(m *model.Proveedor) error
	Delete(ID uuid.UUID) error

	GetByID(ID uuid.UUID) (model.Proveedor, error)
	GetAll() (model.Proveedores, error)
}

type Storage interface {
	Create(m *model.Proveedor) error
	Update(m *model.Proveedor) error
	Delete(ID uuid.UUID) error

	GetByID(ID uuid.UUID) (model.Proveedor, error)
	GetAll() (model.Proveedores, error)
}
