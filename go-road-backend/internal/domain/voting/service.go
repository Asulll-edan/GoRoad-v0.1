package voting

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	CreateVoting(ctx context.Context, roomID, userID uuid.UUID, title, description, votingType string, answers []string) (*Voting, error)
	ListVotings(ctx context.Context, roomID uuid.UUID, cursor string, limit int) ([]Voting, string, bool, error)
	GetVoting(ctx context.Context, id uuid.UUID) (*Voting, error)
	SubmitVote(ctx context.Context, votingID, answerID, userID uuid.UUID) error
	CloseVoting(ctx context.Context, id, userID uuid.UUID) error
	GetResults(ctx context.Context, id uuid.UUID) ([]VoteResult, error)
}
