// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: debates.sql

package database

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const createComment = `-- name: CreateComment :one
INSERT INTO comments (debate_id, parent_comment_id, user_id, content)
VALUES ($1, $2, $3, $4)
RETURNING id, debate_id, parent_comment_id, user_id, content, created_at, updated_at
`

type CreateCommentParams struct {
	DebateID        sql.NullInt32
	ParentCommentID sql.NullInt32
	UserID          sql.NullInt32
	Content         string
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error) {
	row := q.db.QueryRowContext(ctx, createComment,
		arg.DebateID,
		arg.ParentCommentID,
		arg.UserID,
		arg.Content,
	)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.DebateID,
		&i.ParentCommentID,
		&i.UserID,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createDebate = `-- name: CreateDebate :one
INSERT INTO debates (match_id, debate_type, headline, description, ai_generated)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, match_id, debate_type, headline, description, ai_generated, deleted_at, created_at, updated_at
`

type CreateDebateParams struct {
	MatchID     string
	DebateType  string
	Headline    string
	Description sql.NullString
	AiGenerated sql.NullBool
}

func (q *Queries) CreateDebate(ctx context.Context, arg CreateDebateParams) (Debate, error) {
	row := q.db.QueryRowContext(ctx, createDebate,
		arg.MatchID,
		arg.DebateType,
		arg.Headline,
		arg.Description,
		arg.AiGenerated,
	)
	var i Debate
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.DebateType,
		&i.Headline,
		&i.Description,
		&i.AiGenerated,
		&i.DeletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createDebateAnalytics = `-- name: CreateDebateAnalytics :one
INSERT INTO debate_analytics (debate_id, total_votes, total_comments, engagement_score)
VALUES ($1, $2, $3, $4)
ON CONFLICT (debate_id) 
DO UPDATE SET 
    total_votes = $2,
    total_comments = $3,
    engagement_score = $4,
    updated_at = CURRENT_TIMESTAMP
RETURNING id, debate_id, total_votes, total_comments, engagement_score, created_at, updated_at
`

type CreateDebateAnalyticsParams struct {
	DebateID        sql.NullInt32
	TotalVotes      sql.NullInt32
	TotalComments   sql.NullInt32
	EngagementScore sql.NullString
}

func (q *Queries) CreateDebateAnalytics(ctx context.Context, arg CreateDebateAnalyticsParams) (DebateAnalytic, error) {
	row := q.db.QueryRowContext(ctx, createDebateAnalytics,
		arg.DebateID,
		arg.TotalVotes,
		arg.TotalComments,
		arg.EngagementScore,
	)
	var i DebateAnalytic
	err := row.Scan(
		&i.ID,
		&i.DebateID,
		&i.TotalVotes,
		&i.TotalComments,
		&i.EngagementScore,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createDebateCard = `-- name: CreateDebateCard :one
INSERT INTO debate_cards (debate_id, stance, title, description, ai_generated)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, debate_id, stance, title, description, ai_generated, created_at, updated_at
`

type CreateDebateCardParams struct {
	DebateID    sql.NullInt32
	Stance      string
	Title       string
	Description sql.NullString
	AiGenerated sql.NullBool
}

func (q *Queries) CreateDebateCard(ctx context.Context, arg CreateDebateCardParams) (DebateCard, error) {
	row := q.db.QueryRowContext(ctx, createDebateCard,
		arg.DebateID,
		arg.Stance,
		arg.Title,
		arg.Description,
		arg.AiGenerated,
	)
	var i DebateCard
	err := row.Scan(
		&i.ID,
		&i.DebateID,
		&i.Stance,
		&i.Title,
		&i.Description,
		&i.AiGenerated,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createVote = `-- name: CreateVote :one
INSERT INTO votes (debate_card_id, user_id, vote_type, emoji)
VALUES ($1, $2, $3, $4)
ON CONFLICT (debate_card_id, user_id, vote_type, emoji) 
DO UPDATE SET emoji = $4, created_at = CURRENT_TIMESTAMP
RETURNING id, debate_card_id, user_id, vote_type, emoji, created_at
`

type CreateVoteParams struct {
	DebateCardID sql.NullInt32
	UserID       sql.NullInt32
	VoteType     string
	Emoji        sql.NullString
}

func (q *Queries) CreateVote(ctx context.Context, arg CreateVoteParams) (Vote, error) {
	row := q.db.QueryRowContext(ctx, createVote,
		arg.DebateCardID,
		arg.UserID,
		arg.VoteType,
		arg.Emoji,
	)
	var i Vote
	err := row.Scan(
		&i.ID,
		&i.DebateCardID,
		&i.UserID,
		&i.VoteType,
		&i.Emoji,
		&i.CreatedAt,
	)
	return i, err
}

const deleteComment = `-- name: DeleteComment :exec
DELETE FROM comments WHERE id = $1
`

func (q *Queries) DeleteComment(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteComment, id)
	return err
}

const deleteDebate = `-- name: DeleteDebate :exec
DELETE FROM debates WHERE id = $1
`

func (q *Queries) DeleteDebate(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteDebate, id)
	return err
}

const deleteDebateCard = `-- name: DeleteDebateCard :exec
DELETE FROM debate_cards WHERE id = $1
`

func (q *Queries) DeleteDebateCard(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteDebateCard, id)
	return err
}

const deleteVote = `-- name: DeleteVote :exec
DELETE FROM votes WHERE debate_card_id = $1 AND user_id = $2 AND vote_type = $3
`

type DeleteVoteParams struct {
	DebateCardID sql.NullInt32
	UserID       sql.NullInt32
	VoteType     string
}

func (q *Queries) DeleteVote(ctx context.Context, arg DeleteVoteParams) error {
	_, err := q.db.ExecContext(ctx, deleteVote, arg.DebateCardID, arg.UserID, arg.VoteType)
	return err
}

const getComment = `-- name: GetComment :one
SELECT 
    c.id, c.debate_id, c.parent_comment_id, c.user_id, c.content, c.created_at, c.updated_at,
    u.firstname,
    u.lastname
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.id = $1
`

type GetCommentRow struct {
	ID              int32
	DebateID        sql.NullInt32
	ParentCommentID sql.NullInt32
	UserID          sql.NullInt32
	Content         string
	CreatedAt       sql.NullTime
	UpdatedAt       sql.NullTime
	Firstname       string
	Lastname        string
}

func (q *Queries) GetComment(ctx context.Context, id int32) (GetCommentRow, error) {
	row := q.db.QueryRowContext(ctx, getComment, id)
	var i GetCommentRow
	err := row.Scan(
		&i.ID,
		&i.DebateID,
		&i.ParentCommentID,
		&i.UserID,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Firstname,
		&i.Lastname,
	)
	return i, err
}

const getCommentCount = `-- name: GetCommentCount :one
SELECT COUNT(*) FROM comments WHERE debate_id = $1
`

func (q *Queries) GetCommentCount(ctx context.Context, debateID sql.NullInt32) (int64, error) {
	row := q.db.QueryRowContext(ctx, getCommentCount, debateID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getComments = `-- name: GetComments :many
SELECT 
    c.id, c.debate_id, c.parent_comment_id, c.user_id, c.content, c.created_at, c.updated_at,
    u.firstname,
    u.lastname
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.debate_id = $1
ORDER BY c.created_at ASC
`

type GetCommentsRow struct {
	ID              int32
	DebateID        sql.NullInt32
	ParentCommentID sql.NullInt32
	UserID          sql.NullInt32
	Content         string
	CreatedAt       sql.NullTime
	UpdatedAt       sql.NullTime
	Firstname       string
	Lastname        string
}

func (q *Queries) GetComments(ctx context.Context, debateID sql.NullInt32) ([]GetCommentsRow, error) {
	rows, err := q.db.QueryContext(ctx, getComments, debateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCommentsRow
	for rows.Next() {
		var i GetCommentsRow
		if err := rows.Scan(
			&i.ID,
			&i.DebateID,
			&i.ParentCommentID,
			&i.UserID,
			&i.Content,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Firstname,
			&i.Lastname,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDebate = `-- name: GetDebate :one
SELECT id, match_id, debate_type, headline, description, ai_generated, deleted_at, created_at, updated_at FROM debates WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) GetDebate(ctx context.Context, id int32) (Debate, error) {
	row := q.db.QueryRowContext(ctx, getDebate, id)
	var i Debate
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.DebateType,
		&i.Headline,
		&i.Description,
		&i.AiGenerated,
		&i.DeletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getDebateAnalytics = `-- name: GetDebateAnalytics :one
SELECT id, debate_id, total_votes, total_comments, engagement_score, created_at, updated_at FROM debate_analytics WHERE debate_id = $1
`

func (q *Queries) GetDebateAnalytics(ctx context.Context, debateID sql.NullInt32) (DebateAnalytic, error) {
	row := q.db.QueryRowContext(ctx, getDebateAnalytics, debateID)
	var i DebateAnalytic
	err := row.Scan(
		&i.ID,
		&i.DebateID,
		&i.TotalVotes,
		&i.TotalComments,
		&i.EngagementScore,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getDebateCard = `-- name: GetDebateCard :one
SELECT id, debate_id, stance, title, description, ai_generated, created_at, updated_at FROM debate_cards WHERE id = $1
`

func (q *Queries) GetDebateCard(ctx context.Context, id int32) (DebateCard, error) {
	row := q.db.QueryRowContext(ctx, getDebateCard, id)
	var i DebateCard
	err := row.Scan(
		&i.ID,
		&i.DebateID,
		&i.Stance,
		&i.Title,
		&i.Description,
		&i.AiGenerated,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getDebateCards = `-- name: GetDebateCards :many
SELECT id, debate_id, stance, title, description, ai_generated, created_at, updated_at FROM debate_cards WHERE debate_id = $1 ORDER BY stance
`

func (q *Queries) GetDebateCards(ctx context.Context, debateID sql.NullInt32) ([]DebateCard, error) {
	rows, err := q.db.QueryContext(ctx, getDebateCards, debateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DebateCard
	for rows.Next() {
		var i DebateCard
		if err := rows.Scan(
			&i.ID,
			&i.DebateID,
			&i.Stance,
			&i.Title,
			&i.Description,
			&i.AiGenerated,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDebatesByMatch = `-- name: GetDebatesByMatch :many
SELECT id, match_id, debate_type, headline, description, ai_generated, deleted_at, created_at, updated_at FROM debates 
WHERE match_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
`

func (q *Queries) GetDebatesByMatch(ctx context.Context, matchID string) ([]Debate, error) {
	rows, err := q.db.QueryContext(ctx, getDebatesByMatch, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Debate
	for rows.Next() {
		var i Debate
		if err := rows.Scan(
			&i.ID,
			&i.MatchID,
			&i.DebateType,
			&i.Headline,
			&i.Description,
			&i.AiGenerated,
			&i.DeletedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDebatesByType = `-- name: GetDebatesByType :many
SELECT id, match_id, debate_type, headline, description, ai_generated, deleted_at, created_at, updated_at FROM debates 
WHERE debate_type = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
`

func (q *Queries) GetDebatesByType(ctx context.Context, debateType string) ([]Debate, error) {
	rows, err := q.db.QueryContext(ctx, getDebatesByType, debateType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Debate
	for rows.Next() {
		var i Debate
		if err := rows.Scan(
			&i.ID,
			&i.MatchID,
			&i.DebateType,
			&i.Headline,
			&i.Description,
			&i.AiGenerated,
			&i.DeletedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTopDebates = `-- name: GetTopDebates :many
SELECT 
    d.id, d.match_id, d.debate_type, d.headline, d.description, d.ai_generated, d.deleted_at, d.created_at, d.updated_at,
    da.total_votes,
    da.total_comments,
    da.engagement_score
FROM debates d
LEFT JOIN debate_analytics da ON d.id = da.debate_id
WHERE d.deleted_at IS NULL
ORDER BY da.engagement_score DESC NULLS LAST
LIMIT $1
`

type GetTopDebatesRow struct {
	ID              int32
	MatchID         string
	DebateType      string
	Headline        string
	Description     sql.NullString
	AiGenerated     sql.NullBool
	DeletedAt       sql.NullTime
	CreatedAt       sql.NullTime
	UpdatedAt       sql.NullTime
	TotalVotes      sql.NullInt32
	TotalComments   sql.NullInt32
	EngagementScore sql.NullString
}

func (q *Queries) GetTopDebates(ctx context.Context, limit int32) ([]GetTopDebatesRow, error) {
	rows, err := q.db.QueryContext(ctx, getTopDebates, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTopDebatesRow
	for rows.Next() {
		var i GetTopDebatesRow
		if err := rows.Scan(
			&i.ID,
			&i.MatchID,
			&i.DebateType,
			&i.Headline,
			&i.Description,
			&i.AiGenerated,
			&i.DeletedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.TotalVotes,
			&i.TotalComments,
			&i.EngagementScore,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserVote = `-- name: GetUserVote :one
SELECT id, debate_card_id, user_id, vote_type, emoji, created_at FROM votes WHERE debate_card_id = $1 AND user_id = $2 AND vote_type = $3
`

type GetUserVoteParams struct {
	DebateCardID sql.NullInt32
	UserID       sql.NullInt32
	VoteType     string
}

func (q *Queries) GetUserVote(ctx context.Context, arg GetUserVoteParams) (Vote, error) {
	row := q.db.QueryRowContext(ctx, getUserVote, arg.DebateCardID, arg.UserID, arg.VoteType)
	var i Vote
	err := row.Scan(
		&i.ID,
		&i.DebateCardID,
		&i.UserID,
		&i.VoteType,
		&i.Emoji,
		&i.CreatedAt,
	)
	return i, err
}

const getVoteCounts = `-- name: GetVoteCounts :many
SELECT 
    debate_card_id,
    vote_type,
    emoji,
    COUNT(*) as count
FROM votes 
WHERE debate_card_id = ANY($1::int[])
GROUP BY debate_card_id, vote_type, emoji
`

type GetVoteCountsRow struct {
	DebateCardID sql.NullInt32
	VoteType     string
	Emoji        sql.NullString
	Count        int64
}

func (q *Queries) GetVoteCounts(ctx context.Context, dollar_1 []int32) ([]GetVoteCountsRow, error) {
	rows, err := q.db.QueryContext(ctx, getVoteCounts, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetVoteCountsRow
	for rows.Next() {
		var i GetVoteCountsRow
		if err := rows.Scan(
			&i.DebateCardID,
			&i.VoteType,
			&i.Emoji,
			&i.Count,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getVotesByCard = `-- name: GetVotesByCard :many
SELECT id, debate_card_id, user_id, vote_type, emoji, created_at FROM votes WHERE debate_card_id = $1
`

func (q *Queries) GetVotesByCard(ctx context.Context, debateCardID sql.NullInt32) ([]Vote, error) {
	rows, err := q.db.QueryContext(ctx, getVotesByCard, debateCardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Vote
	for rows.Next() {
		var i Vote
		if err := rows.Scan(
			&i.ID,
			&i.DebateCardID,
			&i.UserID,
			&i.VoteType,
			&i.Emoji,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const restoreDebate = `-- name: RestoreDebate :exec
UPDATE debates SET deleted_at = NULL WHERE id = $1
`

func (q *Queries) RestoreDebate(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, restoreDebate, id)
	return err
}

const softDeleteDebate = `-- name: SoftDeleteDebate :exec
UPDATE debates SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1
`

func (q *Queries) SoftDeleteDebate(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, softDeleteDebate, id)
	return err
}

const updateComment = `-- name: UpdateComment :one
UPDATE comments 
SET content = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, debate_id, parent_comment_id, user_id, content, created_at, updated_at
`

type UpdateCommentParams struct {
	ID      int32
	Content string
}

func (q *Queries) UpdateComment(ctx context.Context, arg UpdateCommentParams) (Comment, error) {
	row := q.db.QueryRowContext(ctx, updateComment, arg.ID, arg.Content)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.DebateID,
		&i.ParentCommentID,
		&i.UserID,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateDebate = `-- name: UpdateDebate :one
UPDATE debates 
SET headline = $2, description = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, match_id, debate_type, headline, description, ai_generated, deleted_at, created_at, updated_at
`

type UpdateDebateParams struct {
	ID          int32
	Headline    string
	Description sql.NullString
}

func (q *Queries) UpdateDebate(ctx context.Context, arg UpdateDebateParams) (Debate, error) {
	row := q.db.QueryRowContext(ctx, updateDebate, arg.ID, arg.Headline, arg.Description)
	var i Debate
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.DebateType,
		&i.Headline,
		&i.Description,
		&i.AiGenerated,
		&i.DeletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateDebateAnalytics = `-- name: UpdateDebateAnalytics :one
UPDATE debate_analytics 
SET total_votes = $2, total_comments = $3, engagement_score = $4, updated_at = CURRENT_TIMESTAMP
WHERE debate_id = $1
RETURNING id, debate_id, total_votes, total_comments, engagement_score, created_at, updated_at
`

type UpdateDebateAnalyticsParams struct {
	DebateID        sql.NullInt32
	TotalVotes      sql.NullInt32
	TotalComments   sql.NullInt32
	EngagementScore sql.NullString
}

func (q *Queries) UpdateDebateAnalytics(ctx context.Context, arg UpdateDebateAnalyticsParams) (DebateAnalytic, error) {
	row := q.db.QueryRowContext(ctx, updateDebateAnalytics,
		arg.DebateID,
		arg.TotalVotes,
		arg.TotalComments,
		arg.EngagementScore,
	)
	var i DebateAnalytic
	err := row.Scan(
		&i.ID,
		&i.DebateID,
		&i.TotalVotes,
		&i.TotalComments,
		&i.EngagementScore,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateDebateCard = `-- name: UpdateDebateCard :one
UPDATE debate_cards 
SET title = $2, description = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, debate_id, stance, title, description, ai_generated, created_at, updated_at
`

type UpdateDebateCardParams struct {
	ID          int32
	Title       string
	Description sql.NullString
}

func (q *Queries) UpdateDebateCard(ctx context.Context, arg UpdateDebateCardParams) (DebateCard, error) {
	row := q.db.QueryRowContext(ctx, updateDebateCard, arg.ID, arg.Title, arg.Description)
	var i DebateCard
	err := row.Scan(
		&i.ID,
		&i.DebateID,
		&i.Stance,
		&i.Title,
		&i.Description,
		&i.AiGenerated,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
