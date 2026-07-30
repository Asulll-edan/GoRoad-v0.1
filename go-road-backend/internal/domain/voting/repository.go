package voting

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	CreateVoting(ctx context.Context, v *Voting) error
	FindVotingByID(ctx context.Context, id uuid.UUID) (*Voting, error)
	ListVotingsByRoomID(ctx context.Context, roomID uuid.UUID, cursor string, limit int) ([]Voting, string, bool, error)
	UpdateVoting(ctx context.Context, v *Voting) error

	AddAnswer(ctx context.Context, a *VotingAnswer) error
	ListAnswers(ctx context.Context, votingID uuid.UUID) ([]VotingAnswer, error)

	SubmitVote(ctx context.Context, vote *Vote) error
	GetVote(ctx context.Context, votingID, answerID, userID uuid.UUID) (*Vote, error)
	GetResults(ctx context.Context, votingID uuid.UUID) ([]VoteResult, error)
	CloseVoting(ctx context.Context, id uuid.UUID) error
}

type VoteResult struct {
	AnswerID   uuid.UUID `json:"answer_id"`
	Label      string    `json:"label"`
	VoteCount  int64     `json:"vote_count"`
	Percentage float64   `json:"percentage"`
}
