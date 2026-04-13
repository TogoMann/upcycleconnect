package item

import (
	db "backend/internal/database"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]Item, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, owner_id, container_id, site_id, material_type, physical_state, status, weight, created_at FROM item")
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Item])
}

func (r *Repository) UpdateStatus(id pgtype.Int8, status ItemStatus) error {
	tag, err := r.db.Exec(db.Ctx, "UPDATE item SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("item not found")
	}
	return nil
}

func (r *Repository) CreateAccessCode(itemId pgtype.Int8, userId pgtype.Int8, code string) error {
	_, err := r.db.Exec(db.Ctx,
		"INSERT INTO container_access (item_id, user_id, access_code, expires_at) VALUES ($1, $2, $3, NOW() + INTERVAL '24 hours')",
		itemId, userId, code)
	return err
}

func (r *Repository) GetAccessCode(itemId pgtype.Int8) (string, error) {
	var code string
	err := r.db.QueryRow(db.Ctx, "SELECT access_code FROM container_access WHERE item_id = $1 AND used_at IS NULL ORDER BY created_at DESC LIMIT 1", itemId).Scan(&code)
	return code, err
}

func (r *Repository) Collect(id pgtype.Int8, proId pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "UPDATE item SET status = 'collected', owner_id = $1 WHERE id = $2", proId, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("item not found")
	}
	return nil
}

func (r *Repository) MarkCodeUsed(code string) error {
	tag, err := r.db.Exec(db.Ctx, "UPDATE container_access SET used_at = NOW() WHERE access_code = $1", code)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid or already used code")
	}
	return nil
}
