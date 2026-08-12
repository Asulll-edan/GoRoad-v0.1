package postgres

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	domain "go-road-backend/internal/domain/voting"
)

type votingRepository struct {
	db *Database
}

func NewVotingRepository(db *Database) domain.Repository {
	return &votingRepository{db: db}
}

func (r *votingRepository) CreateVoting(ctx context.Context, v *domain.Voting) error {
	return r.db.WithContext(ctx).Create(v).Error
}

func (r *votingRepository) FindVotingByID(ctx context.Context, id uuid.UUID) (*domain.Voting, error) {
	var v domain.Voting
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&v).Error
	return &v, err
}

func (r *votingRepository) ListVotingsByRoomID(ctx context.Context, roomID uuid.UUID, cursor string, limit int) ([]domain.Voting, string, bool, error) {
	var votings []domain.Voting
	db := r.db.WithContext(ctx).Where("room_id = ?", roomID)
	if cursor != "" {
		data, _ := base64.URLEncoding.DecodeString(cursor)
		var c struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}
		json.Unmarshal(data, &c)
		db = db.Where("(created_at, id) < (?, ?)", c.CreatedAt, c.ID)
	}
	db = db.Order("created_at DESC, id DESC").Limit(limit + 1)
	if err := db.Find(&votings).Error; err != nil {
		return nil, "", false, err
	}
	hasMore := len(votings) > limit
	if hasMore {
		votings = votings[:limit]
	}
	var nextCursor string
	if len(votings) > 0 {
		last := votings[len(votings)-1]
		cData, _ := json.Marshal(struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}{CreatedAt: last.CreatedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cData)
	}
	return votings, nextCursor, hasMore, nil
}

func (r *votingRepository) UpdateVoting(ctx context.Context, v *domain.Voting) error {
	return r.db.WithContext(ctx).Save(v).Error
}

func (r *votingRepository) AddAnswer(ctx context.Context, a *domain.VotingAnswer) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *votingRepository) ListAnswers(ctx context.Context, votingID uuid.UUID) ([]domain.VotingAnswer, error) {
	var answers []domain.VotingAnswer
	err := r.db.WithContext(ctx).Where("voting_id = ?", votingID).
		Order("order_index ASC").Find(&answers).Error
	return answers, err
}

func (r *votingRepository) SubmitVote(ctx context.Context, vote *domain.Vote) error {
	return r.db.WithContext(ctx).Create(vote).Error
}

func (r *votingRepository) GetVote(ctx context.Context, votingID, answerID, userID uuid.UUID) (*domain.Vote, error) {
	var vote domain.Vote
	err := r.db.WithContext(ctx).Where("voting_id = ? AND answer_id = ? AND user_id = ?",
		votingID, answerID, userID).First(&vote).Error
	return &vote, err
}

func (r *votingRepository) GetResults(ctx context.Context, votingID uuid.UUID) ([]domain.VoteResult, error) {
	var results []domain.VoteResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT va.id as answer_id, va.label,
			COUNT(vv.id) as vote_count,
			CASE WHEN COUNT(*) OVER() > 0
				THEN ROUND(COUNT(vv.id)::decimal / SUM(COUNT(vv.id)) OVER() * 100, 1)
				ELSE 0
			END as percentage
		FROM voting_answers va
		LEFT JOIN voting_votes vv ON vv.answer_id = va.id AND vv.voting_id = ?
		WHERE va.voting_id = ?
		GROUP BY va.id, va.label
		ORDER BY va.order_index`, votingID, votingID).Scan(&results).Error
	return results, err
}

func (r *votingRepository) CloseVoting(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.Voting{}).Where("id = ?", id).
		Update("status", "closed").Error
}
