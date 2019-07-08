package task

import "context"

type Task struct {
	ID          int
	Title       string
	Description string
}

type Service interface {
	Create(context.Context, *Task) error
	GetList(context.Context) []*Task
	GetByID(ctx context.Context, ID int) (*Task, error)
	Update(ctx context.Context, t *Task) error
	Delete(ctx context.Context, ID int) error
}

type Repository interface {
	Create(context.Context, *Task) error
	GetList(context.Context) []*Task
	GetByID(ctx context.Context, ID int) (*Task, error)
	Update(ctx context.Context, t *Task) error
	Delete(ctx context.Context, ID int) error
}

func NewService(repo Repository) Service {
	return &service{repo}
}

type service struct {
	repo Repository
}

func (s *service) Create(ctx context.Context, t *Task) error {
	return s.repo.Create(ctx, t)
}

func (s *service) GetList(ctx context.Context) []*Task {
	return s.repo.GetList(ctx)
}

func (s *service) GetByID(ctx context.Context, ID int) (*Task, error) {
	return s.repo.GetByID(ctx, ID)
}

func (s *service) Update(ctx context.Context, t *Task) error {
	return s.repo.Update(ctx, t)
}

func (s *service) Delete(ctx context.Context, ID int) error {
	return s.repo.Delete(ctx, ID)
}
