package service

import (
	"LinkChecker/internal/models"
	"context"
	"reflect"
	"testing"
)

type fakeRepo struct {
	id     int
	groups map[int]models.LinksGroup
}

func newFakeRepo(startID int) *fakeRepo {
	return &fakeRepo{
		id:     startID,
		groups: make(map[int]models.LinksGroup),
	}
}

func (r *fakeRepo) GetID(ctx context.Context) (int, error) {
	id := r.id
	r.id++
	return id, nil
}

func (r *fakeRepo) SaveGroup(ctx context.Context, group models.LinksGroup) error {
	r.groups[group.ID] = group
	return nil
}

func (r *fakeRepo) GetGroups(ctx context.Context, ids []int) ([]models.LinksGroup, error) {
	var res []models.LinksGroup
	for _, id := range ids {
		if group, ok := r.groups[id]; ok {
			res = append(res, group)
		}
	}
	return res, nil
}

type fakeChecker struct {
	statuses map[string]models.LinkStatus
}

func (c *fakeChecker) Check(ctx context.Context, url string) models.LinkStatus {
	if status, ok := c.statuses[url]; ok {
		return status
	}

	return models.StatusNotAvailable
}

func TestLinksServiceCheckLinks(t *testing.T) {
	ctx := context.Background()

	repo := newFakeRepo(1)
	chk := &fakeChecker{
		statuses: map[string]models.LinkStatus{
			"ya.ru": models.StatusAvailable,
			"fake":  models.StatusNotAvailable,
		},
	}

	svc := NewLinkService(repo, chk)

	input := []string{"ya.ru", "fake"}

	group, err := svc.CheckLinks(ctx, input)
	if err != nil {
		t.Fatalf("CheckLinks error: %v", err)
	}
	if group.ID != 1 {
		t.Errorf("expected group id 1, got %d", group.ID)
	}
	if len(group.Links) != 2 {
		t.Fatalf("expected 2 links, got %d", len(group.Links))
	}

	got := make(map[string]models.LinkStatus)
	for _, link := range group.Links {
		got[link.URL] = models.LinkStatus(link.Status)
	}

	want := map[string]models.LinkStatus{
		"ya.ru": models.StatusAvailable,
		"fake":  models.StatusNotAvailable,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}

	saved, ok := repo.groups[group.ID]
	if !ok {
		t.Fatalf("group %v not saved", group.ID)
	}

	if !reflect.DeepEqual(saved, group) {
		t.Errorf("saved group %v, returned group %v", saved, group)
	}
}

func TestLinksServiceGetGroups(t *testing.T) {
	ctx := context.Background()

	repo := newFakeRepo(1)
	chk := &fakeChecker{statuses: map[string]models.LinkStatus{}}
	svc := NewLinkService(repo, chk)

	group1 := models.LinksGroup{ID: 1, Links: []models.LinkCheck{{URL: "ya.ru", Status: "available"}}}
	group2 := models.LinksGroup{ID: 2, Links: []models.LinkCheck{{URL: "fake", Status: "not_available"}}}

	_ = repo.SaveGroup(ctx, group1)
	_ = repo.SaveGroup(ctx, group2)

	got, err := svc.GetGroups(ctx, []int{1, 2})
	if err != nil {
		t.Fatalf("GetGroups error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(got))
	}

	if !reflect.DeepEqual(got[0], group1) || !reflect.DeepEqual(got[1], group2) {
		t.Fatalf("want %v, got %v", []models.LinksGroup{group1, group2}, got)
	}
}
