package models

type LinkStatus string

const (
	StatusAvailable    LinkStatus = "available"
	StatusNotAvailable LinkStatus = "not available"
)

type LinkCheck struct {
	URL    string `json:"url"`
	Status string `json:"status"`
}

type LinksGroup struct {
	ID    int         `json:"id"`
	Links []LinkCheck `json:"links"`
}
