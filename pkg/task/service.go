package task

import (
	"context"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

type Service interface {
	Create(context.Context, *Task) (*Task, error)
	GetList(context.Context) []*Task
	GetByID(ctx context.Context, ID uuid.UUID) (*Task, error)
	Update(ctx context.Context, t *Task) (*Task, error)
	Delete(ctx context.Context, ID uuid.UUID) error
}

type Repository interface {
	Create(context.Context, *Task) (*Task, error)
	GetList(context.Context) []*Task
	GetByID(ctx context.Context, ID uuid.UUID) (*Task, error)
	Update(ctx context.Context, t *Task) (*Task, error)
	Delete(ctx context.Context, ID uuid.UUID) error
}

func NewService(repo Repository) Service {
	return &service{repo}
}

type service struct {
	repo Repository
}

func (s *service) Create(ctx context.Context, t *Task) (*Task, error) {
	return s.repo.Create(ctx, t)
}

func (s *service) GetList(ctx context.Context) []*Task {
	return s.repo.GetList(ctx)
}

func (s *service) GetByID(ctx context.Context, ID uuid.UUID) (*Task, error) {
	return s.repo.GetByID(ctx, ID)
}

func (s *service) Update(ctx context.Context, t *Task) (*Task, error) {
	return s.repo.Update(ctx, t)
}

func (s *service) Delete(ctx context.Context, ID uuid.UUID) error {
	return s.repo.Delete(ctx, ID)
}
