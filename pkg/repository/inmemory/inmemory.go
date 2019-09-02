package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/tennuem/task/pkg/task"
)

var (
	ErrTaskExist    = errors.New("task exist")
	ErrTaskNotExist = errors.New("task does not exist")
)

func NewRepository() task.Repository {
	return &inmemRepo{
		store: make(map[uuid.UUID]*task.Task, 0),
	}
}

type inmemRepo struct {
	mtx   sync.Mutex
	store map[uuid.UUID]*task.Task
}

func (r *inmemRepo) Create(ctx context.Context, t *task.Task) (*task.Task, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	t.ID = uuid.New()
	r.store[t.ID] = t
	return r.store[t.ID], nil
}

func (r *inmemRepo) GetList(ctx context.Context) []*task.Task {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	var result []*task.Task
	for _, v := range r.store {
		result = append(result, v)
	}
	return result
}

func (r *inmemRepo) GetByID(ctx context.Context, ID uuid.UUID) (*task.Task, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	v, ok := r.store[ID]
	if !ok {
		return nil, ErrTaskNotExist
	}
	return v, nil
}

func (r *inmemRepo) Update(ctx context.Context, t *task.Task) (*task.Task, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if _, ok := r.store[t.ID]; !ok {
		return nil, ErrTaskNotExist
	}
	r.store[t.ID] = t
	return r.store[t.ID], nil
}

func (r *inmemRepo) Delete(ctx context.Context, ID uuid.UUID) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if _, ok := r.store[ID]; !ok {
		return ErrTaskNotExist
	}
	delete(r.store, ID)
	return nil
}
