package health

import "context"

var (
	BuildTime string
	Commit    string
	Version   string
)

type Service interface {
	GetLiveness(ctx context.Context) error
	GetReadiness(ctx context.Context) error
	GetVersion(ctx context.Context) (buildTime, commit, version string)
}

func NewService() Service {
	return &service{}
}

type service struct{}

func (s *service) GetLiveness(_ context.Context) error {
	return nil
}

func (s *service) GetReadiness(_ context.Context) error {
	return nil
}

func (s *service) GetVersion(_ context.Context) (buildTime, commit, version string) {
	return BuildTime, Commit, Version
}
