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

func (r *InMemoryUserRepository) GetAll() ([]model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]model.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}

	return users, nil

}

func (r *InMemoryUserRepository) Update(user model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.users[user.ID]
	if !ok {
		return errors.ErrUserNotFound
	}

	user.UpdatedAt = time.Now()

	r.users[user.ID] = user

	return nil

}

func (r *InMemoryUserRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.users[id]
	if !ok {
		return errors.ErrUserNotFound
	}

	delete(r.users, id)

	return nil
}
