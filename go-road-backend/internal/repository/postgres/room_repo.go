package postgres

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	domain "go-road-backend/internal/domain/room"
)

type roomRepository struct {
	db *Database
}

func NewRoomRepository(db *Database) domain.Repository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *domain.TouringRoom) error {
	return r.db.WithContext(ctx).Create(room).Error
}

func (r *roomRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.TouringRoom, error) {
	var room domain.TouringRoom
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.TouringRoom, error) {
	var rooms []domain.TouringRoom
	err := r.db.WithContext(ctx).Where("id IN ? AND deleted_at IS NULL", ids).Find(&rooms).Error
	return rooms, err
}

func (r *roomRepository) Update(ctx context.Context, room *domain.TouringRoom) error {
	return r.db.WithContext(ctx).Save(room).Error
}

func (r *roomRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.TouringRoom{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *roomRepository) List(ctx context.Context, cursor string, limit int, filters map[string]interface{}) ([]domain.TouringRoom, string, bool, error) {
	var rooms []domain.TouringRoom
	db := r.db.WithContext(ctx).Where("deleted_at IS NULL")

	_ = filters

	if cursor != "" {
		data, err := base64.URLEncoding.DecodeString(cursor)
		if err == nil {
			var c struct {
				ID        string    `json:"id"`
				CreatedAt time.Time `json:"created_at"`
			}
			if err := json.Unmarshal(data, &c); err == nil {
				db = db.Where("(created_at, id) < (?, ?)", c.CreatedAt, c.ID)
			}
		}
	}

	db = db.Order("created_at DESC, id DESC").Limit(limit + 1)
	if err := db.Find(&rooms).Error; err != nil {
		return nil, "", false, err
	}

	hasMore := len(rooms) > limit
	if hasMore {
		rooms = rooms[:limit]
	}

	var nextCursor string
	if len(rooms) > 0 {
		last := rooms[len(rooms)-1]
		cursorData, _ := json.Marshal(struct {
			ID        string    `json:"id"`
			CreatedAt time.Time `json:"created_at"`
		}{ID: last.ID.String(), CreatedAt: last.CreatedAt})
		nextCursor = base64.URLEncoding.EncodeToString(cursorData)
	}

	return rooms, nextCursor, hasMore, nil
}

func (r *roomRepository) Discover(ctx context.Context, cursor string, limit int) ([]domain.TouringRoom, string, bool, error) {
	var rooms []domain.TouringRoom
	db := r.db.WithContext(ctx).Where("is_public = true AND status IN ('planning', 'ready') AND deleted_at IS NULL")

	if cursor != "" {
		data, err := base64.URLEncoding.DecodeString(cursor)
		if err == nil {
			var c struct {
				StartDate time.Time `json:"start_date"`
				ID        string    `json:"id"`
			}
			if err := json.Unmarshal(data, &c); err == nil {
				db = db.Where("(start_date, id) > (?, ?)", c.StartDate, c.ID)
			}
		}
	}

	db = db.Order("start_date ASC, id ASC").Limit(limit + 1)
	if err := db.Find(&rooms).Error; err != nil {
		return nil, "", false, err
	}

	hasMore := len(rooms) > limit
	if hasMore {
		rooms = rooms[:limit]
	}

	var nextCursor string
	if len(rooms) > 0 {
		last := rooms[len(rooms)-1]
		if last.StartDate != nil {
			cursorData, _ := json.Marshal(struct {
				StartDate time.Time `json:"start_date"`
				ID        string    `json:"id"`
			}{StartDate: *last.StartDate, ID: last.ID.String()})
			nextCursor = base64.URLEncoding.EncodeToString(cursorData)
		}
	}

	return rooms, nextCursor, hasMore, nil
}

func (r *roomRepository) AddMember(ctx context.Context, member *domain.RoomMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *roomRepository) RemoveMember(ctx context.Context, roomID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("room_id = ? AND user_id = ?", roomID, userID).Delete(&domain.RoomMember{}).Error
}

func (r *roomRepository) UpdateMemberRole(ctx context.Context, roomID, userID uuid.UUID, role string) error {
	return r.db.WithContext(ctx).Model(&domain.RoomMember{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Update("role", role).Error
}

func (r *roomRepository) ListMembers(ctx context.Context, roomID uuid.UUID, cursor string, limit int) ([]domain.RoomMember, string, bool, error) {
	var members []domain.RoomMember
	db := r.db.WithContext(ctx).Where("room_id = ?", roomID)

	if cursor != "" {
		data, err := base64.URLEncoding.DecodeString(cursor)
		if err == nil {
			var c struct {
				ID       string    `json:"id"`
				JoinedAt time.Time `json:"joined_at"`
			}
			if err := json.Unmarshal(data, &c); err == nil {
				db = db.Where("(joined_at, id) > (?, ?)", c.JoinedAt, c.ID)
			}
		}
	}

	db = db.Order("joined_at ASC, id ASC").Limit(limit + 1)
	if err := db.Find(&members).Error; err != nil {
		return nil, "", false, err
	}

	hasMore := len(members) > limit
	if hasMore {
		members = members[:limit]
	}

	var nextCursor string
	if len(members) > 0 {
		last := members[len(members)-1]
		cursorData, _ := json.Marshal(struct {
			ID       string    `json:"id"`
			JoinedAt time.Time `json:"joined_at"`
		}{ID: last.ID.String(), JoinedAt: last.JoinedAt})
		nextCursor = base64.URLEncoding.EncodeToString(cursorData)
	}

	return members, nextCursor, hasMore, nil
}

func (r *roomRepository) GetMember(ctx context.Context, roomID, userID uuid.UUID) (*domain.RoomMember, error) {
	var member domain.RoomMember
	err := r.db.WithContext(ctx).Where("room_id = ? AND user_id = ?", roomID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *roomRepository) CountMembers(ctx context.Context, roomID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.RoomMember{}).Where("room_id = ?", roomID).Count(&count).Error
	return count, err
}

func (r *roomRepository) AddRoleHistory(ctx context.Context, history *domain.RoomRoleHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}
