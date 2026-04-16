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
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			i.id, i.owner_id, i.container_id, i.site_id, i.material_type, 
			i.physical_state, i.status, i.weight, i.created_at,
			s.type_site as site_type
		FROM item i
		LEFT JOIN site s ON i.site_id = s.id
	`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Item])
}

func (r *Repository) GetById(id pgtype.Int8) (*Item, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			i.id, i.owner_id, i.container_id, i.site_id, i.material_type, 
			i.physical_state, i.status, i.weight, i.created_at,
			s.type_site as site_type
		FROM item i
		LEFT JOIN site s ON i.site_id = s.id
		WHERE i.id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Item])
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository) Create(item Item) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx,
		"INSERT INTO item (owner_id, container_id, site_id, material_type, physical_state, status, weight) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		item.OwnerId, item.ContainerId, item.SiteId, item.MaterialType, item.PhysicalState, item.Status, item.Weight).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Update(id pgtype.Int8, item Item) error {
	tag, err := r.db.Exec(db.Ctx,
		"UPDATE item SET owner_id=$1, container_id=$2, site_id=$3, material_type=$4, physical_state=$5, status=$6, weight=$7 WHERE id=$8",
		item.OwnerId, item.ContainerId, item.SiteId, item.MaterialType, item.PhysicalState, item.Status, item.Weight, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("item not found with id %d", id.Int64)
	}
	return nil
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM item WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("item not found with id %d", id.Int64)
	}
	return nil
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
