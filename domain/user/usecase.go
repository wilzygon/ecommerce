package user

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wilzygon/ecommerce/model"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	storage Storage //puerto de salida
}

func New(s Storage) User {
	return User{storage: s}
}

func (u User) Create(m *model.User) error {
	ID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("%s %w", "uuid.NewUUID()", err)
	}

	m.ID = ID
	password, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s %w", "bcrypt.GenerateFromPassword()", err)
	}

	m.Password = string(password)
	if m.Details == nil {
		m.Details = []byte("{}")
	}

	m.CreatedAt = time.Now().Unix()

	//Vamos a guardar ese usuario en el Storage
	//Este dominio User tiene un Storage
	err = u.storage.Create(m) //Le enviamos el modelo que acabamos de utilizar
	if err != nil {
		return fmt.Errorf("%s %w", "storage.Create()", err)
	}

	m.Password = "" //Si no hubo error ya está guardado, y procedemos a limpiar
	return nil
}

func (u User) GetByEmail(email string) (model.User, error) {
	user, err := u.storage.GetByEmail(email)
	if err != nil {
		return model.User{}, fmt.Errorf("%s %w", "storage.GetByEmail()", err)
	}

	return user, nil
}

func (u User) GetAll() (model.Users, error) {
	users, err := u.storage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("%s %w", "storage.GetAll()", err)
	}

	return users, nil
}

func (u User) Login(email, password string) (model.User, error) {
	m, err := u.GetByEmail(email)
	if err != nil {
		return model.User{}, fmt.Errorf("%s %w", "user.GetByEmail()", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password))
	if err != nil {
		return model.User{}, fmt.Errorf("%s %w", "bcrypt.CompareHashAndPassword()", err)
	}

	m.Password = ""

	return m, nil
}
