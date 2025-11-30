package handlers

import (
	"LinkChecker/internal/models"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeService struct {
	group models.LinksGroup
}

func (s *fakeService) CheckLinks(ctx context.Context, urls []string) (models.LinksGroup, error) {
	return s.group, nil
}

func (s *fakeService) GetGroups(ctx context.Context, ids []int) ([]models.LinksGroup, error) {
	return nil, nil
}

type fakePDF struct{}

func (f *fakePDF) Generate(group []models.LinksGroup) ([]byte, error) {
	return []byte("One, two"), nil
}

func TestCheckLinksHandler(t *testing.T) {
	wantGroup := models.LinksGroup{
		ID: 1,
		Links: []models.LinkCheck{
			{URL: "ya.ru", Status: "available"},
			{URL: "fake", Status: "not available"},
		},
	}

	fSvc := &fakeService{wantGroup}
	h := NewHandler(fSvc, &fakePDF{})

	reqBody := `{"links":["ya.ru","fake"]}`
	req := httptest.NewRequest("POST", "/links/check", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	h.CheckLinks(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("got status %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var parsedJSON struct {
		Links    map[string]string `json:"links"`
		LinksNum int               `json:"links_num"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&parsedJSON); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if parsedJSON.LinksNum != wantGroup.ID {
		t.Fatalf("got links_num %d, want %d", parsedJSON.LinksNum, wantGroup.ID)
	}

	wantLinks := map[string]string{
		"ya.ru": "available",
		"fake":  "not available",
	}

	for url, want := range wantLinks {
		status, ok := parsedJSON.Links[url]
		if !ok {
			t.Errorf("link %s not found", url)
			continue
		}
		if status != want {
			t.Errorf("link %s: got status %s, want %s", url, status, want)
		}
	}
}
