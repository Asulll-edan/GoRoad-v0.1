package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

type Cursor struct {
	ID        string    `json:"id"`
	Value     string    `json:"value,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type PaginationParams struct {
	Cursor string `query:"cursor"`
	Limit  int    `query:"limit"`
	Sort   string `query:"sort"`
	Order  string `query:"order" validate:"oneof=asc desc"`
}

type PaginatedResponse struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

type Meta struct {
	Cursor  string `json:"cursor,omitempty"`
	HasMore bool   `json:"has_more"`
	Total   *int64 `json:"total,omitempty"`
}

func ParsePaginationParams(c fiber.Ctx) PaginationParams {
	params := PaginationParams{
		Cursor: c.Query("cursor"),
		Limit:  DefaultLimit,
		Order:  "desc",
	}

	if limit, err := strconv.Atoi(c.Query("limit")); err == nil && limit > 0 {
		if limit > MaxLimit {
			params.Limit = MaxLimit
		} else {
			params.Limit = limit
		}
	}

	if sort := c.Query("sort"); sort != "" {
		params.Sort = sort
	}

	if order := c.Query("order"); order == "asc" {
		params.Order = "asc"
	}

	return params
}

func EncodeCursor(cursor Cursor) string {
	data, _ := json.Marshal(cursor)
	return base64.URLEncoding.EncodeToString(data)
}

func DecodeCursor(cursorStr string) (*Cursor, error) {
	data, err := base64.URLEncoding.DecodeString(cursorStr)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor: %w", err)
	}

	var cursor Cursor
	if err := json.Unmarshal(data, &cursor); err != nil {
		return nil, fmt.Errorf("invalid cursor data: %w", err)
	}

	return &cursor, nil
}

func ApplyCursorPagination(db *gorm.DB, params PaginationParams, cursorField string, tableName string) (*gorm.DB, *Cursor, error) {
	var cursor *Cursor

	if params.Cursor != "" {
		var err error
		cursor, err = DecodeCursor(params.Cursor)
		if err != nil {
			return db, nil, err
		}

		operator := "<"
		if params.Order == "asc" {
			operator = ">"
		}

		whereClause := fmt.Sprintf("(%s, %s.id) %s (?, ?)", cursorField, tableName, operator)
		db = db.Where(whereClause, cursor.CreatedAt, cursor.ID)
	}

	orderClause := fmt.Sprintf("%s %s, %s.id %s", cursorField, params.Order, tableName, params.Order)
	db = db.Order(orderClause)

	// Fetch limit + 1 to check has_more
	db = db.Limit(params.Limit + 1)

	return db, cursor, nil
}

func BuildPaginatedResponse(data interface{}, itemsCount int, limit int, cursorField string, lastID string, lastTime time.Time) PaginatedResponse {
	meta := Meta{
		HasMore: itemsCount > limit,
	}

	if meta.HasMore && itemsCount > 0 {
		nextCursor := EncodeCursor(Cursor{
			ID:        lastID,
			CreatedAt: lastTime,
		})
		meta.Cursor = nextCursor
	}

	return PaginatedResponse{
		Data: data,
		Meta: meta,
	}
}

func BuildPaginatedWithTotal(data interface{}, itemsCount int, limit int, cursorField string, lastID string, lastTime time.Time, total int64) PaginatedResponse {
	resp := BuildPaginatedResponse(data, itemsCount, limit, cursorField, lastID, lastTime)
	resp.Meta.Total = &total
	return resp
}

// For UUID based cursors
func ApplyCursorPaginationUUID(db *gorm.DB, params PaginationParams, cursorField string) (*gorm.DB, *Cursor, error) {
	var cursor *Cursor

	if params.Cursor != "" {
		var err error
		cursor, err = DecodeCursor(params.Cursor)
		if err != nil {
			return db, nil, err
		}

		operator := "<"
		if params.Order == "asc" {
			operator = ">"
		}

		db = db.Where(fmt.Sprintf("(%s, id) %s (?, ?::uuid)", cursorField, operator), cursor.CreatedAt, cursor.ID)
	}

	db = db.Order(fmt.Sprintf("%s %s, id %s", cursorField, params.Order, params.Order))
	db = db.Limit(params.Limit + 1)

	return db, cursor, nil
}

// For numeric score based cursors (leaderboard)
type ScoreCursor struct {
	Score    float64 `json:"score"`
	MemberID string  `json:"member_id"`
}

func EncodeScoreCursor(score float64, memberID string) string {
	cursor := ScoreCursor{Score: score, MemberID: memberID}
	data, _ := json.Marshal(cursor)
	return base64.URLEncoding.EncodeToString(data)
}

func DecodeScoreCursor(cursorStr string) (*ScoreCursor, error) {
	data, err := base64.URLEncoding.DecodeString(cursorStr)
	if err != nil {
		return nil, err
	}
	var cursor ScoreCursor
	if err := json.Unmarshal(data, &cursor); err != nil {
		return nil, err
	}
	return &cursor, nil
}

// For paginated list response with slice trimming
func TrimPaginatedData(data interface{}, limit int) (interface{}, bool) {
	itemsCount := 0

	switch v := data.(type) {
	case *[]interface{}:
		itemsCount = len(*v)
		if itemsCount > limit {
			*v = (*v)[:limit]
		}
	}

	return data, itemsCount > limit
}

// Helper to get last element fields
func GetLastCursor[T any](items []T, getID func(T) string, getTime func(T) time.Time) (string, time.Time) {
	if len(items) == 0 {
		return "", time.Time{}
	}
	last := items[len(items)-1]
	return getID(last), getTime(last)
}

// Float64 comparison for cursor
const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

// Parse comma-separated include parameter for lazy loading
func ParseIncludeParams(c fiber.Ctx) []string {
	include := c.Query("include")
	if include == "" {
		return nil
	}
	return strings.Split(include, ",")
}
