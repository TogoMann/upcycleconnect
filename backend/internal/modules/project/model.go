package project

import "github.com/jackc/pgx/v5/pgtype"

type ProjectStatus string

const (
	InProgress ProjectStatus = "in progress"
	Done       ProjectStatus = "done"
	Featured   ProjectStatus = "featured"
	Cancelled  ProjectStatus = "cancelled"
)

type Project struct {
	Id          pgtype.Int8      `db:"id" json:"id"`
	ListingId   pgtype.Int8      `db:"listing_id" json:"listing_id"`
	CreatorId   pgtype.Int8      `db:"creator_id" json:"creator_id"`
	Title       string           `db:"title" json:"title"`
	Description string           `db:"description" json:"description"`
	FinalScore  pgtype.Int4      `db:"final_score" json:"final_score"`
	Status      ProjectStatus    `db:"status" json:"status"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
	CompletedAt pgtype.Timestamp `db:"completed_at" json:"completed_at"`
}

type ProjectStep struct {
	Id          pgtype.Int8      `db:"id" json:"id"`
	ProjectId   pgtype.Int8      `db:"project_id" json:"project_id"`
	StepNumber  int32            `db:"step_number" json:"step_number"`
	Description string           `db:"description" json:"description"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
}
