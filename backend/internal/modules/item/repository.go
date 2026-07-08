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
			i.id, i.owner_id, i.locker_id, i.site_id, i.material_type,
			i.physical_state, i.size, i.status, i.weight, i.name, i.description, i.created_at,
			s.type_site as site_type
		FROM item i
		LEFT JOIN site s ON i.site_id = s.id
	`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Item])
}

func (r *Repository) GetAdminDepots() ([]AdminDepot, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			i.id, 
			u.username as utilisateur, 
			i.material_type as objet, 
			TO_CHAR(i.created_at, 'DD/MM/YYYY HH24:MI') as date,
			CASE 
				WHEN i.status = 'validated' OR i.status = 'collected' THEN 'valide'
				ELSE 'en_attente'
			END as statut,
			EXISTS(SELECT 1 FROM locker_access la WHERE la.item_id = i.id) as code_envoye
		FROM item i
		JOIN users u ON i.owner_id = u.id
		ORDER BY i.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[AdminDepot])
}

func (r *Repository) GetByUserId(userId pgtype.Int8) ([]Item, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			i.id, i.owner_id, i.locker_id, i.site_id, 
			COALESCE(i.material_type, '') as material_type,
			COALESCE(i.physical_state::text, '') as physical_state, 
			COALESCE(i.size::text, 'M') as size, 
			COALESCE(i.status::text, 'deposited') as status, 
			i.weight, 
			COALESCE(i.name, '') as name, 
			COALESCE(i.description, '') as description, 
			i.created_at, i.entry_id,
			COALESCE(s.type_site, '') as site_type,
			COALESCE(TO_CHAR(e.schedule, 'YYYY-MM-DD'), '') as schedule_date,
			COALESCE(TO_CHAR(e.start, 'HH24:MI'), '') as schedule_time
		FROM item i
		LEFT JOIN site s ON i.site_id = s.id
		LEFT JOIN entry e ON i.entry_id = e.id
		WHERE i.owner_id = $1
	`, userId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Item])
}

func (r *Repository) GetById(id pgtype.Int8) (*Item, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			i.id, i.owner_id, i.locker_id, i.site_id, 
			COALESCE(i.material_type, '') as material_type,
			COALESCE(i.physical_state::text, '') as physical_state, 
			COALESCE(i.size::text, 'M') as size, 
			COALESCE(i.status::text, 'deposited') as status, 
			i.weight, 
			COALESCE(i.name, '') as name, 
			COALESCE(i.description, '') as description, 
			i.created_at, i.entry_id,
			COALESCE(s.type_site, '') as site_type,
			COALESCE(TO_CHAR(e.schedule, 'YYYY-MM-DD'), '') as schedule_date,
			COALESCE(TO_CHAR(e.start, 'HH24:MI'), '') as schedule_time
		FROM item i
		LEFT JOIN site s ON i.site_id = s.id
		LEFT JOIN entry e ON i.entry_id = e.id
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
	size := item.Size
	if size == "" {
		size = SizeM
	}
	var id int64
	err := r.db.QueryRow(db.Ctx,
		"INSERT INTO item (owner_id, locker_id, site_id, material_type, physical_state, size, status, weight, name, description, entry_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
		item.OwnerId, item.LockerId, item.SiteId, item.MaterialType, item.PhysicalState, size, item.Status, item.Weight, item.Name, item.Description, item.EntryId).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Update(id pgtype.Int8, item Item) error {
	size := item.Size
	if size == "" {
		size = SizeM
	}
	tag, err := r.db.Exec(db.Ctx,
		"UPDATE item SET owner_id=$1, locker_id=$2, site_id=$3, material_type=$4, physical_state=$5, size=$6, status=$7, weight=$8, name=$9, description=$10, entry_id=$11 WHERE id=$12",
		item.OwnerId, item.LockerId, item.SiteId, item.MaterialType, item.PhysicalState, size, item.Status, item.Weight, item.Name, item.Description, item.EntryId, id)
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

func (r *Repository) UpdateLockerStatus(lockerId pgtype.Int8, status string) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE locker SET status = $1 WHERE id = $2", status, lockerId)
	return err
}

func (r *Repository) CreateAccessCode(lockerId pgtype.Int8, itemId pgtype.Int8, userId pgtype.Int8, code string) error {
	_, err := r.db.Exec(db.Ctx,
		"INSERT INTO locker_access (locker_id, item_id, user_id, access_code, expires_at) VALUES ($1, $2, $3, $4, NOW() + INTERVAL '24 hours')",
		lockerId, itemId, userId, code)
	return err
}

func (r *Repository) GetAccessCode(itemId pgtype.Int8) (string, error) {
	var code string
	err := r.db.QueryRow(db.Ctx, "SELECT access_code FROM locker_access WHERE item_id = $1 AND used_at IS NULL ORDER BY created_at DESC LIMIT 1", itemId).Scan(&code)
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
	tag, err := r.db.Exec(db.Ctx, "UPDATE locker_access SET used_at = NOW() WHERE access_code = $1", code)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid or already used code")
	}
	return nil
}

func (r *Repository) CancelDepot(id pgtype.Int8) error {
	var entryId pgtype.Int8
	err := r.db.QueryRow(db.Ctx, "SELECT entry_id FROM item WHERE id = $1", id).Scan(&entryId)
	if err != nil {
		return err
	}

	if entryId.Valid {
		_, err = r.db.Exec(db.Ctx, "DELETE FROM entry WHERE id = $1", entryId)
		if err != nil {
			return err
		}
	}

	tag, err := r.db.Exec(db.Ctx, "UPDATE item SET status = 'cancelled' WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("item not found")
	}
	return nil
}
