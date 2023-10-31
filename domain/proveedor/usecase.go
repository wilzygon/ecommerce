package proveedor

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wilzygon/ecommerce/model"
)

// Proveedor implements UseCase
type Proveedor struct {
	storage Storage
}

// New returns a new Proveedor
func New(s Storage) Proveedor {
	return Proveedor{storage: s}
}

// Create creates a model.Product
func (p Proveedor) Create(m *model.Proveedor) error {
	ID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("%s %w", "uuid.NewUUID()", err)
	}

	m.ID = ID
	/* if len(m.Images) == 0 {
		m.Images = []byte(`[]`)
	}
	if len(m.Features) == 0 {
		m.Features = []byte(`{}`)
	} */

	err = p.storage.Create(m)
	if err != nil {
		return err
	}

	return nil
}

// Update updates a model.Proveedor by id
func (p Proveedor) Update(m *model.Proveedor) error {
	if !m.HasID() {
		return fmt.Errorf("proveedor: %w", model.ErrInvalidID)
	}

	/* if len(m.Images) == 0 {
		m.Images = []byte(`[]`)
	}
	if len(m.Features) == 0 {
		m.Features = []byte(`{}`)
	} */
	m.UpdatedAt = time.Now().Unix()

	err := p.storage.Update(m)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a model.Proveedor by id
func (p Proveedor) Delete(ID uuid.UUID) error {
	err := p.storage.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}

func (p Proveedor) GetByID(ID uuid.UUID) (model.Proveedor, error) {
	proveedor, err := p.storage.GetByID(ID)
	if err != nil {
		return model.Proveedor{}, fmt.Errorf("proveedor: %w", err)
	}

	return proveedor, nil
}

// GetAll returns a model.Proveedores according to filters and sorts
func (p Proveedor) GetAll() (model.Proveedores, error) {
	proveedores, err := p.storage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("proveedor: %w", err)
	}

	return proveedores, nil
}
