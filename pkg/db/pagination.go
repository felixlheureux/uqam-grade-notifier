package db

import (
	"encoding/base64"
	"fmt"
	"strings"
)

var DEFAULT_PAGINATION = Pagination{
	Cursor:         PaginationCursor{},
	Limit:          50,
	OrderDirection: "DESC",
	OrderBy:        "id",
}

type PaginationCursor struct {
	Current string `json:"current"`
	Prev    string `json:"prev"`
	Next    string `json:"next"`
}

type Pagination struct {
	Cursor         PaginationCursor `json:"cursor"`
	Limit          int              `json:"limit"`
	OrderDirection string           `json:"order_direction"`
	OrderBy        string           `json:"order_by"`
}

type Result[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

func (p *Pagination) ApplyDefaults() {
	if p.Limit == 0 {
		p.Limit = DEFAULT_PAGINATION.Limit
	}
	if p.OrderDirection == "" {
		p.OrderDirection = DEFAULT_PAGINATION.OrderDirection
	}
	if p.OrderBy == "" {
		p.OrderBy = DEFAULT_PAGINATION.OrderBy
	}
	// Note: We don't set defaults for Cursor fields as empty strings might be valid
}

func EncodeCursor(id string, direction string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s_%s", direction, id)))
}

func DecodeCursor(cursor string) (string, string, error) {
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return "", "", err
	}
	parts := strings.SplitN(string(b), "_", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid cursor format")
	}
	return parts[1], parts[0], nil
}
