package repository

import "github.com/mwidesott/go-user-service/internal/model"

type UserRepository interface {
	Create(user model.User) error
	GetByID(id int) (*model.User, error)
	GetAll() ([]model.User, error)
	Update(user model.User) error
	Delete(id int) error
}
