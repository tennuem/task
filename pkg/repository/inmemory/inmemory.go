package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/tennuem/task/pkg/task"
)

var (
	ErrTaskExist    = errors.New("task exist")
	ErrTaskNotExist = errors.New("task does not exist")
)

func NewRepository() task.Repository {
	return &inmemRepo{
		store: make(map[int]*task.Task, 0),
	}
}

type inmemRepo struct {
	mtx   sync.Mutex
	store map[int]*task.Task
}

func (r *inmemRepo) Create(ctx context.Context, t *task.Task) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.store[t.ID] = t
	return nil
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

func (r *inmemRepo) GetByID(ctx context.Context, ID int) (*task.Task, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	v, ok := r.store[ID]
	if !ok {
		return nil, ErrTaskNotExist
	}
	return v, nil
}

func (r *inmemRepo) Update(ctx context.Context, t *task.Task) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if _, ok := r.store[t.ID]; !ok {
		return ErrTaskNotExist
	}
	r.store[t.ID] = t
	return nil
}

func (r *inmemRepo) Delete(ctx context.Context, ID int) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	delete(r.store, ID)
	return nil
}
