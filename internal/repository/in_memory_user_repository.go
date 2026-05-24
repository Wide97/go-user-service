package repository

import (
	"sync"
	"time"

	"github.com/mwidesott/go-user-service/internal/errors"
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

func (r *InMemoryUserRepository) Create(u model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.counterId++
	u.ID = r.counterId
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	r.users[u.ID] = u

	return nil

}

func (r *InMemoryUserRepository) GetByID(id int) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, errors.ErrUserNotFound
	}

	return &u, nil

}
