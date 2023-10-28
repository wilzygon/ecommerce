package user

import (
	"github.com/google/uuid"
	"github.com/wilzygon/ecommerce/model"
)

// UserCase es el puerto por donde se van a comunicar los datos de entrada
type UseCase interface {
	Create(m *model.User) error
	GetByID(ID uuid.UUID) (model.User, error)
	GetByEmail(email string) (model.User, error)
	GetAll() (model.Users, error)
}

// Storage lo que necesitamos para los datos de salida, donde lo vamos a almacenar
type Storage interface {
	Create(m *model.User) error
	GetByID(ID uuid.UUID) (model.User, error)
	GetByEmail(email string) (model.User, error)
	GetAll() (model.Users, error)
}
