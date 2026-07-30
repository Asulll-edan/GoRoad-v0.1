package lazy

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

type Loader struct {
	RequestedIncludes map[string]bool
}

func NewLoader(c fiber.Ctx) *Loader {
	include := c.Query("include")
	requested := make(map[string]bool)

	if include != "" {
		for _, item := range strings.Split(include, ",") {
			requested[strings.TrimSpace(item)] = true
		}
	}

	return &Loader{RequestedIncludes: requested}
}

func (l *Loader) ShouldInclude(field string) bool {
	if l == nil || l.RequestedIncludes == nil {
		return false
	}
	return l.RequestedIncludes[field]
}

func (l *Loader) HasAny() bool {
	return l != nil && len(l.RequestedIncludes) > 0
}

// PreloadConfig untuk GORM preload dengan pagination
type PreloadConfig struct {
	Query       string
	Args        []interface{}
	Limit       int
	Order       string
}

func (l *Loader) GetPreloads() map[string]*PreloadConfig {
	preloads := make(map[string]*PreloadConfig)

	if l == nil {
		return preloads
	}

	if l.ShouldInclude("members") {
		preloads["Members"] = &PreloadConfig{
			Limit: 20,
			Order: "joined_at DESC",
		}
	}

	if l.ShouldInclude("active_route") {
		preloads["ActiveRoute"] = nil
	}

	if l.ShouldInclude("weather") {
		preloads["Weather"] = nil
	}

	if l.ShouldInclude("settings") {
		preloads["Settings"] = nil
	}

	if l.ShouldInclude("motors") {
		preloads["Motors"] = &PreloadConfig{
			Limit: 5,
		}
	}

	if l.ShouldInclude("badges") {
		preloads["Badges"] = nil
	}

	if l.ShouldInclude("stats") {
		preloads["Stats"] = nil
	}

	if l.ShouldInclude("comments") {
		preloads["Comments"] = &PreloadConfig{
			Limit: 20,
			Order: "created_at DESC",
		}
	}

	return preloads
}

// FilterFields membatasi fields yang dikembalikan berdasarkan sparse fieldsets
func FilterFields(entity string, requestedFields []string) map[string]bool {
	allowedFields := DefaultFields(entity)
	if len(requestedFields) == 0 {
		return allowedFields
	}

	filtered := make(map[string]bool)
	for _, f := range requestedFields {
		if allowedFields[f] {
			filtered[f] = true
		}
	}
	return filtered
}

func DefaultFields(entity string) map[string]bool {
	defaults := map[string]bool{
		"id": true,
	}

	switch entity {
	case "room":
		defaults["name"] = true
		defaults["status"] = true
		defaults["start_date"] = true
		defaults["end_date"] = true
		defaults["start_location"] = true
		defaults["end_location"] = true
		defaults["distance_km"] = true
		defaults["cover_photo_url"] = true
	case "user":
		defaults["username"] = true
		defaults["full_name"] = true
		defaults["photo_url"] = true
		defaults["riding_skill"] = true
		defaults["riding_role"] = true
	case "member":
		defaults["role"] = true
		defaults["position_in_formation"] = true
		defaults["joined_at"] = true
	default:
		defaults["created_at"] = true
		defaults["updated_at"] = true
	}

	return defaults
}

// Response wrapper untuk data yang di-lazy load
type LazyResponse struct {
	Data    interface{}   `json:"data"`
	Related RelatedData   `json:"_related,omitempty"`
}

type RelatedData struct {
	Members     interface{} `json:"members,omitempty"`
	Routes      interface{} `json:"routes,omitempty"`
	Weather     interface{} `json:"weather,omitempty"`
	Settings    interface{} `json:"settings,omitempty"`
	Comments    interface{} `json:"comments,omitempty"`
	Badges      interface{} `json:"badges,omitempty"`
	Motors      interface{} `json:"motors,omitempty"`
	Stats       interface{} `json:"stats,omitempty"`
}
