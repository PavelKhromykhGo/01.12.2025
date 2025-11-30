package service

import (
	"LinkChecker/internal/models"
	"context"
)

type Repo interface {
	GetID(ctx context.Context) (int, error)
	SaveGroup(ctx context.Context, group models.LinksGroup) error
	GetGroups(ctx context.Context, ids []int) ([]models.LinksGroup, error)
}

type Service interface {
	CheckLinks(ctx context.Context, urls []string) (models.LinksGroup, error)
	GetGroups(ctx context.Context, ids []int) ([]models.LinksGroup, error)
}

type Checker interface {
	Check(ctx context.Context, url string) models.LinkStatus
}

type LinkService struct {
	repo    Repo
	checker Checker
}

func NewLinkService(repo Repo, checker Checker) *LinkService {
	return &LinkService{repo: repo, checker: checker}
}

func (s *LinkService) CheckLinks(ctx context.Context, urls []string) (models.LinksGroup, error) {
	id, err := s.repo.GetID(ctx)
	if err != nil {
		return models.LinksGroup{}, err
	}

	links := make([]models.LinkCheck, 0, len(urls))

	for _, url := range urls {
		status := s.checker.Check(ctx, url)

		links = append(links, models.LinkCheck{URL: url, Status: string(status)})
	}

	group := models.LinksGroup{ID: id, Links: links}

	if err := s.repo.SaveGroup(ctx, group); err != nil {
		return models.LinksGroup{}, err
	}

	return group, nil
}

func (s *LinkService) GetGroups(ctx context.Context, ids []int) ([]models.LinksGroup, error) {
	return s.repo.GetGroups(ctx, ids)
}
