package request

type CreateVotingRequest struct {
	Title       string   `json:"title" validate:"required,max=200"`
	Description string   `json:"description,omitempty"`
	Options     []string `json:"options" validate:"required,min=2,max=10"`
	DurationMin int      `json:"duration_minutes" validate:"min=1,max=1440"`
	IsAnonymous bool     `json:"is_anonymous"`
	MaxChoices  int      `json:"max_choices" validate:"min=1,max=5"`
}

type SubmitVoteRequest struct {
	VotingID string   `json:"voting_id" validate:"required"`
	OptionIDs []string `json:"option_ids" validate:"required,min=1"`
}
