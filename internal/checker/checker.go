package checker

import (
	"LinkChecker/internal/models"
	"context"
	"net/http"
	"strings"
	"time"
)

type HTTPChecker struct {
	client *http.Client
}

// NewHTTPChecker создаёт cheker с заданным тайм‑аутом на запросы.
func NewHTTPChecker(timeout time.Duration) *HTTPChecker {
	return &HTTPChecker{client: &http.Client{
		Timeout: timeout,
	}}
}

// Check выполняет GET запрос к указанному URL и возвращает статус его доступности.
func (c *HTTPChecker) Check(ctx context.Context, url string) models.LinkStatus {
	url = normalizeURL(url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return models.StatusNotAvailable
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return models.StatusNotAvailable
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return models.StatusAvailable
	}
	return models.StatusNotAvailable
}

// normalizeURL гарантирует наличие схемы в URL, если клиент его не передал.
func normalizeURL(url string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	return "http://" + url
}
