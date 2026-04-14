package thread

import "github.com/jackc/pgx/v5/pgtype"

type ThreadCategory string

const (
	Bricolage  ThreadCategory = "Bricolage"
	Textile    ThreadCategory = "Textile"
	Ressources ThreadCategory = "Ressources"
	Debutants  ThreadCategory = "Débutants"
	Communaute ThreadCategory = "Communauté"
)

type Thread struct {
	Id         pgtype.Int8      `db:"id" json:"id"`
	CreatedBy  pgtype.Int8      `db:"created_by" json:"created_by"`
	Category   ThreadCategory   `db:"category" json:"category"`
	Title      string           `db:"title" json:"title"`
	Content    string           `db:"content" json:"content"`
	Upvotes    int32            `db:"upvotes" json:"upvotes"`
	Downvotes  int32            `db:"downvotes" json:"downvotes"`
	CreatedAt  pgtype.Timestamp `db:"created_at" json:"created_at"`
	LastPostAt pgtype.Timestamp `db:"last_post_at" json:"last_post_at"`
}
