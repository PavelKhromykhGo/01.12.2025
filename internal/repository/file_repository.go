package repository

import (
	"LinkChecker/internal/models"
	"context"
	"encoding/json"
	"errors"
	"os"
)

type Repository interface {
	GetID(ctx context.Context) (int, error)
	SaveGroup(ctx context.Context, group models.LinksGroup) error
	GetGroups(ctx context.Context, id []int) ([]models.LinksGroup, error)
}

type FileData struct {
	ID     int                 `json:"id"`
	Groups []models.LinksGroup `json:"groups"`
}

type FileRepo struct {
	path string
	data FileData
}

func NewFileRepo(path string) (*FileRepo, error) {
	r := &FileRepo{
		path: path,
		data: FileData{
			ID:     1,
			Groups: []models.LinksGroup{},
		},
	}
	if err := r.load(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return r, nil
		}
		return nil, err
	}
	return r, nil
}

func (r *FileRepo) GetID(ctx context.Context) (int, error) {
	id := r.data.ID
	r.data.ID++
	return id, r.save()
}

func (r *FileRepo) SaveGroup(ctx context.Context, group models.LinksGroup) error {
	r.data.Groups = append(r.data.Groups, group)
	return r.save()
}

func (r *FileRepo) GetGroups(ctx context.Context, ids []int) ([]models.LinksGroup, error) {
	var res []models.LinksGroup

	for _, id := range ids {
		for _, group := range r.data.Groups {
			if group.ID == id {
				res = append(res, group)
				break
			}
		}
	}
	return res, nil
}

func (r *FileRepo) save() error {
	data, err := json.MarshalIndent(r.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, data, 0644)
}

func (r *FileRepo) load() error {
	data, err := os.ReadFile(r.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &r.data)
}
