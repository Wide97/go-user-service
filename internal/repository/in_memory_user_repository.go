package repository

import (
	"sync"

	"github.com/mwidesott/go-user-service/internal/model"
)

type InMemoryUserRepository struct {
	users     map[int]model.User
	mu        sync.RWMutex
	counterId int
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	r := &InMemoryUserRepository{
		users: make(map[int]model.User),
	}

	return r
}
